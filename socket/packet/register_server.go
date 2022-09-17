package packet

import "github.com/sandertv/gophertunnel/minecraft/protocol"

// RegisterServer is sent by a connection to register itself as a server with the provided address.
type RegisterServer struct {
	// Address is the address of the server in the format ip:port.
	Address string
	// UseRakNet is if a connection should use the RakNet protocol when connecting to the server.
	UseRakNet bool
}

// ID ...
func (pk *RegisterServer) ID() uint16 {
	return IDRegisterServer
}

// Marshal ...
func (pk *RegisterServer) Marshal(w *protocol.Writer) {
	w.String(&pk.Address)
	w.Bool(&pk.UseRakNet)
}

// Unmarshal ...
func (pk *RegisterServer) Unmarshal(r *protocol.Reader) {
	r.String(&pk.Address)
	r.Bool(&pk.UseRakNet)
}
