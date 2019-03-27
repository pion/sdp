package sdp

type Semantic int

const (
	SemanticLP Semantic = iota + 1
	SemanticFID
	SemanticBUNDLE
	SemanticFEC
	SemanticWMS
)

const (
	semanticLPStr     = "LS"
	semanticFIDStr    = "FID"
	semanticBUNDLEStr = "BUNDLE"
	semanticFECStr    = "FEC"
	semanticWMSStr    = "WMS"
)

// NewSemantic defines a procedure for creating a new semantic token from a raw
// string.
func NewSemantic(raw string) Semantic {
	switch raw {
	case semanticLPStr:
		return SemanticLP
	case semanticFIDStr:
		return SemanticFID
	case semanticBUNDLEStr:
		return SemanticBUNDLE
	case semanticFECStr:
		return SemanticFEC
	case semanticWMSStr:
		return SemanticWMS
	default:
		return Semantic(unknown)
	}
}

func (t Semantic) String() string {
	switch t {
	case SemanticLP:
		return semanticLPStr
	case SemanticFID:
		return semanticFIDStr
	case SemanticBUNDLE:
		return semanticBUNDLEStr
	case SemanticFEC:
		return semanticFECStr
	case SemanticWMS:
		return semanticWMSStr
	default:
		return unknownStr
	}
}
