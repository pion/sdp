package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRtcp(t *testing.T) {
	tests := []struct {
		parameter string
		expected  string
	}{
		{exampleAttrRtcp, exampleAttrRtcpLine},
	}

	for i, u := range tests {
		actual := Rtcp{}
		assert.Nil(t, actual.Unmarshal(u.parameter))
		assert.Equal(t, u.expected, actual.Marshal(), "%d: %+v", i, u)
	}
}
