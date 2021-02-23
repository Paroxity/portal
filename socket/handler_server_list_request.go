package socket

import (
	"github.com/paroxity/portal/server"
	"github.com/paroxity/portal/socket/packet"
)

// ServerListRequestHandler is responsible for handling the ServerListRequest packet sent by servers.
type ServerListRequestHandler struct{}

// Handle ...
func (*ServerListRequestHandler) Handle(_ packet.Packet, c *Client) error {
	var servers []packet.ServerEntry

	var entry packet.ServerEntry
	for _, g := range server.Groups() {
		for _, s := range g.Servers() {
			entry = packet.ServerEntry{
				Name:        s.Name(),
				Group:       s.Group(),
				Online:      s.Connected(),
				PlayerCount: int64(s.PlayerCount()),
			}
			servers = append(servers, entry)
		}
	}

	return c.WritePacket(&packet.ServerListResponse{
		Servers: servers,
	})
}
