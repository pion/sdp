package sdp

// CandidateType represents the type of the ICE candidate used.
type CandidateType int

const (
	// CandidateTypeHost indicates that the candidate is of Host type as
	// described in https://tools.ietf.org/html/rfc8445#section-5.1.1.1. A
	// candidate obtained by binding to a specific port from an IP address on
	// the host. This includes IP addresses on physical interfaces and logical
	// ones, such as ones obtained through VPNs.
	CandidateTypeHost CandidateType = iota + 1

	// CandidateTypeSrflx indicates the the candidate is of Server
	// Reflexive type as described
	// https://tools.ietf.org/html/rfc8445#section-5.1.1.2. A candidate type
	// whose IP address and port are a binding allocated by a NAT for an ICE
	// agent after it sends a packet through the NAT to a server, such as a
	// STUN server.
	CandidateTypeSrflx

	// CandidateTypePrflx indicates that the candidate is of Peer
	// Reflexive type. A candidate type whose IP address and port are a binding
	// allocated by a NAT for an ICE agent after it sends a packet through the
	// NAT to its peer.
	CandidateTypePrflx

	// CandidateTypeRelay indicates the the candidate is of Relay type as
	// described in https://tools.ietf.org/html/rfc8445#section-5.1.1.2. A
	// candidate type obtained from a relay server, such as a TURN server.
	CandidateTypeRelay
)

const (
	candidateTypeHostStr  = "host"
	candidateTypeSrflxStr = "srflx"
	candidateTypePrflxStr = "prflx"
	candidateTypeRelayStr = "relay"
)

// CandidateType defines a procedure for creating a new candidate type from a
// raw string.
func NewCandidateType(raw string) CandidateType {
	switch raw {
	case candidateTypeHostStr:
		return CandidateTypeHost
	case candidateTypeSrflxStr:
		return CandidateTypeSrflx
	case candidateTypePrflxStr:
		return CandidateTypePrflx
	case candidateTypeRelayStr:
		return CandidateTypeRelay
	default:
		return CandidateType(unknown)
	}
}

func (t CandidateType) String() string {
	switch t {
	case CandidateTypeHost:
		return candidateTypeHostStr
	case CandidateTypeSrflx:
		return candidateTypeSrflxStr
	case CandidateTypePrflx:
		return candidateTypePrflxStr
	case CandidateTypeRelay:
		return candidateTypeRelayStr
	default:
		return unknownStr
	}
}
