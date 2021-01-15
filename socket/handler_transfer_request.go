package socket

import (
	"fmt"
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
	response := func(status byte, reason string) error {
		return c.WritePacket(&portalpacket.TransferResponse{
			PlayerRuntimeID: pk.PlayerRuntimeID,
			Status:          portalpacket.TransferResponseGroupNotFound,
			Reason:          fmt.Sprintf("Group %s does not exist", pk.Group),
		})
	}

	g, ok := server.GroupFromName(pk.Group)
	if !ok {
		return response(portalpacket.TransferResponseGroupNotFound, fmt.Sprintf("Group %s does not exist", pk.Group))
	}

	srv, ok := g.Server(pk.Server)
	if !ok {
		return response(portalpacket.TransferResponseServerNotFound, fmt.Sprintf("Server %s does not exist in group %s", pk.Server, pk.Group))
	}

	for _, s := range session.All() {
		if s.ServerConn().GameData().EntityRuntimeID == pk.PlayerRuntimeID {
			if s.Server().Address() == srv.Address() {
				return response(portalpacket.TransferResponseAlreadyOnServer, "Player is already on the server")
			}

			if err := s.Transfer(srv); err != nil {
				return response(portalpacket.TransferResponseError, err.Error())
			}
			return response(portalpacket.TransferResponseSuccess, "The player was transferred to the server")
		}
	}

	return response(portalpacket.TransferResponsePlayerNotFound, "Player is not connected to the proxy")
}
