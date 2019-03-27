package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnection(t *testing.T) {
	tests := []struct {
		parameter string
		expected  string
	}{
		{exampleConnection, exampleConnectionLine},
	}

	for i, u := range tests {
		actual := Connection{}
		assert.Nil(t, actual.Unmarshal(u.parameter))
		assert.Equal(t, u.expected, actual.Marshal(), "%d: %+v", i, u)
	}
}
