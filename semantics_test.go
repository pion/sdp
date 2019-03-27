package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSemantic(t *testing.T) {
	tests := []struct {
		value    string
		expected Semantic
	}{
		{"unknown", Semantic(unknown)},
		{"LS", SemanticLP},
		{"FID", SemanticFID},
		{"BUNDLE", SemanticBUNDLE},
		{"FEC", SemanticFEC},
		{"WMS", SemanticWMS},
	}

	for i, u := range tests {
		assert.Equal(t, u.expected, NewSemantic(u.value), "%d: %+v", i, u)
	}
}

func TestSignalingState_String(t *testing.T) {
	tests := []struct {
		actual   Semantic
		expected string
	}{
		{Semantic(unknown), unknownStr},
		{SemanticLP, "LS"},
		{SemanticFID, "FID"},
		{SemanticBUNDLE, "BUNDLE"},
		{SemanticFEC, "FEC"},
		{SemanticWMS, "WMS"},
	}

	for i, u := range tests {
		assert.Equal(t, u.expected, u.actual.String(), "%d: %+v", i, u)
	}
}
