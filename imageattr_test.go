package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImageattr(t *testing.T) {
	tests := []struct {
		parameter string
		expected  string
	}{
		{exampleAttrImageattr, exampleAttrImageattrLine},
	}

	for i, u := range tests {
		actual := ImageAttr{}
		assert.Nil(t, actual.Unmarshal(u.parameter))
		assert.Equal(t, u.expected, actual.Marshal(), "%d: %+v", i, u)
	}
}
