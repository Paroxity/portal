package tcpprotocol

import (
	"github.com/df-mc/dragonfly/server/session"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"net"
)

// Listener implements a Minecraft listener on top of an unspecific net.Listener. It abstracts away the login sequence
// of connecting clients and provides the implements the net.Listener interface to provide a consistent API.
type Listener struct {
	listener net.Listener
}

// Listen announces on the local network address. The network must be "tcp", "tcp4", "tcp6", "unix", "unixpacket"
// A Listener is returned which may be used to accept connections. If the host in the address parameter is empty or a
// literal unspecified IP address, Listen listens on all available unicast and anycast IP addresses of the local system.
func Listen(network, address string) (*Listener, error) {
	l, err := net.Listen(network, address)
	if err != nil {
		return nil, err
	}
	listener := &Listener{
		listener: l,
	}
	return listener, nil
}

// Accept ...
func (l *Listener) Accept() (session.Conn, error) {
	netConn, err := l.listener.Accept()
	if err != nil {
		return nil, err
	}

	conn := newConn(netConn)
	err = conn.identify()
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// Disconnect ...
func (l *Listener) Disconnect(conn session.Conn, message string) error {
	_ = conn.WritePacket(&packet.Disconnect{
		HideDisconnectionScreen: message == "",
		Message:                 message,
	})
	return conn.Close()
}

// Close ...
func (l *Listener) Close() error {
	return l.listener.Close()
}
