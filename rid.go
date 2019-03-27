package sdp

import (
	"fmt"
	"strings"

	"github.com/pions/webrtc/pkg/rtcerr"

	"github.com/pkg/errors"
)

type RID struct {
	Value string
}

func (r *RID) Clone() Attribute {
	return &RID{Value: r.Value}
}

func (r *RID) Unmarshal(raw string) error {
	parts := strings.Split(raw, ":")
	if len(parts) != 2 {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", raw)}, pkgName)
	}

	r.Value = parts[1]
	return nil
}

func (r *RID) Marshal() string {
	return attributeKey + r.Name() + ":" + r.Value + endline
}

func (r *RID) Name() string {
	return AttributeNameRID
}
