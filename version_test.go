package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	tests := []struct {
		parameter string
		expected  string
	}{
		{exampleVersion, exampleVersionLine},
	}

	for i, u := range tests {
		actual := Version{}
		assert.Nil(t, actual.Unmarshal(u.parameter))
		assert.Equal(t, u.expected, actual.Marshal(), "%d: %+v", i, u)
	}
}
