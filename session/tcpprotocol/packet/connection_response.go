package packet

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

const (
	ConnectionResponseSuccess byte = iota
	ConnectionResponseUnsupportedProtocol
)

// ConnectionResponse is sent by the server in response to a ConnectionRequest packet. It contains the response and if
// the player was able to connect to the server successfully.
type ConnectionResponse struct {
	// Response is the response from the server. This can be one of the constants above.
	Response byte
}

// ID ...
func (pk *ConnectionResponse) ID() uint32 {
	return IDConnectionResponse
}

// Marshal ...
func (pk *ConnectionResponse) Marshal(w *protocol.Writer) {
	w.Uint8(&pk.Response)
}

// Unmarshal ...
func (pk *ConnectionResponse) Unmarshal(r *protocol.Reader) {
	r.Uint8(&pk.Response)
}
