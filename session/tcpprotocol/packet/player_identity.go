package packet

import (
	"encoding/json"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/login"
)

// PlayerIdentity is sent by the proxy to give the server information about the client that would usually be sent within
// the login sequence.
type PlayerIdentity struct {
	// IdentityData contains identity data of the player logged in.
	IdentityData login.IdentityData
	// ClientData is a container of client specific data of a Login packet. It holds data such as the skin of a player,
	// but also its language code and device information.
	ClientData login.ClientData
	// EnableClientCache, if set to true, enables the client blob cache for the client. This means that the server will
	// send chunks as blobs, which may be saved by the client so that chunks don't have to be transmitted every time,
	// resulting in less network transmission.
	EnableClientCache bool
	// Address is the address of the player that has joined the server.
	Address string
}

// ID ...
func (pk *PlayerIdentity) ID() uint32 {
	return IDPlayerIdentity
}

// Marshal ...
func (pk *PlayerIdentity) Marshal(w *protocol.Writer) {
	identityData, _ := json.Marshal(pk.IdentityData)
	clientData, _ := json.Marshal(pk.ClientData)
	w.ByteSlice(&identityData)
	w.ByteSlice(&clientData)
	w.Bool(&pk.EnableClientCache)
	w.String(&pk.Address)
}

// Unmarshal ...
func (pk *PlayerIdentity) Unmarshal(r *protocol.Reader) {
	var identityData, clientData []byte
	r.ByteSlice(&identityData)
	r.ByteSlice(&clientData)
	r.Bool(&pk.EnableClientCache)
	r.String(&pk.Address)
	_ = json.Unmarshal(identityData, &pk.IdentityData)
	_ = json.Unmarshal(clientData, &pk.ClientData)
}
