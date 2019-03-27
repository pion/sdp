package sdp

// Protocol indicates the transport protocol type of the ICE candidate.
type Protocol int

const (
	// ProtocolUDP indicates UDP transport.
	ProtocolUDP Protocol = iota + 1

	// ProtocolTCP indicates a TCP transport.
	ProtocolTCP
)

const (
	protocolUDPStr = "UDP"
	protocolTCPStr = "TCP"
)

// NewProtocol defines a procedure for creating a new protocol from a raw
// string.
func NewProtocol(raw string) Protocol {
	switch raw {
	case protocolUDPStr:
		return ProtocolUDP
	case protocolTCPStr:
		return ProtocolTCP
	default:
		return Protocol(unknown)
	}
}

func (t Protocol) String() string {
	switch t {
	case ProtocolUDP:
		return protocolUDPStr
	case ProtocolTCP:
		return protocolTCPStr
	default:
		return unknownStr
	}
}
