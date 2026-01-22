// SPDX-FileCopyrightText: 2026 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithFingerprint(t *testing.T) {
	m := new(MediaDescription)

	assert.Equal(t, []Attribute(nil), m.Attributes)

	m = m.WithFingerprint("testalgorithm", "testfingerprint")

	assert.Equal(t, []Attribute{
		{"fingerprint", "testalgorithm testfingerprint"},
	},
		m.Attributes)
}

func TestMediaDescription_Attribute(t *testing.T) {
	md := &MediaDescription{
		Attributes: []Attribute{
			{Key: "rtcp-mux"},
			{Key: "mid", Value: "video"},
			{Key: "setup", Value: "actpass"},
		},
	}

	t.Run("found", func(t *testing.T) {
		v, ok := md.Attribute("mid")
		assert.True(t, ok)
		assert.Equal(t, "video", v)
	})

	t.Run("not found", func(t *testing.T) {
		v, ok := md.Attribute("nonexistent")
		assert.False(t, ok)
		assert.Equal(t, "", v)
	})
}

func TestRangedPort_String(t *testing.T) {
	t.Run("no range", func(t *testing.T) {
		p := &RangedPort{Value: 5004}
		assert.Equal(t, "5004", p.String())
	})

	t.Run("with range", func(t *testing.T) {
		r := 2
		p := &RangedPort{Value: 5004, Range: &r}
		assert.Equal(t, "5004/2", p.String())
	})
}

func TestRangedPort_marshalInto_RangeBranch(t *testing.T) {
	r := 3
	p := RangedPort{Value: 49170, Range: &r}
	out := p.marshalInto(nil)
	assert.Equal(t, "49170/3", string(out))
}

func TestRangedPort_marshalSize_RangeBranch(t *testing.T) {
	r := 12
	p := RangedPort{Value: 65535, Range: &r}

	gotSize := p.marshalSize()
	wantLen := len((&RangedPort{Value: 65535, Range: &r}).String())
	assert.Equal(t, wantLen, gotSize)
}

func TestMediaName_String(t *testing.T) {
	r := 2
	m := MediaName{
		Media:  "audio",
		Port:   RangedPort{Value: 5004, Range: &r},
		Protos: []string{"UDP", "TLS", "RTP", "SAVPF"},
		Formats: []string{
			"111", "96",
		},
	}
	assert.Equal(t, "audio 5004/2 UDP/TLS/RTP/SAVPF 111 96", m.String())
}
