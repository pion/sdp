package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDirection(t *testing.T) {
	tests := []struct {
		value    string
		expected Direction
	}{
		{"unknown", Direction(unknown)},
		{"sendrecv", DirectionSendRecv},
		{"sendonly", DirectionSendOnly},
		{"recvonly", DirectionRecvOnly},
		{"inactive", DirectionInactive},
	}

	for i, u := range tests {
		assert.Equal(t, u.expected, NewDirection(u.value), "%d: %+v", i, u)
	}
}

func TestDirection_String(t *testing.T) {
	tests := []struct {
		actual   Direction
		expected string
	}{
		{Direction(unknown), unknownStr},
		{DirectionSendRecv, "sendrecv"},
		{DirectionSendOnly, "sendonly"},
		{DirectionRecvOnly, "recvonly"},
		{DirectionInactive, "inactive"},
	}

	for i, u := range tests {
		assert.Equal(t, u.expected, u.actual.String(), "%d: %+v", i, u)
	}
}
