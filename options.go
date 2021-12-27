package portal

import (
	"github.com/paroxity/portal/internal"
	"github.com/paroxity/portal/session"
	"github.com/sandertv/gophertunnel/minecraft"
)

// Options represents the options that control how the proxy should be set up. After the proxy has been
// instantiated, the options below are immutable unless instantiated again.
type Options struct {
	// Logger represents the logger that will be used for the lifetime of the proxy.
	Logger internal.Logger

	// Address is the address that the proxy should run on. It should be in the format of "address:port".
	Address string
	// ListenConfig contains settings that can be changed for the listener. It can be used to change the MOTD
	// and add resource packs etc.
	ListenConfig minecraft.ListenConfig

	// LoadBalancer is the method used to balance load across the servers on the proxy. It can be used to
	// change which servers players connect to when they join the proxy.
	LoadBalancer session.LoadBalancer

	// Whitelist is used to limit the proxy to only allow certain players to join.
	Whitelist session.Whitelist
}
