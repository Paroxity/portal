package socket

import (
	"github.com/paroxity/portal/socket/packet"
)

// TransferRequestHandler is responsible for handling the TransferRequest packet sent by servers.
type TransferRequestHandler struct{ requireAuth }

// Handle ...
func (*TransferRequestHandler) Handle(p packet.Packet, srv Server, c *Client) error {
	pk := p.(*packet.TransferRequest)
	response := func(status byte, error string) error {
		return c.WritePacket(&packet.TransferResponse{
			PlayerUUID: pk.PlayerUUID,
			Status:     status,
			Error:      error,
		})
	}

	targetSrv, ok := srv.ServerRegistry().Server(pk.Server)
	if !ok {
		return response(packet.TransferResponseServerNotFound, "")
	}

	s, ok := srv.SessionStore().Load(pk.PlayerUUID)
	if !ok {
		return response(packet.TransferResponsePlayerNotFound, "")
	}

	if s.Server().Address() == targetSrv.Address() {
		return response(packet.TransferResponseAlreadyOnServer, "")
	}

	if err := s.Transfer(targetSrv); err != nil {
		return response(packet.TransferResponseError, err.Error())
	}

	return response(packet.TransferResponseSuccess, "")
}
