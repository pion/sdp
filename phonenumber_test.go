package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPhoneNumber(t *testing.T) {
	tests := []struct {
		parameter string
		expected  string
	}{
		{examplePhone, examplePhoneLine},
	}

	for i, u := range tests {
		actual := PhoneNumber{}
		assert.Nil(t, actual.Unmarshal(u.parameter))
		assert.Equal(t, u.expected, actual.Marshal(), "%d: %+v", i, u)
	}
}
