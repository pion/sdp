package sdp

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/pions/webrtc/pkg/rtcerr"

	"github.com/pkg/errors"
)

type ExtMap struct {
	Value     int
	Direction Direction
	URI       *url.URL
	ExtAttr   *string
}

func (e *ExtMap) Clone() Attribute {
	return Attribute{Key: "extmap", Value: e.string()}
}

func (e *ExtMap) Unmarshal(raw string) error {
	parts := strings.SplitN(raw, ":", 2)
	if len(parts) != 2 {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", raw)}, pkgName)
	}

	fields := strings.Fields(parts[1])
	if len(parts) < 2 {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", raw)}, pkgName)
	}

	valdir := strings.Split(fields[0], "/")
	value, err := strconv.ParseInt(valdir[0], 10, 64)
	if err != nil {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", valdir[0])}, pkgName)
	}

	var direction Direction
	if len(valdir) == 2 {
		direction = NewDirection(valdir[1])
	}

	uri, err := url.Parse(fields[1])
	if err != nil {
		return errors.Wrap(&rtcerr.UnknownError{Err: err}, pkgName)
	}

	if len(fields) == 3 {
		tmp := fields[2]
		e.ExtAttr = &tmp
	}

	e.Value = int(value)
	e.Direction = direction
	e.URI = uri
	return nil
}

func (e *ExtMap) Marshal() string {
	return attributeKey + e.Name() + ":" + e.string() + endline
}

func (e *ExtMap) string() string {
	output := fmt.Sprintf("%d", e.Value)
	if e.Direction != Direction(unknown) {
		output += "/" + e.Direction.String()
	}

	output += " " + e.URI.String()
	if e.ExtAttr != nil {
		output += " " + *e.ExtAttr
	}

	return output
}

func (e *ExtMap) Name() string {
	return "extmap"
}
