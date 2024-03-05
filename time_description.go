// SPDX-FileCopyrightText: 2023 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

package sdp

import (
	"strconv"
	"strings"
)

// TimeDescription describes "t=", "r=" fields of the session description
// which are used to specify the start and stop times for a session as well as
// repeat intervals and durations for the scheduled session.
type TimeDescription struct {
	// t=<start-time> <stop-time>
	// https://tools.ietf.org/html/rfc4566#section-5.9
	Timing Timing

	// r=<repeat interval> <active duration> <offsets from start-time>
	// https://tools.ietf.org/html/rfc4566#section-5.10
	RepeatTimes []RepeatTime
}

// Timing defines the "t=" field's structured representation for the start and
// stop times.
type Timing struct {
	StartTime uint64
	StopTime  uint64
}

func (t Timing) String() string {
	output := strconv.FormatUint(t.StartTime, 10)
	output += " " + strconv.FormatUint(t.StopTime, 10)
	return output
}

func (t Timing) marshalSize() (size int) {
	return lenUint(t.StartTime) + 1 + lenUint(t.StopTime)
}

// RepeatTime describes the "r=" fields of the session description which
// represents the intervals and durations for repeated scheduled sessions.
type RepeatTime struct {
	Interval int64
	Duration int64
	Offsets  []int64
}

func (r RepeatTime) String() string {
	fields := make([]string, 0)
	fields = append(fields, strconv.FormatInt(r.Interval, 10))
	fields = append(fields, strconv.FormatInt(r.Duration, 10))
	for _, value := range r.Offsets {
		fields = append(fields, strconv.FormatInt(value, 10))
	}

	return strings.Join(fields, " ")
}

func (r RepeatTime) marshalSize() (size int) {
	size = lenInt(r.Interval)
	size += 1 + lenInt(r.Duration)
	for _, o := range r.Offsets {
		size += 1 + lenInt(o)
	}

	return
}
