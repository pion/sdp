// SPDX-FileCopyrightText: 2026 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func getTestSessionDescription() SessionDescription {
	return SessionDescription{
		MediaDescriptions: []*MediaDescription{
			{
				MediaName: MediaName{
					Media: "video",
					Port: RangedPort{
						Value: 51372,
					},
					Protos:  []string{"RTP", "AVP"},
					Formats: []string{"120", "121", "126", "97", "98"},
				},
				Attributes: []Attribute{
					NewAttribute("fmtp:126 profile-level-id=42e01f;level-asymmetry-allowed=1;packetization-mode=1", ""),
					NewAttribute("fmtp:97 profile-level-id=42e01f;level-asymmetry-allowed=1", ""),
					NewAttribute("fmtp:98 profile-level-id=42e01e; packetization-mode=1", ""),
					NewAttribute("fmtp:120 max-fs=12288;max-fr=60", ""),
					NewAttribute("fmtp:121 max-fs=12288;max-fr=60", ""),
					NewAttribute("rtpmap:120 VP8/90000", ""),
					NewAttribute("rtpmap:121 VP9/90000", ""),
					NewAttribute("rtpmap:126 H264/90000", ""),
					NewAttribute("rtpmap:97 H264/90000", ""),
					NewAttribute("rtpmap:98 H264/90000", ""),
					NewAttribute("rtcp-fb:97 ccm fir", ""),
					NewAttribute("rtcp-fb:97 nack", ""),
					NewAttribute("rtcp-fb:97 nack pli", ""),
					NewAttribute("rtcp-fb:* transport-cc", ""),
					NewAttribute("rtcp-fb:* nack", ""),
				},
			},
		},
	}
}

func TestGetPayloadTypeForVP8(t *testing.T) {
	for _, test := range []struct {
		Codec    Codec
		Expected uint8
	}{
		{
			Codec: Codec{
				Name: "VP8",
			},
			Expected: 120,
		},
		{
			Codec: Codec{
				Name: "VP9",
			},
			Expected: 121,
		},
		{
			Codec: Codec{
				Name: "H264",
				Fmtp: "profile-level-id=42e01f;level-asymmetry-allowed=1",
			},
			Expected: 97,
		},
		{
			Codec: Codec{
				Name: "H264",
				Fmtp: "level-asymmetry-allowed=1;profile-level-id=42e01f",
			},
			Expected: 97,
		},
		{
			Codec: Codec{
				Name: "H264",
				Fmtp: "profile-level-id=42e01f;level-asymmetry-allowed=1;packetization-mode=1",
			},
			Expected: 126,
		},
	} {
		sd := getTestSessionDescription()

		actual, err := sd.GetPayloadTypeForCodec(test.Codec)
		assert.NoError(t, err)
		assert.Equal(t, actual, test.Expected)
	}
}

func TestGetCodecForPayloadType(t *testing.T) {
	for _, test := range []struct {
		name        string
		SD          SessionDescription
		PayloadType uint8
		Expected    Codec
	}{
		{
			"vp8",
			getTestSessionDescription(),
			120,
			Codec{
				PayloadType:  120,
				Name:         "VP8",
				ClockRate:    90000,
				Fmtp:         "max-fs=12288;max-fr=60",
				RTCPFeedback: []string{"transport-cc", "nack"},
			},
		},
		{
			"vp9",
			getTestSessionDescription(),
			121,
			Codec{
				PayloadType:  121,
				Name:         "VP9",
				ClockRate:    90000,
				Fmtp:         "max-fs=12288;max-fr=60",
				RTCPFeedback: []string{"transport-cc", "nack"},
			},
		},
		{
			"h264 126",
			getTestSessionDescription(),
			126,
			Codec{
				PayloadType:  126,
				Name:         "H264",
				ClockRate:    90000,
				Fmtp:         "profile-level-id=42e01f;level-asymmetry-allowed=1;packetization-mode=1",
				RTCPFeedback: []string{"transport-cc", "nack"},
			},
		},
		{
			"h264 97",
			getTestSessionDescription(),
			97,
			Codec{
				PayloadType:  97,
				Name:         "H264",
				ClockRate:    90000,
				Fmtp:         "profile-level-id=42e01f;level-asymmetry-allowed=1",
				RTCPFeedback: []string{"ccm fir", "nack", "nack pli", "transport-cc"},
			},
		},
		{
			"h264 98",
			getTestSessionDescription(),
			98,
			Codec{
				PayloadType:  98,
				Name:         "H264",
				ClockRate:    90000,
				Fmtp:         "profile-level-id=42e01e; packetization-mode=1",
				RTCPFeedback: []string{"transport-cc", "nack"},
			},
		},
		{
			"pcmu without rtpmap",
			SessionDescription{
				MediaDescriptions: []*MediaDescription{
					{
						MediaName: MediaName{
							Media:   "audio",
							Protos:  []string{"RTP", "AVP"},
							Formats: []string{"0", "8", "9"},
						},
					},
				},
			},
			0,
			Codec{
				PayloadType: 0,
				Name:        "PCMU",
				ClockRate:   8000,
			},
		},
		{
			"pcma without rtpmap",
			SessionDescription{
				MediaDescriptions: []*MediaDescription{
					{
						MediaName: MediaName{
							Media:   "audio",
							Protos:  []string{"RTP", "AVP"},
							Formats: []string{"0", "8", "9"},
						},
					},
				},
			},
			8,
			Codec{
				PayloadType: 8,
				Name:        "PCMA",
				ClockRate:   8000,
			},
		},
		{
			"g722 without rtpmap",
			SessionDescription{
				MediaDescriptions: []*MediaDescription{
					{
						MediaName: MediaName{
							Media:   "audio",
							Protos:  []string{"RTP", "AVP"},
							Formats: []string{"0", "8", "9"},
						},
					},
				},
			},
			9,
			Codec{
				PayloadType: 9,
				Name:        "G722",
				ClockRate:   8000,
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			actual, err := test.SD.GetCodecForPayloadType(test.PayloadType)
			assert.NoError(t, err)
			assert.Equal(t, actual, test.Expected)
		})
	}
}

func TestGetCodecsForPayloadTypes(t *testing.T) {
	for _, test := range []struct {
		name         string
		SD           SessionDescription
		PayloadTypes []uint8
		Expected     []Codec
	}{
		{
			"vp8-9",
			getTestSessionDescription(),
			[]uint8{120, 121},
			[]Codec{
				{
					PayloadType:  120,
					Name:         "VP8",
					ClockRate:    90000,
					Fmtp:         "max-fs=12288;max-fr=60",
					RTCPFeedback: []string{"transport-cc", "nack"},
				},
				{
					PayloadType:  121,
					Name:         "VP9",
					ClockRate:    90000,
					Fmtp:         "max-fs=12288;max-fr=60",
					RTCPFeedback: []string{"transport-cc", "nack"},
				},
			},
		},
		{
			"pcma without rtpmap",
			SessionDescription{
				MediaDescriptions: []*MediaDescription{
					{
						MediaName: MediaName{
							Media:   "audio",
							Protos:  []string{"RTP", "AVP"},
							Formats: []string{"0", "8"},
						},
					},
				},
			},
			[]uint8{0, 8},
			[]Codec{
				{
					PayloadType: 0,
					Name:        "PCMU",
					ClockRate:   8000,
				},
				{
					PayloadType: 8,
					Name:        "PCMA",
					ClockRate:   8000,
				},
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			actual, err := test.SD.GetCodecsForPayloadTypes(test.PayloadTypes)
			assert.NoError(t, err)
			assert.Equal(t, actual, test.Expected)
		})
	}
}

func TestNewSessionID(t *testing.T) {
	minVal := uint64(0x7FFFFFFFFFFFFFFF)
	maxVal := uint64(0)
	for i := 0; i < 10000; i++ {
		r, err := newSessionID()
		assert.NoError(t, err)
		assert.Lessf(t, r, uint64((1<<64)-1), "Session ID must be less than 2**64-1, got %d", r)
		if r < minVal {
			minVal = r
		}
		if r > maxVal {
			maxVal = r
		}
	}

	assert.Less(t, minVal, uint64(0x1000000000000000), "Value around upper boundary was not generated")
	assert.Greater(t, maxVal, uint64(0x7000000000000000), "Value around lower boundary was not generated")
}

func TestConnectionRole_String(t *testing.T) {
	assert.Equal(t, "active", ConnectionRoleActive.String())
	assert.Equal(t, "passive", ConnectionRolePassive.String())
	assert.Equal(t, "actpass", ConnectionRoleActpass.String())
	assert.Equal(t, "holdconn", ConnectionRoleHoldconn.String())

	var zero ConnectionRole
	assert.Equal(t, "Unknown", zero.String())

	var bogus ConnectionRole = 99
	assert.Equal(t, "Unknown", bogus.String())
}

func TestCodec_String(t *testing.T) {
	c := Codec{
		PayloadType:        111,
		Name:               "opus",
		ClockRate:          48000,
		EncodingParameters: "2",
		Fmtp:               "minptime=10;useinbandfec=1",
		RTCPFeedback:       []string{"nack", "pli"},
	}

	got := c.String()
	assert.Equal(t, "111 opus/48000/2 (minptime=10;useinbandfec=1) [nack, pli]", got)
}

func TestParseRtpmap_NoEncodingParams(t *testing.T) {
	codec, err := parseRtpmap("rtpmap:111 opus/48000")
	assert.NoError(t, err)

	assert.Equal(t, uint8(111), codec.PayloadType)
	assert.Equal(t, "opus", codec.Name)
	assert.Equal(t, uint32(48000), codec.ClockRate)
	assert.Equal(t, "", codec.EncodingParameters)
}

func TestParseRtpmap_WithEncodingParams(t *testing.T) {
	codec, err := parseRtpmap("rtpmap:96 MP4A-LATM/44100/2")
	assert.NoError(t, err)

	assert.Equal(t, uint8(96), codec.PayloadType)
	assert.Equal(t, "MP4A-LATM", codec.Name)
	assert.Equal(t, uint32(44100), codec.ClockRate)
	assert.Equal(t, "2", codec.EncodingParameters)
}

func TestParseRtpmap_Error_MissingSpace(t *testing.T) {
	// missing space
	_, err := parseRtpmap("rtpmap:111")
	assert.ErrorIs(t, err, errExtractCodecRtpmap)
}

func TestParseRtpmap_Error_MissingColon(t *testing.T) {
	// missing colon
	_, err := parseRtpmap("rtpmap111 opus/48000")
	assert.ErrorIs(t, err, errExtractCodecRtpmap)
}

func TestParseRtpmap_Error_NonNumericPayload(t *testing.T) {
	// non-numeric
	_, err := parseRtpmap("rtpmap:xx opus/48000")
	assert.ErrorIs(t, err, errExtractCodecRtpmap)
}

func TestParseRtpmap_Error_NonNumericClockRate(t *testing.T) {
	_, err := parseRtpmap("rtpmap:111 opus/notanumber")
	assert.ErrorIs(t, err, errExtractCodecRtpmap)
}

func TestParseFmtp_Error_MissingSpace(t *testing.T) {
	_, err := parseFmtp("fmtp:111")
	assert.ErrorIs(t, err, errExtractCodecFmtp)
}

func TestParseFmtp_Error_MissingColon(t *testing.T) {
	_, err := parseFmtp("fmtp111 a=b;c=d")
	assert.ErrorIs(t, err, errExtractCodecFmtp)
}

func TestParseFmtp_Error_NonNumericPayload(t *testing.T) {
	_, err := parseFmtp("fmtp:xx profile-level-id=42e01f;packetization-mode=1")
	assert.ErrorIs(t, err, errExtractCodecFmtp)
}

func TestEquivalentFmtp_MismatchAfterSortAndTrim(t *testing.T) {
	want := "profile-level-id=42e01f; packetization-mode=1"
	got := "packetization-mode=0;  profile-level-id=42e01f"

	assert.False(t, equivalentFmtp(want, got))
}

func TestParseRtcpFb_MissingSpace(t *testing.T) {
	c, wildcard, err := parseRtcpFb("rtcp-fb:97")
	assert.ErrorIs(t, err, errExtractCodecRtcpFb)
	assert.False(t, wildcard)
	assert.Equal(t, uint8(0), c.PayloadType)
	assert.Empty(t, c.RTCPFeedback)
}

func TestParseRtcpFb_MissingColon(t *testing.T) {
	c, wildcard, err := parseRtcpFb("rtcp-fb97 nack")
	assert.ErrorIs(t, err, errExtractCodecRtcpFb)
	assert.False(t, wildcard)
	assert.Equal(t, uint8(0), c.PayloadType)
	assert.Empty(t, c.RTCPFeedback)
}

func TestParseRtcpFb_NonNumeric(t *testing.T) {
	c, wildcard, err := parseRtcpFb("rtcp-fb:xx nack")
	assert.Error(t, err)
	assert.False(t, wildcard)
	assert.Equal(t, uint8(0), c.PayloadType)
	assert.Empty(t, c.RTCPFeedback)
}

func TestBuildCodecMap_RtcpFbError(t *testing.T) {
	sd := SessionDescription{
		MediaDescriptions: []*MediaDescription{
			{
				Attributes: []Attribute{
					// non-numeric should return an error
					NewAttribute("rtcp-fb:xx nack", ""),
				},
			},
		},
	}

	codecs := sd.buildCodecMap()

	// the three static codecs should be present, unchanged.
	if assert.Len(t, codecs, 3) {
		if c, ok := codecs[0]; assert.True(t, ok) {
			assert.Equal(t, uint8(0), c.PayloadType)
			assert.Equal(t, "PCMU", c.Name)
			assert.Equal(t, uint32(8000), c.ClockRate)
			assert.Empty(t, c.RTCPFeedback)
		}
		if c, ok := codecs[8]; assert.True(t, ok) {
			assert.Equal(t, uint8(8), c.PayloadType)
			assert.Equal(t, "PCMA", c.Name)
			assert.Equal(t, uint32(8000), c.ClockRate)
			assert.Empty(t, c.RTCPFeedback)
		}
		if c, ok := codecs[9]; assert.True(t, ok) {
			assert.Equal(t, uint8(9), c.PayloadType)
			assert.Equal(t, "G722", c.Name)
			assert.Equal(t, uint32(8000), c.ClockRate)
			assert.Empty(t, c.RTCPFeedback)
		}
	}
}

func TestCodecsMatch_MiddleFalse_ClockRateMismatch(t *testing.T) {
	expected := Codec{ClockRate: 44100}
	actual := Codec{Name: "opus", ClockRate: 48000}
	assert.False(t, codecsMatch(expected, actual))
}

func TestCodecsMatch_MiddleFalse_EncodingParamsMismatch(t *testing.T) {
	expected := Codec{EncodingParameters: "1"}
	actual := Codec{EncodingParameters: "2"}
	assert.False(t, codecsMatch(expected, actual))
}

func TestGetCodecForPayloadType_Error_NotFound(t *testing.T) {
	var sd SessionDescription
	_, err := sd.GetCodecForPayloadType(42)
	assert.ErrorIs(t, err, errPayloadTypeNotFound)
}

func TestGetPayloadTypeForCodec_Error_NotFound(t *testing.T) {
	var sd SessionDescription
	_, err := sd.GetPayloadTypeForCodec(Codec{Name: "doesnotexist"})
	assert.ErrorIs(t, err, errCodecNotFound)
}

func TestLexer_HandleType_ElseIfErrorFromReadType(t *testing.T) {
	// should cause a syntaxError (key='a', err=syntaxError) because second byte != '='
	l := &lexer{baseLexer: baseLexer{value: "a-"}}

	called := false
	fn := func(key byte) stateFn {
		called = true // should not be called because handleType returns early on err != nil

		return func(*lexer) (stateFn, error) { return nil, nil }
	}

	st, err := l.handleType(fn)
	assert.Nil(t, st)
	assert.False(t, called)

	var se syntaxError
	assert.ErrorAs(t, err, &se)
}

func TestLexer_HandleType_SyntaxErrorWhenFnReturnsNil(t *testing.T) {
	// valid type "a=" so readType returns nil error so fn returns nil
	l := &lexer{baseLexer: baseLexer{value: "a="}}

	fn := func(key byte) stateFn {
		assert.Equal(t, byte('a'), key)

		return nil // should trigger final syntaxError
	}

	st, err := l.handleType(fn)
	assert.Nil(t, st)

	var se syntaxError
	assert.ErrorAs(t, err, &se)
}
