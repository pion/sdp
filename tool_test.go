package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTool(t *testing.T) {
	tests := []struct {
		parameter string
		expected  string
	}{
		{exampleAttrTool, exampleAttrToolLine},
	}

	for i, u := range tests {
		actual := Tool{}
		assert.Nil(t, actual.Unmarshal(u.parameter))
		assert.Equal(t, u.expected, actual.Marshal(), "%d: %+v", i, u)
	}
}
