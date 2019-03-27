package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFramerate(t *testing.T) {
	tests := []struct {
		parameter string
		expected  string
	}{
		{exampleAttrFramerate, exampleAttrFramerateLine},
	}

	for i, u := range tests {
		actual := Framerate{}
		assert.Nil(t, actual.Unmarshal(u.parameter))
		assert.Equal(t, u.expected, actual.Marshal(), "%d: %+v", i, u)
	}
}
