package server

import (
	"github.com/paroxity/portal/socket/packet"
	"io"
)

// Client represents a client connected over the TCP socket system.
type Client interface {
	io.Closer

	// Name returns the name of the Client. This must be unique for all clients as it is used for
	// identification by the proxy and other clients.
	Name() string
	// ReadPacket reads a packet from the Client and returns it. It also returns an error in case one
	// occurred whilst reading the packet.
	ReadPacket() (packet.Packet, error)
	// ReadPacket writes a packet to the Client. It returns an error in case one occurred whilst writing the
	// packet.
	WritePacket(pk packet.Packet) error
}
