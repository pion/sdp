package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRtcpMux(t *testing.T) {
	tests := []struct {
		parameter string
		expected  string
	}{
		{exampleAttrRtcpMux, exampleAttrRtcpMuxLine},
	}

	for i, u := range tests {
		actual := RtcpMux{}
		assert.Nil(t, actual.Unmarshal(u.parameter))
		assert.Equal(t, u.expected, actual.Marshal(), "%d: %+v", i, u)
	}
}
