package sdp

type Rotation int

const (
	RotationPortrait Rotation = iota + 1
	RotationLandscape
	RotationSeascape
)

const (
	rotationPortraitStr  = "portrait"
	rotationLandscapeStr = "landscape"
	rotationSeascapeStr  = "seascape"
)

// NewRotation defines a procedure for creating a new rotation from a raw
// string.
func NewRotation(raw string) Rotation {
	switch raw {
	case rotationPortraitStr:
		return RotationPortrait
	case rotationLandscapeStr:
		return RotationLandscape
	case rotationSeascapeStr:
		return RotationSeascape
	default:
		return Rotation(unknown)
	}
}

func (o Rotation) String() string {
	switch o {
	case RotationPortrait:
		return rotationPortraitStr
	case RotationLandscape:
		return rotationLandscapeStr
	case RotationSeascape:
		return rotationSeascapeStr
	default:
		return unknownStr
	}
}
