package sdp

import (
	"fmt"
	"strings"

	"github.com/pions/webrtc/pkg/rtcerr"

	"github.com/pkg/errors"
)

type Type struct {
	Value string
}

func (t *Type) Clone() Attribute {
	return &Type{Value: t.Value}
}

func (t *Type) Unmarshal(raw string) error {
	parts := strings.Split(raw, ":")
	if len(parts) != 2 {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", raw)}, pkgName)
	}

	t.Value = parts[1]
	return nil
}

func (t *Type) Marshal() string {
	return attributeKey + t.Name() + ":" + t.Value + endline
}

func (t *Type) Name() string {
	return AttributeNameType
}
