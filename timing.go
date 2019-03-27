package sdp

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pions/webrtc/pkg/rtcerr"

	"github.com/pkg/errors"
)

// Timing defines the "t=" field's structured representation for the start and
// stop times.
type Timing struct {
	StartTime uint64
	StopTime  uint64
}

func (t *Timing) Unmarshal(raw string) error {
	fields := strings.Fields(raw)
	if len(fields) < 2 {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("t=%v", fields)}, pkgName)
	}

	startTime, err := strconv.ParseUint(fields[0], 10, 64)
	if err != nil {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", fields[1])}, pkgName)
	}

	stopTime, err := strconv.ParseUint(fields[1], 10, 64)
	if err != nil {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", fields[1])}, pkgName)
	}

	t.StartTime = startTime
	t.StopTime = stopTime
	return nil
}

func (t Timing) Marshal() string {
	return timingKey + t.String() + endline
}

func (t Timing) String() string {
	output := strconv.FormatUint(t.StartTime, 10)
	output += " " + strconv.FormatUint(t.StopTime, 10)
	return output
}
