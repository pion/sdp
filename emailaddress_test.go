package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmailAddress(t *testing.T) {
	tests := []struct {
		parameter string
		expected  string
	}{
		{exampleEmail, exampleEmailLine},
	}

	for i, u := range tests {
		actual := EmailAddress{}
		assert.Nil(t, actual.Unmarshal(u.parameter))
		assert.Equal(t, u.expected, actual.Marshal(), "%d: %+v", i, u)
	}
}
