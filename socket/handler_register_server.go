package socket

import (
	"github.com/paroxity/portal/server"
	"github.com/paroxity/portal/socket/packet"
)

// RegisterServerHandler is responsible for handling the RegisterServer packet sent by servers.
type RegisterServerHandler struct{ requireAuth }

// Handle ...
func (*RegisterServerHandler) Handle(p packet.Packet, srv Server, c *Client) error {
	pk := p.(*packet.RegisterServer)
	srv.ServerRegistry().AddServer(server.New(c.Name(), pk.Address))
	srv.Logger().Debugf("socket connection \"%s\" has registered itself as a server with the address \"%s\"", c.Name(), pk.Address)
	return nil
}
