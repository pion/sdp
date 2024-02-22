// SPDX-FileCopyrightText: 2023 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

package sdp

import (
	"errors"
	"reflect"
	"testing"
)

func getTestSessionDescription() SessionDescription {
	return SessionDescription{
		MediaDescriptions: []MediaDescription{
			{
				MediaName: MediaName{
					Media: kVideo,
					Port: RangedPort{
						Value: 51372,
					},
					Protos:  [][]byte{[]byte("RTP"), []byte("AVP")},
					Formats: [][]byte{[]byte("120"), []byte("121"), []byte("126"), []byte("97")},
				},
				Attributes: []Attribute{
					NewAttribute([]byte("fmtp"), []byte("126 profile-level-id=42e01f;level-asymmetry-allowed=1;packetization-mode=1")),
					NewAttribute([]byte("fmtp"), []byte("97 profile-level-id=42e01f;level-asymmetry-allowed=1")),
					NewAttribute([]byte("fmtp"), []byte("120 max-fs=12288;max-fr=60")),
					NewAttribute([]byte("fmtp"), []byte("121 max-fs=12288;max-fr=60")),
					NewAttribute([]byte("rtpmap"), []byte("120 VP8/90000")),
					NewAttribute([]byte("rtpmap"), []byte("121 VP9/90000")),
					NewAttribute([]byte("rtpmap"), []byte("126 H264/90000")),
					NewAttribute([]byte("rtpmap"), []byte("97 H264/90000")),
					NewAttribute([]byte("rtcp-fb"), []byte("97 ccm fir")),
					NewAttribute([]byte("rtcp-fb"), []byte("97 nack")),
					NewAttribute([]byte("rtcp-fb"), []byte("97 nack pli")),
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
				Name: []byte("VP8"),
			},
			Expected: 120,
		},
		{
			Codec: Codec{
				Name: []byte("VP9"),
			},
			Expected: 121,
		},
		{
			Codec: Codec{
				Name: []byte("H264"),
				Fmtp: []byte("profile-level-id=42e01f;level-asymmetry-allowed=1"),
			},
			Expected: 97,
		},
		{
			Codec: Codec{
				Name: []byte("H264"),
				Fmtp: []byte("level-asymmetry-allowed=1;profile-level-id=42e01f"),
			},
			Expected: 97,
		},
		{
			Codec: Codec{
				Name: []byte("H264"),
				Fmtp: []byte("profile-level-id=42e01f;level-asymmetry-allowed=1;packetization-mode=1"),
			},
			Expected: 126,
		},
	} {
		sd := getTestSessionDescription()

		actual, err := sd.GetPayloadTypeForCodec(test.Codec)
		if got, want := err, error(nil); !errors.Is(got, want) {
			t.Fatalf("GetPayloadTypeForCodec(): err=%v, want=%v", got, want)
		}

		if actual != test.Expected {
			t.Errorf("error:\n\nEXPECTED:\n%v\nACTUAL:\n%v", test.Expected, actual)
		}
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
				PayloadType: 120,
				Name:        []byte("VP8"),
				ClockRate:   90000,
				Fmtp:        []byte("max-fs=12288;max-fr=60"),
			},
		},
		{
			"vp9",
			getTestSessionDescription(),
			121,
			Codec{
				PayloadType: 121,
				Name:        []byte("VP9"),
				ClockRate:   90000,
				Fmtp:        []byte("max-fs=12288;max-fr=60"),
			},
		},
		{
			"h264 126",
			getTestSessionDescription(),
			126,
			Codec{
				PayloadType: 126,
				Name:        []byte("H264"),
				ClockRate:   90000,
				Fmtp:        []byte("profile-level-id=42e01f;level-asymmetry-allowed=1;packetization-mode=1"),
			},
		},
		{
			"h264 97",
			getTestSessionDescription(),
			97,
			Codec{
				PayloadType:  97,
				Name:         []byte("H264"),
				ClockRate:    90000,
				Fmtp:         []byte("profile-level-id=42e01f;level-asymmetry-allowed=1"),
				RTCPFeedback: [][]byte{[]byte("ccm fir"), []byte("nack"), []byte("nack pli")},
			},
		},
		{
			"pcmu without rtpmap",
			SessionDescription{
				MediaDescriptions: []MediaDescription{
					{
						MediaName: MediaName{
							Media:   []byte("audio"),
							Protos:  [][]byte{[]byte("RTP"), []byte("AVP")},
							Formats: [][]byte{[]byte("0"), []byte("8")},
						},
					},
				},
			},
			0,
			Codec{
				PayloadType: 0,
				Name:        []byte("PCMU"),
				ClockRate:   8000,
			},
		},
		{
			"pcma without rtpmap",
			SessionDescription{
				MediaDescriptions: []MediaDescription{
					{
						MediaName: MediaName{
							Media:   []byte("audio"),
							Protos:  [][]byte{[]byte("RTP"), []byte("AVP")},
							Formats: [][]byte{[]byte("0"), []byte("8")},
						},
					},
				},
			},
			8,
			Codec{
				PayloadType: 8,
				Name:        []byte("PCMA"),
				ClockRate:   8000,
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			actual, err := test.SD.GetCodecForPayloadType(test.PayloadType)
			if got, want := err, error(nil); !errors.Is(got, want) {
				t.Fatalf("GetCodecForPayloadType(): err=%v, want=%v", got, want)
			}

			if !reflect.DeepEqual(actual, test.Expected) {
				t.Errorf("error:\n\nEXPECTED:\n%v\nACTUAL:\n%v", test.Expected, actual)
			}
		})
	}
}

func TestNewSessionID(t *testing.T) {
	min := uint64(0x7FFFFFFFFFFFFFFF)
	max := uint64(0)
	for i := 0; i < 10000; i++ {
		r, err := newSessionID()
		if err != nil {
			t.Fatal(err)
		}
		if r > (1<<63)-1 {
			t.Fatalf("Session ID must be less than 2**64-1, got %d", r)
		}
		if r < min {
			min = r
		}
		if r > max {
			max = r
		}
	}
	if min > 0x1000000000000000 {
		t.Error("Value around lower boundary was not generated")
	}
	if max < 0x7000000000000000 {
		t.Error("Value around upper boundary was not generated")
	}
}
