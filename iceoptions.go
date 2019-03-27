package sdp

import (
	"fmt"
	"strings"

	"github.com/pions/webrtc/pkg/rtcerr"

	"github.com/pkg/errors"
)

// IceOptions is defined in https://tools.ietf.org/html/rfc5245.
type IceOptions struct {
	Value string
}

func (i *IceOptions) Clone() Attribute {
	return &IceOptions{Value: i.Value}
}

func (i *IceOptions) Unmarshal(raw string) error {
	parts := strings.Split(raw, ":")
	if len(parts) != 2 {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", raw)}, pkgName)
	}

	i.Value = parts[1]
	return nil
}

func (i *IceOptions) Marshal() string {
	return attributeKey + i.Name() + ":" + i.Value + endline
}

func (i *IceOptions) Name() string {
	return AttributeNameIceOptions
}
