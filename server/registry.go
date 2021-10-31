package server

import (
	"strings"
	"sync"
)

// Registry represents a registry which stores the severs registered on the proxy.
type Registry struct {
	mu      sync.Mutex
	servers map[string]*Server
}

// NewDefaultRegistry creates a new Registry and returns it.
func NewDefaultRegistry() *Registry {
	return &Registry{servers: make(map[string]*Server)}
}

// Server attempts to find a server from its name, and returns the server and if it was found or not.
func (r *Registry) Server(name string) (*Server, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	srv, ok := r.servers[strings.ToLower(name)]
	return srv, ok
}

// Servers returns a slice of all the available servers on the proxy.
func (r *Registry) Servers() (all []*Server) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, srv := range r.servers {
		all = append(all, srv)
	}
	return
}

// AddServer adds a server to the register.
func (r *Registry) AddServer(srv *Server) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.servers[strings.ToLower(srv.Name())] = srv
}

// RemoveServer removes a server from the register.
func (r *Registry) RemoveServer(srv *Server) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.servers, strings.ToLower(srv.Name()))
}
