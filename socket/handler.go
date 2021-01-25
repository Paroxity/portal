package socket

import (
	"github.com/paroxity/portal/socket/packet"
)

// PacketHandler represents a type which handles a specific packet coming from a client.
type PacketHandler interface {
	// Handle is responsible for handling an incoming packet for the client.
	Handle(p packet.Packet, c *Client) error
}

var handlers = make(map[uint16]PacketHandler)

// RegisterHandler registers a PacketHandler for the provided packet ID. Handlers do not stack, meaning
// registering multiple handlers for the same id will override the previous one.
func RegisterHandler(id uint16, h PacketHandler) {
	handlers[id] = h
}

func init() {
	RegisterHandler(packet.IDAuthRequest, &AuthRequestHandler{})
	RegisterHandler(packet.IDTransferRequest, &TransferRequestHandler{})
	RegisterHandler(packet.IDPlayerInfoRequest, &PlayerInfoRequest{})
}
