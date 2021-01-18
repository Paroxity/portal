package socket

import (
	"github.com/paroxity/portal/server"
	"github.com/paroxity/portal/session"
	portalpacket "github.com/paroxity/portal/socket/packet"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// TransferRequestHandler is responsible for handling the TransferRequest packet sent by servers.
type TransferRequestHandler struct{}

// Handle ...
func (*TransferRequestHandler) Handle(p packet.Packet, c *Client) error {
	pk := p.(*portalpacket.TransferRequest)
	response := func(status byte, error string) error {
		return c.WritePacket(&portalpacket.TransferResponse{
			PlayerUUID: pk.PlayerUUID,
			Status:     status,
			Error:      error,
		})
	}

	g, ok := server.GroupFromName(pk.Group)
	if !ok {
		return response(portalpacket.TransferResponseGroupNotFound, "")
	}

	srv, ok := g.Server(pk.Server)
	if !ok {
		return response(portalpacket.TransferResponseServerNotFound, "")
	}

	s, ok := session.Lookup(pk.PlayerUUID)
	if !ok {
		return response(portalpacket.TransferResponsePlayerNotFound, "")
	}

	if s.Server().Address() == srv.Address() {
		return response(portalpacket.TransferResponseAlreadyOnServer, "")
	}

	if err := s.Transfer(srv); err != nil {
		return response(portalpacket.TransferResponseError, err.Error())
	}

	return response(portalpacket.TransferResponseSuccess, "")
}
