package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIceUfrag(t *testing.T) {
	tests := []struct {
		parameter string
		expected  string
	}{
		{exampleAttrIceUfrag, exampleAttrIceUfragLine},
	}

	for i, u := range tests {
		actual := IceUfrag{}
		assert.Nil(t, actual.Unmarshal(u.parameter))
		assert.Equal(t, u.expected, actual.Marshal(), "%d: %+v", i, u)
	}
}
