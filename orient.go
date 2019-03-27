package sdp

import (
	"fmt"
	"strings"

	"github.com/pions/webrtc/pkg/rtcerr"

	"github.com/pkg/errors"
)

type Orient struct {
	Value Rotation
}

func (o *Orient) Clone() Attribute {
	return &Orient{Value: o.Value}
}

func (o *Orient) Unmarshal(raw string) error {
	parts := strings.Split(raw, ":")
	if len(parts) != 2 {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", raw)}, pkgName)
	}

	rotation := NewRotation(strings.TrimSpace(parts[1]))
	if rotation == Rotation(unknown) {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", raw)}, pkgName)
	}

	o.Value = rotation
	return nil
}

func (o *Orient) Marshal() string {
	return attributeKey + o.Name() + ":" + o.Value.String() + endline
}

func (o *Orient) Name() string {
	return AttributeNameOrient
}
