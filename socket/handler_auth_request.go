package socket

import (
	"bytes"
	"github.com/paroxity/portal/server"
	"github.com/paroxity/portal/socket/packet"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	_ "unsafe"
)

// AuthRequestHandler is responsible for handling the AuthRequest packet sent by servers.
type AuthRequestHandler struct{}

// Handle ...
func (*AuthRequestHandler) Handle(p packet.Packet, srv Server, c *Client) error {
	pk := p.(*packet.AuthRequest)

	if pk.Secret != srv.Secret() {
		srv.Logger().Errorf("failed socket authentication attempt from \"%s\": incorrect secret provided", pk.Name)
		return c.WritePacket(&packet.AuthResponse{Status: packet.AuthResponseIncorrectSecret})
	}

	data := bytes.NewBuffer(pk.ExtraData)
	r := protocol.NewReader(data, 0)
	switch pk.Type {
	case packet.ClientTypeServer:
		var address string
		r.String(&address)

		s, ok := srv.ServerRegistry().Server(pk.Name)
		if !ok {
			s = server.New(pk.Name, address)
			srv.ServerRegistry().AddServer(s)
		} else if s.Connected() {
			srv.Logger().Errorf("failed socket authentication attempt from \"%s\": server is already connected", pk.Name)
			return c.WritePacket(&packet.AuthResponse{Status: packet.AuthResponseInvalidData})
		}

		c.name = pk.Name
		c.clientType = pk.Type
		c.extraData["address"] = address

		server_setConn(s, c)
	default:
		return c.WritePacket(&packet.AuthResponse{Status: packet.AuthResponseUnknownType})
	}

	c.Authenticate()
	srv.Logger().Infof("socket connection \"%s\" successfully authenticated", pk.Name)
	return c.WritePacket(&packet.AuthResponse{Status: packet.AuthResponseSuccess})
}

// RequiresAuth ...
func (*AuthRequestHandler) RequiresAuth() bool {
	return false
}

//go:linkname server_setConn github.com/paroxity/portal/server.(*Server).setConn
//noinspection ALL
func server_setConn(s *server.Server, c server.Client)
