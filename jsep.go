package sdp

import (
	"net"
)

func NewJsepSessionDescription() *SessionDescription {
	return &SessionDescription{
		Version: Version{Value: 0},
		Origin: Origin{
			Username:       "-",
			NetworkType:    "IN",
			AddressType:    "IP4",
			UnicastAddress: "0.0.0.0",
		},
		Session: Session{Value: "-"},
		TimeDescriptions: []TimeDescription{
			{
				Timing: Timing{
					StartTime: 0,
					StopTime:  0,
				},
				RepeatTimes: nil,
			},
		},
		Attributes: SessionAttributes{
			&IceOptions{Value: "trickle"},
		},
	}
}

func NewJsepRtpDescription() *MediaDescription {
	return &MediaDescription{
		Media: Media{
			Protos: []string{"UDP", "TLS", "RTP", "SAVPF"},
		},
		Connection: &Connection{
			NetworkType: "IN",
			AddressType: "IP4",
			Address: &Address{
				IP: net.ParseIP("0.0.0.0"),
			},
		},
	}
}

func NewJsepSctpDescription() *MediaDescription {
	return &MediaDescription{
		Media: Media{
			Type:    MediaTypeApplication,
			Protos:  []string{"UDP", "TLS", "RTP", "SAVPF"},
			Formats: []string{"webrtc-datachannel"},
		},
		Connection: &Connection{
			NetworkType: "IN",
			AddressType: "IP4",
			Address: &Address{
				IP: net.ParseIP("0.0.0.0"),
			},
		},
	}
}
