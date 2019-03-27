package sdp

import (
	"fmt"
	"strings"

	"github.com/pions/webrtc/pkg/rtcerr"

	"github.com/pkg/errors"
)

type SctpPort struct {
	Value uint16
}

func (s *SctpPort) Clone() Attribute {
	return &SctpPort{Value: s.Value}
}

func (s *SctpPort) Unmarshal(raw string) error {
	parts := strings.Split(raw, ":")
	if len(parts) != 2 {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", raw)}, pkgName)
	}

	port, err := parsePort(parts[1])
	if err != nil {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", parts[1])}, pkgName)
	}

	s.Value = uint16(port)
	return nil
}

func (s *SctpPort) Marshal() string {
	return attributeKey + s.Name() + ":" + fmt.Sprintf("%d", s.Value) + endline
}

func (s *SctpPort) Name() string {
	return AttributeNameSctpPort
}
