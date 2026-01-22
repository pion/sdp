// SPDX-FileCopyrightText: 2026 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTiming_String(t *testing.T) {
	assert.Equal(t, "0 0", Timing{StartTime: 0, StopTime: 0}.String())
	assert.Equal(t, "12345 67890", Timing{StartTime: 12345, StopTime: 67890}.String())
}

func TestRepeatTime_String(t *testing.T) {
	assert.Equal(t, "3600 900", RepeatTime{Interval: 3600, Duration: 900}.String())
	assert.Equal(
		t,
		"604800 3600 -60 0 60",
		RepeatTime{Interval: 604800, Duration: 3600, Offsets: []int64{-60, 0, 60}}.String(),
	)
}
