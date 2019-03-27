package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLang(t *testing.T) {
	tests := []struct {
		parameter string
		expected  string
	}{
		{exampleAttrLang, exampleAttrLangLine},
	}

	for i, u := range tests {
		actual := Lang{}
		assert.Nil(t, actual.Unmarshal(u.parameter))
		assert.Equal(t, u.expected, actual.Marshal(), "%d: %+v", i, u)
	}
}
