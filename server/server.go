package server

import (
	"fmt"
	"go.uber.org/atomic"
)

type Server struct {
	name    string
	group   string
	address string

	playerCount atomic.Int64
}

func New(name, group, address string) (*Server, error) {
	s := &Server{
		name:    name,
		group:   group,
		address: address,
	}
	g, ok := GroupFromName(group)
	if !ok {
		return nil, fmt.Errorf("group %s not found", group)
	}
	g.AddServer(s)
	return s, nil
}

func (s *Server) Name() string {
	return s.name
}

func (s *Server) Address() string {
	return s.address
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
