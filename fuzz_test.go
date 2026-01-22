// SPDX-FileCopyrightText: 2026 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func FuzzUnmarshal(f *testing.F) {
	f.Add("")
	f.Add(CanonicalUnmarshalSDP)
	f.Fuzz(func(t *testing.T, data string) {
		// Check that unmarshalling any byte slice does not panic.
		var sd SessionDescription
		if err := sd.UnmarshalString(data); err != nil {
			return
		}
		// Check that we can marshal anything we unmarshalled.
		_, err := sd.Marshal()
		assert.NoError(t, err)
	})
}
