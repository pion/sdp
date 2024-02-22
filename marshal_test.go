// SPDX-FileCopyrightText: 2023 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

package sdp

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	CanonicalMarshalSDP = "v=0\r\n" +
		"o=jdoe 2890844526 2890842807 IN IP4 10.47.16.5\r\n" +
		"s=SDP Seminar\r\n" +
		"i=A Seminar on the session description protocol\r\n" +
		"u=http://www.example.com/seminars/sdp.pdf\r\n" +
		"e=j.doe@example.com (Jane Doe)\r\n" +
		"p=+1 617 555-6011\r\n" +
		"c=IN IP4 224.2.17.12/127\r\n" +
		"b=X-YZ:128\r\n" +
		"b=AS:12345\r\n" +
		"t=2873397496 2873404696\r\n" +
		"t=3034423619 3042462419\r\n" +
		"r=604800 3600 0 90000\r\n" +
		"z=2882844526 -3600 2898848070 0\r\n" +
		"k=prompt\r\n" +
		"a=candidate:0 1 UDP 2113667327 203.0.113.1 54400 typ host\r\n" +
		"a=recvonly\r\n" +
		"m=audio 49170 RTP/AVP 0\r\n" +
		"i=Vivamus a posuere nisl\r\n" +
		"c=IN IP4 203.0.113.1\r\n" +
		"b=X-YZ:128\r\n" +
		"k=prompt\r\n" +
		"a=sendrecv\r\n" +
		"m=video 51372 RTP/AVP 99\r\n" +
		"a=rtpmap:99 h263-1998/90000\r\n"
)

func TestMarshalCanonical(t *testing.T) {
	sd := &SessionDescription{
		Version: 0,
		Origin: Origin{
			Username:       []byte("jdoe"),
			SessionID:      uint64(2890844526),
			SessionVersion: uint64(2890842807),
			NetworkType:    []byte("IN"),
			AddressType:    []byte("IP4"),
			UnicastAddress: []byte("10.47.16.5"),
		},
		SessionName:        []byte("SDP Seminar"),
		SessionInformation: Information("A Seminar on the session description protocol"),
		URI:                []byte("http://www.example.com/seminars/sdp.pdf"),
		EmailAddress:       EmailAddress("j.doe@example.com (Jane Doe)"),
		PhoneNumber:        PhoneNumber("+1 617 555-6011"),
		ConnectionInformation: ConnectionInformation{
			NetworkType: []byte("IN"),
			AddressType: []byte("IP4"),
			Address: Address{
				Address: []byte("224.2.17.12"),
				TTL:     127,
			},
		},
		Bandwidth: []Bandwidth{
			{
				Experimental: true,
				Type:         []byte("YZ"),
				Bandwidth:    128,
			},
			{
				Type:      []byte("AS"),
				Bandwidth: 12345,
			},
		},
		TimeDescriptions: []TimeDescription{
			{
				Timing: Timing{
					StartTime: 2873397496,
					StopTime:  2873404696,
				},
				RepeatTimes: nil,
			},
			{
				Timing: Timing{
					StartTime: 3034423619,
					StopTime:  3042462419,
				},
				RepeatTimes: []RepeatTime{
					{
						Interval: 604800,
						Duration: 3600,
						Offsets:  []int64{0, 90000},
					},
				},
			},
		},
		TimeZones: []TimeZone{
			{
				AdjustmentTime: 2882844526,
				Offset:         -3600,
			},
			{
				AdjustmentTime: 2898848070,
				Offset:         0,
			},
		},
		EncryptionKey: EncryptionKey("prompt"),
		Attributes: []Attribute{
			NewAttribute([]byte("candidate"), []byte("0 1 UDP 2113667327 203.0.113.1 54400 typ host")),
			NewAttribute([]byte("recvonly"), nil),
		},
		MediaDescriptions: []MediaDescription{
			{
				MediaName: MediaName{
					Media: []byte("audio"),
					Port: RangedPort{
						Value: 49170,
					},
					Protos:  [][]byte{[]byte("RTP"), []byte("AVP")},
					Formats: [][]byte{[]byte("0")},
				},
				MediaTitle: Information("Vivamus a posuere nisl"),
				ConnectionInformation: ConnectionInformation{
					NetworkType: []byte("IN"),
					AddressType: []byte("IP4"),
					Address: Address{
						Address: []byte("203.0.113.1"),
					},
				},
				Bandwidth: []Bandwidth{
					{
						Experimental: true,
						Type:         []byte("YZ"),
						Bandwidth:    128,
					},
				},
				EncryptionKey: EncryptionKey("prompt"),
				Attributes: []Attribute{
					NewAttribute([]byte("sendrecv"), nil),
				},
			},
			{
				MediaName: MediaName{
					Media: []byte("video"),
					Port: RangedPort{
						Value: 51372,
					},
					Protos:  [][]byte{[]byte("RTP"), []byte("AVP")},
					Formats: [][]byte{[]byte("99")},
				},
				Attributes: []Attribute{
					NewAttribute([]byte("rtpmap"), []byte("99 h263-1998/90000")),
				},
			},
		},
	}

	actual, err := sd.Marshal()
	require.NoError(t, err)
	require.Equal(t, CanonicalMarshalSDP, string(actual))
}

// var sink []byte

func BenchmarkMarshal(b *testing.B) {
	b.ReportAllocs()
	var sd SessionDescription
	err := sd.Unmarshal([]byte(BigSDP))
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		// sink, err = sd.Marshal()
		_, err = sd.Marshal()
		if err != nil {
			b.Fatal(err)
		}
	}
}
