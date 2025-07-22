// SPDX-FileCopyrightText: 2023 The Pion community <https://pion.ly>
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
							Formats: []string{"0", "8"},
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
							Formats: []string{"0", "8"},
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
