package socket

import (
	"github.com/paroxity/portal/socket/packet"
	_ "unsafe"
)

// AuthRequestHandler is responsible for handling the AuthRequest packet sent by servers.
type AuthRequestHandler struct{}

// Handle ...
func (*AuthRequestHandler) Handle(p packet.Packet, srv Server, c *Client) error {
	pk := p.(*packet.AuthRequest)

	if c.Authenticated() {
		return nil
	}

	if pk.Protocol != packet.ProtocolVersion {
		srv.Logger().Errorf("failed socket authentication attempt from \"%s\": unsupported protocol version %d", pk.Name, pk.Protocol)
		return c.WritePacket(&packet.AuthResponse{Status: packet.AuthResponseUnsupportedProtocol})
	}
	if pk.Secret != srv.Secret() {
		srv.Logger().Errorf("failed socket authentication attempt from \"%s\": incorrect secret provided", pk.Name)
		return c.WritePacket(&packet.AuthResponse{Status: packet.AuthResponseIncorrectSecret})
	}
	_, ok := srv.Client(pk.Name)
	if ok {
		srv.Logger().Errorf("failed socket authentication attempt from \"%s\": a connection already exists with this name", pk.Name)
		return c.WritePacket(&packet.AuthResponse{Status: packet.AuthResponseAlreadyConnected})
	}

	srv.Authenticate(c, pk.Name)
	srv.Logger().Debugf("socket connection \"%s\" successfully authenticated", pk.Name)
	return c.WritePacket(&packet.AuthResponse{Status: packet.AuthResponseSuccess})
}

// RequiresAuth ...
func (*AuthRequestHandler) RequiresAuth() bool {
	return false
}
