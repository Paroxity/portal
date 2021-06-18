package packet

import (
	"github.com/google/uuid"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

// UpdatePlayerLatency is sent by the proxy to update a player's latency on the server they are connected to.
type UpdatePlayerLatency struct {
	// PlayerUUID is the UUID of the player the latency belongs to.
	PlayerUUID uuid.UUID
	// Latency is the latency of the player's connection to the proxy in milliseconds.
	Latency int64
}

// ID ...
func (*UpdatePlayerLatency) ID() uint16 {
	return IDUpdatePlayerLatency
}

// Marshal ...
func (pk *UpdatePlayerLatency) Marshal(w *protocol.Writer) {
	w.UUID(&pk.PlayerUUID)
	w.Int64(&pk.Latency)
}

// Unmarshal ...
func (pk *UpdatePlayerLatency) Unmarshal(r *protocol.Reader) {
	r.UUID(&pk.PlayerUUID)
	r.Int64(&pk.Latency)
}
