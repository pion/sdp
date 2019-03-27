package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestType(t *testing.T) {
	tests := []struct {
		parameter string
		expected  string
	}{
		{exampleAttrType, exampleAttrTypeLine},
	}

	for i, u := range tests {
		actual := Type{}
		assert.Nil(t, actual.Unmarshal(u.parameter))
		assert.Equal(t, u.expected, actual.Marshal(), "%d: %+v", i, u)
	}
}
