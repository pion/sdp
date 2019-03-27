package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMid(t *testing.T) {
	tests := []struct {
		parameter string
		expected  string
	}{
		{exampleAttrMid, exampleAttrMidLine},
	}

	for i, u := range tests {
		actual := MID{}
		assert.Nil(t, actual.Unmarshal(u.parameter))
		assert.Equal(t, u.expected, actual.Marshal(), "%d: %+v", i, u)
	}
}
