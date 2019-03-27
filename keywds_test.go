package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeywords(t *testing.T) {
	tests := []struct {
		parameter string
		expected  string
	}{
		{exampleAttrKeywords, exampleAttrKeywordsLine},
	}

	for i, u := range tests {
		actual := Keywds{}
		assert.Nil(t, actual.Unmarshal(u.parameter))
		assert.Equal(t, u.expected, actual.Marshal(), "%d: %+v", i, u)
	}
}
