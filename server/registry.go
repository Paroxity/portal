package server

import (
	"strings"
	"sync"
)

// Registry represents a register which stores the servers available on the proxy.
type Registry interface {
	// Server attempts to find a server from its name, and returns the server and if it was found or not.
	Server(name string) (*Server, bool)
	// Servers returns a slice of all of the available servers on the proxy.
	Servers() []*Server
	// AddServer attempts to add a server to the register.
	AddServer(srv *Server)
	// RemoveServer attempts to remove a server from the register.
	RemoveServer(srv *Server)
}

// DefaultRegistry represents a server registry with basic behaviour.
type DefaultRegistry struct {
	servers map[string]*Server
	mu      sync.Mutex
}

// NewDefaultRegistry creates a new DefaultRegistry and returns it.
func NewDefaultRegistry() *DefaultRegistry {
	return &DefaultRegistry{servers: make(map[string]*Server)}
}

// Server ...
func (r *DefaultRegistry) Server(name string) (*Server, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	srv, ok := r.servers[strings.ToLower(name)]
	return srv, ok
}

// Servers ...
func (r *DefaultRegistry) Servers() (all []*Server) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, srv := range r.servers {
		all = append(all, srv)
	}
	return
}

// AddServer ...
func (r *DefaultRegistry) AddServer(srv *Server) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.servers[strings.ToLower(srv.Name())] = srv
}

// RemoveServer ...
func (r *DefaultRegistry) RemoveServer(srv *Server) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.servers, strings.ToLower(srv.Name()))
}
