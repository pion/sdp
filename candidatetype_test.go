package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCandidateType(t *testing.T) {
	tests := []struct {
		value    string
		expected CandidateType
	}{
		{"unknown", CandidateType(unknown)},
		{"host", CandidateTypeHost},
		{"srflx", CandidateTypeSrflx},
		{"prflx", CandidateTypePrflx},
		{"relay", CandidateTypeRelay},
	}

	for i, u := range tests {
		assert.Equal(t, u.expected, NewCandidateType(u.value), "%d: %+v", i, u)
	}
}

func TestIceCandidateType_String(t *testing.T) {
	tests := []struct {
		actual   CandidateType
		expected string
	}{
		{CandidateType(unknown), unknownStr},
		{CandidateTypeHost, "host"},
		{CandidateTypeSrflx, "srflx"},
		{CandidateTypePrflx, "prflx"},
		{CandidateTypeRelay, "relay"},
	}

	for i, u := range tests {
		assert.Equal(t, u.expected, u.actual.String(), "%d: %+v", i, u)
	}
}
