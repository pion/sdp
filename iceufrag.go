package sdp

import (
	"fmt"
	"strings"

	"github.com/pions/webrtc/pkg/rtcerr"

	"github.com/pkg/errors"
)

// IceUfrag is defined in https://tools.ietf.org/html/rfc5245.
type IceUfrag struct {
	Value string
}

func (i *IceUfrag) Clone() Attribute {
	return &IceUfrag{Value: i.Value}
}

func (i *IceUfrag) Unmarshal(raw string) error {
	parts := strings.Split(raw, ":")
	if len(parts) != 2 {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", raw)}, pkgName)
	}

	i.Value = parts[1]
	return nil
}

func (i *IceUfrag) Marshal() string {
	return attributeKey + i.Name() + ":" + i.Value + endline
}

func (i *IceUfrag) Name() string {
	return AttributeNameIceUfrag
}
