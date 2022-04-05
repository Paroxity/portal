package packet

// ProtocolVersion is the protocol version supported by the proxy. It will only accept clients that match this version,
// and it should be incremented every time the protocol changes.
const ProtocolVersion = 1

const (
	IDAuthRequest uint16 = iota
	IDAuthResponse
	IDRegisterServer
	IDTransferRequest
	IDTransferResponse
	IDPlayerInfoRequest
	IDPlayerInfoResponse
	IDServerListRequest
	IDServerListResponse
	IDFindPlayerRequest
	IDFindPlayerResponse
	IDUpdatePlayerLatency
)
