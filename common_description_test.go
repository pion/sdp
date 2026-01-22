// SPDX-FileCopyrightText: 2026 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInformation_String(t *testing.T) {
	i := Information("About this session")
	assert.Equal(t, "About this session", i.String())
}

func TestConnectionInformation_String(t *testing.T) {
	t.Run("without address", func(t *testing.T) {
		c := ConnectionInformation{NetworkType: "IN", AddressType: "IP4", Address: nil}
		assert.Equal(t, "IN IP4", c.String())
	})

	t.Run("with address + TTL + Range", func(t *testing.T) {
		ttl, rg := 127, 3
		addr := &Address{Address: "224.2.17.12", TTL: &ttl, Range: &rg}
		c := ConnectionInformation{NetworkType: "IN", AddressType: "IP4", Address: addr}
		assert.Equal(t, "IN IP4 224.2.17.12/127/3", c.String())
	})
}

func TestAddress_String_Variants(t *testing.T) {
	t.Run("only address", func(t *testing.T) {
		a := &Address{Address: "239.255.255.250"}
		assert.Equal(t, "239.255.255.250", a.String())
	})

	t.Run("TTL only", func(t *testing.T) {
		ttl := 5
		a := &Address{Address: "239.255.255.250", TTL: &ttl}
		assert.Equal(t, "239.255.255.250/5", a.String())
	})

	t.Run("Range only", func(t *testing.T) {
		rg := 7
		a := &Address{Address: "239.255.255.250", Range: &rg}
		assert.Equal(t, "239.255.255.250/7", a.String())
	})

	t.Run("TTL and Range", func(t *testing.T) {
		ttl, rg := 5, 7
		a := &Address{Address: "239.255.255.250", TTL: &ttl, Range: &rg}
		assert.Equal(t, "239.255.255.250/5/7", a.String())
	})
}

func TestAddress_marshal_RangePaths(t *testing.T) {
	t.Run("Range only size matches", func(t *testing.T) {
		rg := 9
		a := Address{Address: "a", Range: &rg}
		out := a.marshalInto(nil)
		want := "a/9"
		assert.Equal(t, want, string(out))
		assert.Equal(t, len(want), a.marshalSize())
	})

	t.Run("TTL and Range size matches", func(t *testing.T) {
		ttl, rg := 2, 11
		a := Address{Address: "addr", TTL: &ttl, Range: &rg}
		s := (&a).String()
		assert.Equal(t, "addr/2/11", s)
		assert.Equal(t, len(s), a.marshalSize())
	})
}

func TestBandwidth_String(t *testing.T) {
	t.Run("standard", func(t *testing.T) {
		b := Bandwidth{Experimental: false, Type: "AS", Bandwidth: 512}
		assert.Equal(t, "AS:512", b.String())
	})

	t.Run("experimental", func(t *testing.T) {
		b := Bandwidth{Experimental: true, Type: "AS", Bandwidth: 512}
		assert.Equal(t, "X-AS:512", b.String())
	})
}

func TestEncryptionKey_String(t *testing.T) {
	e := EncryptionKey("clear:hunter2")
	assert.Equal(t, "clear:hunter2", e.String())
}

func TestAttribute_IsICECandidate(t *testing.T) {
	assert.True(t, Attribute{Key: "candidate"}.IsICECandidate())
	assert.False(t, Attribute{Key: "Candidate"}.IsICECandidate())
	assert.False(t, Attribute{Key: "ice-candidate"}.IsICECandidate())
	assert.False(t, Attribute{Key: ""}.IsICECandidate())
}
