package packet

import "github.com/sandertv/gophertunnel/minecraft/protocol"

const (
	AuthResponseSuccess byte = iota
	AuthResponseIncorrectSecret
	AuthResponseUnknownType
	AuthResponseInvalidData
)

// AuthResponse is sent by the proxy in response to AuthRequest. It tells the client if the authentication
// request was successful or not.
type AuthResponse struct {
	// Status is the response status from authentication. The possible values for this can be found above.
	Status byte
}

func (*AuthResponse) ID() uint16 {
	return IDAuthResponse
}

func (pk *AuthResponse) Marshal(w *protocol.Writer) {
	w.Uint8(&pk.Status)
}

func (pk *AuthResponse) Unmarshal(r *protocol.Reader) {
	r.Uint8(&pk.Status)
}
