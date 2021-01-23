package packet

import (
	"github.com/google/uuid"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

const (
	PlayerInfoResponseSuccess byte = iota
	PlayerInfoResponsePlayerNotFound
)

type PlayerInfoResponse struct {
	PlayerUUID uuid.UUID
	Status     byte
	XUID       string
	Address    string
}

func (*PlayerInfoResponse) ID() uint32 {
	return IDPlayerInfoResponse
}

func (pk *PlayerInfoResponse) Marshal(w *protocol.Writer) {
	w.UUID(&pk.PlayerUUID)
	w.Uint8(&pk.Status)
	w.String(&pk.XUID)
	w.String(&pk.Address)
}

func (pk *PlayerInfoResponse) Unmarshal(r *protocol.Reader) {
	r.UUID(&pk.PlayerUUID)
	r.Uint8(&pk.Status)
	r.String(&pk.XUID)
	r.String(&pk.Address)
}
