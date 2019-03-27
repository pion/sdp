package sdp

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSendrecv(t *testing.T) {
	tests := []struct {
		parameter string
		expected  string
	}{
		{exampleAttrSendrecv, exampleAttrSendrecvLine},
	}

	for i, u := range tests {
		actual := SendRecv{}
		assert.Nil(t, actual.Unmarshal(u.parameter))
		assert.Equal(t, u.expected, actual.Marshal(), "%d: %+v", i, u)
	}
}
