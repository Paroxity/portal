package packet

import (
	"github.com/google/uuid"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

// TransferRequest is sent by a server to request the transfer of a player.
type TransferRequest struct {
	// PlayerUUID is the UUID of the player to be transferred.
	PlayerUUID uuid.UUID
	// Server is the name of the server in the group to transfer to.
	Server string
}

// ID ...
func (*TransferRequest) ID() uint16 {
	return IDTransferRequest
}

// Marshal ...
func (pk *TransferRequest) Marshal(w *protocol.Writer) {
	w.UUID(&pk.PlayerUUID)
	w.String(&pk.Server)
}

// Unmarshal ...
func (pk *TransferRequest) Unmarshal(r *protocol.Reader) {
	r.UUID(&pk.PlayerUUID)
	r.String(&pk.Server)
}
