package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRtcpFb(t *testing.T) {
	tests := []struct {
		parameter string
		expected  string
	}{
		{exampleAttrRtcpFb1, exampleAttrRtcpFb1Line},
		{exampleAttrRtcpFb2, exampleAttrRtcpFb2Line},
		{exampleAttrRtcpFb3, exampleAttrRtcpFb3Line},
	}

	for i, u := range tests {
		actual := RtcpFb{}
		assert.Nil(t, actual.Unmarshal(u.parameter))
		assert.Equal(t, u.expected, actual.Marshal(), "%d: %+v", i, u)
	}
}
