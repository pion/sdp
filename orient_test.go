package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrient(t *testing.T) {
	tests := []struct {
		parameter string
		expected  string
	}{
		{exampleAttrOrient1, exampleAttrOrient1Line},
		{exampleAttrOrient2, exampleAttrOrient2Line},
		{exampleAttrOrient3, exampleAttrOrient3Line},
	}

	for i, u := range tests {
		actual := Orient{}
		assert.Nil(t, actual.Unmarshal(u.parameter))
		assert.Equal(t, u.expected, actual.Marshal(), "%d: %+v", i, u)
	}
}
