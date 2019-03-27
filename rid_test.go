package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRid(t *testing.T) {
	tests := []struct {
		parameter string
		expected  string
	}{
		{exampleAttrRid, exampleAttrRidLine},
	}

	for i, u := range tests {
		actual := RID{}
		assert.Nil(t, actual.Unmarshal(u.parameter))
		assert.Equal(t, u.expected, actual.Marshal(), "%d: %+v", i, u)
	}
}
