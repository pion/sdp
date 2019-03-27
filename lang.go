package sdp

import (
	"fmt"
	"strings"

	"github.com/pions/webrtc/pkg/rtcerr"

	"github.com/pkg/errors"
)

type Lang struct {
	Value string
}

func (l *Lang) Clone() Attribute {
	return &Lang{Value: l.Value}
}

func (l *Lang) Unmarshal(raw string) error {
	parts := strings.Split(raw, ":")
	if len(parts) != 2 {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", raw)}, pkgName)
	}

	l.Value = parts[1]
	return nil
}

func (l *Lang) Marshal() string {
	return attributeKey + l.Name() + ":" + l.Value + endline
}

func (l *Lang) Name() string {
	return AttributeNameLang
}
