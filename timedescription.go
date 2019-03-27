package sdp

// TimeDescription describes "t=", "r=" fields of the session description
// which are used to specify the start and stop times for a session as well as
// repeat intervals and durations for the scheduled session.
//
// https://tools.ietf.org/html/rfc4566#section-5.9
//
// https://tools.ietf.org/html/rfc4566#section-5.10
type TimeDescription struct {
	// t=<start-time> <stop-time>
	Timing Timing

	// r=<repeat interval> <active duration> <offsets from start-time>
	RepeatTimes RepeatTimes
}

func (t *TimeDescription) Clone() *TimeDescription {
	timeDesc := &TimeDescription{}
	timeDesc.Timing = t.Timing

	if len(t.RepeatTimes) > 0 {
		timeDesc.RepeatTimes = *t.RepeatTimes.Clone()
	}

	return timeDesc
}
