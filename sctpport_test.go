package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSctpPort(t *testing.T) {
	tests := []struct {
		parameter string
		expected  string
	}{
		{exampleAttrSctpPort, exampleAttrSctpPortLine},
	}

	for i, u := range tests {
		actual := SctpPort{}
		assert.Nil(t, actual.Unmarshal(u.parameter))
		assert.Equal(t, u.expected, actual.Marshal(), "%d: %+v", i, u)
	}
}
