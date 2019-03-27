package sdp

import (
	"fmt"
	"strings"

	"github.com/pions/webrtc/pkg/rtcerr"

	"github.com/pkg/errors"
)

type MsID struct {
	Value string
}

func (m *MsID) Clone() Attribute {
	return &MsID{Value: m.Value}
}

func (m *MsID) Unmarshal(raw string) error {
	parts := strings.Split(raw, ":")
	if len(parts) != 2 {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", raw)}, pkgName)
	}

	m.Value = parts[1]
	return nil
}

func (m *MsID) Marshal() string {
	return attributeKey + m.Name() + ":" + m.Value + endline
}

func (m *MsID) Name() string {
	return AttributeNameMsID
}
