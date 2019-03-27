package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIceMismatch(t *testing.T) {
	tests := []struct {
		parameter string
		expected  string
	}{
		{exampleAttrIceMismatch, exampleAttrIceMismatchLine},
	}

	for i, u := range tests {
		actual := IceMismatch{}
		assert.Nil(t, actual.Unmarshal(u.parameter))
		assert.Equal(t, u.expected, actual.Marshal(), "%d: %+v", i, u)
	}
}
