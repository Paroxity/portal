package server

import (
	"go.uber.org/atomic"
)

type Server struct {
	name    string
	group   string
	address string

	connection Client

	playerCount atomic.Int64
}

func New(name, group, address string) *Server {
	s := &Server{
		name:    name,
		group:   group,
		address: address,
	}

	return s
}

func (s *Server) Name() string {
	return s.name
}

func (s *Server) Group() string {
	return s.group
}

func (s *Server) Address() string {
	return s.address
}

func (s *Server) Connected() bool {
	return s.connection != nil
}

func (s *Server) Conn() Client {
	return s.connection
}

func (s *Server) setConn(c Client) {
	s.connection = c
}

func (s *Server) IncrementPlayerCount() {
	s.playerCount.Add(1)
}

func (s *Server) DecrementPlayerCount() {
	s.playerCount.Sub(1)
}

func (s *Server) PlayerCount() int {
	return int(s.playerCount.Load())
}
