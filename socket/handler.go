package socket

import (
	portalpacket "github.com/paroxity/portal/socket/packet"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// PacketHandler represents a type which handles a specific packet coming from a client.
type PacketHandler interface {
	// Handle is responsible for handling an incoming packet for the client.
	Handle(p packet.Packet, c *Client) error
}

var handlers = make(map[uint32]PacketHandler)

// RegisterHandler registers a PacketHandler for the provided packet ID. Handlers do not stack, meaning
// registering multiple handlers for the same id will override the previous one.
func RegisterHandler(id uint32, h PacketHandler) {
	handlers[id] = h
}

func init() {
	RegisterHandler(portalpacket.IDAuthRequest, &AuthRequestHandler{})
	RegisterHandler(portalpacket.IDTransferRequest, &TransferRequestHandler{})
}
