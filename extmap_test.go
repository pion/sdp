package sdp

import (
	"net/url"
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

func TestTransportCCExtMap(t *testing.T) {
	//a=extmap:<value>["/"<direction>] <URI> <extensionattributes>
	//a=extmap:3 http://www.ietf.org/id/draft-holmer-rmcat-transport-wide-cc-extensions-01
	uri, _ := url.Parse("http://www.ietf.org/id/draft-holmer-rmcat-transport-wide-cc-extensions-01")
	e := ExtMap{
		Value: 3,
		URI:   uri,
	}

	if e.Marshal() == "3 http://www.ietf.org/id/draft-holmer-rmcat-transport-wide-cc-extensions-01" {
		t.Error("TestTransportCC failed")
	}
}
