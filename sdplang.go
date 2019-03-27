package sdp

import (
	"fmt"
	"strings"

	"github.com/pions/webrtc/pkg/rtcerr"

	"github.com/pkg/errors"
)

//noinspection GoNameStartsWithPackageName
type SdpLang struct {
	Value string
}

func (s *SdpLang) Clone() Attribute {
	return &SdpLang{Value: s.Value}
}

func (s *SdpLang) Unmarshal(raw string) error {
	parts := strings.Split(raw, ":")
	if len(parts) != 2 {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", raw)}, pkgName)
	}

	s.Value = parts[1]
	return nil
}

func (s *SdpLang) Marshal() string {
	return attributeKey + s.Name() + ":" + s.Value + endline
}

func (s *SdpLang) Name() string {
	return AttributeNameSdpLang
}
