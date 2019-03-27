package sdp

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pions/webrtc/pkg/rtcerr"

	"github.com/pkg/errors"
)

type MaxMessageSize struct {
	Value uint64
}

func (m *MaxMessageSize) Clone() Attribute {
	return &MaxMessageSize{Value: m.Value}
}

func (m *MaxMessageSize) Unmarshal(raw string) error {
	parts := strings.Split(raw, ":")
	if len(parts) != 2 {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", raw)}, pkgName)
	}

	value, err := strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", parts[1])}, pkgName)
	}

	m.Value = value
	return nil
}

func (m *MaxMessageSize) Marshal() string {
	return attributeKey + m.Name() + ":" + fmt.Sprintf("%d", m.Value) + endline
}

func (m *MaxMessageSize) Name() string {
	return AttributeNameMaxMessageSize
}
