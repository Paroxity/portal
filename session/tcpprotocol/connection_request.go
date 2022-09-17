package tcpprotocol

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

// ConnectionRequest is sent by the proxy to request a connection for a player who is attempting to join the server.
type ConnectionRequest struct {
	// ProtocolVersion is the protocol version of the TCP protocol used by the proxy.
	ProtocolVersion uint32
}

// ID ...
func (pk *ConnectionRequest) ID() uint32 {
	return IDConnectionRequest
}

// Marshal ...
func (pk *ConnectionRequest) Marshal(w *protocol.Writer) {
	w.Uint32(&pk.ProtocolVersion)
}

// Unmarshal ...
func (pk *ConnectionRequest) Unmarshal(r *protocol.Reader) {
	r.Uint32(&pk.ProtocolVersion)
}
