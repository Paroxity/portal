package tcpprotocol

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol/login"
	"net"
)

// Dialer allows specifying specific settings for connection to a Minecraft server. The zero value of Dialer is used for
// the package level Dial function.
type Dialer struct {
	// IdentityData is the identity data used to login to the server with. It includes the username, UUID and XUID of the
	// player. The IdentityData object is obtained using Minecraft auth if Email and Password are set. If not, the object
	// provided here is used, or a default one if left empty.
	IdentityData login.IdentityData
	// ClientData is the client data used to login to the server with. It includes fields such as the skin, locale and
	// UUIDs unique to the client. If empty, a default is sent produced using defaultClientData().
	ClientData login.ClientData
	// EnableClientCache, if set to true, enables the client blob cache for the client. This means that the server will
	// send chunks as blobs, which may be saved by the client so that chunks don't have to be transmitted every time,
	// resulting in less network transmission.
	EnableClientCache bool
}

// Dial dials a Minecraft connection to the address passed over the network passed. The network is typically "tcp". A
// Conn is returned which may be used to receive packets from and send packets to. A zero value of a Dialer struct is
// used to initiate the connection. A custom Dialer may be used to specify additional behaviour.
func Dial(network, address, playerAddress string) (*Conn, error) {
	var d Dialer
	return d.Dial(network, address, playerAddress)
}

// Dial dials a Minecraft connection to the address passed over the network passed. The network is typically
// "raknet". A Conn is returned which may be used to receive packets from and send packets to.
func (d Dialer) Dial(network, address, playerAddress string) (*Conn, error) {
	netConn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}

	conn := newConn(netConn)
	conn.identityData = d.IdentityData
	conn.clientData = d.ClientData
	conn.enableClientCache = d.EnableClientCache

	err = conn.login(playerAddress)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
