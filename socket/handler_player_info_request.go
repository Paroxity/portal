package socket

import (
	"github.com/paroxity/portal/session"
	"github.com/paroxity/portal/socket/packet"
)

type PlayerInfoRequest struct{}

func (*PlayerInfoRequest) Handle(p packet.Packet, c *Client) error {
	pk := p.(*packet.PlayerInfoRequest)
	response := func(status byte, xuid string, address string) error {
		return c.WritePacket(&packet.PlayerInfoResponse{
			PlayerUUID: pk.PlayerUUID,
			Status:     status,
			XUID:       xuid,
			Address:    address,
		})
	}

	s, ok := session.Lookup(pk.PlayerUUID)
	if !ok {
		return response(packet.PlayerInfoResponsePlayerNotFound, "", "")
	}

	return response(packet.PlayerInfoResponseSuccess, s.Conn().IdentityData().XUID, s.Conn().RemoteAddr().String())
}