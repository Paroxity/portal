package session

import (
	"errors"
	"github.com/google/uuid"
	"github.com/paroxity/portal/query"
	"github.com/paroxity/portal/server"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"go.uber.org/atomic"
	"sync"
	"time"
)

var (
	emptyChunkData = make([]byte, 257)

	sessions sync.Map
)

// Session stores the data for an active session on the proxy.
type Session struct {
	conn       *minecraft.Conn
	translator *translator

	server         *server.Server
	serverConn     *minecraft.Conn
	tempServerConn *minecraft.Conn

	uuid uuid.UUID

	transferring atomic.Bool
	closed       atomic.Bool
}

// Lookup attempts to find a Session with the provided UUID.
func Lookup(v uuid.UUID) (*Session, bool) {
	s, ok := sessions.Load(v)
	if !ok {
		return nil, false
	}
	return s.(*Session), true
}

// New creates a new Session with the provided connection and target server.
func New(conn *minecraft.Conn) error {
	s := &Session{
		conn:       conn,
		translator: newTranslator(conn.GameData()),
		uuid:       uuid.MustParse(conn.IdentityData().Identity),
	}

	srv := LoadBalancer()(s)
	s.server = srv

	srvConn, err := s.dial(srv)
	if err != nil {
		return err
	}
	s.serverConn = srvConn
	if err := s.login(); err != nil {
		return err
	}
	handlePackets(s)
	sessions.Store(s.UUID(), s)
	srv.IncrementPlayerCount()
	query.IncrementPlayerCount()
	return nil
}

func (s *Session) dial(srv *server.Server) (*minecraft.Conn, error) {
	i := s.conn.IdentityData()
	i.XUID = ""
	return minecraft.Dialer{
		ClientData:   s.conn.ClientData(),
		IdentityData: i,
	}.Dial("raknet", srv.Address())
}

// login performs the initial login sequence for the session.
func (s *Session) login() error {
	var g sync.WaitGroup
	g.Add(2)
	var loginErr error = nil
	go func() {
		if err := s.conn.StartGameTimeout(s.serverConn.GameData(), time.Minute); err != nil {
			loginErr = err
		}
		g.Done()
	}()
	go func() {
		if err := s.serverConn.DoSpawnTimeout(time.Minute); err != nil {
			loginErr = err
		}
		g.Done()
	}()
	g.Wait()
	return loginErr
}

// Conn returns the active connection for the session.
func (s *Session) Conn() *minecraft.Conn {
	return s.conn
}

// ServerConn returns the connection for the session's current server.
func (s *Session) ServerConn() *minecraft.Conn {
	return s.serverConn
}

// UUID returns the UUID from the session's connection.
func (s *Session) UUID() uuid.UUID {
	return s.uuid
}

func (s *Session) Transfer(srv *server.Server) error {
	if !s.transferring.CAS(false, true) {
		return errors.New("already being transferred")
	}

	conn, err := s.dial(srv)
	if err != nil {
		return err
	}
	if err := conn.DoSpawnTimeout(time.Minute); err != nil {
		return err
	}
	s.tempServerConn = conn

	pos := s.conn.GameData().PlayerPosition
	_ = s.conn.WritePacket(&packet.ChangeDimension{
		Dimension: packet.DimensionNether,
		Position:  pos,
	})

	// TODO: Clear inventory, scoreboard & entities

	chunkX := int32(pos.X()) >> 4
	chunkZ := int32(pos.Z()) >> 4
	for x := int32(-1); x <= 1; x++ {
		for z := int32(-1); z <= 1; z++ {
			_ = s.conn.WritePacket(&packet.LevelChunk{
				ChunkX:        chunkX + x,
				ChunkZ:        chunkZ + z,
				SubChunkCount: 0,
				RawPayload:    emptyChunkData,
			})
		}
	}

	s.server.DecrementPlayerCount()
	s.server = srv
	s.server.IncrementPlayerCount()

	return nil
}

// Transferring returns if the session is currently transferring to a different server or not.
func (s *Session) Transferring() bool {
	return s.transferring.Load()
}

// SetTransferring sets if the session is transferring to a different server.
func (s *Session) SetTransferring(v bool) {
	s.transferring.Store(v)
}

// Close closes the session and any linked connections/counters.
func (s *Session) Close() {
	if !s.closed.CAS(false, true) {
		return
	}

	_ = s.conn.Close()
	_ = s.serverConn.Close()

	s.server.DecrementPlayerCount()
	query.DecrementPlayerCount()
}
