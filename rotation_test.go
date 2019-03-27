package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRotation(t *testing.T) {
	tests := []struct {
		value    string
		expected Rotation
	}{
		{"unknown", Rotation(unknown)},
		{"portrait", RotationPortrait},
		{"landscape", RotationLandscape},
		{"seascape", RotationSeascape},
	}

	for i, u := range tests {
		assert.Equal(t, u.expected, NewRotation(u.value), "%d: %+v", i, u)
	}
}

func TestRotation_String(t *testing.T) {
	tests := []struct {
		actual   Rotation
		expected string
	}{
		{Rotation(unknown), unknownStr},
		{RotationPortrait, "portrait"},
		{RotationLandscape, "landscape"},
		{RotationSeascape, "seascape"},
	}

	for i, u := range tests {
		assert.Equal(t, u.expected, u.actual.String(), "%d: %+v", i, u)
	}
}
