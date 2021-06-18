package packet

const (
	IDAuthRequest uint16 = iota
	IDAuthResponse
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
