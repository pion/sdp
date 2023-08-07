// SPDX-FileCopyrightText: 2023 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

package sdp

import "testing"

func FuzzUnmarshal(f *testing.F) {
	f.Add([]byte{})
	f.Add([]byte(CanonicalUnmarshalSDP))
	f.Fuzz(func(t *testing.T, data []byte) {
		// Check that unmarshalling any byte slice does not panic.
		var sd SessionDescription
		if err := sd.Unmarshal(data); err != nil {
			return
		}
		// Check that we can marshal anything we unmarshalled.
		_, err := sd.Marshal()
		if err != nil {
			t.Fatalf("failed to marshal")
		}
	})
}
