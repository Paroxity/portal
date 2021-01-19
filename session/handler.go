package session

import (
	"github.com/paroxity/portal/event"
	"github.com/paroxity/portal/server"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// Handler handles events that are called by a player's session.
type Handler interface {
	// HandleClientBoundPacket handles a packet that's sent by the session's connected server. ctx.Cancel()
	// may be called to cancel the packet.
	HandleClientBoundPacket(ctx *event.Context, pk packet.Packet)
	// HandleServerBoundPacket handles a packet that's sent by the session. ctx.Cancel() may be called to
	// cancel the packet.
	HandleServerBoundPacket(ctx *event.Context, pk packet.Packet)
	// HandleTransfer handles a session being transferred to another server. ctx.Cancel() may be called to
	// cancel the transfer.
	HandleTransfer(ctx *event.Context, svr *server.Server)
}

// NopHandler implements the Handler interface but does not execute any code when an event is called. The
// default handler of sessions is set to NopHandler.
// Users may embed NopHandler to avoid having to implement each method.
type NopHandler struct{}

// Compile time check to make sure NopHandler implements Handler.
var _ Handler = (*NopHandler)(nil)

// HandleClientBoundPacket ...
func (NopHandler) HandleClientBoundPacket(*event.Context, packet.Packet) {}

// HandleServerBoundPacket ...
func (NopHandler) HandleServerBoundPacket(*event.Context, packet.Packet) {}

// HandleTransfer ...
func (NopHandler) HandleTransfer(*event.Context, *server.Server) {}
