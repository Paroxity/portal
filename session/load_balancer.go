package session

import (
	"errors"
	"github.com/paroxity/portal/server"
	"strings"
)

var (
	loadBalancers map[string]func(session *Session) *server.Server
	loadBalancer  string
)

// LoadBalancer returns the selected load balancer method.
func LoadBalancer() func(session *Session) *server.Server {
	return loadBalancers[loadBalancer]
}

// RegisterLoadBalancer registers a load balancer method to be used when a session joins.
func RegisterLoadBalancer(name string, f func(session *Session) *server.Server) {
	loadBalancers[strings.ToLower(name)] = f
}

// SetLoadBalancer sets the default load balancer to be used when a player joins the proxy. If the provided
// load balancer does not exist, an error is returned.
func SetLoadBalancer(name string) error {
	if _, ok := loadBalancers[strings.ToLower(name)]; !ok {
		return errors.New("load balancer not registered")
	}
	loadBalancer = strings.ToLower(name)
	return nil
}

func init() {
	RegisterLoadBalancer("split", func(session *Session) *server.Server {
		var srv *server.Server
		for _, s := range server.DefaultGroup().Servers() {
			if srv == nil || srv.PlayerCount() > s.PlayerCount() {
				srv = s
			}
		}
		return srv
	})
	if err := SetLoadBalancer("split"); err != nil {
		panic(err)
	}
}
