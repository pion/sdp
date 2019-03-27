package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEndOfCandidates(t *testing.T) {
	tests := []struct {
		parameter string
		expected  string
	}{
		{exampleAttrEndOfCandidates, exampleAttrEndOfCandidatesLine},
	}

	for i, u := range tests {
		actual := EndOfCandidates{}
		assert.Nil(t, actual.Unmarshal(u.parameter))
		assert.Equal(t, u.expected, actual.Marshal(), "%d: %+v", i, u)
	}
}
