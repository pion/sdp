package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMaxMessageSize(t *testing.T) {
	tests := []struct {
		parameter string
		expected  string
	}{
		{exampleAttrMaxMessageSize, exampleAttrMaxMessageSizeLine},
	}

	for i, u := range tests {
		actual := MaxMessageSize{}
		assert.Nil(t, actual.Unmarshal(u.parameter))
		assert.Equal(t, u.expected, actual.Marshal(), "%d: %+v", i, u)
	}
}
