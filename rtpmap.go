package sdp

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pions/webrtc/pkg/rtcerr"

	"github.com/pkg/errors"
)

type RtpMap struct {
	Payload   int
	Encoding  string
	ClockRate int
	Params    []int
}

func (r *RtpMap) Clone() Attribute {
	rtpmap := &RtpMap{}
	rtpmap.Payload = r.Payload
	rtpmap.Encoding = r.Encoding
	rtpmap.ClockRate = r.ClockRate
	rtpmap.Params = append([]int(nil), r.Params...)
	return rtpmap
}

func (r *RtpMap) Unmarshal(raw string) error {
	parts := strings.Split(raw, ":")
	if len(parts) != 2 {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", raw)}, pkgName)
	}

	split := strings.Fields(parts[1])
	if len(split) != 2 {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", raw)}, pkgName)
	}

	payload, err := strconv.ParseInt(split[0], 10, 64)
	if err != nil {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", split[0])}, pkgName)
	}

	extras := strings.Split(split[1], "/")
	if len(extras) < 2 {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", raw)}, pkgName)
	}

	clockRate, err := strconv.ParseInt(extras[1], 10, 64)
	if err != nil {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", split[0])}, pkgName)
	}

	var params []int
	for _, param := range extras[2:] {
		value, err := strconv.ParseInt(param, 10, 64)
		if err != nil {
			return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", param)}, pkgName)
		}
		params = append(params, int(value))
	}

	r.Payload = int(payload)
	r.Encoding = extras[0]
	r.ClockRate = int(clockRate)
	r.Params = params
	return nil
}

func (r *RtpMap) Marshal() string {
	return attributeKey + r.Name() + ":" + r.string() + endline
}

func (r *RtpMap) string() string {
	var extras []string
	extras = append(extras, r.Encoding)
	extras = append(extras, strconv.Itoa(r.ClockRate))
	for _, param := range r.Params {
		extras = append(extras, strconv.Itoa(param))
	}

	return fmt.Sprintf(
		"%d %v",
		r.Payload,
		strings.Join(extras, "/"),
	)
}

func (r *RtpMap) Name() string {
	return AttributeNameRtpMap
}
