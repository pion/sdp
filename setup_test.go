package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetup(t *testing.T) {
	tests := []struct {
		parameter string
		expected  string
	}{
		{exampleAttrSetup, exampleAttrSetupLine},
	}

	for i, u := range tests {
		actual := Setup{}
		assert.Nil(t, actual.Unmarshal(u.parameter))
		assert.Equal(t, u.expected, actual.Marshal(), "%d: %+v", i, u)
	}
}
