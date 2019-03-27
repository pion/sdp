package sdp

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/pions/webrtc/pkg/rtcerr"

	"github.com/pkg/errors"
)

//ExtMap represents the activation of a single RTP header extension
type ExtMap struct {
	Value     int
	Direction Direction
	URI       *url.URL
	ExtAttr   *string
}

//Clone converts this object to an Attribute
func (e *ExtMap) Clone() Attribute {
	return Attribute{Key: "extmap", Value: e.string()}
}

//Unmarshal creates an Extmap from a string
func (e *ExtMap) Unmarshal(raw string) error {
	parts := strings.SplitN(raw, ":", 2)
	if len(parts) != 2 {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", raw)}, pkgName)
	}

	fields := strings.Fields(parts[1])
	if len(fields) < 2 {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", raw)}, pkgName)
	}

	valdir := strings.Split(fields[0], "/")
	value, err := strconv.ParseInt(valdir[0], 10, 64)
	if (value < 1) || (value > 246) {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", valdir[0])}, pkgName+": extmap key must be in the range 1-256")
	}
	if err != nil {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", valdir[0])}, pkgName)
	}

	var direction Direction
	if len(valdir) == 2 {
		direction, err = NewDirection(valdir[1])
		if err != nil {
			return errors.Wrap(&rtcerr.SyntaxError{Err: err}, pkgName)
		}
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

//Marshal creates a string from an ExtMap
func (e *ExtMap) Marshal() string {
	return attributeKey + e.Name() + ":" + e.string() + endline
}

func (e *ExtMap) string() string {
	output := fmt.Sprintf("%d", e.Value)
	dirstring := e.Direction.String()
	if dirstring != directionUnknownStr {
		output += "/" + dirstring
	}

	if e.URI != nil {
		output += " " + e.URI.String()
	}

	if e.ExtAttr != nil {
		output += " " + *e.ExtAttr
	}

	return output
}

//Name returns the constant name of this object
func (e *ExtMap) Name() string {
	return "extmap"
}
