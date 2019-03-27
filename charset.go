package sdp

import (
	"fmt"
	"strings"

	"github.com/pions/webrtc/pkg/rtcerr"

	"github.com/pkg/errors"
)

type Charset struct {
	Value string
}

func (c *Charset) Clone() Attribute {
	return &Charset{Value: c.Value}
}

func (c *Charset) Unmarshal(raw string) error {
	parts := strings.Split(raw, ":")
	if len(parts) != 2 {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", raw)}, pkgName)
	}

	c.Value = parts[1]
	return nil
}

func (c *Charset) Marshal() string {
	return attributeKey + c.Name() + ":" + c.Value + endline
}

func (c *Charset) Name() string {
	return AttributeNameCharset
}
