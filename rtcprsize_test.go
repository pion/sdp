package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRtcpRsize(t *testing.T) {
	tests := []struct {
		parameter string
		expected  string
	}{
		{exampleAttrRtcpRsize, exampleAttrRtcpRsizeLine},
	}

	for i, u := range tests {
		actual := RtcpRsize{}
		assert.Nil(t, actual.Unmarshal(u.parameter))
		assert.Equal(t, u.expected, actual.Marshal(), "%d: %+v", i, u)
	}
}
