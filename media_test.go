package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMediaName(t *testing.T) {
	tests := []struct {
		parameter string
		expected  string
	}{
		{exampleMedia1, exampleMedia1Line},
		{exampleMedia2, exampleMedia2Line},
	}

	for i, u := range tests {
		actual := Media{}
		assert.Nil(t, actual.Unmarshal(u.parameter))
		assert.Equal(t, u.expected, actual.Marshal(), "%d: %+v", i, u)
	}
}
