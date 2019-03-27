package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTiming(t *testing.T) {
	tests := []struct {
		parameter string
		expected  string
	}{
		{exampleTiming1, exampleTiming1Line},
	}

	for i, u := range tests {
		actual := Timing{}
		assert.Nil(t, actual.Unmarshal(u.parameter))
		assert.Equal(t, u.expected, actual.Marshal(), "%d: %+v", i, u)
	}
}
