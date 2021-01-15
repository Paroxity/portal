package packet

import "github.com/sandertv/gophertunnel/minecraft/protocol"

const (
	TransferResponseSuccess = iota
	TransferResponseGroupNotFound
	TransferResponseServerNotFound
	TransferResponseAlreadyOnServer
	TransferResponsePlayerNotFound
)

// TransferResponse is sent by the proxy in response to a transfer request.
type TransferResponse struct {
	// PlayerRuntimeID is the entity runtime ID of the player being transferred.
	PlayerRuntimeID uint64
	// Status is the response status from the transfer. The possible values for this can be found above.
	Status byte
	// Reason is the reason behind the Status provided.
	Reason string
}

// ID ...
func (*TransferResponse) ID() uint32 {
	return IDTransferResponse
}

// Marshal ...
func (pk *TransferResponse) Marshal(w *protocol.Writer) {
	w.Varuint64(&pk.PlayerRuntimeID)
	w.Uint8(&pk.Status)
	w.String(&pk.Reason)
}

// Unmarshal ...
func (pk *TransferResponse) Unmarshal(r *protocol.Reader) {
	r.Varuint64(&pk.PlayerRuntimeID)
	r.Uint8(&pk.Status)
	r.String(&pk.Reason)
}
