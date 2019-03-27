package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIceOptions(t *testing.T) {
	tests := []struct {
		parameter string
		expected  string
	}{
		{exampleAttrIceOptions, exampleAttrIceOptionsLine},
	}

	for i, u := range tests {
		actual := IceOptions{}
		assert.Nil(t, actual.Unmarshal(u.parameter))
		assert.Equal(t, u.expected, actual.Marshal(), "%d: %+v", i, u)
	}
}
