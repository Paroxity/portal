package session

import (
	"errors"
	"github.com/google/uuid"
	"github.com/paroxity/portal/event"
	"github.com/paroxity/portal/query"
	"github.com/paroxity/portal/server"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"github.com/sirupsen/logrus"
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
	*translator

	conn *minecraft.Conn

	hMutex sync.RWMutex
	// h holds the current handler of the session.
	h Handler

	serverMu       sync.RWMutex
	server         *server.Server
	serverConn     *minecraft.Conn
	tempServerConn *minecraft.Conn

	entityMu sync.Mutex
	entities map[int64]struct{}

	uuid uuid.UUID

	transferring atomic.Bool
	once         sync.Once
}

// All returns all of the connected sessions on the proxy.
func All() []*Session {
	var s []*Session
	sessions.Range(func(_, v interface{}) bool {
		s = append(s, v.(*Session))
		return true
	})
	return s
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
func New(conn *minecraft.Conn) (*Session, error) {
	s := &Session{
		translator: newTranslator(conn.GameData()),
		conn:       conn,

		entities: make(map[int64]struct{}),

		h:    NopHandler{},
		uuid: uuid.MustParse(conn.IdentityData().Identity),
	}

	srv := LoadBalancer()(s)
	s.server = srv

	srvConn, err := s.dial(srv)
	if err != nil {
		return nil, err
	}

	s.serverConn = srvConn
	if err := s.login(); err != nil {
		return nil, err
	}
	handlePackets(s)
	sessions.Store(s.UUID(), s)
	srv.IncrementPlayerCount()
	query.IncrementPlayerCount()
	return s, nil
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
func (s *Session) login() (err error) {
	var g sync.WaitGroup
	g.Add(2)
	go func() {
		err = s.conn.StartGameTimeout(s.ServerConn().GameData(), time.Minute)
		g.Done()
	}()
	go func() {
		err = s.ServerConn().DoSpawnTimeout(time.Minute)
		g.Done()
	}()
	g.Wait()
	return
}

// Conn returns the active connection for the session.
func (s *Session) Conn() *minecraft.Conn {
	return s.conn
}

// Server returns the server the session is currently connected to.
func (s *Session) Server() *server.Server {
	s.serverMu.RLock()
	defer s.serverMu.RUnlock()
	return s.server
}

// ServerConn returns the connection for the session's current server.
func (s *Session) ServerConn() *minecraft.Conn {
	s.serverMu.RLock()
	defer s.serverMu.RUnlock()
	return s.serverConn
}

// UUID returns the UUID from the session's connection.
func (s *Session) UUID() uuid.UUID {
	return s.uuid
}

func (s *Session) Handle(h Handler) {
	s.hMutex.Lock()
	defer s.hMutex.Unlock()

	if h == nil {
		h = NopHandler{}
	}
	s.h = h
}

func (s *Session) Transfer(srv *server.Server) (err error) {
	if !s.transferring.CAS(false, true) {
		return errors.New("already being transferred")
	}

	logrus.Infof("Transferring %s to server %s in group %s\n", s.conn.IdentityData().DisplayName, srv.Name(), srv.Group())

	ctx := event.C()
	s.handler().HandleTransfer(ctx, srv)

	ctx.Continue(func() {
		conn, err := s.dial(srv)
		if err != nil {
			return
		}
		if err = conn.DoSpawnTimeout(time.Minute); err != nil {
			return
		}

		s.serverMu.Lock()
		s.tempServerConn = conn
		s.serverMu.Unlock()

		pos := s.conn.GameData().PlayerPosition
		_ = s.conn.WritePacket(&packet.ChangeDimension{
			Dimension: packet.DimensionNether,
			Position:  pos,
		})

		// TODO: Clear inventory & scoreboard

		s.entityMu.Lock()
		e := s.entities
		s.entities = map[int64]struct{}{}
		s.entityMu.Unlock()

		for id := range e {
			_ = s.conn.WritePacket(&packet.RemoveActor{EntityUniqueID: id})
		}

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

		s.serverMu.Lock()
		s.server.DecrementPlayerCount()
		s.server = srv
		s.server.IncrementPlayerCount()
		s.serverMu.Unlock()
	})

	ctx.Stop(func() {
		s.setTransferring(false)
	})

	return
}

// Transferring returns if the session is currently transferring to a different server or not.
func (s *Session) Transferring() bool {
	return s.transferring.Load()
}

// setTransferring sets if the session is transferring to a different server.
func (s *Session) setTransferring(v bool) {
	s.transferring.Store(v)
}

func (s *Session) handler() Handler {
	s.hMutex.RLock()
	handler := s.h
	s.hMutex.RUnlock()
	return handler
}

func (s *Session) addLoadedEntity(eid int64) {
	s.entityMu.Lock()
	s.entities[eid] = struct{}{}
	s.entityMu.Unlock()
}

func (s *Session) removeLoadedEntity(eid int64) {
	s.entityMu.Lock()
	delete(s.entities, eid)
	s.entityMu.Unlock()
}

// Close closes the session and any linked connections/counters.
func (s *Session) Close() {
	s.once.Do(func() {
		_ = s.conn.Close()
		_ = s.ServerConn().Close()

		s.entityMu.Lock()
		s.entities = map[int64]struct{}{}
		s.entityMu.Unlock()

		s.server.DecrementPlayerCount()
		query.DecrementPlayerCount()
	})
}
