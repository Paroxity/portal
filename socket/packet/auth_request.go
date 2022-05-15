package packet

import "github.com/sandertv/gophertunnel/minecraft/protocol"

// AuthRequest is sent by a connection to authenticate with the proxy.
type AuthRequest struct {
	// Protocol is the protocol version supported by the client. It must match the proxy version otherwise the client
	// cannot authenticate.
	Protocol uint32
	// Secret is the secret key to authenticate with. It must match the configured key in the proxy otherwise
	// the client will not be authenticated.
	Secret string
	// Name is the name of the client that is being authenticated. The name must be different to existing
	// connections.
	Name string
}

// ID ...
func (pk *AuthRequest) ID() uint16 {
	return IDAuthRequest
}

// Marshal ...
func (pk *AuthRequest) Marshal(w *protocol.Writer) {
	w.Uint32(&pk.Protocol)
	w.String(&pk.Secret)
	w.String(&pk.Name)
}

// Unmarshal ...
func (pk *AuthRequest) Unmarshal(r *protocol.Reader) {
	r.Uint32(&pk.Protocol)
	r.String(&pk.Secret)
	r.String(&pk.Name)
}
