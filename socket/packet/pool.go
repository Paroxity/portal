package packet

// Register registers a function that returns a packet for a specific ID. Packets with this ID coming in from
// connections will resolve to the packet returned by the function passed.
func Register(id uint32, pk func() Packet) {
	registeredPackets[id] = pk
}

// registeredPackets holds packets registered by the user.
var registeredPackets = map[uint32]func() Packet{}

// Pool is a map holding packets indexed by a packet ID.
type Pool map[uint32]Packet

// NewPool returns a new pool with all supported packets sent. Packets may be retrieved from it simply by
// indexing it with the packet ID.
func NewPool() Pool {
	p := Pool{}
	for id, pk := range registeredPackets {
		p[id] = pk()
	}
	return p
}

func init() {
	packets := map[uint32]func() Packet{
		IDAuthRequest:        func() Packet { return &AuthRequest{} },
		IDAuthResponse:       func() Packet { return &AuthResponse{} },
		IDTransferRequest:    func() Packet { return &TransferRequest{} },
		IDTransferResponse:   func() Packet { return &TransferResponse{} },
		IDPlayerInfoRequest:  func() Packet { return &PlayerInfoRequest{} },
		IDPlayerInfoResponse: func() Packet { return &PlayerInfoResponse{} },
	}
	for id, pk := range packets {
		Register(id, pk)
	}
}
