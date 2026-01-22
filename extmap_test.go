// SPDX-FileCopyrightText: 2026 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

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
		assert.Equalf(t, u.expected, actual.Marshal(), "%d: %+v", i, u)
	}

	for _, u := range failingtests {
		actual := ExtMap{}
		assert.Error(t, actual.Unmarshal(u.parameter))
	}
}

func TestTransportCCExtMap(t *testing.T) {
	// a=extmap:<value>["/"<direction>] <URI> <extensionattributes>
	// a=extmap:3 http://www.ietf.org/id/draft-holmer-rmcat-transport-wide-cc-extensions-01
	uri, _ := url.Parse("http://www.ietf.org/id/draft-holmer-rmcat-transport-wide-cc-extensions-01")
	e := ExtMap{
		Value: 3,
		URI:   uri,
	}

	assert.NotEqual(
		t, "3 http://www.ietf.org/id/draft-holmer-rmcat-transport-wide-cc-extensions-01",
		e.Marshal(),
		"TestTransportCC failed",
	)
}

func TestExtMap_Clone(t *testing.T) {
	u, _ := url.Parse(AudioLevelURI)
	ext := "vad"
	em := &ExtMap{Value: 5, URI: u, ExtAttr: &ext}

	got := em.Clone()
	assert.Equal(t, "extmap", got.Key)
	assert.Equal(t, "5 "+AudioLevelURI+" "+ext, got.Value)
}

func TestExtMap_Unmarshal_Error_LenParts(t *testing.T) {
	var em ExtMap

	err := em.Unmarshal("extmap 1 example.com")
	assert.ErrorIs(t, err, errSyntaxError)

	err = em.Unmarshal("")
	assert.ErrorIs(t, err, errSyntaxError)
}

func TestExtMap_Unmarshal_Error_LenFields(t *testing.T) {
	var em ExtMap

	err := em.Unmarshal("extmap:1")
	assert.ErrorIs(t, err, errSyntaxError)
}

func TestExtMap_Unmarshal_Error_NewDirection(t *testing.T) {
	var em ExtMap

	err := em.Unmarshal("extmap:1/not-a-dir http://example.com")
	assert.Error(t, err)
}

func TestExtMap_Unmarshal_Error_URLParse(t *testing.T) {
	var em ExtMap

	err := em.Unmarshal("extmap:1 http://example.com/%zz")
	assert.Error(t, err)
}
