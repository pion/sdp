package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoteCandidates(t *testing.T) {
	tests := []struct {
		parameter string
		expected  string
	}{
		{exampleAttrRemoteCandidates, exampleAttrRemoteCandidatesLine},
	}

	for i, u := range tests {
		actual := RemoteCandidates{}
		assert.Nil(t, actual.Unmarshal(u.parameter))
		assert.Equal(t, u.expected, actual.Marshal(), "%d: %+v", i, u)
	}
}
