package sdp

import (
	"fmt"
	"strings"

	"github.com/pions/webrtc/pkg/rtcerr"

	"github.com/pkg/errors"
)

type Simulcast struct {
	Value string
}

func (s *Simulcast) Clone() Attribute {
	return &RID{Value: s.Value}
}

func (s *Simulcast) Unmarshal(raw string) error {
	parts := strings.Split(raw, ":")
	if len(parts) != 2 {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", raw)}, pkgName)
	}

	s.Value = parts[1]
	return nil
}

func (s *Simulcast) Marshal() string {
	return attributeKey + s.Name() + ":" + s.Value + endline
}

func (s *Simulcast) Name() string {
	return AttributeNameSimulcast
}
