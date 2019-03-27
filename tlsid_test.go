package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTlsID(t *testing.T) {
	tests := []struct {
		parameter string
		expected  string
	}{
		{exampleAttrTlsId, exampleAttrTlsIdLine},
	}

	for i, u := range tests {
		actual := TlsID{}
		assert.Nil(t, actual.Unmarshal(u.parameter))
		assert.Equal(t, u.expected, actual.Marshal(), "%d: %+v", i, u)
	}
}
