package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtmap(t *testing.T) {
	passingtests := []struct {
		parameter string
		expected  string
	}{
		{exampleAttrExtmap1, exampleAttrExtmap1Line},
		{exampleAttrExtmap2, exampleAttrExtmap2Line},
	}
	failingtests := []struct {
		parameter string
		expected  string
	}{
		{failingAttrExtmap1, failingAttrExtmap1Line},
		{failingAttrExtmap2, failingAttrExtmap2Line},
	}

	for i, u := range passingtests {
		actual := ExtMap{}
		assert.NoError(t, actual.Unmarshal(u.parameter))
		assert.Equal(t, u.expected, actual.Marshal(), "%d: %+v", i, u)
	}

	for _, u := range failingtests {
		actual := ExtMap{}
		assert.Error(t, actual.Unmarshal(u.parameter))
	}
}
