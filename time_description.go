// SPDX-FileCopyrightText: 2023 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

package sdp

import (
	"strconv"
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

func (t Timing) ByteLen() int {
	return uintLen(t.StartTime) + uintLen(t.StopTime) + 1
}

func (t Timing) MarshalAppend(b []byte) []byte {
	b = growByteSlice(b, t.ByteLen())
	b = strconv.AppendUint(b, t.StartTime, 10)
	b = append(b, ' ')
	b = strconv.AppendUint(b, t.StopTime, 10)
	return b
}

func (t Timing) String() string {
	output := strconv.FormatUint(t.StartTime, 10)
	output += " " + strconv.FormatUint(t.StopTime, 10)
	return output
}

// RepeatTime describes the "r=" fields of the session description which
// represents the intervals and durations for repeated scheduled sessions.
type RepeatTime struct {
	Interval int64
	Duration int64
	Offsets  []int64
}

func (r RepeatTime) ByteLen() int {
	n := uintLen(uint64(r.Interval)) + uintLen(uint64(r.Duration)) + 1
	for _, o := range r.Offsets {
		n += uintLen(uint64(o)) + 1
	}
	return n
}

func (r RepeatTime) MarshalAppend(b []byte) []byte {
	b = growByteSlice(b, r.ByteLen())
	b = strconv.AppendInt(b, r.Interval, 10)
	b = append(b, ' ')
	b = strconv.AppendInt(b, r.Duration, 10)
	for _, o := range r.Offsets {
		b = append(b, ' ')
		b = strconv.AppendUint(b, uint64(o), 10)
	}
	return b
}

func (r RepeatTime) String() string {
	return string(r.MarshalAppend(nil))
}
