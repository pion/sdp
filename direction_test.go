package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDirection(t *testing.T) {
	passingtests := []struct {
		value    string
		expected Direction
	}{
		{"sendrecv", DirectionSendRecv},
		{"sendonly", DirectionSendOnly},
		{"recvonly", DirectionRecvOnly},
		{"inactive", DirectionInactive},
	}
	failingtests := []string{
		"",
		"notadirection",
	}

	for i, u := range passingtests {
		dir, err := NewDirection(u.value)
		assert.NoError(t, err)
		assert.Equal(t, u.expected, dir, "%d: %+v", i, u)
	}
	for _, u := range failingtests {
		_, err := NewDirection(u)
		assert.Error(t, err)
	}
}

func TestDirection_String(t *testing.T) {
	tests := []struct {
		actual   Direction
		expected string
	}{
		{Direction(unknown), directionUnknownStr},
		{DirectionSendRecv, "sendrecv"},
		{DirectionSendOnly, "sendonly"},
		{DirectionRecvOnly, "recvonly"},
		{DirectionInactive, "inactive"},
	}

	for i, u := range tests {
		assert.Equal(t, u.expected, u.actual.String(), "%d: %+v", i, u)
	}
}
