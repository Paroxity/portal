package server

import (
	"go.uber.org/atomic"
)

// Server represents a server connected to the proxy which players can join and play on.
type Server struct {
	name      string
	address   string
	useRakNet bool

	playerCount atomic.Int64
}

// New creates a new Server with the provided name, group and address as well as if the connection should use the RakNet
// protocol or not.
func New(name, address string, useRakNet bool) *Server {
	s := &Server{
		name:      name,
		address:   address,
		useRakNet: useRakNet,
	}

	return s
}

// Name returns the name the server was registered with.
func (s *Server) Name() string {
	return s.name
}

// Address returns the IP address the server was registered with. This should also contain the port separated
// by a colon. E.g. "127.0.0.1:19132".
func (s *Server) Address() string {
	return s.address
}

// UseRakNet returns if a connection should use the RakNet protocol when connecting to the server.
func (s *Server) UseRakNet() bool {
	return s.useRakNet
}

// IncrementPlayerCount increments the player count of the server.
func (s *Server) IncrementPlayerCount() {
	s.playerCount.Add(1)
}

// DecrementPlayerCount decreases the player count of the server.
func (s *Server) DecrementPlayerCount() {
	s.playerCount.Sub(1)
}

// PlayerCount returns the player count of the server controlled by the IncrementPlayerCount and
// DecrementPlayerCount functions above.
func (s *Server) PlayerCount() int {
	return int(s.playerCount.Load())
}
