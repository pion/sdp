// SPDX-FileCopyrightText: 2026 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	exampleAttrExtmap1     = "extmap:1 http://example.com/082005/ext.htm#ttime"
	exampleAttrExtmap1Line = exampleAttrExtmap1
	exampleAttrExtmap2     = "extmap:2/sendrecv http://example.com/082005/ext.htm#xmeta short"
	exampleAttrExtmap2Line = exampleAttrExtmap2
	failingAttrExtmap1     = "extmap:257/sendrecv http://example.com/082005/ext.htm#xmeta short"
	failingAttrExtmap1Line = attributeKey + failingAttrExtmap1
	failingAttrExtmap2     = "extmap:2/blorg http://example.com/082005/ext.htm#xmeta short"
	failingAttrExtmap2Line = attributeKey + failingAttrExtmap2
)

func TestSessionDescription_Attribute(t *testing.T) {
	sd := &SessionDescription{
		Attributes: []Attribute{
			{Key: "ice-options", Value: "trickle"},
			{Key: "mid", Value: "video"},
		},
	}

	t.Run("found", func(t *testing.T) {
		v, ok := sd.Attribute("mid")
		assert.True(t, ok)
		assert.Equal(t, "video", v)
	})

	t.Run("not found", func(t *testing.T) {
		v, ok := sd.Attribute("does-not-exist")
		assert.False(t, ok)
		assert.Equal(t, "", v)
	})
}

func TestVersion_String(t *testing.T) {
	var v Version = 0
	assert.Equal(t, "0", v.String())

	v = 2
	assert.Equal(t, "2", v.String())
}

func TestOrigin_String(t *testing.T) {
	o := Origin{
		Username:       "alice",
		SessionID:      12345,
		SessionVersion: 678,
		NetworkType:    "IN",
		AddressType:    "IP4",
		UnicastAddress: "111.1.111.1",
	}
	assert.Equal(t, "alice 12345 678 IN IP4 111.1.111.1", o.String())
}

func TestSessionName_String(t *testing.T) {
	sn := SessionName("My Session")
	assert.Equal(t, "My Session", sn.String())

	empty := SessionName("")
	assert.Equal(t, "", empty.String())
}

func TestEmailAddress_String(t *testing.T) {
	e := EmailAddress("user@pion.com")
	assert.Equal(t, "user@pion.com", e.String())
}

func TestPhoneNumber_String(t *testing.T) {
	p := PhoneNumber("+1 111 1111")
	assert.Equal(t, "+1 111 1111", p.String())
}

func TestTimeZone_String(t *testing.T) {
	z := TimeZone{AdjustmentTime: 3600, Offset: -1800}
	assert.Equal(t, "3600 -1800", z.String())

	z = TimeZone{AdjustmentTime: 0, Offset: 0}
	assert.Equal(t, "0 0", z.String())
}
