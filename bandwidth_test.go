package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBandwidth(t *testing.T) {
	tests := []struct {
		parameter string
		expected  string
	}{
		{exampleBandwidth1, exampleBandwidth1Line},
		{exampleBandwidth2, exampleBandwidth2Line},
	}

	for i, u := range tests {
		actual := Bandwidth{}
		assert.Nil(t, actual.Unmarshal(u.parameter))
		assert.Equal(t, u.expected, actual.Marshal(), "%d: %+v", i, u)
	}
}
