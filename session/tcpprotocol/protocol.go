package tcpprotocol

import (
	packet2 "github.com/paroxity/portal/session/tcpprotocol/packet"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// ProtocolVersion is the current supported version of the protocol. If a server is using an outdated version of the
// protocol, players will be unable to connect. This constant gets updated every time the protocol is changed.
const ProtocolVersion = 1

func init() {
	packet.Register(packet2.IDPlayerIdentity, func() packet.Packet { return &packet2.PlayerIdentity{} })
	packet.Register(packet2.IDConnectionRequest, func() packet.Packet { return &packet2.ConnectionRequest{} })
	packet.Register(packet2.IDConnectionResponse, func() packet.Packet { return &packet2.ConnectionResponse{} })
}
