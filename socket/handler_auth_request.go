package socket

import (
	"bytes"
	"github.com/paroxity/portal/config"
	"github.com/paroxity/portal/server"
	"github.com/paroxity/portal/socket/packet"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sirupsen/logrus"
	_ "unsafe"
)

// AuthRequestHandler is responsible for handling the AuthRequest packet sent by servers.
type AuthRequestHandler struct{}

// Handle ...
func (*AuthRequestHandler) Handle(p packet.Packet, c *Client) error {
	pk := p.(*packet.AuthRequest)

	if pk.Secret != config.SocketSecret() {
		logrus.Errorf("Failed socket authentication attempt from \"%s\": Incorrect secret provided", pk.Name)
		return c.WritePacket(&packet.AuthResponse{Status: packet.AuthResponseIncorrectSecret})
	}

	data := bytes.NewBuffer(pk.ExtraData)
	r := protocol.NewReader(data, 0)
	switch pk.Type {
	case packet.ClientTypeServer:
		var group, address string
		r.String(&group)
		r.String(&address)

		g, ok := server.GroupFromName(group)
		if !ok {
			logrus.Errorf("Failed socket authentication attempt from \"%s\": Group \"%s\" not found", pk.Name, group)
			return c.WritePacket(&packet.AuthResponse{Status: packet.AuthResponseInvalidData})
		}

		s, ok := g.Server(pk.Name)
		if !ok {
			s = server.New(pk.Name, g.Name(), address)
			g.AddServer(s)
		} else if s.Connected() {
			logrus.Errorf("Failed socket authentication attempt from \"%s\": Server is already connected\n", pk.Name)
			return c.WritePacket(&packet.AuthResponse{Status: packet.AuthResponseInvalidData})
		}

		c.name = pk.Name
		c.clientType = pk.Type
		c.extraData["address"] = address
		c.extraData["group"] = g.Name()

		server_setConn(s, c)
	default:
		return c.WritePacket(&packet.AuthResponse{Status: packet.AuthResponseUnknownType})
	}

	logrus.Infof("Socket connection \"%s\" successfully authenticated\n", pk.Name)
	return c.WritePacket(&packet.AuthResponse{Status: packet.AuthResponseSuccess})
}

//go:linkname server_setConn github.com/paroxity/portal/server.(*Server).setConn
//noinspection ALL
func server_setConn(s *server.Server, c server.Client)
