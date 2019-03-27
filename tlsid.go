package sdp

import (
	"fmt"
	"github.com/pions/webrtc/pkg/rtcerr"
	"strings"

	"github.com/pkg/errors"
)

type TlsID struct {
	Value string
}

func (t *TlsID) Clone() Attribute {
	return &Lang{Value: t.Value}
}

func (t *TlsID) Unmarshal(raw string) error {
	parts := strings.Split(raw, ":")
	if len(parts) != 2 {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", raw)}, pkgName)
	}

	t.Value = parts[1]
	return nil
}

func (t *TlsID) Marshal() string {
	return attributeKey + t.Name() + ":" + t.Value + endline
}

func (t *TlsID) Name() string {
	return AttributeNameTlsID
}
