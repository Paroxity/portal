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

	g, ok := server.GroupFromName(pk.Group)
	if !ok {
		return c.WritePacket(&portalpacket.TransferResponse{
			Status: portalpacket.TransferResponseGroupNotFound,
			Reason: fmt.Sprintf("Group %s does not exist", pk.Group),
		})
	}

	srv, ok := g.Server(pk.Server)
	if !ok {
		return c.WritePacket(&portalpacket.TransferResponse{
			Status: portalpacket.TransferResponseServerNotFound,
			Reason: fmt.Sprintf("Server %s does not exist in group %s", pk.Server, pk.Group),
		})
	}
	if srv == nil {
		// TODO: Send response saying no servers in group
		return nil
	}

	for _, s := range session.All() {
		if s.ServerConn().GameData().EntityRuntimeID == pk.PlayerRuntimeID {
			if s.Server().Address() == srv.Address() {
				return c.WritePacket(&portalpacket.TransferResponse{
					Status: portalpacket.TransferResponseAlreadyOnServer,
					Reason: "Player is already on the server",
				})
			}

			if err := s.Transfer(srv); err != nil {
				// TODO: Send response saying error
				return nil
			}
			return c.WritePacket(&portalpacket.TransferResponse{
				Status: portalpacket.TransferResponseSuccess,
				Reason: "The player was transferred to the server",
			})
		}
	}

	return c.WritePacket(&portalpacket.TransferResponse{
		Status: portalpacket.TransferResponsePlayerNotFound,
		Reason: "Player is not connected to the proxy",
	})
}
