package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTimeZones(t *testing.T) {
	tests := []struct {
		parameter string
		expected  string
	}{
		{exampleTimeZones1, exampleTimeZones2Line},
	}

	for i, u := range tests {
		actual := TimeZones{}
		assert.Nil(t, actual.Unmarshal(u.parameter))
		assert.Equal(t, u.expected, actual.Marshal(), "%d: %+v", i, u)
	}
}
