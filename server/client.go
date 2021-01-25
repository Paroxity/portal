package server

import (
	"github.com/paroxity/portal/socket/packet"
	"io"
)

// Client represents a client connected over the TCP socket system.
type Client interface {
	io.Closer

	Name() string
	ReadPacket() (packet.Packet, error)
	WritePacket(pk packet.Packet) error
}
