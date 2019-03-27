package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrigin(t *testing.T) {
	tests := []struct {
		parameter string
		expected  string
	}{
		{exampleOrigin, exampleOriginLine},
	}

	for i, u := range tests {
		actual := Origin{}
		assert.Nil(t, actual.Unmarshal(u.parameter))
		assert.Equal(t, u.expected, actual.Marshal(), "%d: %+v", i, u)
	}
}
