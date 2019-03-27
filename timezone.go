package sdp

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pions/webrtc/pkg/rtcerr"

	"github.com/pkg/errors"
)

// TimeZone defines the structured object for "z=" line which describes
// repeated sessions scheduling.
type TimeZone struct {
	AdjustmentTime uint64
	Offset         int64
}

func (z *TimeZone) Clone() *TimeZone {
	return &TimeZone{
		AdjustmentTime: z.AdjustmentTime,
		Offset:         z.Offset,
	}
}

func (z *TimeZone) Unmarshal(raw string) error {
	fields := strings.Fields(raw)
	if len(fields)%2 != 0 {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("z=%v", fields)}, pkgName)
	}

	adjustmentTime, err := strconv.ParseUint(fields[0], 10, 64)
	if err != nil {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", fields[0])}, pkgName)
	}

	offset, err := parseTimeUnits(fields[1])
	if err != nil {
		return err
	}

	z.AdjustmentTime = adjustmentTime
	z.Offset = offset
	return nil
}

func (z *TimeZone) Marshal() string {
	return fmt.Sprintf("%d %d", z.AdjustmentTime, z.Offset)
}
