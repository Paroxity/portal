package session

import (
	"errors"
	"github.com/google/uuid"
	"github.com/paroxity/portal/event"
	"github.com/paroxity/portal/internal"
	"github.com/paroxity/portal/server"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"github.com/scylladb/go-set/b16set"
	"github.com/scylladb/go-set/i32set"
	"github.com/scylladb/go-set/i64set"
	"github.com/scylladb/go-set/strset"
	"go.uber.org/atomic"
	"sync"
	"time"
)

var (
	emptyChunkData = make([]byte, 257)
)

// Session stores the data for an active session on the proxy.
type Session struct {
	*translator

	log   internal.Logger
	conn  *minecraft.Conn
	store *Store

	hMutex sync.RWMutex
	// h holds the current handler of the session.
	h Handler

	serverMu       sync.RWMutex
	server         *server.Server
	serverConn     *minecraft.Conn
	tempServerConn *minecraft.Conn

	entities    *i64set.Set
	playerList  *b16set.Set
	effects     *i32set.Set
	bossBars    *i64set.Set
	scoreboards *strset.Set

	uuid uuid.UUID

	transferring atomic.Bool
	postTransfer atomic.Bool
	once         sync.Once
}

// New creates a new Session with the provided connection.
func New(conn *minecraft.Conn, store *Store, loadBalancer LoadBalancer, log internal.Logger) (_ *Session, err error) {
	s := &Session{
		log:   log,
		conn:  conn,
		store: store,

		entities:    i64set.New(),
		playerList:  b16set.New(),
		effects:     i32set.New(),
		bossBars:    i64set.New(),
		scoreboards: strset.New(),

		h:    NopHandler{},
		uuid: uuid.MustParse(conn.IdentityData().Identity),
	}

	store.Store(s)
	defer func() {
		if err != nil {
			store.Delete(s.UUID())
		}
	}()

	srv := loadBalancer.FindServer(s)
	if srv == nil {
		return nil, errors.New("load balancer did not return a server for the player to join")
	}
	s.server = srv

	srvConn, err := s.dial(srv)
	if err != nil {
		return nil, err
	}

	s.serverConn = srvConn
	if err = s.login(); err != nil {
		_ = srvConn.Close()

		return nil, err
	}

	s.translator = newTranslator(conn.GameData())

	handlePackets(s)
	srv.IncrementPlayerCount()
	return s, nil
}

// dial dials a new connection to the provided server. It then returns the connection between the proxy and
// that server, along with any error that may have occurred.
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

// Handle sets the handler for the current session which can be used to handle different events from the
// session. If the handler is nil, a NopHandler is used instead.
func (s *Session) Handle(h Handler) {
	s.hMutex.Lock()
	defer s.hMutex.Unlock()

	if h == nil {
		h = NopHandler{}
	}
	s.h = h
}

// Transfer transfers the session to the provided server, returning any error that may have occurred during
// the initial transfer.
func (s *Session) Transfer(srv *server.Server) (err error) {
	if !s.transferring.CAS(false, true) {
		return errors.New("already being transferred")
	}

	s.log.Infof("%s is being transferred from %s to %s", s.conn.IdentityData().DisplayName, s.Server().Name(), srv.Name())

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
		_ = s.conn.WritePacket(&packet.StopSound{StopAll: true})

		chunkX := int32(pos.X()) >> 4
		chunkZ := int32(pos.Z()) >> 4
		for x := int32(-1); x <= 1; x++ {
			for z := int32(-1); z <= 1; z++ {
				_ = s.conn.WritePacket(&packet.LevelChunk{
					Position:      protocol.ChunkPos{chunkX + x, chunkZ + z},
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

// handler() returns the handler connected to the session.
func (s *Session) handler() Handler {
	s.hMutex.RLock()
	handler := s.h
	s.hMutex.RUnlock()
	return handler
}

// Close closes the session and any linked connections/counters.
func (s *Session) Close() {
	s.once.Do(func() {
		s.handler().HandleQuit()
		s.Handle(NopHandler{})

		s.store.Delete(s.UUID())

		_ = s.conn.Close()
		_ = s.ServerConn().Close()
		if s.tempServerConn != nil {
			_ = s.tempServerConn.Close()
		}

		s.Server().DecrementPlayerCount()
	})
}

// clearEntities flushes the entities map and despawns the entities for the client.
func (s *Session) clearEntities() {
	s.entities.Each(func(id int64) bool {
		_ = s.conn.WritePacket(&packet.RemoveActor{EntityUniqueID: id})
		return true
	})

	s.entities.Clear()
}

// clearPlayerList flushes the playerList map and removes all the entries for the client.
func (s *Session) clearPlayerList() {
	var entries = make([]protocol.PlayerListEntry, s.playerList.Size())
	s.playerList.Each(func(uid [16]byte) bool {
		entries = append(entries, protocol.PlayerListEntry{UUID: uid})
		return true
	})

	_ = s.conn.WritePacket(&packet.PlayerList{ActionType: packet.PlayerListActionRemove, Entries: entries})

	s.playerList.Clear()
}

// clearEffects flushes the effects map and removes all the effects for the client.
func (s *Session) clearEffects() {
	s.effects.Each(func(i int32) bool {
		_ = s.conn.WritePacket(&packet.MobEffect{
			EntityRuntimeID: s.originalRuntimeID,
			Operation:       packet.MobEffectRemove,
			EffectType:      i,
		})
		return true
	})

	s.effects.Clear()
}

// clearBossBars clears all the boss bars currently visible the client.
func (s *Session) clearBossBars() {
	s.bossBars.Each(func(b int64) bool {
		_ = s.conn.WritePacket(&packet.BossEvent{
			BossEntityUniqueID: b,
			EventType:          packet.BossEventHide,
		})
		return true
	})

	s.bossBars.Clear()
}

// clearScoreboard clears the current scoreboard visible by the client.
func (s *Session) clearScoreboard() {
	s.scoreboards.Each(func(sb string) bool {
		_ = s.conn.WritePacket(&packet.RemoveObjective{ObjectiveName: sb})
		return true
	})

	s.scoreboards.Clear()
}
