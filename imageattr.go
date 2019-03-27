package sdp

import (
	"fmt"
	"strings"

	"github.com/pions/webrtc/pkg/rtcerr"

	"github.com/pkg/errors"
)

type ImageAttr struct {
	Value string
}

func (i *ImageAttr) Clone() Attribute {
	return &ImageAttr{Value: i.Value}
}

func (i *ImageAttr) Unmarshal(raw string) error {
	parts := strings.Split(raw, ":")
	if len(parts) != 2 {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", raw)}, pkgName)
	}

	i.Value = parts[1]
	return nil
}

func (i *ImageAttr) Marshal() string {
	return attributeKey + i.Name() + ":" + i.Value + endline
}

func (i *ImageAttr) Name() string {
	return AttributeNameImageAttr
}
