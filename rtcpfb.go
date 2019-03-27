package sdp

import (
	"fmt"
	"strings"

	"github.com/pions/webrtc/pkg/rtcerr"

	"github.com/pkg/errors"
)

type RtcpFb struct {
	Value string
}

func (r *RtcpFb) Clone() Attribute {
	return &RtcpFb{Value: r.Value}
}

func (r *RtcpFb) Unmarshal(raw string) error {
	parts := strings.Split(raw, ":")
	if len(parts) != 2 {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", raw)}, pkgName)
	}

	r.Value = parts[1]
	return nil
}

func (r *RtcpFb) Marshal() string {
	return attributeKey + r.Name() + ":" + r.Value + endline
}

func (r *RtcpFb) Name() string {
	return AttributeNameRtcpFb
}
