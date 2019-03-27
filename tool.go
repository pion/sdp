package sdp

import (
	"fmt"
	"strings"

	"github.com/pions/webrtc/pkg/rtcerr"

	"github.com/pkg/errors"
)

type Tool struct {
	Value string
}

func (t *Tool) Clone() Attribute {
	return &Tool{Value: t.Value}
}

func (t *Tool) Unmarshal(raw string) error {
	parts := strings.Split(raw, ":")
	if len(parts) != 2 {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", raw)}, pkgName)
	}

	t.Value = parts[1]
	return nil
}

func (t *Tool) Marshal() string {
	return attributeKey + t.Name() + ":" + t.Value + endline
}

func (t *Tool) Name() string {
	return AttributeNameTool
}
