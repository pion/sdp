// SPDX-FileCopyrightText: 2023 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

package sdp

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewJSEPSessionDescription(t *testing.T) {
	t.Run("Without Identity", func(t *testing.T) {
		sd, err := NewJSEPSessionDescription(false)
		assert.NoError(t, err)
		assert.NotNil(t, sd)
		assert.Zero(t, sd.Version)
		assert.Equal(t, "-", sd.Origin.Username)
		assert.Equal(t, "IN", sd.Origin.NetworkType)
		assert.Equal(t, "IP4", sd.Origin.AddressType)
		assert.Equal(t, "0.0.0.0", sd.Origin.UnicastAddress)
		assert.Equal(t, SessionName("-"), sd.SessionName)
		assert.Len(t, sd.TimeDescriptions, 1)
		assert.Zero(t, sd.TimeDescriptions[0].Timing.StartTime)
		assert.Zero(t, sd.TimeDescriptions[0].Timing.StopTime)
		assert.Empty(t, sd.Attributes)
	})

	t.Run("With Identity", func(t *testing.T) {
		sd, err := NewJSEPSessionDescription(true)
		assert.NoError(t, err)
		assert.NotNil(t, sd)
		assert.Len(t, sd.Attributes, 1)
		assert.Equal(t, AttrKeyIdentity, sd.Attributes[0].Key)
	})
}

func TestSessionDescriptionAttributes(t *testing.T) {
	t.Run("WithPropertyAttribute", func(t *testing.T) {
		sd, err := NewJSEPSessionDescription(false)
		assert.NoError(t, err)
		sd = sd.WithPropertyAttribute(AttrKeyRTCPMux)
		assert.Len(t, sd.Attributes, 1)
		assert.Equal(t, AttrKeyRTCPMux, sd.Attributes[0].Key)
	})

	t.Run("WithValueAttribute", func(t *testing.T) {
		sd, err := NewJSEPSessionDescription(false)
		assert.NoError(t, err)
		sd = sd.WithValueAttribute(AttrKeyMID, "video")
		assert.Len(t, sd.Attributes, 1)
		assert.Equal(t, AttrKeyMID, sd.Attributes[0].Key)
		assert.Equal(t, "video", sd.Attributes[0].Value)
	})

	t.Run("WithICETrickleAdvertised", func(t *testing.T) {
		sd, err := NewJSEPSessionDescription(false)
		assert.NoError(t, err)
		sd = sd.WithICETrickleAdvertised()
		assert.Len(t, sd.Attributes, 1)
		assert.Equal(t, AttrKeyICEOptions, sd.Attributes[0].Key)
		assert.Equal(t, "trickle", sd.Attributes[0].Value)
	})

	t.Run("WithICERenomination", func(t *testing.T) {
		sd, err := NewJSEPSessionDescription(false)
		assert.NoError(t, err)
		sd = sd.WithICETrickleAdvertised().WithICERenomination()
		assert.Len(t, sd.Attributes, 1)
		assert.Equal(t, AttrKeyICEOptions, sd.Attributes[0].Key)
		assert.Equal(t, "trickle renomination", sd.Attributes[0].Value)
	})

	t.Run("WithFingerprint", func(t *testing.T) {
		sd, err := NewJSEPSessionDescription(false)
		assert.NoError(t, err)
		sd = sd.WithFingerprint("sha-256", "test-fingerprint")
		assert.Len(t, sd.Attributes, 1)
		assert.Equal(t, "fingerprint", sd.Attributes[0].Key)
		assert.Equal(t, "sha-256 test-fingerprint", sd.Attributes[0].Value)
	})
}

func TestSessionDescription_ICEOptions_Combined(t *testing.T) {
	t.Run("WithICETrickleAdvertised and WithICERenominationAdvertised", func(t *testing.T) {
		sd, err := NewJSEPSessionDescription(false)
		assert.NoError(t, err)

		sd = sd.WithICETrickleAdvertised().WithICERenomination()

		iceOptionsCount := 0
		var iceOptionsValue string
		for _, attr := range sd.Attributes {
			if attr.Key == AttrKeyICEOptions {
				iceOptionsCount++
				iceOptionsValue = attr.Value
			}
		}

		assert.Equal(t, 1, iceOptionsCount, "Should have exactly one ice-options attribute")
		assert.Equal(t, "trickle renomination", iceOptionsValue, "Should combine both values with space")
	})

	t.Run("WithICERenominationAdvertised and WithICETrickleAdvertised (reverse order)", func(t *testing.T) {
		sd, err := NewJSEPSessionDescription(false)
		assert.NoError(t, err)

		sd = sd.WithICERenomination().WithICETrickleAdvertised()

		iceOptionsCount := 0
		var iceOptionsValue string
		for _, attr := range sd.Attributes {
			if attr.Key == AttrKeyICEOptions {
				iceOptionsCount++
				iceOptionsValue = attr.Value
			}
		}

		assert.Equal(t, 1, iceOptionsCount, "Should have exactly one ice-options attribute")
		assert.Equal(t, "renomination trickle", iceOptionsValue, "Should combine both values with space")
	})
}

func TestNewJSEPMediaDescription(t *testing.T) {
	md := NewJSEPMediaDescription("video", []string{"96", "97"})
	assert.NotNil(t, md)
	assert.Equal(t, "video", md.MediaName.Media)
	assert.Equal(t, int(9), md.MediaName.Port.Value)
	assert.Equal(t, []string{"UDP", "TLS", "RTP", "SAVPF"}, md.MediaName.Protos)
	assert.Equal(t, "IN", md.ConnectionInformation.NetworkType)
	assert.Equal(t, "IP4", md.ConnectionInformation.AddressType)
	assert.Equal(t, "0.0.0.0", md.ConnectionInformation.Address.Address)
}

func TestMediaDescriptionAttributes(t *testing.T) {
	md := NewJSEPMediaDescription("audio", nil)

	t.Run("WithPropertyAttribute", func(t *testing.T) {
		md = md.WithPropertyAttribute(AttrKeyRTCPMux)
		assert.Len(t, md.Attributes, 1)
		assert.Equal(t, AttrKeyRTCPMux, md.Attributes[0].Key)
	})

	t.Run("WithValueAttribute", func(t *testing.T) {
		md = md.WithValueAttribute(AttrKeyMID, "audio")
		assert.Len(t, md.Attributes, 2)
		assert.Equal(t, AttrKeyMID, md.Attributes[1].Key)
		assert.Equal(t, "audio", md.Attributes[1].Value)
	})

	t.Run("WithFingerprint", func(t *testing.T) {
		md = md.WithFingerprint("sha-256", "test-fingerprint")
		assert.Len(t, md.Attributes, 3)
		assert.Equal(t, "fingerprint", md.Attributes[2].Key)
		assert.Equal(t, "sha-256 test-fingerprint", md.Attributes[2].Value)
	})

	t.Run("WithICECredentials", func(t *testing.T) {
		md = md.WithICECredentials("test-ufrag", "test-pwd")
		assert.Len(t, md.Attributes, 5)
		assert.Equal(t, "ice-ufrag", md.Attributes[3].Key)
		assert.Equal(t, "test-ufrag", md.Attributes[3].Value)
		assert.Equal(t, "ice-pwd", md.Attributes[4].Key)
		assert.Equal(t, "test-pwd", md.Attributes[4].Value)
	})
}

func TestMediaDescriptionCodec(t *testing.T) {
	md := NewJSEPMediaDescription("audio", nil)

	t.Run("WithCodec", func(t *testing.T) {
		md = md.WithCodec(111, "opus", 48000, 2, "minptime=10;useinbandfec=1")
		assert.Len(t, md.MediaName.Formats, 1)
		assert.Equal(t, "111", md.MediaName.Formats[0])
		assert.Len(t, md.Attributes, 2)
		assert.Equal(t, "rtpmap", md.Attributes[0].Key)
		assert.Equal(t, "111 opus/48000/2", md.Attributes[0].Value)
		assert.Equal(t, "fmtp", md.Attributes[1].Key)
		assert.Equal(t, "111 minptime=10;useinbandfec=1", md.Attributes[1].Value)
	})

	t.Run("WithMediaSource", func(t *testing.T) {
		md = md.WithMediaSource(1234567890, "test-cname", "test-stream", "test-label")
		assert.Len(t, md.Attributes, 6)
		assert.Equal(t, "ssrc", md.Attributes[2].Key)
		assert.Equal(t, "1234567890 cname:test-cname", md.Attributes[2].Value)
		assert.Equal(t, "ssrc", md.Attributes[3].Key)
		assert.Equal(t, "1234567890 msid:test-stream test-label", md.Attributes[3].Value)
		assert.Equal(t, "ssrc", md.Attributes[4].Key)
		assert.Equal(t, "1234567890 mslabel:test-stream", md.Attributes[4].Value)
		assert.Equal(t, "ssrc", md.Attributes[5].Key)
		assert.Equal(t, "1234567890 label:test-label", md.Attributes[5].Value)
	})
}

func Test_extMapURI_TransportCC(t *testing.T) {
	m := extMapURI()
	u, ok := m[ExtMapValueTransportCC]
	assert.True(t, ok)
	assert.Equal(t, "http://www.ietf.org/id/draft-holmer-rmcat-transport-wide-cc-extensions-01", u)
}

func TestSessionDescription_WithMedia_Appends(t *testing.T) {
	sd, err := NewJSEPSessionDescription(false)
	assert.NoError(t, err)

	md := NewJSEPMediaDescription("audio", nil)
	prev := len(sd.MediaDescriptions)

	ret := sd.WithMedia(md)
	assert.Same(t, sd, ret)
	assert.Equal(t, prev+1, len(sd.MediaDescriptions))
	assert.Equal(t, md, sd.MediaDescriptions[len(sd.MediaDescriptions)-1])
}

func TestMediaDescription_WithExtMap_AddsPropertyAttribute(t *testing.T) {
	md := NewJSEPMediaDescription("audio", nil)

	u, _ := url.Parse(extMapURI()[ExtMapValueTransportCC])
	em := ExtMap{Value: ExtMapValueTransportCC, URI: u}

	ret := md.WithExtMap(em)
	assert.Same(t, md, ret)
	if assert.Len(t, md.Attributes, 1) {
		assert.Equal(t, "extmap:3 "+u.String(), md.Attributes[0].Key)
		assert.Empty(t, md.Attributes[0].Value)
	}
}

func TestMediaDescription_WithTransportCCExtMap_AddsExpectedAttribute(t *testing.T) {
	md := NewJSEPMediaDescription("audio", nil)

	ret := md.WithTransportCCExtMap()
	assert.Same(t, md, ret)
	if assert.Len(t, md.Attributes, 1) {
		want := "extmap:3 " + extMapURI()[ExtMapValueTransportCC]
		assert.Equal(t, want, md.Attributes[0].Key)
		assert.Empty(t, md.Attributes[0].Value)
	}
}
