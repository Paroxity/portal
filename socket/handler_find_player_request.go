package socket

import (
	"github.com/paroxity/portal/session"
	"github.com/paroxity/portal/socket/packet"
)

// FindPlayerRequestHandler is responsible for handling the FindPlayerRequest packet sent by servers.
type FindPlayerRequestHandler struct{}

// Handle ...
func (*FindPlayerRequestHandler) Handle(p packet.Packet, c *Client) error {
	pk := p.(*packet.FindPlayerRequest)
	s, ok := session.Lookup(pk.PlayerUUID)
	if !ok {
		s, ok = session.LookupByName(pk.PlayerName)
		if !ok {
			return c.WritePacket(&packet.FindPlayerResponse{
				PlayerUUID: pk.PlayerUUID,
				PlayerName: pk.PlayerName,
				Online:     false,
			})
		}
	}

	return c.WritePacket(&packet.FindPlayerResponse{
		PlayerUUID: s.UUID(),
		PlayerName: s.Conn().IdentityData().DisplayName,
		Online:     true,
		Group:      s.Server().Group(),
		Server:     s.Server().Name(),
	})
}
