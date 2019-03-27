package sdp

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pions/webrtc/pkg/rtcerr"

	"github.com/pkg/errors"
)

type MaxPtime struct {
	Value uint64
}

func (m *MaxPtime) Clone() Attribute {
	return &MaxPtime{Value: m.Value}
}

func (m *MaxPtime) Unmarshal(raw string) error {
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

func (m *MaxPtime) Marshal() string {
	return attributeKey + m.Name() + ":" + fmt.Sprintf("%d", m.Value) + endline
}

func (m *MaxPtime) Name() string {
	return AttributeNameMaxPtime
}
