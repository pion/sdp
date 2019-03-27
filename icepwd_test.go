package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIcePwd(t *testing.T) {
	tests := []struct {
		parameter string
		expected  string
	}{
		{exampleAttrIcePwd, exampleAttrIcePwdLine},
	}

	for i, u := range tests {
		actual := IcePwd{}
		assert.Nil(t, actual.Unmarshal(u.parameter))
		assert.Equal(t, u.expected, actual.Marshal(), "%d: %+v", i, u)
	}
}
