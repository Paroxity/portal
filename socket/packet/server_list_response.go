package packet

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

// ServerListResponse is sent by the proxy in response to ServerListRequest. It sends list of all
// the servers connected to the proxy.
type ServerListResponse struct {
	// Servers represents all the servers connected to the proxy.
	Servers []ServerEntry
}

// ServerEntry represents server connected the proxy.
type ServerEntry struct {
	// Name is name of the server.
	Name string
	// PlayerCount returns player count of the server.
	PlayerCount int64
}

// ID ...
func (*ServerListResponse) ID() uint16 {
	return IDServerListResponse
}

// Marshal ...
func (pk *ServerListResponse) Marshal(w *protocol.Writer) {
	l := uint32(len(pk.Servers))
	w.Uint32(&l)

	for _, s := range pk.Servers {
		w.String(&s.Name)
		w.Int64(&s.PlayerCount)
	}
}

// Unmarshal ...
func (pk *ServerListResponse) Unmarshal(r *protocol.Reader) {
	var l uint32
	r.Uint32(&l)

	pk.Servers = make([]ServerEntry, l)
	for i := uint32(0); i < l; i++ {
		r.String(&pk.Servers[i].Name)
		r.Int64(&pk.Servers[i].PlayerCount)
	}
}
