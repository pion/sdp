package sdp

//Direction is a marker for transmission directon of an endpoint
type Direction int

const (
	//DirectionSendRecv is for bidirectional communication
	DirectionSendRecv Direction = iota + 1
	//DirectionSendOnly is for outgoing communication
	DirectionSendOnly
	//DirectionRecvOnly is for incoming communication
	DirectionRecvOnly
	//DirectionInactive is for no communication
	DirectionInactive
)

const (
	directionSendRecvStr = "sendrecv"
	directionSendOnlyStr = "sendonly"
	directionRecvOnlyStr = "recvonly"
	directionInactiveStr = "inactive"
)

// NewDirection defines a procedure for creating a new direction from a raw
// string.
func NewDirection(raw string) Direction {
	switch raw {
	case directionSendRecvStr:
		return DirectionSendRecv
	case directionSendOnlyStr:
		return DirectionSendOnly
	case directionRecvOnlyStr:
		return DirectionRecvOnly
	case directionInactiveStr:
		return DirectionInactive
	default:
		return Direction(unknown)
	}
}

func (t Direction) String() string {
	switch t {
	case DirectionSendRecv:
		return directionSendRecvStr
	case DirectionSendOnly:
		return directionSendOnlyStr
	case DirectionRecvOnly:
		return directionRecvOnlyStr
	case DirectionInactive:
		return directionInactiveStr
	default:
		return unknownStr
	}
}
