package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPtime(t *testing.T) {
	tests := []struct {
		parameter string
		expected  string
	}{
		{exampleAttrPtime, exampleAttrPtimeLine},
	}

	for i, u := range tests {
		actual := Ptime{}
		assert.Nil(t, actual.Unmarshal(u.parameter))
		assert.Equal(t, u.expected, actual.Marshal(), "%d: %+v", i, u)
	}
}
