package sdp

import (
	"fmt"
	"strconv"

	"github.com/pions/webrtc/pkg/rtcerr"
	"github.com/pkg/errors"
)

// Version describes the value provided by the "v=" field which gives
// the version of the Session Description protocol.
type Version struct {
	Value int
}

func (v *Version) Clone() *Version {
	return &Version{Value: v.Value}
}

func (v *Version) Unmarshal(raw string) error {
	version, err := strconv.ParseInt(raw, 10, 32)
	if err != nil {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", version)}, pkgName)
	}

	// As off the latest draft of the rfc this value is required to be 0.
	// https://tools.ietf.org/html/draft-ietf-rtcweb-jsep-24#section-5.8.1
	if version != 0 {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", version)}, pkgName)
	}

	v.Value = int(version)
	return nil
}

func (v *Version) Marshal() string {
	return versionKey + strconv.Itoa(v.Value) + endline
}
