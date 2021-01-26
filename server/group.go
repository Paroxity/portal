package server

import (
	"strings"
	"sync"
)

// Group represents a group of servers on the proxy. Groups are used to keep servers organized and in one
// place on the proxy, making it easier for servers to separate different parts of their network.
type Group struct {
	name string

	servers   map[string]*Server
	serversMu sync.Mutex
}

// NewGroup creates a new group with the given name and returns it.
func NewGroup(name string) *Group {
	g := &Group{
		name: name,

		servers: make(map[string]*Server),
	}
	return g
}

// Name returns the name the group was registered with.
func (g *Group) Name() string {
	return g.name
}

// AddServer adds the provided server to the group, overriding existing ones if the name conflicts.
func (g *Group) AddServer(s *Server) {
	g.serversMu.Lock()
	g.servers[strings.ToLower(s.Name())] = s
	g.serversMu.Unlock()
}

// RemoveServer removes the provided server name from the group.
func (g *Group) RemoveServer(name string) {
	g.serversMu.Lock()
	delete(g.servers, strings.ToLower(name))
	g.serversMu.Unlock()
}

// Server attempts to find the server by the name provided. It returns the Server and a bool to say if the
// Server was found or not.
func (g *Group) Server(name string) (*Server, bool) {
	g.serversMu.Lock()
	defer g.serversMu.Unlock()
	s, ok := g.servers[strings.ToLower(name)]
	return s, ok
}

// Servers returns all of the servers connected to the group.
func (g *Group) Servers() map[string]*Server {
	return g.servers
}
