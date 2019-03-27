package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMaxptime(t *testing.T) {
	tests := []struct {
		parameter string
		expected  string
	}{
		{exampleAttrMaxptime, exampleAttrMaxptimeLine},
	}

	for i, u := range tests {
		actual := MaxPtime{}
		assert.Nil(t, actual.Unmarshal(u.parameter))
		assert.Equal(t, u.expected, actual.Marshal(), "%d: %+v", i, u)
	}
}
