package sdp

type Direction int

const (
	DirectionSendRecv Direction = iota + 1
	DirectionSendOnly
	DirectionRecvOnly
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
