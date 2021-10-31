package socket

import (
	"github.com/paroxity/portal/socket/packet"
)

// FindPlayerRequestHandler is responsible for handling the FindPlayerRequest packet sent by servers.
type FindPlayerRequestHandler struct{ requireAuth }

// Handle ...
func (*FindPlayerRequestHandler) Handle(p packet.Packet, srv Server, c *Client) error {
	pk := p.(*packet.FindPlayerRequest)
	s, ok := srv.SessionStore().Load(pk.PlayerUUID)
	if !ok {
		s, ok = srv.SessionStore().LoadFromName(pk.PlayerName)
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
		Server:     s.Server().Name(),
	})
}
