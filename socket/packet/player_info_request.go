package packet

import (
	"github.com/google/uuid"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

type PlayerInfoRequest struct {
	PlayerUUID uuid.UUID
}

func (*PlayerInfoRequest) ID() uint32 {
	return IDPlayerInfoRequest
}

func (pk *PlayerInfoRequest) Marshal(w *protocol.Writer) {
	w.UUID(&pk.PlayerUUID)
}

func (pk *PlayerInfoRequest) Unmarshal(r *protocol.Reader) {
	r.UUID(&pk.PlayerUUID)
}
