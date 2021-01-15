package packet

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

// Transfer request is sent by a server to request the transfer of a player.
type TransferRequest struct {
	// PlayerRuntimeID is the entity runtime ID of the player to be transferred.
	PlayerRuntimeID uint64
	// Group is the name of the group to transfer the player to.
	Group string
	// Server is the name of the server in the group to transfer to.
	Server string
}

// ID ...
func (*TransferRequest) ID() uint32 {
	return IDTransferRequest
}

// Marshal ...
func (pk *TransferRequest) Marshal(w *protocol.Writer) {
	w.Varuint64(&pk.PlayerRuntimeID)
	w.String(&pk.Group)
	w.String(&pk.Server)
}

// Unmarshal ...
func (pk *TransferRequest) Unmarshal(r *protocol.Reader) {
	r.Varuint64(&pk.PlayerRuntimeID)
	r.String(&pk.Group)
	r.String(&pk.Server)
}
