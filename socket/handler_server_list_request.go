package socket

import (
	"github.com/paroxity/portal/server"
	"github.com/paroxity/portal/socket/packet"
)

// ServerListRequestHandler is responsible for handling the ServerList packet sent by servers.
type ServerListRequestHandler struct{}

// Handle ...
func (*ServerListRequestHandler) Handle(_ packet.Packet, c *Client) error {
	servers := make([]packet.ServerEntry, 0)
	for _, g := range server.Groups() {
		for _, s := range g.Servers() {
			name := s.Name()
			group := s.Group()
			isOnline := s.Connected()
			playerCount := uint16(s.PlayerCount())

			servers = append(servers, packet.ServerEntry{
				Name:        &name,
				Group:       &group,
				IsOnline:    &isOnline,
				PlayerCount: &playerCount,
			})
		}
	}

	return c.WritePacket(&packet.ServerListResponse{
		Servers: servers,
	})
}
