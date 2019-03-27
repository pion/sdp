package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMsid(t *testing.T) {
	tests := []struct {
		parameter string
		expected  string
	}{
		{exampleAttrMsid, exampleAttrMsidLine},
	}

	for i, u := range tests {
		actual := MsID{}
		assert.Nil(t, actual.Unmarshal(u.parameter))
		assert.Equal(t, u.expected, actual.Marshal(), "%d: %+v", i, u)
	}
}
