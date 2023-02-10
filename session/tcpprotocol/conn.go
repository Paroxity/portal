package tcpprotocol

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"github.com/google/uuid"
	"github.com/klauspost/compress/snappy"
	packet2 "github.com/paroxity/portal/session/tcpprotocol/packet"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/login"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"go.uber.org/atomic"
	"net"
	"sync"
	"time"
)

// Conn represents a player's connection to a server that is using the TCP protocol instead of RakNet.
type Conn struct {
	conn net.Conn
	pool packet.Pool

	identityData login.IdentityData
	clientData   login.ClientData
	gameData     minecraft.GameData

	enableClientCache bool

	sendMu   sync.Mutex
	hdr      *packet.Header
	buf      *bytes.Buffer
	shieldID atomic.Int32

	spawn chan struct{}
}

// newConn attempts to create a new TCP-based connection to the provided server using the provided client data. If
// successful the connection will be returned, otherwise an error will be returned instead.
func newConn(netConn net.Conn) *Conn {
	conn := &Conn{
		conn:  netConn,
		pool:  packet.NewPool(),
		buf:   bytes.NewBuffer(make([]byte, 0, 4096)),
		hdr:   &packet.Header{},
		spawn: make(chan struct{}, 1),
	}
	return conn
}

// login attempts to connect to a server by identifying the player and waiting to receive the StartGame packet. An error
// is returned if it failed to connect.
func (conn *Conn) login(playerAddress string) error {
	err := conn.WritePacket(&packet2.ConnectionRequest{
		ProtocolVersion: ProtocolVersion,
	})
	if err != nil {
		return err
	}

	connectionResponse, err := expect(conn, packet2.IDConnectionResponse)
	if err != nil {
		return err
	}
	switch connectionResponse.(*packet2.ConnectionResponse).Response {
	case packet2.ConnectionResponseUnsupportedProtocol:
		return fmt.Errorf("unsupported protocol version %d", ProtocolVersion)
	}

	err = conn.WritePacket(&packet2.PlayerIdentity{
		IdentityData:      conn.identityData,
		ClientData:        conn.clientData,
		EnableClientCache: conn.enableClientCache,
		Address:           playerAddress,
	})
	if err != nil {
		return err
	}

	startGame, err := expect(conn, packet.IDStartGame)
	if err != nil {
		return err
	}
	startGamePacket := startGame.(*packet.StartGame)
	conn.gameData = minecraft.GameData{
		Difficulty:                   startGamePacket.Difficulty,
		WorldName:                    startGamePacket.WorldName,
		EntityUniqueID:               startGamePacket.EntityUniqueID,
		EntityRuntimeID:              startGamePacket.EntityRuntimeID,
		PlayerGameMode:               startGamePacket.PlayerGameMode,
		BaseGameVersion:              startGamePacket.BaseGameVersion,
		PlayerPosition:               startGamePacket.PlayerPosition,
		Pitch:                        startGamePacket.Pitch,
		Yaw:                          startGamePacket.Yaw,
		Dimension:                    startGamePacket.Dimension,
		WorldSpawn:                   startGamePacket.WorldSpawn,
		EditorWorld:                  startGamePacket.EditorWorld,
		GameRules:                    startGamePacket.GameRules,
		Time:                         startGamePacket.Time,
		ServerBlockStateChecksum:     startGamePacket.ServerBlockStateChecksum,
		CustomBlocks:                 startGamePacket.Blocks,
		Items:                        startGamePacket.Items,
		PlayerMovementSettings:       startGamePacket.PlayerMovementSettings,
		WorldGameMode:                startGamePacket.WorldGameMode,
		ServerAuthoritativeInventory: startGamePacket.ServerAuthoritativeInventory,
		Experiments:                  startGamePacket.Experiments,
	}
	return nil
}

// identify attempts to identify the player attempting to connect to the server through the PlayerIdentity packet. An
// error is returned if it failed to identify.
func (conn *Conn) identify() error {
	connectionRequest, err := expect(conn, packet2.IDConnectionRequest)
	if err != nil {
		return err
	}

	response := packet2.ConnectionResponseSuccess
	if connectionRequest.(*packet2.ConnectionRequest).ProtocolVersion != ProtocolVersion {
		response = packet2.ConnectionResponseUnsupportedProtocol
	}
	err = conn.WritePacket(&packet2.ConnectionResponse{
		Response: response,
	})
	if err != nil {
		return err
	}

	playerIdentity, err := expect(conn, packet2.IDPlayerIdentity)
	if err != nil {
		return err
	}
	playerIdentityPacket := playerIdentity.(*packet2.PlayerIdentity)
	conn.identityData = playerIdentityPacket.IdentityData
	conn.clientData = playerIdentityPacket.ClientData
	conn.enableClientCache = playerIdentityPacket.EnableClientCache
	return nil
}

// Close ...
func (conn *Conn) Close() error {
	return conn.conn.Close()
}

// IdentityData ...
func (conn *Conn) IdentityData() login.IdentityData {
	return conn.identityData
}

// ClientData ...
func (conn *Conn) ClientData() login.ClientData {
	return conn.clientData
}

// ClientCacheEnabled ...
func (conn *Conn) ClientCacheEnabled() bool {
	return conn.enableClientCache
}

// ChunkRadius ...
func (conn *Conn) ChunkRadius() int {
	//TODO implement me
	return 8
}

// Latency ...
func (conn *Conn) Latency() time.Duration {
	//TODO implement me
	return time.Millisecond * 20
}

// Flush ...
func (conn *Conn) Flush() error {
	return nil
}

// RemoteAddr ...
func (conn *Conn) RemoteAddr() net.Addr {
	return conn.conn.RemoteAddr()
}

// GameData ...
func (conn *Conn) GameData() minecraft.GameData {
	return conn.gameData
}

// DoSpawnTimeout ...
func (conn *Conn) DoSpawnTimeout(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	select {
	case <-ctx.Done():
		return fmt.Errorf("spawn timeout")
	case <-conn.spawn:
		return nil
	}
}

// StartGameContext ...
func (conn *Conn) StartGameContext(_ context.Context, data minecraft.GameData) error {
	for _, item := range data.Items {
		if item.Name == "minecraft:shield" {
			conn.shieldID.Store(int32(item.RuntimeID))
			break
		}
	}
	return conn.WritePacket(&packet.StartGame{
		Difficulty:                   data.Difficulty,
		EntityUniqueID:               data.EntityUniqueID,
		EntityRuntimeID:              data.EntityRuntimeID,
		PlayerGameMode:               data.PlayerGameMode,
		PlayerPosition:               data.PlayerPosition,
		Pitch:                        data.Pitch,
		Yaw:                          data.Yaw,
		Dimension:                    data.Dimension,
		WorldSpawn:                   data.WorldSpawn,
		EditorWorld:                  data.EditorWorld,
		GameRules:                    data.GameRules,
		Time:                         data.Time,
		Blocks:                       data.CustomBlocks,
		Items:                        data.Items,
		AchievementsDisabled:         true,
		Generator:                    1,
		EducationFeaturesEnabled:     true,
		MultiPlayerGame:              true,
		MultiPlayerCorrelationID:     uuid.Must(uuid.NewRandom()).String(),
		CommandsEnabled:              true,
		WorldName:                    data.WorldName,
		LANBroadcastEnabled:          true,
		PlayerMovementSettings:       data.PlayerMovementSettings,
		WorldGameMode:                data.WorldGameMode,
		ServerAuthoritativeInventory: data.ServerAuthoritativeInventory,
		Experiments:                  data.Experiments,
		BaseGameVersion:              data.BaseGameVersion,
		GameVersion:                  protocol.CurrentVersion,
	})
}

// ReadPacket ...
func (conn *Conn) ReadPacket() (pk packet.Packet, err error) {
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
func (conn *Conn) WritePacket(pk packet.Packet) error {
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

func expect(conn *Conn, expectedID uint32) (packet.Packet, error) {
	pk, err := conn.ReadPacket()
	if err != nil {
		return pk, err
	}
	if pk.ID() != expectedID {
		return nil, fmt.Errorf("received unexpected packet %T", pk)
	}
	return pk, nil
}
