package socket

import (
	"github.com/paroxity/portal/internal"
	"github.com/paroxity/portal/server"
	"github.com/paroxity/portal/session"
	"github.com/paroxity/portal/socket/packet"
	"net"
	"strings"
	"sync"
)

type Server interface {
	// Listen starts listening for connections on an address.
	Listen() error

	// Logger returns the logger attached to the socket server.
	Logger() internal.Logger

	// Secret returns the secret required for connections to authenticate.
	Secret() string

	// Clients returns all the clients that are connected to the socket server.
	Clients() []*Client
	// Client attempts to return a client from the provided name, case-sensitive.
	Client(name string) (*Client, bool)

	// SessionStore returns the store used to hold the open sessions on the proxy.
	SessionStore() *session.Store
	// ServerRegistry returns the registry used to store available servers on the proxy.
	ServerRegistry() *server.Registry
}

// DefaultServer represents a basic TCP socket server implementation. It allows external connections to
// connect and authenticate to be able to communicate with the proxy.
type DefaultServer struct {
	log internal.Logger

	addr   string
	secret string

	listener  net.Listener
	clientsMu sync.RWMutex
	clients   map[string]*Client

	sessionStore   *session.Store
	serverRegistry *server.Registry
}

// NewDefaultServer creates a new default server to be used for accepting socket connections.
func NewDefaultServer(addr, secret string, sessionStore *session.Store, serverRegistry *server.Registry, log internal.Logger) *DefaultServer {
	return &DefaultServer{
		log: log,

		addr:   addr,
		secret: secret,

		clients: make(map[string]*Client),

		sessionStore:   sessionStore,
		serverRegistry: serverRegistry,
	}
}

// Listen ...
func (s *DefaultServer) Listen() error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	s.log.Infof("socket server listening on %s\n", s.addr)
	s.listener = listener

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				s.log.Infof("socket server unable to accept connection: %v", err)
				continue
			}
			s.log.Debugf("socket server accepted a new connection")

			go s.handleClient(NewClient(conn, s.log))
		}
	}()
	return nil
}

// handleClient handles a client that has been accepted from the server.
func (s *DefaultServer) handleClient(c *Client) {
	defer s.handleClientDisconnect(c)
	s.clientsMu.Lock()
	s.clients[c.Name()] = c
	s.clientsMu.Unlock()

	for {
		pk, err := c.ReadPacket()
		if err != nil {
			if containsAny(err.Error(), "EOF", "closed") {
				return
			}
			s.log.Errorf("socket server unable to read packet: %v", err)
			continue
		}

		h, ok := handlers[pk.ID()]
		if ok {
			if !c.Authenticated() && h.RequiresAuth() {
				_ = c.WritePacket(&packet.AuthResponse{Status: packet.AuthResponseUnauthenticated})
				s.log.Debugf("received packet %T from unauthenticated client", pk)
				return
			}
			if err := h.Handle(pk, s, c); err != nil {
				s.log.Errorf("socket server unable to handle packet: %v", err)
			}
		} else {
			if c.name == "" {
				s.log.Debugf("unhandled packet %T from unauthenticated socket connection", pk)
			} else {
				s.log.Debugf("unhandled packet %T from %s socket connection", pk, c.name)
			}
		}
	}
}

func (s *DefaultServer) handleClientDisconnect(c *Client) {
	s.clientsMu.Lock()
	defer s.clientsMu.Unlock()
	delete(s.clients, c.Name())
	s.log.Debugf("socket connection \"%s\" closed", c.name)
	srv, ok := s.serverRegistry.Server(c.Name())
	if ok {
		s.serverRegistry.RemoveServer(srv)
		s.log.Debugf("removed server for socket connection \"%s\"", c.Name())
	}
}

// Logger ...
func (s *DefaultServer) Logger() internal.Logger {
	return s.log
}

// Secret ...
func (s *DefaultServer) Secret() string {
	return s.secret
}

// Clients ...
func (s *DefaultServer) Clients() (clients []*Client) {
	s.clientsMu.RLock()
	defer s.clientsMu.RUnlock()
	for _, client := range s.clients {
		clients = append(clients, client)
	}
	return
}

// Client ...
func (s *DefaultServer) Client(name string) (*Client, bool) {
	s.clientsMu.RLock()
	defer s.clientsMu.RUnlock()
	client, ok := s.clients[name]
	return client, ok
}

// SessionStore ...
func (s *DefaultServer) SessionStore() *session.Store {
	return s.sessionStore
}

// ServerRegistry ...
func (s *DefaultServer) ServerRegistry() *server.Registry {
	return s.serverRegistry
}

// containsAny checks if the string contains any of the provided sub strings.
func containsAny(s string, subs ...string) bool {
	for _, sub := range subs {
		if strings.Contains(s, sub) {
			return true
		}
	}

	return false
}
