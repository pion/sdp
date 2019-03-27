package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTimeZone(t *testing.T) {
	tests := []struct {
		parameter string
		expected  string
	}{
		{exampleTimeZone1, exampleTimeZone3},
		{exampleTimeZone2, exampleTimeZone2},
	}

	for i, u := range tests {
		actual := TimeZone{}
		assert.Nil(t, actual.Unmarshal(u.parameter))
		assert.Equal(t, u.expected, actual.Marshal(), "%d: %+v", i, u)
	}
}
