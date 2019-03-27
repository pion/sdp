package sdp

import (
	"fmt"
	"strings"

	"github.com/pions/webrtc/pkg/rtcerr"

	"github.com/pkg/errors"
)

// IcePwd is defined in https://tools.ietf.org/html/rfc5245.
type IcePwd struct {
	Value string
}

func (i *IcePwd) Clone() Attribute {
	return &IcePwd{Value: i.Value}
}

func (i *IcePwd) Unmarshal(raw string) error {
	parts := strings.Split(raw, ":")
	if len(parts) != 2 {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", raw)}, pkgName)
	}

	i.Value = parts[1]
	return nil
}

func (i *IcePwd) Marshal() string {
	return attributeKey + i.Name() + ":" + i.Value + endline
}

func (i *IcePwd) Name() string {
	return AttributeNameIcePwd
}
