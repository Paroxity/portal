package packet

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

// ServerListResponse is sent by the proxy in response to ServerListResponse. It sends list of all
// the servers connected to the proxy (including offline servers)
type ServerListResponse struct {
	// Servers represents all the servers connected to the proxy
	Servers []ServerEntry
}

type ServerEntry struct {
	// Name is name of the server
	Name        string
	// Group is group of the server
	Group       string
	// IsOnline returns if the server is currently online
	IsOnline    bool
	// PlayerCount is count of online players connected to that server through proxy
	PlayerCount uint16
}

// ID ...
func (*ServerListResponse) ID() uint16 {
	return IDServerListResponse
}

// Marshal ...
func (pk *ServerListResponse) Marshal(w *protocol.Writer) {
	l := uint16(len(pk.Servers))
	w.Uint16(&l)

	for _, s := range pk.Servers {
		w.String(&s.Name)
		w.String(&s.Group)
		w.Bool(&s.IsOnline)
		w.Uint16(&s.PlayerCount)
	}
}

// Unmarshal ...
func (pk *ServerListResponse) Unmarshal(r *protocol.Reader) {
	var l uint16
	r.Uint16(&l)

	pk.Servers = make([]ServerEntry, l)
	for i := uint16(0); i < l; i++ {
		entry := ServerEntry{}
		r.String(&entry.Name)
		r.String(&entry.Group)
		r.Bool(&entry.IsOnline)
		r.Uint16(&entry.PlayerCount)

		pk.Servers[i] = entry
	}
}

