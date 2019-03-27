package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIdentity(t *testing.T) {
	tests := []struct {
		parameter string
		expected  string
	}{
		{exampleAttrIdentity1, exampleAttrIdentity1Line},
		{exampleAttrIdentity2, exampleAttrIdentity2Line},
	}

	for i, u := range tests {
		actual := Identity{}
		assert.Nil(t, actual.Unmarshal(u.parameter))
		assert.Equal(t, u.expected, actual.Marshal(), "%d: %+v", i, u)
	}
}
