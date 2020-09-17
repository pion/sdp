package sdp

import (
	"errors"
	"reflect"
	"testing"
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
					Formats: []string{"120", "121", "126", "97"},
				},
				Attributes: []Attribute{
					NewAttribute("fmtp:126 profile-level-id=42e01f;level-asymmetry-allowed=1;packetization-mode=1", ""),
					NewAttribute("fmtp:97 profile-level-id=42e01f;level-asymmetry-allowed=1", ""),
					NewAttribute("fmtp:120 max-fs=12288;max-fr=60", ""),
					NewAttribute("fmtp:121 max-fs=12288;max-fr=60", ""),
					NewAttribute("rtpmap:120 VP8/90000", ""),
					NewAttribute("rtpmap:121 VP9/90000", ""),
					NewAttribute("rtpmap:126 H264/90000", ""),
					NewAttribute("rtpmap:97 H264/90000", ""),
					NewAttribute("rtcp-fb:97 ccm fir", ""),
					NewAttribute("rtcp-fb:97 nack", ""),
					NewAttribute("rtcp-fb:97 nack pli", ""),
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
		PayloadType uint8
		Expected    Codec
	}{
		{
			PayloadType: 120,
			Expected: Codec{
				PayloadType: 120,
				Name:        "VP8",
				ClockRate:   90000,
				Fmtp:        "max-fs=12288;max-fr=60",
			},
		},
		{
			PayloadType: 121,
			Expected: Codec{
				PayloadType: 121,
				Name:        "VP9",
				ClockRate:   90000,
				Fmtp:        "max-fs=12288;max-fr=60",
			},
		},
		{
			PayloadType: 126,
			Expected: Codec{
				PayloadType: 126,
				Name:        "H264",
				ClockRate:   90000,
				Fmtp:        "profile-level-id=42e01f;level-asymmetry-allowed=1;packetization-mode=1",
			},
		},
		{
			PayloadType: 97,
			Expected: Codec{
				PayloadType:  97,
				Name:         "H264",
				ClockRate:    90000,
				Fmtp:         "profile-level-id=42e01f;level-asymmetry-allowed=1",
				RTCPFeedback: []string{"ccm fir", "nack", "nack pli"},
			},
		},
	} {
		sd := getTestSessionDescription()

		actual, err := sd.GetCodecForPayloadType(test.PayloadType)
		if got, want := err, error(nil); !errors.Is(got, want) {
			t.Fatalf("GetCodecForPayloadType(): err=%v, want=%v", got, want)
		}

		if !reflect.DeepEqual(actual, test.Expected) {
			t.Errorf("error:\n\nEXPECTED:\n%v\nACTUAL:\n%v", test.Expected, actual)
		}
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
