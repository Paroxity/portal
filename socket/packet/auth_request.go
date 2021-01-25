package packet

import "github.com/sandertv/gophertunnel/minecraft/protocol"

const (
	ClientTypeServer = iota
)

// AuthRequest is sent by a connection to authenticate with the proxy.
type AuthRequest struct {
	// Type is the type of connection being made. The different types can be seen above.
	Type byte
	// Secret is the secret key to authenticate with. It must match the configured key in the proxy otherwise
	// the client will not be authenticated.
	Secret string
	// Name is the name of the client that is being authenticated. The name must be different to existing
	// connections.
	Name string
	// ExtraData contains extra data linked to the connection. Different types of connections will require
	// different data to be provided.
	ExtraData []byte
}

// ID ...
func (pk *AuthRequest) ID() uint16 {
	return IDAuthRequest
}

// Marshal ...
func (pk *AuthRequest) Marshal(w *protocol.Writer) {
	w.Uint8(&pk.Type)
	w.String(&pk.Secret)
	w.String(&pk.Name)
	w.Bytes(&pk.ExtraData)
}

// Unmarshal ...
func (pk *AuthRequest) Unmarshal(r *protocol.Reader) {
	r.Uint8(&pk.Type)
	r.String(&pk.Secret)
	r.String(&pk.Name)
	r.Bytes(&pk.ExtraData)
}
