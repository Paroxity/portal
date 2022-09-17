package tcpprotocol

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
	w.String(&pk.Address)
}

// Unmarshal ...
func (pk *PlayerIdentity) Unmarshal(r *protocol.Reader) {
	var identityData, clientData []byte
	r.ByteSlice(&identityData)
	r.ByteSlice(&clientData)
	r.String(&pk.Address)
	_ = json.Unmarshal(identityData, &pk.IdentityData)
	_ = json.Unmarshal(clientData, &pk.ClientData)
}
