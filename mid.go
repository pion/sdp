package sdp

import (
	"fmt"
	"strings"

	"github.com/pions/webrtc/pkg/rtcerr"

	"github.com/pkg/errors"
)

// MID is defined in https://tools.ietf.org/html/rfc5888.
type MID struct {
	Value string
}

func (m *MID) Clone() Attribute {
	return &MID{Value: m.Value}
}

func (m *MID) Unmarshal(raw string) error {
	parts := strings.Split(raw, ":")
	if len(parts) != 2 {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", raw)}, pkgName)
	}

	m.Value = parts[1]
	return nil
}

func (m MID) Marshal() string {
	return attributeKey + m.Name() + ":" + m.Value + endline
}

func (m *MID) Name() string {
	return AttributeNameMID
}
