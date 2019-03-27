package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSdplang(t *testing.T) {
	tests := []struct {
		parameter string
		expected  string
	}{
		{exampleAttrSdplang, exampleAttrSdplangLine},
	}

	for i, u := range tests {
		actual := SdpLang{}
		assert.Nil(t, actual.Unmarshal(u.parameter))
		assert.Equal(t, u.expected, actual.Marshal(), "%d: %+v", i, u)
	}
}
