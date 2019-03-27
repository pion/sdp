package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCharset(t *testing.T) {
	tests := []struct {
		parameter string
		expected  string
	}{
		{exampleAttrCharset, exampleAttrCharsetLine},
	}

	for i, u := range tests {
		actual := Charset{}
		assert.Nil(t, actual.Unmarshal(u.parameter))
		assert.Equal(t, u.expected, actual.Marshal(), "%d: %+v", i, u)
	}
}
