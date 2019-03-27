package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecvonly(t *testing.T) {
	tests := []struct {
		parameter string
		expected  string
	}{
		{exampleAttrRecvonly, exampleAttrRecvonlyLine},
	}

	for i, u := range tests {
		actual := RecvOnly{}
		assert.Nil(t, actual.Unmarshal(u.parameter))
		assert.Equal(t, u.expected, actual.Marshal(), "%d: %+v", i, u)
	}
}
