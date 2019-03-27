package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCategory(t *testing.T) {
	tests := []struct {
		parameter string
		expected  string
	}{
		{exampleAttrCategory, exampleAttrCategoryLine},
	}

	for i, u := range tests {
		actual := Category{}
		assert.Nil(t, actual.Unmarshal(u.parameter))
		assert.Equal(t, u.expected, actual.Marshal(), "%d: %+v", i, u)
	}
}
