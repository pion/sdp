// SPDX-FileCopyrightText: 2023 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithFingerprint(t *testing.T) {
	m := new(MediaDescription)

	assert.EqualValues(t, []Attribute(nil), m.Attributes)

	m = m.WithFingerprint([]byte("testalgorithm"), []byte("testfingerprint"))

	assert.EqualValues(t, []Attribute{
		{[]byte("fingerprint"), []byte("testalgorithm testfingerprint")},
	},
		m.Attributes)
}
