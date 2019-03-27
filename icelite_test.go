package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIceLite(t *testing.T) {
	tests := []struct {
		parameter string
		expected  string
	}{
		{exampleAttrIceLite, exampleAttrIceLiteLine},
	}

	for i, u := range tests {
		actual := IceLite{}
		assert.Nil(t, actual.Unmarshal(u.parameter))
		assert.Equal(t, u.expected, actual.Marshal(), "%d: %+v", i, u)
	}
}
