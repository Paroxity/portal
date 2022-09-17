package session

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"github.com/klauspost/compress/snappy"
	"github.com/paroxity/portal/session/tcpprotocol"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/login"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"go.uber.org/atomic"
	"io"
	"net"
	"sync"
	"time"
)

// ServerConn represents a connection that can be used between the proxy and a server to communicate on behalf of a client.
type ServerConn interface {
	io.Closer
	// GameData returns specific game data set to the connection for the player to be initialised with. This data is
	// obtained from the server during the login process.
	GameData() minecraft.GameData
	// DoSpawnTimeout starts the game for the client in the server with a timeout after which an error is returned if the
	// client has not yet spawned by that time. DoSpawnTimeout will start the spawning sequence using the game data found
	// in conn.GameData(), which was sent earlier by the server.
	DoSpawnTimeout(timeout time.Duration) error
	// ReadPacket reads a packet from the Conn, depending on the packet ID that is found in front of the packet data. If
	// a read deadline is set, an error is returned if the deadline is reached before any packet is received. ReadPacket
	// must not be called on multiple goroutines simultaneously. If the packet read was not implemented, a *packet.Unknown
	// is returned, containing the raw payload of the packet read.
	ReadPacket() (packet.Packet, error)
	// WritePacket encodes the packet passed and writes it to the Conn. The encoded data is buffered until the next 20th
	// of a second, after which the data is flushed and sent over the connection.
	WritePacket(packet.Packet) error
}

// TCPConn represents a player's connection to a server that is using the TCP protocol instead of RakNet.
type TCPConn struct {
	conn net.Conn
	pool packet.Pool

	identityData login.IdentityData
	clientData   login.ClientData
	gameData     minecraft.GameData

	sendMu   sync.Mutex
	hdr      *packet.Header
	buf      *bytes.Buffer
	shieldID atomic.Int32

	spawn chan struct{}
}

// NewTCPConn attempts to create a new TCP-based connection to the provided server using the provided client data. If
// successful the connection will be returned, otherwise an error will be returned instead.
func NewTCPConn(address, playerAddress string, identityData login.IdentityData, clientData login.ClientData) (*TCPConn, error) {
	conn := &TCPConn{
		identityData: identityData,
		clientData:   clientData,
		pool:         packet.NewPool(),
		buf:          bytes.NewBuffer(make([]byte, 0, 4096)),
		hdr:          &packet.Header{},
		spawn:        make(chan struct{}, 1),
	}
	err := conn.dial(address, playerAddress)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// dial attempts to dial a connection to the provided address for the player. An error is returned if it failed to dial.
func (conn *TCPConn) dial(address, playerAddress string) error {
	tcpConn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	conn.conn = tcpConn
	err = conn.WritePacket(&tcpprotocol.PlayerIdentity{
		IdentityData: conn.identityData,
		ClientData:   conn.clientData,
		Address:      playerAddress,
	})
	if err != nil {
		return err
	}
	pk, err := conn.ReadPacket()
	if err != nil {
		return err
	}
	startGame, ok := pk.(*packet.StartGame)
	if !ok {
		return fmt.Errorf("expected start game packet, got %T (%d)", pk, pk.ID())
	}
	conn.gameData = minecraft.GameData{
		Difficulty:                   startGame.Difficulty,
		WorldName:                    startGame.WorldName,
		EntityUniqueID:               startGame.EntityUniqueID,
		EntityRuntimeID:              startGame.EntityRuntimeID,
		PlayerGameMode:               startGame.PlayerGameMode,
		BaseGameVersion:              startGame.BaseGameVersion,
		PlayerPosition:               startGame.PlayerPosition,
		Pitch:                        startGame.Pitch,
		Yaw:                          startGame.Yaw,
		Dimension:                    startGame.Dimension,
		WorldSpawn:                   startGame.WorldSpawn,
		EditorWorld:                  startGame.EditorWorld,
		GameRules:                    startGame.GameRules,
		Time:                         startGame.Time,
		ServerBlockStateChecksum:     startGame.ServerBlockStateChecksum,
		CustomBlocks:                 startGame.Blocks,
		Items:                        startGame.Items,
		PlayerMovementSettings:       startGame.PlayerMovementSettings,
		WorldGameMode:                startGame.WorldGameMode,
		ServerAuthoritativeInventory: startGame.ServerAuthoritativeInventory,
		Experiments:                  startGame.Experiments,
	}
	return nil
}

// Close ...
func (conn *TCPConn) Close() error {
	return conn.conn.Close()
}

// GameData ...
func (conn *TCPConn) GameData() minecraft.GameData {
	return conn.gameData
}

// DoSpawnTimeout ...
func (conn *TCPConn) DoSpawnTimeout(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	select {
	case <-ctx.Done():
		return fmt.Errorf("spawn timeout")
	case <-conn.spawn:
		return nil
	}
}

// ReadPacket ...
func (conn *TCPConn) ReadPacket() (pk packet.Packet, err error) {
	var l uint32
	if err := binary.Read(conn.conn, binary.LittleEndian, &l); err != nil {
		return nil, err
	}

	data := make([]byte, l)
	read, err := conn.conn.Read(data)
	if err != nil {
		return nil, err
	}
	if read != int(l) {
		return nil, fmt.Errorf("expected %v bytes, got %v", l, read)
	}

	decoded, err := snappy.Decode(nil, data)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(decoded)
	header := &packet.Header{}
	if err := header.Read(buf); err != nil {
		return nil, err
	}

	pkFunc, ok := conn.pool[header.PacketID]
	if !ok {
		return nil, fmt.Errorf("unknown packet %v", header.PacketID)
	}

	defer func() {
		if recoveredErr := recover(); recoveredErr != nil {
			err = fmt.Errorf("%T: %w", pk, recoveredErr.(error))
		}
	}()
	pk = pkFunc()
	pk.Unmarshal(protocol.NewReader(buf, 0))
	if buf.Len() > 0 {
		return nil, fmt.Errorf("still have %v bytes unread", buf.Len())
	}

	if _, ok := pk.(*packet.StartGame); ok {
		close(conn.spawn)
	}

	return pk, nil
}

// WritePacket ...
func (conn *TCPConn) WritePacket(pk packet.Packet) error {
	conn.sendMu.Lock()
	conn.hdr.PacketID = pk.ID()
	_ = conn.hdr.Write(conn.buf)

	pk.Marshal(protocol.NewWriter(conn.buf, conn.shieldID.Load()))

	data := conn.buf.Bytes()
	conn.buf.Reset()
	conn.sendMu.Unlock()

	encoded := snappy.Encode(nil, data)

	buf := bytes.NewBuffer(make([]byte, 0, 4+len(encoded)))

	if err := binary.Write(buf, binary.LittleEndian, int32(len(encoded))); err != nil {
		return err
	}
	if _, err := buf.Write(encoded); err != nil {
		return err
	}

	if _, err := conn.conn.Write(buf.Bytes()); err != nil {
		return err
	}

	return nil
}
