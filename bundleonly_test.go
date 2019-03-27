package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBundleOnly(t *testing.T) {
	tests := []struct {
		parameter string
		expected  string
	}{
		{exampleAttrBundleOnly, exampleAttrBundleOnlyLine},
	}

	for i, u := range tests {
		actual := BundleOnly{}
		assert.Nil(t, actual.Unmarshal(u.parameter))
		assert.Equal(t, u.expected, actual.Marshal(), "%d: %+v", i, u)
	}
}
