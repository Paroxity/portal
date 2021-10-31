package socket

import (
	"github.com/paroxity/portal/socket/packet"
)

// PlayerInfoRequestHandler is responsible for handling the PlayerInfoRequest packet sent by servers.
type PlayerInfoRequestHandler struct{ requireAuth }

// Handle ...
func (*PlayerInfoRequestHandler) Handle(p packet.Packet, srv Server, c *Client) error {
	pk := p.(*packet.PlayerInfoRequest)
	response := func(status byte, xuid string, address string) error {
		return c.WritePacket(&packet.PlayerInfoResponse{
			PlayerUUID: pk.PlayerUUID,
			Status:     status,
			XUID:       xuid,
			Address:    address,
		})
	}

	s, ok := srv.SessionStore().Load(pk.PlayerUUID)
	if !ok {
		return response(packet.PlayerInfoResponsePlayerNotFound, "", "")
	}

	return response(packet.PlayerInfoResponseSuccess, s.Conn().IdentityData().XUID, s.Conn().RemoteAddr().String())
}
