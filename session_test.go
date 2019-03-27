package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSession(t *testing.T) {
	tests := []struct {
		parameter string
		expected  string
	}{
		{exampleSession, exampleSessionLine},
	}

	for i, u := range tests {
		actual := Session{}
		assert.Nil(t, actual.Unmarshal(u.parameter))
		assert.Equal(t, u.expected, actual.Marshal(), "%d: %+v", i, u)
	}
}
