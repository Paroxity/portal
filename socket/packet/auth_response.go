package packet

import "github.com/sandertv/gophertunnel/minecraft/protocol"

const (
	AuthResponseSuccess byte = iota
	AuthResponseUnsupportedProtocol
	AuthResponseIncorrectSecret
	AuthResponseAlreadyConnected
	AuthResponseUnauthenticated
)

// AuthResponse is sent by the proxy in response to AuthRequest. It tells the client if the authentication
// request was successful or not.
type AuthResponse struct {
	// Protocol is the protocol version supported by the socket server.
	Protocol uint32
	// Status is the response status from authentication. The possible values for this can be found above.
	Status byte
}

// ID ...
func (*AuthResponse) ID() uint16 {
	return IDAuthResponse
}

// Marshal ...
func (pk *AuthResponse) Marshal(w *protocol.Writer) {
	w.Uint32(&pk.Protocol)
	w.Uint8(&pk.Status)
}

// Unmarshal ...
func (pk *AuthResponse) Unmarshal(r *protocol.Reader) {
	r.Uint32(&pk.Protocol)
	r.Uint8(&pk.Status)
}
