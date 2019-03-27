package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncryptionKey(t *testing.T) {
	tests := []struct {
		parameter string
		expected  string
	}{
		{exampleEncryptionKey, exampleEncryptionKeyLine},
	}

	for i, u := range tests {
		actual := EncryptionKey{}
		assert.Nil(t, actual.Unmarshal(u.parameter))
		assert.Equal(t, u.expected, actual.Marshal(), "%d: %+v", i, u)
	}
}
