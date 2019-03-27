package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepeatTime(t *testing.T) {
	tests := []struct {
		parameter string
		expected  string
	}{
		{exampleRepeatTime1, exampleRepeatTime1Line},
		{exampleRepeatTime2, exampleRepeatTime3Line},
	}

	for i, u := range tests {
		actual := RepeatTime{}
		assert.Nil(t, actual.Unmarshal(u.parameter))
		assert.Equal(t, u.expected, actual.Marshal(), "%d: %+v", i, u)
	}
}
