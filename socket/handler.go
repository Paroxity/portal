package socket

import (
	"github.com/paroxity/portal/socket/packet"
)

// PacketHandler represents a type which handles a specific packet coming from a client.
type PacketHandler interface {
	// Handle is responsible for handling an incoming packet for the client.
	Handle(p packet.Packet, src Server, c *Client) error
	// RequiresAuth returns if the client must be authenticated in order for the handler to be triggered.
	RequiresAuth() bool
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
	RegisterHandler(packet.IDPlayerInfoRequest, &PlayerInfoRequestHandler{})
	RegisterHandler(packet.IDServerListRequest, &ServerListRequestHandler{})
	RegisterHandler(packet.IDFindPlayerRequest, &FindPlayerRequestHandler{})
}

// requireAuth implements the RequiresAuth() method and always returns true.
type requireAuth struct{}

// RequiresAuth ...
func (*requireAuth) RequiresAuth() bool {
	return true
}
