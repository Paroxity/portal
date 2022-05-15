package session

import (
	"github.com/paroxity/portal/server"
)

// LoadBalancer represents a load balancer which helps balance the load of players on the proxy.
type LoadBalancer interface {
	// FindServer finds a server for the session to connect to when they first join. If nil is returned, the
	// player is kicked from the proxy.
	FindServer(session *Session) *server.Server
}

// SplitLoadBalancer attempts to split players evenly across all the servers.
type SplitLoadBalancer struct {
	registry *server.Registry
}

// NewSplitLoadBalancer creates a "split" load balancer with the provided server registry.
func NewSplitLoadBalancer(registry *server.Registry) *SplitLoadBalancer {
	return &SplitLoadBalancer{registry: registry}
}

// FindServer ...
func (b *SplitLoadBalancer) FindServer(*Session) (srv *server.Server) {
	for _, s := range b.registry.Servers() {
		if srv == nil || srv.PlayerCount() > s.PlayerCount() {
			srv = s
		}
	}
	return srv
}
