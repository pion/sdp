package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithFingerprint(t *testing.T) {
	m := new(MediaDescription)

	assert.Equal(t, []Attribute(nil), m.Attributes)

	m = m.WithFingerprint("testalgorithm", "testfingerprint")

	assert.Equal(t, []Attribute{
		{"fingerprint", "testalgorithm testfingerprint"},
	},
		m.Attributes)
}
