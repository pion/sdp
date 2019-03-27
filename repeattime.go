package sdp

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pions/webrtc/pkg/rtcerr"

	"github.com/pkg/errors"
)

// RepeatTime describes the "r=" fields of the session description which
// represents the intervals and durations for repeated scheduled sessions.
type RepeatTime struct {
	Interval int64
	Duration int64
	Offsets  []int64
}

func (r *RepeatTime) Clone() *RepeatTime {
	repeatTime := &RepeatTime{}
	repeatTime.Interval = r.Interval
	repeatTime.Duration = r.Duration
	repeatTime.Offsets = append([]int64(nil), r.Offsets...)
	return repeatTime
}

func (r *RepeatTime) Unmarshal(raw string) error {
	fields := strings.Fields(raw)
	if len(fields) < 3 {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("r=%v", fields)}, pkgName)
	}

	interval, err := parseTimeUnits(fields[0])
	if err != nil {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", fields)}, pkgName)
	}

	duration, err := parseTimeUnits(fields[1])
	if err != nil {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", fields)}, pkgName)
	}

	var offsets []int64
	for i := 2; i < len(fields); i++ {
		offset, err := parseTimeUnits(fields[i])
		if err != nil {
			return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", fields)}, pkgName)
		}
		offsets = append(offsets, offset)
	}

	r.Interval = interval
	r.Duration = duration
	r.Offsets = offsets
	return nil
}

func (r *RepeatTime) Marshal() string {
	return repeatTimeKey + r.string() + endline
}

func (r *RepeatTime) string() string {
	fields := make([]string, 0)
	fields = append(fields, strconv.FormatInt(r.Interval, 10))
	fields = append(fields, strconv.FormatInt(r.Duration, 10))
	for _, value := range r.Offsets {
		fields = append(fields, strconv.FormatInt(value, 10))
	}

	output := strings.Join(fields, " ")
	return output
}
