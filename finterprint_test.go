package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFingerprint(t *testing.T) {
	tests := []struct {
		parameter string
		expected  string
	}{
		{exampleAttrFingerprint, exampleAttrFingerprintLine},
	}

	for i, u := range tests {
		actual := Fingerprint{}
		assert.Nil(t, actual.Unmarshal(u.parameter))
		assert.Equal(t, u.expected, actual.Marshal(), "%d: %+v", i, u)
	}
}
