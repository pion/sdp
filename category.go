package sdp

import (
	"fmt"
	"strings"

	"github.com/pions/webrtc/pkg/rtcerr"

	"github.com/pkg/errors"
)

type Category struct {
	Path []string
}

func (c *Category) Clone() Attribute {
	cat := &Category{}
	cat.Path = append([]string(nil), c.Path...)
	return cat
}

func (c *Category) Unmarshal(raw string) error {
	parts := strings.Split(raw, ":")
	if len(parts) != 2 {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", raw)}, pkgName)
	}

	c.Path = strings.Split(parts[1], ".")
	return nil
}

func (c *Category) Marshal() string {
	return attributeKey + c.Name() + ":" + strings.Join(c.Path, ".") + endline
}

func (c *Category) Name() string {
	return AttributeNameCategory
}
