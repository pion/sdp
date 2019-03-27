package sdp

import (
	"fmt"
	"strings"

	"github.com/pions/webrtc/pkg/rtcerr"

	"github.com/pkg/errors"
)

type Rtcp struct {
	Value string
}

func (i *Rtcp) Clone() Attribute {
	return &Rtcp{Value: i.Value}
}

func (i *Rtcp) Unmarshal(raw string) error {
	parts := strings.SplitN(raw, ":", 2)
	if len(parts) != 2 {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", raw)}, pkgName)
	}

	i.Value = parts[1]
	return nil
}

func (i *Rtcp) Marshal() string {
	return attributeKey + i.Name() + ":" + i.Value + endline
}

func (i *Rtcp) Name() string {
	return AttributeNameRtcp
}
