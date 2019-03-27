package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuality(t *testing.T) {
	tests := []struct {
		parameter string
		expected  string
	}{
		{exampleAttrQuality, exampleAttrQualityLine},
	}

	for i, u := range tests {
		actual := Quality{}
		assert.Nil(t, actual.Unmarshal(u.parameter))
		assert.Equal(t, u.expected, actual.Marshal(), "%d: %+v", i, u)
	}
}
