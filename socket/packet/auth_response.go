package packet

import "github.com/sandertv/gophertunnel/minecraft/protocol"

const (
	AuthResponseSuccess = iota
	AuthResponseIncorrectSecret
	AuthResponseUnknownType
	AuthResponseInvalidData
)

// AuthResponse is sent by the proxy in response to AuthRequest. It tells the client if the authentication
// request was successful or not.
type AuthResponse struct {
	// Status is the response status from authentication. The possible values for this can be found above.
	Status byte
	// Reason is the reason behind the Status provided.
	Reason string
}

func (*AuthResponse) ID() uint32 {
	return IDAuthResponse
}

func (pk *AuthResponse) Marshal(w *protocol.Writer) {
	w.Uint8(&pk.Status)
	w.String(&pk.Reason)
}

func (pk *AuthResponse) Unmarshal(r *protocol.Reader) {
	r.Uint8(&pk.Status)
	r.String(&pk.Reason)
}
