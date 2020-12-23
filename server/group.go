package server

import (
	"strings"
	"sync"
)

type Group struct {
	name string

	servers   map[string]*Server
	serversMu sync.Mutex
}

func NewGroup(name string) *Group {
	g := &Group{
		name: name,

		servers: make(map[string]*Server),
	}
	return g
}

func (g *Group) Name() string {
	return g.name
}

func (g *Group) AddServer(s *Server) {
	g.serversMu.Lock()
	g.servers[strings.ToLower(s.Name())] = s
	g.serversMu.Unlock()
}

func (g *Group) RemoveServer(name string) {
	g.serversMu.Lock()
	delete(g.servers, strings.ToLower(name))
	g.serversMu.Unlock()
}

func (g *Group) Server(name string) (*Server, bool) {
	g.serversMu.Lock()
	defer g.serversMu.Unlock()
	s, ok := g.servers[strings.ToLower(name)]
	return s, ok
}

func (g *Group) Servers() map[string]*Server {
	return g.servers
}
