package sdp

import (
	"fmt"
	"strings"

	"github.com/pions/webrtc/pkg/rtcerr"

	"github.com/pkg/errors"
)

type TimeZones []TimeZone

func (t *TimeZones) Clone() *TimeZones {
	timeZones := &TimeZones{}
	for _, timeZone := range *t {
		*timeZones = append(*timeZones, *timeZone.Clone())
	}
	return timeZones
}

func (t *TimeZones) Unmarshal(raw string) error {
	// These fields are transimitted in pairs
	// z=<adjustment time> <offset> <adjustment time> <offset> ....
	// so we are making sure that there are actually multiple of 2 total.
	fields := strings.Fields(raw)
	if len(fields)%2 != 0 {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("z=%v", fields)}, pkgName)
	}

	var timeZones TimeZones
	for i := 0; i < len(fields); i += 2 {
		timeZone := TimeZone{}
		if err := timeZone.Unmarshal(strings.Join(fields[i:i+2], " ")); err != nil {
			return err
		}
		timeZones = append(timeZones, timeZone)
	}

	*t = timeZones
	return nil
}

func (t *TimeZones) Marshal() string {
	return timeZonesKey + t.string() + endline
}

func (t *TimeZones) string() string {
	rawTimeZones := make([]string, 0)
	for _, z := range *t {
		rawTimeZones = append(rawTimeZones, z.Marshal())
	}
	return strings.Join(rawTimeZones, " ")
}
