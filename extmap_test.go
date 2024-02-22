// SPDX-FileCopyrightText: 2023 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

package sdp

import (
	"testing"

	"github.com/stretchr/testify/require"
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
		require.NoError(t, actual.Unmarshal([]byte(u.parameter)))
		require.Equal(t, u.expected, actual.Marshal(), "%d: %+v", i, u)
	}

	for _, u := range failingtests {
		actual := ExtMap{}
		require.Error(t, actual.Unmarshal([]byte(u.parameter)))
	}
}

func TestTransportCCExtMap(t *testing.T) {
	// a=extmap:<value>["/"<direction>] <URI> <extensionattributes>
	// a=extmap:3 http://www.ietf.org/id/draft-holmer-rmcat-transport-wide-cc-extensions-01
	e := ExtMap{
		Value: 3,
		URI:   []byte("http://www.ietf.org/id/draft-holmer-rmcat-transport-wide-cc-extensions-01"),
	}

	require.EqualValues(t, e.Marshal(), "3 http://www.ietf.org/id/draft-holmer-rmcat-transport-wide-cc-extensions-01")
}
