package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtmap(t *testing.T) {
	tests := []struct {
		parameter string
		expected  string
	}{
		{exampleAttrExtmap1, exampleAttrExtmap1Line},
		{exampleAttrExtmap2, exampleAttrExtmap2Line},
	}

	for i, u := range tests {
		actual := ExtMap{}
		assert.Nil(t, actual.Unmarshal(u.parameter))
		assert.Equal(t, u.expected, actual.Marshal(), "%d: %+v", i, u)
	}
}
