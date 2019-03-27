package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendonly(t *testing.T) {
	tests := []struct {
		parameter string
		expected  string
	}{
		{exampleAttrSendonly, exampleAttrSendonlyLine},
	}

	for i, u := range tests {
		actual := SendOnly{}
		assert.Nil(t, actual.Unmarshal(u.parameter))
		assert.Equal(t, u.expected, actual.Marshal(), "%d: %+v", i, u)
	}
}
