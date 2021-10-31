package socket

import (
	"github.com/paroxity/portal/socket/packet"
)

// ServerListRequestHandler is responsible for handling the ServerListRequest packet sent by servers.
type ServerListRequestHandler struct{ requireAuth }

// Handle ...
func (*ServerListRequestHandler) Handle(_ packet.Packet, srv Server, c *Client) error {
	var servers []packet.ServerEntry

	for _, s := range srv.ServerRegistry().Servers() {
		entry := packet.ServerEntry{
			Name:        s.Name(),
			Online:      s.Connected(),
			PlayerCount: int64(s.PlayerCount()),
		}
		servers = append(servers, entry)
	}

	return c.WritePacket(&packet.ServerListResponse{
		Servers: servers,
	})
}
