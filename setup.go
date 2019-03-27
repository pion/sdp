package sdp

import (
	"fmt"
	"strings"

	"github.com/pions/webrtc/pkg/rtcerr"

	"github.com/pkg/errors"
)

type Setup struct {
	Value ConnectionRole
}

func (s *Setup) Clone() Attribute {
	return &Setup{Value: s.Value}
}

func (s *Setup) Unmarshal(raw string) error {
	parts := strings.Split(raw, ":")
	if len(parts) != 2 {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", raw)}, pkgName)
	}

	role := NewConnectionRole(strings.TrimSpace(parts[1]))
	if role == ConnectionRole(unknown) {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", raw)}, pkgName)
	}

	s.Value = role
	return nil
}

func (s *Setup) Marshal() string {
	return attributeKey + s.Name() + ":" + s.Value.String() + endline
}

func (s *Setup) Name() string {
	return AttributeNameSetup
}
