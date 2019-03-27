package sdp

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pions/webrtc/pkg/rtcerr"

	"github.com/pkg/errors"
)

type Ssrc struct {
	ID   int
	Attr string
}

func (s *Ssrc) Clone() Attribute {
	ssrc := &Ssrc{}
	ssrc.ID = s.ID
	ssrc.Attr = s.Attr
	return ssrc
}

func (s *Ssrc) Unmarshal(raw string) error {
	parts := strings.SplitN(raw, ":", 2)
	if len(parts) != 2 {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", raw)}, pkgName)
	}

	split := strings.SplitN(parts[1], " ", 2)
	if len(parts) != 2 {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", raw)}, pkgName)
	}

	id, err := strconv.ParseInt(split[0], 10, 64)
	if err != nil {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", split[0])}, pkgName)
	}

	s.ID = int(id)
	s.Attr = split[1]
	return nil
}

func (s *Ssrc) Marshal() string {
	return attributeKey + s.Name() + ":" + fmt.Sprintf("%d %v", s.ID, s.Attr) + endline
}

func (s *Ssrc) Name() string {
	return AttributeNameSsrc
}
