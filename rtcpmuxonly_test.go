package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRtcpOnly(t *testing.T) {
	tests := []struct {
		parameter string
		expected  string
	}{
		{exampleAttrRtcpMuxOnly, exampleAttrRtcpMuxOnlyLine},
	}

	for i, u := range tests {
		actual := RtcpMuxOnly{}
		assert.Nil(t, actual.Unmarshal(u.parameter))
		assert.Equal(t, u.expected, actual.Marshal(), "%d: %+v", i, u)
	}
}
