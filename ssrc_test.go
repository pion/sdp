package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSsrc(t *testing.T) {
	tests := []struct {
		parameter string
		expected  string
	}{
		{exampleAttrSsrc1, exampleAttrSsrc1Line},
		{exampleAttrSsrc2, exampleAttrSsrc2Line},
		{exampleAttrSsrc3, exampleAttrSsrc3Line},
		{exampleAttrSsrc4, exampleAttrSsrc4Line},
	}

	for i, u := range tests {
		actual := Ssrc{}
		assert.Nil(t, actual.Unmarshal(u.parameter))
		assert.Equal(t, u.expected, actual.Marshal(), "%d: %+v", i, u)
	}
}
