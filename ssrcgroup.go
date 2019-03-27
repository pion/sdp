package sdp

import (
	"fmt"
	"strings"

	"github.com/pions/webrtc/pkg/rtcerr"

	"github.com/pkg/errors"
)

// SsrcGroup is defined in https://tools.ietf.org/html/rfc5576.
type SsrcGroup struct {
	Semantic Semantic
	SSRCs    []string
}

func (s *SsrcGroup) Clone() Attribute {
	group := &SsrcGroup{}
	group.Semantic = s.Semantic
	group.SSRCs = append([]string(nil), s.SSRCs...)
	return group
}

func (s *SsrcGroup) Unmarshal(raw string) error {
	parts := strings.Split(raw, ":")
	if len(parts) != 2 {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", raw)}, pkgName)
	}

	split := strings.Fields(parts[1])
	if len(parts) < 1 {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", raw)}, pkgName)
	}

	semantic := NewSemantic(split[0])
	if semantic != SemanticFEC && semantic != SemanticFID {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", split[0])}, pkgName)
	}

	s.Semantic = semantic
	s.SSRCs = append(s.SSRCs, split[1:]...)
	return nil
}

func (s *SsrcGroup) Marshal() string {
	return attributeKey + s.Name() + ":" + s.string() + endline
}

func (s *SsrcGroup) string() string {
	SSRCs := strings.Join(s.SSRCs, " ")
	if len(s.SSRCs) > 0 {
		SSRCs = " " + SSRCs
	}

	return fmt.Sprintf(
		"%v%v",
		s.Semantic.String(),
		SSRCs,
	)
}

func (s *SsrcGroup) Name() string {
	return AttributeNameSsrcGroup
}

// FirstSSRC returns the first ssrc-id
func (s *SsrcGroup) FirstSSRC() *string {
	if len(s.SSRCs) > 0 {
		tmp := s.SSRCs[0]
		return &tmp
	}
	return nil
}

func (s *SsrcGroup) HasSSRC(ssrc string) int {
	for i, each := range s.SSRCs {
		if each == ssrc {
			return i
		}
	}
	return -1
}

func (s *SsrcGroup) AddSSRC(ssrc string) int {
	if i := s.HasSSRC(ssrc); i != -1 {
		return i
	}

	i := len(s.SSRCs)
	s.SSRCs = append(s.SSRCs, ssrc)
	return i
}

func (s *SsrcGroup) RemoveSSRC(ssrc string) bool {
	if i := s.HasSSRC(ssrc); i != -1 {
		copy(s.SSRCs[i:], s.SSRCs[i+1:])
		s.SSRCs[len(s.SSRCs)-1] = ""
		s.SSRCs = s.SSRCs[:len(s.SSRCs)-1]
		return true
	}
	return false
}
