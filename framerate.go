package sdp

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pions/webrtc/pkg/rtcerr"

	"github.com/pkg/errors"
)

type Framerate struct {
	Value float64
}

func (f *Framerate) Clone() Attribute {
	return &Framerate{Value: f.Value}
}

func (f *Framerate) Unmarshal(raw string) error {
	parts := strings.Split(raw, ":")
	if len(parts) != 2 {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", raw)}, pkgName)
	}

	value, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", parts[1])}, pkgName)
	}

	f.Value = value
	return nil
}

func (f *Framerate) Marshal() string {
	return attributeKey + f.Name() + ":" + f.string() + endline
}

func (f *Framerate) string() string {
	return strconv.FormatFloat(f.Value, 'f', 1, 64)
}

func (f *Framerate) Name() string {
	return AttributeNameFramerate
}
