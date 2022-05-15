package packet

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

// ServerListRequest is sent by the client to request list of all the servers
// connected to portal proxy (including offline servers).
type ServerListRequest struct{}

// ID ...
func (*ServerListRequest) ID() uint16 {
	return IDServerListRequest
}

// Marshal ...
func (pk *ServerListRequest) Marshal(*protocol.Writer) {}

// Unmarshal ...
func (pk *ServerListRequest) Unmarshal(*protocol.Reader) {}
