package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRtpmap(t *testing.T) {
	tests := []struct {
		parameter string
		expected  string
	}{
		{exampleAttrRtpmap1, exampleAttrRtpmap1Line},
		{exampleAttrRtpmap2, exampleAttrRtpmap2Line},
		{exampleAttrRtpmap3, exampleAttrRtpmap3Line},
	}

	for i, u := range tests {
		actual := RtpMap{}
		assert.Nil(t, actual.Unmarshal(u.parameter))
		assert.Equal(t, u.expected, actual.Marshal(), "%d: %+v", i, u)
	}
}
