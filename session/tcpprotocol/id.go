package tcpprotocol

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"math"
)

// ProtocolVersion is the current supported version of the protocol. If a server is using an outdated version of the
// protocol, players will be unable to connect. This constant gets updated every time the protocol is changed.
const ProtocolVersion = 1

const (
	IDConnectionRequest uint32 = math.MaxUint32 - iota
	IDConnectionResponse
	IDPlayerIdentity
)

func init() {
	packet.Register(IDPlayerIdentity, func() packet.Packet { return &PlayerIdentity{} })
	packet.Register(IDConnectionRequest, func() packet.Packet { return &ConnectionRequest{} })
	packet.Register(IDConnectionResponse, func() packet.Packet { return &ConnectionResponse{} })
}
