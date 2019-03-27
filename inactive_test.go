package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInactive(t *testing.T) {
	tests := []struct {
		parameter string
		expected  string
	}{
		{exampleAttrInactive, exampleAttrInactiveLine},
	}

	for i, u := range tests {
		actual := Inactive{}
		assert.Nil(t, actual.Unmarshal(u.parameter))
		assert.Equal(t, u.expected, actual.Marshal(), "%d: %+v", i, u)
	}
}
