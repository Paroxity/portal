package packet

// Register registers a function that returns a packet for a specific ID. Packets with this ID coming in from
// connections will resolve to the packet returned by the function passed.
func Register(id uint16, pk func() Packet) {
	registeredPackets[id] = pk
}

// registeredPackets holds packets registered by the user.
var registeredPackets = map[uint16]func() Packet{}

// Pool is a map holding packets indexed by a packet ID.
type Pool map[uint16]Packet

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
	packets := map[uint16]func() Packet{
		IDAuthRequest:         func() Packet { return &AuthRequest{} },
		IDAuthResponse:        func() Packet { return &AuthResponse{} },
		IDRegisterServer:      func() Packet { return &RegisterServer{} },
		IDTransferRequest:     func() Packet { return &TransferRequest{} },
		IDTransferResponse:    func() Packet { return &TransferResponse{} },
		IDPlayerInfoRequest:   func() Packet { return &PlayerInfoRequest{} },
		IDPlayerInfoResponse:  func() Packet { return &PlayerInfoResponse{} },
		IDServerListRequest:   func() Packet { return &ServerListRequest{} },
		IDServerListResponse:  func() Packet { return &ServerListResponse{} },
		IDFindPlayerRequest:   func() Packet { return &FindPlayerRequest{} },
		IDFindPlayerResponse:  func() Packet { return &FindPlayerResponse{} },
		IDUpdatePlayerLatency: func() Packet { return &UpdatePlayerLatency{} },
	}
	for id, pk := range packets {
		Register(id, pk)
	}
}
