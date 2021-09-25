package portal

import (
	"fmt"
	"github.com/paroxity/portal/internal"
	"github.com/paroxity/portal/server"
	"github.com/paroxity/portal/session"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sirupsen/logrus"
)

// Portal represents the proxy and controls its functionality.
type Portal struct {
	log internal.Logger

	address      string
	listenConfig minecraft.ListenConfig
	listener     *minecraft.Listener

	sessionStore   session.Store
	serverRegistry server.Registry
	loadBalancer   session.LoadBalancer
}

// New instantiates portal using the provided options and returns it. If some options are not set, default
// values will be used in replacement.
func New(opts Options) *Portal {
	if opts.Logger == nil {
		opts.Logger = logrus.New()
	}
	if opts.SessionStore == nil {
		opts.SessionStore = session.NewDefaultStore()
	}
	if opts.ServerRegistry == nil {
		opts.ServerRegistry = server.NewDefaultRegistry()
	}
	if opts.LoadBalancer == nil {
		opts.LoadBalancer = session.NewSplitLoadBalancer(opts.ServerRegistry)
	}
	return &Portal{
		log: opts.Logger,

		address:      opts.Address,
		listenConfig: opts.ListenConfig,

		sessionStore:   opts.SessionStore,
		serverRegistry: opts.ServerRegistry,
		loadBalancer:   opts.LoadBalancer,
	}
}

// Logger returns the global logger used by the proxy.
func (p *Portal) Logger() internal.Logger {
	return p.log
}

// SessionStore returns the session store provided to portal. It is used to store all of the open sessions.
func (p *Portal) SessionStore() session.Store {
	return p.sessionStore
}

// ServerRegistry returns the server registry provided to portal. It is used to store the available servers.
func (p *Portal) ServerRegistry() server.Registry {
	return p.serverRegistry
}

// Listen starts to listen on the set address and allows connections from minecraft clients. An error is
// returned if the listener failed to listen.
func (p *Portal) Listen() error {
	l, err := p.listenConfig.Listen("raknet", p.address)
	if err != nil {
		return err
	}
	p.listener = l
	return nil
}

// Accept accepts a fully connected (on Minecraft layer) connection which is ready to receive and send
// packets. If the listener is closed or the player failed to spawn in then an error will be returned.
func (p *Portal) Accept() (*session.Session, error) {
	if p.listener == nil {
		return nil, fmt.Errorf("no active listener")
	}
	conn, err := p.listener.Accept()
	if err != nil {
		return nil, err
	}
	return session.New(conn.(*minecraft.Conn), p.sessionStore, p.loadBalancer, p.log)
}

// Disconnect disconnects a Minecraft Conn passed by first sending a disconnect with the message passed, and
// closing the connection after. If the message passed is empty, the client will be immediately sent to the
// player list instead of a disconnect screen.
func (p *Portal) Disconnect(conn *minecraft.Conn, message string) error {
	if p.listener == nil {
		return fmt.Errorf("no listener to disconnect connection")
	}
	return p.listener.Disconnect(conn, message)
}
