package packet

import (
	"github.com/google/uuid"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

const (
	TransferResponseSuccess byte = iota
	TransferResponseServerNotFound
	TransferResponseAlreadyOnServer
	TransferResponsePlayerNotFound
	TransferResponseError
)

// TransferResponse is sent by the proxy in response to a transfer request.
type TransferResponse struct {
	// PlayerUUID is the UUID of the player being transferred.
	PlayerUUID uuid.UUID
	// Status is the response status from the transfer. The possible values for this can be found above.
	Status byte
	// Error is the error message when the Status field is TransferResponseError.
	Error string
}

// ID ...
func (*TransferResponse) ID() uint16 {
	return IDTransferResponse
}

// Marshal ...
func (pk *TransferResponse) Marshal(w *protocol.Writer) {
	w.UUID(&pk.PlayerUUID)
	w.Uint8(&pk.Status)
	if pk.Status == TransferResponseError {
		w.String(&pk.Error)
	}
}

// Unmarshal ...
func (pk *TransferResponse) Unmarshal(r *protocol.Reader) {
	r.UUID(&pk.PlayerUUID)
	r.Uint8(&pk.Status)
	if pk.Status == TransferResponseError {
		r.String(&pk.Error)
	}
}
