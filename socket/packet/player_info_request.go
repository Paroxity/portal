package packet

import (
	"github.com/google/uuid"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

// PlayerInfoRequest is sent by a connection to request information such as the XUID and IP address of a
// player connected to the proxy.
type PlayerInfoRequest struct {
	// PlayerUUID is the UUID of the player to get information about.
	PlayerUUID uuid.UUID
}

// ID ...
func (*PlayerInfoRequest) ID() uint16 {
	return IDPlayerInfoRequest
}

// Marshal ...
func (pk *PlayerInfoRequest) Marshal(w *protocol.Writer) {
	w.UUID(&pk.PlayerUUID)
}

// Unmarshal ...
func (pk *PlayerInfoRequest) Unmarshal(r *protocol.Reader) {
	r.UUID(&pk.PlayerUUID)
}
