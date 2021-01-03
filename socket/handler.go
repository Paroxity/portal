package socket

import "github.com/sandertv/gophertunnel/minecraft/protocol/packet"

// packetHandler represents a type which handles a specific packet coming from a client.
type packetHandler interface {
	// Handle is responsible for handling an incoming packet for the client.
	Handle(p packet.Packet, c *Client) error
}
