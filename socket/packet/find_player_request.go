package packet

import (
	"github.com/google/uuid"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

// FindPlayerRequest is sent by a connection to find the server the request player is currently on.
type FindPlayerRequest struct {
	// PlayerUUID is the UUID of the player to find.
	PlayerUUID uuid.UUID
	// PlayerName is the name of the player to find.
	PlayerName string
}

// ID ...
func (*FindPlayerRequest) ID() uint16 {
	return IDFindPlayerRequest
}

// Marshal ...
func (pk *FindPlayerRequest) Marshal(w *protocol.Writer) {
	w.UUID(&pk.PlayerUUID)
	w.String(&pk.PlayerName)
}

// Unmarshal ...
func (pk *FindPlayerRequest) Unmarshal(r *protocol.Reader) {
	r.UUID(&pk.PlayerUUID)
	r.String(&pk.PlayerName)
}
