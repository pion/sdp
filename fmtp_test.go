package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFmtp(t *testing.T) {
	tests := []struct {
		parameter string
		expected  string
	}{
		{exampleAttrFmtp1, exampleAttrFmtp1Line},
		{exampleAttrFmtp2, exampleAttrFmtp2Line},
	}

	for i, u := range tests {
		actual := Fmtp{}
		assert.Nil(t, actual.Unmarshal(u.parameter))
		assert.Equal(t, u.expected, actual.Marshal(), "%d: %+v", i, u)
	}
}
