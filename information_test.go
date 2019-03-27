package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInformation(t *testing.T) {
	tests := []struct {
		parameter string
		expected  string
	}{
		{exampleInfo1, exampleInfo1Line},
	}

	for i, u := range tests {
		actual := Information{}
		assert.Nil(t, actual.Unmarshal(u.parameter))
		assert.Equal(t, u.expected, actual.Marshal(), "%d: %+v", i, u)
	}
}
