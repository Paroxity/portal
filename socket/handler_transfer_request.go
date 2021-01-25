package socket

import (
	"github.com/paroxity/portal/server"
	"github.com/paroxity/portal/session"
	"github.com/paroxity/portal/socket/packet"
)

// TransferRequestHandler is responsible for handling the TransferRequest packet sent by servers.
type TransferRequestHandler struct{}

// Handle ...
func (*TransferRequestHandler) Handle(p packet.Packet, c *Client) error {
	pk := p.(*packet.TransferRequest)
	response := func(status byte, error string) error {
		return c.WritePacket(&packet.TransferResponse{
			PlayerUUID: pk.PlayerUUID,
			Status:     status,
			Error:      error,
		})
	}

	g, ok := server.GroupFromName(pk.Group)
	if !ok {
		return response(packet.TransferResponseGroupNotFound, "")
	}

	srv, ok := g.Server(pk.Server)
	if !ok {
		return response(packet.TransferResponseServerNotFound, "")
	}

	s, ok := session.Lookup(pk.PlayerUUID)
	if !ok {
		return response(packet.TransferResponsePlayerNotFound, "")
	}

	if s.Server().Address() == srv.Address() {
		return response(packet.TransferResponseAlreadyOnServer, "")
	}

	if err := s.Transfer(srv); err != nil {
		return response(packet.TransferResponseError, err.Error())
	}

	return response(packet.TransferResponseSuccess, "")
}
