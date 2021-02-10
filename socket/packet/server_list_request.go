package packet

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

type ServerListRequest struct {}

// ID ...
func (*ServerListRequest) ID() uint16 {
	return IDServerListRequest
}

// Marshal ...
func (pk *ServerListRequest) Marshal(w *protocol.Writer) {}

// Unmarshal ...
func (pk *ServerListRequest) Unmarshal(r *protocol.Reader) {}
