package packet

import (
	"github.com/google/uuid"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

const (
	PlayerInfoResponseSuccess byte = iota
	PlayerInfoResponsePlayerNotFound
)

// PlayerInfoResponse is sent by the proxy in response to PlayerInfoRequest to tell the connection the XUID
// and IP address of the requested player.
type PlayerInfoResponse struct {
	// PlayerUUID is the UUID of the player the information belongs to.
	PlayerUUID uuid.UUID
	// Status is the response status from fetching the player information. The possible values for this can
	// be found above.
	Status byte
	// XUID is the Xbox Unique Identifier of the requested player.
	XUID string
	// Address is the IP address of the requested player. This can be IPv4 or IPv6 depending on which address
	// they join with.
	Address string
}

// ID ...
func (*PlayerInfoResponse) ID() uint16 {
	return IDPlayerInfoResponse
}

// Marshal ...
func (pk *PlayerInfoResponse) Marshal(w *protocol.Writer) {
	w.UUID(&pk.PlayerUUID)
	w.Uint8(&pk.Status)
	w.String(&pk.XUID)
	w.String(&pk.Address)
}

// Unmarshal ...
func (pk *PlayerInfoResponse) Unmarshal(r *protocol.Reader) {
	r.UUID(&pk.PlayerUUID)
	r.Uint8(&pk.Status)
	r.String(&pk.XUID)
	r.String(&pk.Address)
}
