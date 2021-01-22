package socket

import (
	"github.com/paroxity/portal/session"
	portalpacket "github.com/paroxity/portal/socket/packet"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type PlayerInfoRequest struct{}

func (*PlayerInfoRequest) Handle(p packet.Packet, c *Client) error {
	pk := p.(*portalpacket.PlayerInfoRequest)
	response := func(status byte, xuid string, address string) error {
		return c.WritePacket(&portalpacket.PlayerInfoResponse{
			PlayerUUID: pk.PlayerUUID,
			Status:     status,
			XUID:       xuid,
			Address:    address,
		})
	}

	s, ok := session.Lookup(pk.PlayerUUID)
	if !ok {
		return response(portalpacket.PlayerInfoResponsePlayerNotFound, "", "")
	}

	return response(portalpacket.PlayerInfoResponseSuccess, s.Conn().IdentityData().XUID, s.Conn().RemoteAddr().String())
}
