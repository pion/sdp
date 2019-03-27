package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimulcast(t *testing.T) {
	tests := []struct {
		parameter string
		expected  string
	}{
		{exampleAttrSimulcast, exampleAttrSimulcastLine},
	}

	for i, u := range tests {
		actual := Simulcast{}
		assert.Nil(t, actual.Unmarshal(u.parameter))
		assert.Equal(t, u.expected, actual.Marshal(), "%d: %+v", i, u)
	}
}
