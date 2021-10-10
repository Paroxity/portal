package server

import (
	"github.com/paroxity/portal/socket/packet"
)

// Client represents a client connected over the TCP socket system.
type Client interface {
	// Name returns the name of the Client. This must be unique for all clients as it is used for
	// identification by the proxy and other clients.
	Name() string
	// ReadPacket reads a packet from the Client and returns it. It also returns an error in case one
	// occurred whilst reading the packet.
	ReadPacket() (packet.Packet, error)
	// WritePacket writes a packet to the Client. It returns an error in case one occurred whilst writing the
	// packet.
	WritePacket(pk packet.Packet) error
	// Close is called when the client is disconnected from the socket server. The registry is provided in
	// case it needs to remove any servers etc.
	Close(registry *Registry) error
}
