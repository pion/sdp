package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCandidate(t *testing.T) {
	tests := []struct {
		parameter string
		expected  string
	}{
		{exampleAttrCandidate1, exampleAttrCandidate1Line},
		{exampleAttrCandidate2, exampleAttrCandidate2Line},
		{exampleAttrCandidate3, exampleAttrCandidate3Line},
	}

	for i, u := range tests {
		actual := Candidate{}
		assert.Nil(t, actual.Unmarshal(u.parameter))
		assert.Equal(t, u.expected, actual.Marshal(), "%d: %+v", i, u)
	}
}
