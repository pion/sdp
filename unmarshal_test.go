// SPDX-FileCopyrightText: 2026 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

package sdp

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	BaseSDP = "v=0\r\n" +
		"o=jdoe 2890844526 2890842807 IN IP4 10.47.16.5\r\n" +
		"s=SDP Seminar\r\n"

	SessionInformationSDP = BaseSDP +
		"i=A Seminar on the session description protocol\r\n" +
		"t=3034423619 3042462419\r\n"

	// https://tools.ietf.org/html/rfc4566#section-5
	// Parsers SHOULD be tolerant and also accept records terminated
	// with a single newline character.
	SessionInformationSDPLFOnly = "v=0\n" +
		"o=jdoe 2890844526 2890842807 IN IP4 10.47.16.5\n" +
		"s=SDP Seminar\n" +
		"i=A Seminar on the session description protocol\n" +
		"t=3034423619 3042462419\n"

	// SessionInformationSDPCROnly = "v=0\r" +
	// 	"o=jdoe 2890844526 2890842807 IN IP4 10.47.16.5\r" +
	// 	"s=SDP Seminar\r"
	// 	"i=A Seminar on the session description protocol\r" +
	// 	"t=3034423619 3042462419\r"

	// Other SDP parsers (e.g. one in VLC media player) allow
	// empty lines.
	SessionInformationSDPExtraCRLF = "v=0\r\n" +
		"o=jdoe 2890844526 2890842807 IN IP4 10.47.16.5\r\n" +
		"\r\n" +
		"s=SDP Seminar\r\n" +
		"\r\n" +
		"i=A Seminar on the session description protocol\r\n" +
		"\r\n" +
		"t=3034423619 3042462419\r\n" +
		"\r\n"

	URISDP = BaseSDP +
		"u=http://www.example.com/seminars/sdp.pdf\r\n" +
		"t=3034423619 3042462419\r\n"

	EmailAddressSDP = BaseSDP +
		"e=j.doe@example.com (Jane Doe)\r\n" +
		"t=3034423619 3042462419\r\n"

	PhoneNumberSDP = BaseSDP +
		"p=+1 617 555-6011\r\n" +
		"t=3034423619 3042462419\r\n"

	SessionConnectionInformationSDP = BaseSDP +
		"c=IN IP4 224.2.17.12/127\r\n" +
		"t=3034423619 3042462419\r\n"

	SessionBandwidthSDP = BaseSDP +
		"b=X-YZ:128\r\n" +
		"b=AS:12345\r\n" +
		"t=3034423619 3042462419\r\n"

	TimingSDP = BaseSDP +
		"t=2873397496 2873404696\r\n"

	// Short hand time notation is converted into NTP timestamp format in
	// seconds. Because of that unittest comparisons will fail as the same time
	// will be expressed in different units.
	RepeatTimesSDP = TimingSDP +
		"r=604800 3600 0 90000\r\n" +
		"r=3d 2h 0 21h\r\n"

	RepeatTimesSDPExpected = TimingSDP +
		"r=604800 3600 0 90000\r\n" +
		"r=259200 7200 0 75600\r\n"

	RepeatTimesSDPExtraCRLF = RepeatTimesSDPExpected +
		"\r\n"

	// The expected value looks a bit different for the same reason as mentioned
	// above regarding RepeatTimes.
	TimeZonesSDP = TimingSDP +
		"r=2882844526 -1h 2898848070 0\r\n"

	TimeZonesSDPExpected = TimingSDP +
		"r=2882844526 -3600 2898848070 0\r\n"

	TimeZonesSDP2 = TimingSDP +
		"z=2882844526 -3600 2898848070 0\r\n"

	TimeZonesSDP2ExtraCRLF = TimeZonesSDP2 +
		"\r\n"

	SessionEncryptionKeySDP = TimingSDP +
		"k=prompt\r\n"

	SessionEncryptionKeySDPExtraCRLF = SessionEncryptionKeySDP +
		"\r\n"

	SessionAttributesSDP = TimingSDP +
		"a=rtpmap:96 opus/48000\r\n"

	MediaNameSDP = TimingSDP +
		"m=video 51372 RTP/AVP 99\r\n" +
		"m=audio 54400 RTP/SAVPF 0 96\r\n" +
		"m=message 5028 TCP/MSRP *\r\n"

	MediaNameSDPExtraCRLF = MediaNameSDP +
		"\r\n"

	MediaTitleSDP = MediaNameSDP +
		"i=Vivamus a posuere nisl\r\n"

	MediaConnectionInformationSDP = MediaNameSDP +
		"c=IN IP4 203.0.113.1\r\n"

	MediaConnectionInformationSDPExtraCRLF = MediaConnectionInformationSDP +
		"\r\n"

	MediaDescriptionOutOfOrderSDP = MediaNameSDP +
		"a=rtpmap:99 h263-1998/90000\r\n" +
		"a=candidate:0 1 UDP 2113667327 203.0.113.1 54400 typ host\r\n" +
		"c=IN IP4 203.0.113.1\r\n" +
		"i=Vivamus a posuere nisl\r\n"

	MediaDescriptionOutOfOrderSDPActual = MediaNameSDP +
		"i=Vivamus a posuere nisl\r\n" +
		"c=IN IP4 203.0.113.1\r\n" +
		"a=rtpmap:99 h263-1998/90000\r\n" +
		"a=candidate:0 1 UDP 2113667327 203.0.113.1 54400 typ host\r\n"

	MediaBandwidthSDP = MediaNameSDP +
		"b=X-YZ:128\r\n" +
		"b=AS:12345\r\n" +
		"b=TIAS:12345\r\n" +
		"b=RS:12345\r\n" +
		"b=RR:12345\r\n"

	MediaEncryptionKeySDP = MediaNameSDP +
		"k=prompt\r\n"

	MediaEncryptionKeySDPExtraCRLF = MediaEncryptionKeySDP +
		"\r\n"

	MediaAttributesSDP = MediaNameSDP +
		"a=rtpmap:99 h263-1998/90000\r\n" +
		"a=candidate:0 1 UDP 2113667327 203.0.113.1 54400 typ host\r\n" +
		"a=rtcp-fb:97 ccm fir\r\n" +
		"a=rtcp-fb:97 nack\r\n" +
		"a=rtcp-fb:97 nack pli\r\n"

	MediaBfcpSDP = TimingSDP +
		"m=application 3238 UDP/BFCP *\r\n" +
		"a=sendrecv\r\n" +
		"a=setup:actpass\r\n" +
		"a=connection:new\r\n" +
		"a=floorctrl:c-s\r\n"

	MediaCubeSDP = TimingSDP +
		"m=application 2455 UDP/UDT/IX *\r\n" +
		"a=ixmap:0 ping\r\n" +
		"a=ixmap:2 xccp\r\n"

	MediaTCPMRCPv2 = TimingSDP +
		"m=application 1544 TCP/MRCPv2 1\r\n"

	MediaTCPTLSMRCPv2 = TimingSDP +
		"m=application 1544 TCP/TLS/MRCPv2 1\r\n"

	MediaFECSDP = TimingSDP +
		"m=application 50178 UDP/FEC *\r\n" +
		"a=fec-repair-flow:encoding-id=8; fssi=E:1234,S:0,m:7\r\n" +
		"a=repair-window:500ms\r\n"

	CanonicalUnmarshalSDP = "v=0\r\n" +
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

func TestRoundTrip(t *testing.T) {
	for _, test := range []struct {
		Name   string
		SDP    string
		Actual string
	}{
		{
			Name:   "SessionInformationSDPLFOnly",
			SDP:    SessionInformationSDPLFOnly,
			Actual: SessionInformationSDP,
		},
		// {
		// 	Name:   "SessionInformationSDPCROnly",
		// 	SDP:    SessionInformationSDPCROnly,
		// 	Actual: SessionInformationSDPBaseSDP,
		// },
		{
			Name:   "SessionInformationSDPExtraCRLF",
			SDP:    SessionInformationSDPExtraCRLF,
			Actual: SessionInformationSDP,
		},
		{
			Name: "SessionInformation",
			SDP:  SessionInformationSDP,
		},
		{
			Name: "URI",
			SDP:  URISDP,
		},
		{
			Name: "EmailAddress",
			SDP:  EmailAddressSDP,
		},
		{
			Name: "PhoneNumber",
			SDP:  PhoneNumberSDP,
		},
		{
			Name:   "RepeatTimesSDPExtraCRLF",
			SDP:    RepeatTimesSDPExtraCRLF,
			Actual: RepeatTimesSDPExpected,
		},
		{
			Name: "SessionConnectionInformation",
			SDP:  SessionConnectionInformationSDP,
		},
		{
			Name: "SessionBandwidth",
			SDP:  SessionBandwidthSDP,
		},
		{
			Name: "SessionEncryptionKey",
			SDP:  SessionEncryptionKeySDP,
		},
		{
			Name:   "SessionEncryptionKeyExtraCRLF",
			SDP:    SessionEncryptionKeySDPExtraCRLF,
			Actual: SessionEncryptionKeySDP,
		},
		{
			Name: "SessionAttributes",
			SDP:  SessionAttributesSDP,
		},
		{
			Name:   "TimeZonesSDP2ExtraCRLF",
			SDP:    TimeZonesSDP2ExtraCRLF,
			Actual: TimeZonesSDP2,
		},
		{
			Name: "MediaName",
			SDP:  MediaNameSDP,
		},
		{
			Name:   "MediaNameExtraCRLF",
			SDP:    MediaNameSDPExtraCRLF,
			Actual: MediaNameSDP,
		},
		{
			Name: "MediaTitle",
			SDP:  MediaTitleSDP,
		},
		{
			Name: "MediaConnectionInformation",
			SDP:  MediaConnectionInformationSDP,
		},
		{
			Name:   "MediaConnectionInformationExtraCRLF",
			SDP:    MediaConnectionInformationSDPExtraCRLF,
			Actual: MediaConnectionInformationSDP,
		},
		{
			Name:   "MediaDescriptionOutOfOrder",
			SDP:    MediaDescriptionOutOfOrderSDP,
			Actual: MediaDescriptionOutOfOrderSDPActual,
		},
		{
			Name: "MediaBandwidth",
			SDP:  MediaBandwidthSDP,
		},
		{
			Name: "MediaEncryptionKey",
			SDP:  MediaEncryptionKeySDP,
		},
		{
			Name:   "MediaEncryptionKeyExtraCRLF",
			SDP:    MediaEncryptionKeySDPExtraCRLF,
			Actual: MediaEncryptionKeySDP,
		},
		{
			Name: "MediaAttributes",
			SDP:  MediaAttributesSDP,
		},
		{
			Name: "CanonicalUnmarshal",
			SDP:  CanonicalUnmarshalSDP,
		},
		{
			Name: "MediaBfcpSDP",
			SDP:  MediaBfcpSDP,
		},
		{
			Name: "MediaCubeSDP",
			SDP:  MediaCubeSDP,
		},
		{
			Name: "MediaTCPMRCPv2",
			SDP:  MediaTCPMRCPv2,
		},
		{
			Name: "MediaTCPTLSMRCPv2",
			SDP:  MediaTCPTLSMRCPv2,
		},
		{
			Name: "MediaFEC",
			SDP:  MediaFECSDP,
		},
	} {
		test := test
		t.Run(test.Name, func(t *testing.T) {
			sd := &SessionDescription{}

			err := sd.UnmarshalString(test.SDP)
			assert.NoError(t, err)

			actual, err := sd.Marshal()
			assert.NoError(t, err)

			want := test.SDP
			if test.Actual != "" {
				want = test.Actual
			}

			assert.Equal(t, want, string(actual))
		})
	}
}

func TestUnmarshalRepeatTimes(t *testing.T) {
	sd := &SessionDescription{}
	assert.NoError(t, sd.UnmarshalString(RepeatTimesSDP))

	actual, err := sd.Marshal()
	assert.NoError(t, err)
	assert.Equal(t, RepeatTimesSDPExpected, string(actual))

	err = sd.UnmarshalString(TimingSDP + "r=\r\n")
	assert.ErrorIs(t, err, errSDPInvalidValue)
}

func TestUnmarshalTimeZones(t *testing.T) {
	sd := &SessionDescription{}
	assert.NoError(t, sd.UnmarshalString(TimeZonesSDP))

	actual, err := sd.Marshal()
	assert.NoError(t, err)
	assert.Equal(t, TimeZonesSDPExpected, string(actual))
}

func TestUnmarshalNonNilAddress(t *testing.T) {
	in := "v=0\r\no=0 0 0 IN IP4 0\r\ns=0\r\nc=IN IP4\r\nt=0 0\r\n"
	var sd SessionDescription
	err := sd.UnmarshalString(in)
	assert.NoError(t, err)

	out, err := sd.Marshal()
	assert.NoError(t, err)
	assert.Equal(t, in, string(out))
}

func TestUnmarshalZeroValues(t *testing.T) {
	in := "v=0\r\no=0 0 0 IN IP4 0\r\ns=\r\nt=0 0\r\n"
	var sd SessionDescription
	assert.NoError(t, sd.UnmarshalString(in))

	out, err := sd.Marshal()
	assert.NoError(t, err)
	assert.Equal(t, in, string(out))
}

func TestUnmarshalPortRange(t *testing.T) {
	for _, test := range []struct {
		In          string
		ExpectError error
	}{
		{
			In:          SessionAttributesSDP + "m=video -1 RTP/AVP 99\r\n",
			ExpectError: errSDPInvalidPortValue,
		},
		{
			In:          SessionAttributesSDP + "m=video 65536 RTP/AVP 99\r\n",
			ExpectError: errSDPInvalidPortValue,
		},
		{
			In:          SessionAttributesSDP + "m=video 0 RTP/AVP 99\r\n",
			ExpectError: nil,
		},
		{
			In:          SessionAttributesSDP + "m=video 65535 RTP/AVP 99\r\n",
			ExpectError: nil,
		},
		{
			In:          SessionAttributesSDP + "m=video --- RTP/AVP 99\r\n",
			ExpectError: errSDPInvalidPortValue,
		},
	} {
		var sd SessionDescription
		err := sd.UnmarshalString(test.In)
		if test.ExpectError != nil {
			assert.ErrorIs(t, err, test.ExpectError)
		} else {
			assert.NoError(t, err)
		}
	}
}

func BenchmarkUnmarshal(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var sd SessionDescription
		err := sd.UnmarshalString(CanonicalUnmarshalSDP)
		assert.NoError(b, err)
	}
}

func TestUnmarshalOriginIncomplete(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Origin
	}{
		{
			name:  "missing unicast address - Uniview camera case",
			input: "v=0\r\no=- 1001 1 IN IP4\r\ns=VCP IPC Realtime stream\r\nt=0 0\r\n",
			expected: Origin{
				Username:       "-",
				SessionID:      1001,
				SessionVersion: 1,
				NetworkType:    "IN",
				AddressType:    "IP4",
				UnicastAddress: "0.0.0.0",
			},
		},
		{
			name:  "missing address type and address",
			input: "v=0\r\no=- 1001 1 IN\r\ns=Test Stream\r\nt=0 0\r\n",
			expected: Origin{
				Username:       "-",
				SessionID:      1001,
				SessionVersion: 1,
				NetworkType:    "IN",
				AddressType:    "IP4",
				UnicastAddress: "0.0.0.0",
			},
		},
		{
			name:  "IPv6 missing address",
			input: "v=0\r\no=- 1001 1 IN IP6\r\ns=Test Stream\r\nt=0 0\r\n",
			expected: Origin{
				Username:       "-",
				SessionID:      1001,
				SessionVersion: 1,
				NetworkType:    "IN",
				AddressType:    "IP6",
				UnicastAddress: "::",
			},
		},
		{
			name:  "complete origin line - should work as before",
			input: "v=0\r\no=jdoe 2890844526 2890842807 IN IP4 10.47.16.5\r\ns=SDP Seminar\r\nt=3034423619 3042462419\r\n",
			expected: Origin{
				Username:       "jdoe",
				SessionID:      2890844526,
				SessionVersion: 2890842807,
				NetworkType:    "IN",
				AddressType:    "IP4",
				UnicastAddress: "10.47.16.5",
			},
		},
		{
			name:  "empty address field",
			input: "v=0\r\no=- 1001 1 IN IP4 \r\ns=Test\r\nt=0 0\r\n",
			expected: Origin{
				Username:       "-",
				SessionID:      1001,
				SessionVersion: 1,
				NetworkType:    "IN",
				AddressType:    "IP4",
				UnicastAddress: "0.0.0.0",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var sd SessionDescription
			err := sd.UnmarshalString(test.input)
			assert.NoError(t, err)
			assert.Equal(t, test.expected, sd.Origin)
		})
	}
}

func TestUnmarshalOriginInvalidFields(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "invalid network type",
			input: "v=0\r\no=- 1001 1 INVALID IP4 10.0.0.1\r\ns=Test\r\nt=0 0\r\n",
		},
		{
			name:  "invalid address type",
			input: "v=0\r\no=- 1001 1 IN INVALID 10.0.0.1\r\ns=Test\r\nt=0 0\r\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var sd SessionDescription
			err := sd.UnmarshalString(test.input)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "invalid value")
		})
	}
}

// Test edge cases for 100% coverage.
func TestUnmarshalOriginEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
	}{
		{
			name:        "missing mandatory username",
			input:       "v=0\r\no=\r\ns=Test\r\nt=0 0\r\n",
			expectError: true,
		},
		{
			name:        "missing mandatory session ID",
			input:       "v=0\r\no=user\r\ns=Test\r\nt=0 0\r\n",
			expectError: true,
		},
		{
			name:        "missing mandatory network type",
			input:       "v=0\r\no=user 1001 1\r\ns=Test\r\nt=0 0\r\n",
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var sd SessionDescription
			err := sd.UnmarshalString(test.input)
			if test.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUnmarshalString_ErrSDPCacheInvalid(t *testing.T) {
	origNew := unmarshalCachePool.New
	t.Cleanup(func() { unmarshalCachePool.New = origNew })

	// ensure there are no cached values.
	unmarshalCachePool.New = nil
	for v := unmarshalCachePool.Get(); v != nil; v = unmarshalCachePool.Get() {
		// discard
	}

	unmarshalCachePool.New = func() any { return 123 }

	var sd SessionDescription
	err := sd.UnmarshalString("")
	assert.ErrorIs(t, err, errSDPCacheInvalid)
}

func TestUnmarshal_DelegatesToUnmarshalString(t *testing.T) {
	in := []byte("v=0\r\no=0 0 0 IN IP4 0\r\ns=0\r\nt=0 0\r\n")
	var sd SessionDescription
	assert.NoError(t, sd.Unmarshal(in))

	out, err := sd.Marshal()
	assert.NoError(t, err)
	assert.Equal(t, string(in), string(out))
}

func TestS1_SyntaxError(t *testing.T) {
	l := &lexer{baseLexer: baseLexer{value: "a="}} // not 'v'
	st, err := s1(l)
	assert.Nil(t, st)

	var se syntaxError
	assert.ErrorAs(t, err, &se)
}

func TestS2_SyntaxError(t *testing.T) {
	l := &lexer{baseLexer: baseLexer{value: "a="}} // not 'o'
	st, err := s2(l)
	assert.Nil(t, st)

	var se syntaxError
	assert.ErrorAs(t, err, &se)
}

func TestS3_SyntaxError(t *testing.T) {
	l := &lexer{baseLexer: baseLexer{value: "a="}} // not 's'
	st, err := s3(l)
	assert.Nil(t, st)

	var se syntaxError
	assert.ErrorAs(t, err, &se)
}

func TestS4_SyntaxError(t *testing.T) {
	l := &lexer{baseLexer: baseLexer{value: "a="}}

	st, err := s4(l)
	assert.Nil(t, st)

	var se syntaxError
	assert.ErrorAs(t, err, &se)
}

func TestS5_SyntaxError(t *testing.T) {
	l := &lexer{baseLexer: baseLexer{value: "a="}}

	st, err := s5(l)
	assert.Nil(t, st)

	var se syntaxError
	assert.ErrorAs(t, err, &se)
}

func TestS6_KeyC_UnmarshalSessionConnectionInformation(t *testing.T) {
	lex := &lexer{
		desc:      &SessionDescription{},
		baseLexer: baseLexer{value: "c=IN IP4 111.1.111.1\r\n"},
	}

	st, err := s6(lex)
	assert.NoError(t, err)

	if assert.NotNil(t, st) {
		_, err = st(lex)
		assert.NoError(t, err)

		if assert.NotNil(t, lex.desc.ConnectionInformation) {
			ci := lex.desc.ConnectionInformation
			assert.Equal(t, "IN", ci.NetworkType)
			assert.Equal(t, "IP4", ci.AddressType)

			if assert.NotNil(t, ci.Address) {
				assert.Equal(t, "111.1.111.1", ci.Address.Address)
			}
		}
	}
}

func TestS6_KeyB_UnmarshalSessionBandwidth(t *testing.T) {
	lex := &lexer{
		desc:      &SessionDescription{},
		baseLexer: baseLexer{value: "b=AS:123\r\n"},
	}

	st, err := s6(lex)
	assert.NoError(t, err)

	if assert.NotNil(t, st) {
		_, err = st(lex)
		assert.NoError(t, err)

		if assert.Len(t, lex.desc.Bandwidth, 1) {
			bw := lex.desc.Bandwidth[0]
			assert.False(t, bw.Experimental)
			assert.Equal(t, "AS", bw.Type)
			assert.Equal(t, uint64(123), bw.Bandwidth)
		}
	}
}

func TestS6_SyntaxError(t *testing.T) {
	l := &lexer{baseLexer: baseLexer{value: "a="}}

	st, err := s6(l)
	assert.Nil(t, st)

	var se syntaxError
	assert.ErrorAs(t, err, &se)
}

func TestS7_KeyE_UnmarshalEmail(t *testing.T) {
	lex := &lexer{
		desc:      &SessionDescription{},
		baseLexer: baseLexer{value: "e=abc.Def@example.com (abc Def)\r\n"},
	}

	st, err := s7(lex)
	assert.NoError(t, err)

	if assert.NotNil(t, st) {
		_, err = st(lex)
		assert.NoError(t, err)

		if assert.NotNil(t, lex.desc.EmailAddress) {
			assert.Equal(t, "abc.Def@example.com (abc Def)", string(*lex.desc.EmailAddress))
		}
	}
}

func TestS7_KeyP_UnmarshalPhone(t *testing.T) {
	lex := &lexer{
		desc:      &SessionDescription{},
		baseLexer: baseLexer{value: "p=+1 111 111-1111\r\n"},
	}

	st, err := s7(lex)
	assert.NoError(t, err)

	if assert.NotNil(t, st) {
		_, err = st(lex)
		assert.NoError(t, err)

		if assert.NotNil(t, lex.desc.PhoneNumber) {
			assert.Equal(t, "+1 111 111-1111", string(*lex.desc.PhoneNumber))
		}
	}
}

func TestS7_KeyC_UnmarshalSessionConnectionInformation(t *testing.T) {
	lex := &lexer{
		desc:      &SessionDescription{},
		baseLexer: baseLexer{value: "c=IN IP4 111.1.111.1\r\n"},
	}

	st, err := s7(lex)
	assert.NoError(t, err)

	if assert.NotNil(t, st) {
		_, err = st(lex)
		assert.NoError(t, err)

		if assert.NotNil(t, lex.desc.ConnectionInformation) {
			ci := lex.desc.ConnectionInformation
			assert.Equal(t, "IN", ci.NetworkType)
			assert.Equal(t, "IP4", ci.AddressType)

			if assert.NotNil(t, ci.Address) {
				assert.Equal(t, "111.1.111.1", ci.Address.Address)
			}
		}
	}
}

func TestS7_KeyB_UnmarshalSessionBandwidth(t *testing.T) {
	lex := &lexer{
		desc:      &SessionDescription{},
		baseLexer: baseLexer{value: "b=AS:123\r\n"},
	}

	st, err := s7(lex)
	assert.NoError(t, err)

	if assert.NotNil(t, st) {
		_, err = st(lex)
		assert.NoError(t, err)

		if assert.Len(t, lex.desc.Bandwidth, 1) {
			bw := lex.desc.Bandwidth[0]
			assert.False(t, bw.Experimental)
			assert.Equal(t, "AS", bw.Type)
			assert.Equal(t, uint64(123), bw.Bandwidth)
		}
	}
}

func TestS7_SyntaxError(t *testing.T) {
	l := &lexer{baseLexer: baseLexer{value: "a="}}

	st, err := s7(l)
	assert.Nil(t, st)

	var se syntaxError
	assert.ErrorAs(t, err, &se)
}

func TestS8_KeyB_UnmarshalSessionBandwidth(t *testing.T) {
	lex := &lexer{
		desc:      &SessionDescription{},
		baseLexer: baseLexer{value: "b=AS:123\r\n"},
	}

	st, err := s8(lex)
	assert.NoError(t, err)

	if assert.NotNil(t, st) {
		_, err = st(lex)
		assert.NoError(t, err)

		if assert.Len(t, lex.desc.Bandwidth, 1) {
			bw := lex.desc.Bandwidth[0]
			assert.False(t, bw.Experimental)
			assert.Equal(t, "AS", bw.Type)
			assert.Equal(t, uint64(123), bw.Bandwidth)
		}
	}
}

func TestS8_SyntaxError(t *testing.T) {
	l := &lexer{baseLexer: baseLexer{value: "a="}}

	st, err := s8(l)
	assert.Nil(t, st)

	var se syntaxError
	assert.ErrorAs(t, err, &se)
}

func TestS9_SyntaxError(t *testing.T) {
	l := &lexer{baseLexer: baseLexer{value: "e="}}

	st, err := s9(l)
	assert.Nil(t, st)

	var se syntaxError
	assert.ErrorAs(t, err, &se)
}

func TestS10_KeyP_UnmarshalPhone(t *testing.T) {
	lex := &lexer{
		desc:      &SessionDescription{},
		baseLexer: baseLexer{value: "p=+1 111 111-1111\r\n"},
	}

	st, err := s10(lex)
	assert.NoError(t, err)

	if assert.NotNil(t, st) {
		_, err = st(lex)
		assert.NoError(t, err)

		if assert.NotNil(t, lex.desc.PhoneNumber) {
			assert.Equal(t, "+1 111 111-1111", string(*lex.desc.PhoneNumber))
		}
	}
}

func TestS10_KeyC_UnmarshalSessionConnectionInformation(t *testing.T) {
	lex := &lexer{
		desc:      &SessionDescription{},
		baseLexer: baseLexer{value: "c=IN IP4 111.1.111.1\r\n"},
	}

	st, err := s10(lex)
	assert.NoError(t, err)

	if assert.NotNil(t, st) {
		_, err = st(lex)
		assert.NoError(t, err)

		if assert.NotNil(t, lex.desc.ConnectionInformation) {
			ci := lex.desc.ConnectionInformation
			assert.Equal(t, "IN", ci.NetworkType)
			assert.Equal(t, "IP4", ci.AddressType)

			if assert.NotNil(t, ci.Address) {
				assert.Equal(t, "111.1.111.1", ci.Address.Address)
			}
		}
	}
}

func TestS10_KeyB_UnmarshalSessionBandwidth(t *testing.T) {
	lex := &lexer{
		desc:      &SessionDescription{},
		baseLexer: baseLexer{value: "b=AS:123\r\n"},
	}

	st, err := s10(lex)
	assert.NoError(t, err)

	if assert.NotNil(t, st) {
		_, err = st(lex)
		assert.NoError(t, err)

		if assert.Len(t, lex.desc.Bandwidth, 1) {
			bw := lex.desc.Bandwidth[0]
			assert.False(t, bw.Experimental)
			assert.Equal(t, "AS", bw.Type)
			assert.Equal(t, uint64(123), bw.Bandwidth)
		}
	}
}

func TestS10_SyntaxError(t *testing.T) {
	l := &lexer{baseLexer: baseLexer{value: "a="}}

	st, err := s10(l)
	assert.Nil(t, st)

	var se syntaxError
	assert.ErrorAs(t, err, &se)
}

func TestS11_SyntaxError(t *testing.T) {
	l := &lexer{baseLexer: baseLexer{value: "t="}}

	st, err := s11(l)
	assert.Nil(t, st)

	var se syntaxError
	assert.ErrorAs(t, err, &se)
}

func TestS12_SyntaxError(t *testing.T) {
	l := &lexer{baseLexer: baseLexer{value: "u="}}

	st, err := s12(l)
	assert.Nil(t, st)

	var se syntaxError
	assert.ErrorAs(t, err, &se)
}

func TestS13_KeyA_UnmarshalSessionAttribute(t *testing.T) {
	lex := &lexer{
		desc:      &SessionDescription{},
		cache:     &unmarshalCache{},
		baseLexer: baseLexer{value: "a=recvonly\r\n"},
	}

	st, err := s13(lex)
	assert.NoError(t, err)

	if assert.NotNil(t, st) {
		_, err = st(lex)
		assert.NoError(t, err)

		attrs := lex.cache.cloneSessionAttributes()

		if assert.Len(t, attrs, 1) {
			assert.Equal(t, "recvonly", attrs[0].Key)
			assert.Equal(t, "", attrs[0].Value)
		}
	}
}

func TestS13_KeyM_UnmarshalMediaDescription(t *testing.T) {
	lex := &lexer{
		desc:      &SessionDescription{},
		cache:     &unmarshalCache{},
		baseLexer: baseLexer{value: "m=audio 49170 RTP/AVP 0\r\n"},
	}

	st, err := s13(lex)
	assert.NoError(t, err)

	if assert.NotNil(t, st) {
		_, err = st(lex)
		assert.NoError(t, err)

		if assert.Len(t, lex.desc.MediaDescriptions, 1) {
			md := lex.desc.MediaDescriptions[0]
			assert.Equal(t, "audio", md.MediaName.Media)
			assert.Equal(t, 49170, md.MediaName.Port.Value)
			assert.Equal(t, []string{"RTP", "AVP"}, md.MediaName.Protos)
			assert.Equal(t, []string{"0"}, md.MediaName.Formats)
		}
	}
}

func TestS13_SyntaxError(t *testing.T) {
	l := &lexer{baseLexer: baseLexer{value: "t="}}

	st, err := s13(l)
	assert.Nil(t, st)

	var se syntaxError
	assert.ErrorAs(t, err, &se)
}

func TestS14_KeyK_UnmarshalMediaEncryptionKey(t *testing.T) {
	lex := &lexer{
		desc: &SessionDescription{
			MediaDescriptions: []*MediaDescription{{}},
		},
		baseLexer: baseLexer{value: "k=prompt\r\n"},
	}

	st, err := s14(lex)
	assert.NoError(t, err)

	if assert.NotNil(t, st) {
		_, err = st(lex)
		assert.NoError(t, err)

		md := lex.desc.MediaDescriptions[len(lex.desc.MediaDescriptions)-1]

		if assert.NotNil(t, md.EncryptionKey) {
			assert.Equal(t, "prompt", string(*md.EncryptionKey))
		}
	}
}

func TestS14_KeyB_UnmarshalMediaBandwidth(t *testing.T) {
	lex := &lexer{
		desc: &SessionDescription{
			MediaDescriptions: []*MediaDescription{{}},
		},
		baseLexer: baseLexer{value: "b=AS:123\r\n"},
	}

	st, err := s14(lex)
	assert.NoError(t, err)

	if assert.NotNil(t, st) {
		_, err = st(lex)
		assert.NoError(t, err)

		md := lex.desc.MediaDescriptions[len(lex.desc.MediaDescriptions)-1]

		if assert.Len(t, md.Bandwidth, 1) {
			bw := md.Bandwidth[0]
			assert.False(t, bw.Experimental)
			assert.Equal(t, "AS", bw.Type)
			assert.Equal(t, uint64(123), bw.Bandwidth)
		}
	}
}

func TestS14_KeyI_UnmarshalMediaTitle(t *testing.T) {
	lex := &lexer{
		desc: &SessionDescription{
			MediaDescriptions: []*MediaDescription{{}},
		},
		baseLexer: baseLexer{value: "i=My Title\r\n"},
	}

	st, err := s14(lex)
	assert.NoError(t, err)

	if assert.NotNil(t, st) {
		_, err = st(lex)
		assert.NoError(t, err)

		md := lex.desc.MediaDescriptions[len(lex.desc.MediaDescriptions)-1]

		if assert.NotNil(t, md.MediaTitle) {
			assert.Equal(t, "My Title", string(*md.MediaTitle))
		}
	}
}

func TestS14_SyntaxError(t *testing.T) {
	l := &lexer{baseLexer: baseLexer{value: "t="}}

	st, err := s14(l)
	assert.Nil(t, st)

	var se syntaxError
	assert.ErrorAs(t, err, &se)
}

func TestS15_KeyA_UnmarshalMediaAttribute(t *testing.T) {
	lex := &lexer{
		desc: &SessionDescription{
			MediaDescriptions: []*MediaDescription{{}}, // need an existing media section
		},
		cache:     &unmarshalCache{},
		baseLexer: baseLexer{value: "a=rtpmap:96 opus/48000\r\n"},
	}

	st, err := s15(lex)
	assert.NoError(t, err)
	if assert.NotNil(t, st) {
		_, err = st(lex) // run unmarshalMediaAttribute
		assert.NoError(t, err)

		attrs := lex.cache.cloneMediaAttributes()
		if assert.Len(t, attrs, 1) {
			assert.Equal(t, "rtpmap", attrs[0].Key)
			assert.Equal(t, "96 opus/48000", attrs[0].Value)
		}
	}
}

func TestS15_KeyC_UnmarshalMediaConnectionInformation(t *testing.T) {
	lex := &lexer{
		desc: &SessionDescription{
			MediaDescriptions: []*MediaDescription{{}},
		},
		baseLexer: baseLexer{value: "c=IN IP4 203.0.113.1\r\n"},
	}

	st, err := s15(lex)
	assert.NoError(t, err)

	if assert.NotNil(t, st) {
		_, err = st(lex)
		assert.NoError(t, err)

		md := lex.desc.MediaDescriptions[len(lex.desc.MediaDescriptions)-1]

		if assert.NotNil(t, md.ConnectionInformation) {
			ci := md.ConnectionInformation
			assert.Equal(t, "IN", ci.NetworkType)
			assert.Equal(t, "IP4", ci.AddressType)

			if assert.NotNil(t, ci.Address) {
				assert.Equal(t, "203.0.113.1", ci.Address.Address)
			}
		}
	}
}

func TestS15_KeyM_UnmarshalMediaDescription(t *testing.T) {
	lex := &lexer{
		desc:      &SessionDescription{},
		cache:     &unmarshalCache{},
		baseLexer: baseLexer{value: "m=audio 49170 RTP/AVP 0\r\n"},
	}

	st, err := s15(lex)
	assert.NoError(t, err)

	if assert.NotNil(t, st) {
		_, err = st(lex)
		assert.NoError(t, err)

		if assert.Len(t, lex.desc.MediaDescriptions, 1) {
			md := lex.desc.MediaDescriptions[0]
			assert.Equal(t, "audio", md.MediaName.Media)
			assert.Equal(t, 49170, md.MediaName.Port.Value)
			assert.Equal(t, []string{"RTP", "AVP"}, md.MediaName.Protos)
			assert.Equal(t, []string{"0"}, md.MediaName.Formats)
		}
	}
}

func TestS15_SyntaxError(t *testing.T) {
	l := &lexer{baseLexer: baseLexer{value: "t="}}

	st, err := s15(l)
	assert.Nil(t, st)

	var se syntaxError
	assert.ErrorAs(t, err, &se)
}

func TestS16_KeyA_UnmarshalMediaAttribute(t *testing.T) {
	lex := &lexer{
		desc:      &SessionDescription{},
		cache:     &unmarshalCache{},
		baseLexer: baseLexer{value: "a=rtpmap:96 opus/48000\r\n"},
	}

	st, err := s16(lex)
	assert.NoError(t, err)

	if assert.NotNil(t, st) {
		_, err = st(lex)
		assert.NoError(t, err)

		attrs := lex.cache.cloneMediaAttributes()

		if assert.Len(t, attrs, 1) {
			assert.Equal(t, "rtpmap", attrs[0].Key)
			assert.Equal(t, "96 opus/48000", attrs[0].Value)
		}
	}
}

func TestS16_KeyK_UnmarshalMediaEncryptionKey(t *testing.T) {
	lex := &lexer{
		desc: &SessionDescription{
			MediaDescriptions: []*MediaDescription{{}},
		},
		baseLexer: baseLexer{value: "k=prompt\r\n"},
	}

	st, err := s16(lex)
	assert.NoError(t, err)

	if assert.NotNil(t, st) {
		_, err = st(lex)
		assert.NoError(t, err)

		md := lex.desc.MediaDescriptions[len(lex.desc.MediaDescriptions)-1]

		if assert.NotNil(t, md.EncryptionKey) {
			assert.Equal(t, "prompt", string(*md.EncryptionKey))
		}
	}
}

func TestS16_KeyB_UnmarshalMediaBandwidth(t *testing.T) {
	lex := &lexer{
		desc: &SessionDescription{
			MediaDescriptions: []*MediaDescription{{}},
		},
		baseLexer: baseLexer{value: "b=AS:123\r\n"},
	}

	st, err := s16(lex)
	assert.NoError(t, err)

	if assert.NotNil(t, st) {
		_, err = st(lex)
		assert.NoError(t, err)

		md := lex.desc.MediaDescriptions[len(lex.desc.MediaDescriptions)-1]

		if assert.Len(t, md.Bandwidth, 1) {
			bw := md.Bandwidth[0]
			assert.False(t, bw.Experimental)
			assert.Equal(t, "AS", bw.Type)
			assert.Equal(t, uint64(123), bw.Bandwidth)
		}
	}
}

func TestS16_KeyI_UnmarshalMediaTitle(t *testing.T) {
	lex := &lexer{
		desc: &SessionDescription{
			MediaDescriptions: []*MediaDescription{{}},
		},
		baseLexer: baseLexer{value: "i=My Title\r\n"},
	}

	st, err := s16(lex)
	assert.NoError(t, err)

	if assert.NotNil(t, st) {
		_, err = st(lex)
		assert.NoError(t, err)

		md := lex.desc.MediaDescriptions[len(lex.desc.MediaDescriptions)-1]

		if assert.NotNil(t, md.MediaTitle) {
			assert.Equal(t, "My Title", string(*md.MediaTitle))
		}
	}
}

func TestS16_KeyM_UnmarshalMediaDescription(t *testing.T) {
	lex := &lexer{
		desc:      &SessionDescription{},
		cache:     &unmarshalCache{},
		baseLexer: baseLexer{value: "m=audio 49170 RTP/AVP 0\r\n"},
	}

	st, err := s16(lex)
	assert.NoError(t, err)

	if assert.NotNil(t, st) {
		_, err = st(lex)
		assert.NoError(t, err)

		if assert.Len(t, lex.desc.MediaDescriptions, 1) {
			md := lex.desc.MediaDescriptions[0]
			assert.Equal(t, "audio", md.MediaName.Media)
			assert.Equal(t, 49170, md.MediaName.Port.Value)
			assert.Equal(t, []string{"RTP", "AVP"}, md.MediaName.Protos)
			assert.Equal(t, []string{"0"}, md.MediaName.Formats)
		}
	}
}

func TestS16_SyntaxError(t *testing.T) {
	l := &lexer{baseLexer: baseLexer{value: "u="}}

	st, err := s16(l)
	assert.Nil(t, st)

	var se syntaxError
	assert.ErrorAs(t, err, &se)
}

func TestUnmarshalProtocolVersion_Error_ReadUint64Field(t *testing.T) {
	// non-numeric
	l := &lexer{baseLexer: baseLexer{value: "x\r\n"}}

	st, err := unmarshalProtocolVersion(l)
	assert.Nil(t, st)

	var se syntaxError
	assert.ErrorAs(t, err, &se)
}

func TestUnmarshalProtocolVersion_Error_InvalidNonZeroVersion(t *testing.T) {
	// version must be 0
	l := &lexer{baseLexer: baseLexer{value: "1\r\n"}}

	st, err := unmarshalProtocolVersion(l)
	assert.Nil(t, st)
	assert.ErrorIs(t, err, errSDPInvalidValue)
}

func TestUnmarshalOrigin_Error_ReadUsernameField(t *testing.T) {
	l := &lexer{
		desc:      &SessionDescription{},
		baseLexer: baseLexer{value: ""},
	}

	st, err := unmarshalOrigin(l)
	assert.Nil(t, st)
	assert.ErrorIs(t, err, io.EOF)
}

func TestUnmarshalOrigin_Error_ReadNetworkTypeField(t *testing.T) {
	// missing NetworkType
	l := &lexer{
		desc:      &SessionDescription{},
		baseLexer: baseLexer{value: "test 1 1"},
	}

	st, err := unmarshalOrigin(l)
	assert.Nil(t, st)
	assert.ErrorIs(t, err, io.EOF)
}

func TestUnmarshalOrigin_Error_ReadUint64_SessionID(t *testing.T) {
	// non-numeric sessionID
	l := &lexer{
		desc:      &SessionDescription{},
		baseLexer: baseLexer{value: "test NaN 1 IN IP4 11.1.1.1\r\n"},
	}

	st, err := unmarshalOrigin(l)
	assert.Nil(t, st)

	var se syntaxError
	assert.ErrorAs(t, err, &se)
}

func TestUnmarshalOrigin_Error_ReadUint64_SessionVersion(t *testing.T) {
	// non-numeric sessionVersion
	l := &lexer{
		desc:      &SessionDescription{},
		baseLexer: baseLexer{value: "test 1 NaN IN IP4 11.1.1.1\r\n"},
	}

	st, err := unmarshalOrigin(l)
	assert.Nil(t, st)

	var se syntaxError
	assert.ErrorAs(t, err, &se)
}

func TestUnmarshalOrigin_Error_InvalidNetworkType(t *testing.T) {
	// invalid network type
	l := &lexer{
		desc:      &SessionDescription{},
		baseLexer: baseLexer{value: "test 1 1 INVALID IP4 11.1.1.1\r\n"},
	}

	st, err := unmarshalOrigin(l)
	assert.Nil(t, st)
	assert.ErrorIs(t, err, errSDPInvalidValue)
}

func TestUnmarshalOrigin_Error_HandleAddressType_Propagates(t *testing.T) {
	// missing AddressType
	l := &lexer{
		desc:      &SessionDescription{},
		baseLexer: baseLexer{value: "test 1 1 IN"},
	}

	st, err := unmarshalOrigin(l)
	assert.Nil(t, st)
	assert.ErrorIs(t, err, io.EOF)
}

func TestUnmarshalOrigin_Error_HandleUnicastAddress_Propagates(t *testing.T) {
	// missing UnicastAddress
	l := &lexer{
		desc:      &SessionDescription{},
		baseLexer: baseLexer{value: "test 1 1 IN IP4"},
	}

	st, err := unmarshalOrigin(l)
	assert.Nil(t, st)
	assert.ErrorIs(t, err, io.EOF)
}

func TestHandleAddressType_ReturnsUnderlyingError(t *testing.T) {
	l := &lexer{
		desc:      &SessionDescription{},
		baseLexer: baseLexer{value: ""},
	}

	err := handleAddressType(l)
	assert.Error(t, err)
	assert.ErrorIs(t, err, io.EOF)
}

func TestHandleUnicastAddress_ReturnsUnderlyingError(t *testing.T) {
	l := &lexer{
		desc:      &SessionDescription{},
		baseLexer: baseLexer{value: ""},
	}

	err := handleUnicastAddress(l)
	assert.Error(t, err)
	assert.ErrorIs(t, err, io.EOF)
}

func TestUnmarshalSessionName_Error_ReadLine(t *testing.T) {
	l := &lexer{desc: &SessionDescription{}, baseLexer: baseLexer{value: ""}}
	st, err := unmarshalSessionName(l)

	assert.Nil(t, st)
	assert.ErrorIs(t, err, io.EOF)
}

func TestUnmarshalSessionInformation_Error_ReadLine(t *testing.T) {
	l := &lexer{desc: &SessionDescription{}, baseLexer: baseLexer{value: ""}}
	st, err := unmarshalSessionInformation(l)

	assert.Nil(t, st)
	assert.ErrorIs(t, err, io.EOF)
}

func TestUnmarshalURI_Error_ReadLine(t *testing.T) {
	l := &lexer{desc: &SessionDescription{}, baseLexer: baseLexer{value: ""}}
	st, err := unmarshalURI(l)

	assert.Nil(t, st)
	assert.ErrorIs(t, err, io.EOF)
}

func TestUnmarshalURI_Error_Parse(t *testing.T) {
	l := &lexer{desc: &SessionDescription{}, baseLexer: baseLexer{value: "%zz\r\n"}}
	st, err := unmarshalURI(l)

	assert.Nil(t, st)
	assert.Error(t, err)
}

func TestUnmarshalEmail_Error_ReadLine(t *testing.T) {
	l := &lexer{desc: &SessionDescription{}, baseLexer: baseLexer{value: ""}}
	st, err := unmarshalEmail(l)

	assert.Nil(t, st)
	assert.ErrorIs(t, err, io.EOF)
}

func TestUnmarshalPhone_Error_ReadLine(t *testing.T) {
	l := &lexer{desc: &SessionDescription{}, baseLexer: baseLexer{value: ""}}
	st, err := unmarshalPhone(l)

	assert.Nil(t, st)
	assert.ErrorIs(t, err, io.EOF)
}

func TestUnmarshalSessionConnectionInformation_Error_FromInner(t *testing.T) {
	l := &lexer{desc: &SessionDescription{}, baseLexer: baseLexer{value: ""}}
	st, err := unmarshalSessionConnectionInformation(l)

	assert.Nil(t, st)
	assert.ErrorIs(t, err, io.EOF)
}

func TestUnmarshalConnectionInformation_ErrInvalidNetworkType(t *testing.T) {
	l := &lexer{baseLexer: baseLexer{value: "INVALID IP4 111.1.111.1\r\n"}}

	ci, err := l.unmarshalConnectionInformation()
	assert.Nil(t, ci)
	assert.ErrorIs(t, err, errSDPInvalidValue)
}

func TestUnmarshalConnectionInformation_ErrReadAddressType(t *testing.T) {
	// missing AddressType token
	l := &lexer{baseLexer: baseLexer{value: "IN"}}

	ci, err := l.unmarshalConnectionInformation()
	assert.Nil(t, ci)
	assert.ErrorIs(t, err, io.EOF)
}

func TestUnmarshalConnectionInformation_ErrInvalidAddressType(t *testing.T) {
	l := &lexer{baseLexer: baseLexer{value: "IN INVALID 111.1.111.1\r\n"}}

	ci, err := l.unmarshalConnectionInformation()
	assert.Nil(t, ci)
	assert.ErrorIs(t, err, errSDPInvalidValue)
}

func TestUnmarshalConnectionInformation_ErrReadAddress(t *testing.T) {
	// missing address token
	l := &lexer{baseLexer: baseLexer{value: "IN IP4"}}

	ci, err := l.unmarshalConnectionInformation()
	assert.Nil(t, ci)
	assert.ErrorIs(t, err, io.EOF)
}

func TestUnmarshalSessionBandwidth_Error_ReadLine(t *testing.T) {
	l := &lexer{
		desc:      &SessionDescription{},
		baseLexer: baseLexer{value: ""},
	}

	st, err := unmarshalSessionBandwidth(l)
	assert.Nil(t, st)
	assert.ErrorIs(t, err, io.EOF)
}

func TestUnmarshalSessionBandwidth_Error_InvalidBandwidthValue(t *testing.T) {
	l := &lexer{
		desc:      &SessionDescription{},
		baseLexer: baseLexer{value: "bad\r\n"},
	}

	st, err := unmarshalSessionBandwidth(l)
	assert.Nil(t, st)
	assert.ErrorIs(t, err, errSDPInvalidValue)
}

func TestUnmarshalBandwidth_InvalidType(t *testing.T) {
	bw, err := unmarshalBandwidth("ZZ:123")
	assert.Nil(t, bw)
	assert.ErrorIs(t, err, errSDPInvalidValue)
}

func TestUnmarshalBandwidth_InvalidNumeric(t *testing.T) {
	bw, err := unmarshalBandwidth("AS:notanumber")
	assert.Nil(t, bw)
	assert.ErrorIs(t, err, errSDPInvalidNumericValue)
}

func TestUnmarshalTiming_Error_StartTime(t *testing.T) {
	// non-numeric start time
	l := &lexer{
		desc:      &SessionDescription{},
		baseLexer: baseLexer{value: "NaN 0\r\n"},
	}

	st, err := unmarshalTiming(l)
	assert.Nil(t, st)

	var se syntaxError
	assert.ErrorAs(t, err, &se)
}

func TestUnmarshalTiming_Error_StopTime(t *testing.T) {
	// non-numeric stop time
	l := &lexer{
		desc:      &SessionDescription{},
		baseLexer: baseLexer{value: "123 NaN\r\n"},
	}

	st, err := unmarshalTiming(l)
	assert.Nil(t, st)

	var se syntaxError
	assert.ErrorAs(t, err, &se)
}

func TestUnmarshalRepeatTimes_Error_FirstFieldRead(t *testing.T) {
	// no tokens
	l := &lexer{
		desc: &SessionDescription{
			TimeDescriptions: []TimeDescription{{}},
		},
		baseLexer: baseLexer{value: ""},
	}

	st, err := unmarshalRepeatTimes(l)
	assert.Nil(t, st)
	assert.ErrorIs(t, err, io.EOF)
}

func TestUnmarshalRepeatTimes_Error_SecondFieldRead(t *testing.T) {
	// missing duration
	l := &lexer{
		desc: &SessionDescription{
			TimeDescriptions: []TimeDescription{{}},
		},
		baseLexer: baseLexer{value: "604800"},
	}

	st, err := unmarshalRepeatTimes(l)
	assert.Nil(t, st)
	assert.ErrorIs(t, err, io.EOF)
}

func TestUnmarshalRepeatTimes_Error_DurationParse(t *testing.T) {
	// invalid duration token
	l := &lexer{
		desc: &SessionDescription{
			TimeDescriptions: []TimeDescription{{}},
		},
		baseLexer: baseLexer{value: "604800 bad\r\n"},
	}

	st, err := unmarshalRepeatTimes(l)
	assert.Nil(t, st)
	assert.ErrorIs(t, err, errSDPInvalidValue)
}

func TestUnmarshalRepeatTimes_Error_OffsetParse(t *testing.T) {
	// invalid offset
	l := &lexer{
		desc: &SessionDescription{
			TimeDescriptions: []TimeDescription{{}},
		},
		baseLexer: baseLexer{value: "604800 3600 nope\r\n"},
	}

	st, err := unmarshalRepeatTimes(l)
	assert.Nil(t, st)
	assert.ErrorIs(t, err, errSDPInvalidValue)
}

func TestUnmarshalRepeatTimes_Error_ReadFieldInsideLoop(t *testing.T) {
	l := &lexer{
		desc: &SessionDescription{
			TimeDescriptions: []TimeDescription{{}},
		},
		baseLexer: baseLexer{value: "604800 3600"},
	}

	st, err := unmarshalRepeatTimes(l)
	assert.Nil(t, st)
	assert.ErrorIs(t, err, io.EOF)
}

func TestUnmarshalTimeZones_Error_ReadUint64Field(t *testing.T) {
	// non-numeric starting token
	l := &lexer{
		desc:      &SessionDescription{},
		baseLexer: baseLexer{value: "NaN"},
	}

	st, err := unmarshalTimeZones(l)
	assert.Nil(t, st)

	var se syntaxError
	assert.ErrorAs(t, err, &se)
}

func TestUnmarshalTimeZones_Error_ReadField(t *testing.T) {
	// no space/offset token
	l := &lexer{
		desc:      &SessionDescription{},
		baseLexer: baseLexer{value: "123"},
	}

	st, err := unmarshalTimeZones(l)
	assert.Nil(t, st)
	assert.ErrorIs(t, err, io.EOF)
}

func TestUnmarshalTimeZones_Error_ParseOffset(t *testing.T) {
	// invalid offset token invalid
	l := &lexer{
		desc:      &SessionDescription{},
		baseLexer: baseLexer{value: "123 bad\r\n"},
	}

	st, err := unmarshalTimeZones(l)
	assert.Nil(t, st)
	assert.ErrorIs(t, err, errSDPInvalidValue)
}

func TestUnmarshalSessionEncryptionKey_Error_ReadLine(t *testing.T) {
	l := &lexer{
		desc:      &SessionDescription{},
		baseLexer: baseLexer{value: ""},
	}

	st, err := unmarshalSessionEncryptionKey(l)
	assert.Nil(t, st)
	assert.ErrorIs(t, err, io.EOF)
}

func TestUnmarshalSessionAttribute_Error_ReadLine(t *testing.T) {
	l := &lexer{
		desc:      &SessionDescription{},
		cache:     &unmarshalCache{},
		baseLexer: baseLexer{value: ""},
	}

	st, err := unmarshalSessionAttribute(l)
	assert.Nil(t, st)
	assert.ErrorIs(t, err, io.EOF)
}

func TestUnmarshalMediaDescription_Error_ReadMediaField(t *testing.T) {
	// no tokens
	l := &lexer{
		desc:      &SessionDescription{},
		baseLexer: baseLexer{value: ""},
	}

	st, err := unmarshalMediaDescription(l)
	assert.Nil(t, st)
	assert.ErrorIs(t, err, io.EOF)
}

func TestUnmarshalMediaDescription_Error_InvalidMediaToken(t *testing.T) {
	// media token is not in allowed set
	l := &lexer{
		desc:      &SessionDescription{},
		baseLexer: baseLexer{value: "data 9 RTP/AVP 0\r\n"},
	}

	st, err := unmarshalMediaDescription(l)
	assert.Nil(t, st)
	assert.ErrorIs(t, err, errSDPInvalidValue)
}

func TestUnmarshalMediaDescription_Error_ReadPortField(t *testing.T) {
	// no port token
	l := &lexer{
		desc:      &SessionDescription{},
		baseLexer: baseLexer{value: "audio"},
	}

	st, err := unmarshalMediaDescription(l)
	assert.Nil(t, st)
	assert.ErrorIs(t, err, io.EOF)
}

func TestUnmarshalMediaDescription_Error_PortRangeInvalid(t *testing.T) {
	// has invalid range part in port token
	l := &lexer{
		desc:      &SessionDescription{},
		baseLexer: baseLexer{value: "audio 123/abc RTP/AVP 0\r\n"},
	}

	st, err := unmarshalMediaDescription(l)
	assert.Nil(t, st)
	assert.ErrorIs(t, err, errSDPInvalidValue)
}

func TestUnmarshalMediaDescription_Error_ReadProtoField(t *testing.T) {
	// but no proto token
	l := &lexer{
		desc:      &SessionDescription{},
		baseLexer: baseLexer{value: "audio 9"},
	}

	st, err := unmarshalMediaDescription(l)
	assert.Nil(t, st)
	assert.ErrorIs(t, err, io.EOF)
}

func TestUnmarshalMediaDescription_Error_ReadFieldInFormatsLoop(t *testing.T) {
	// no newline or fmt tokens
	l := &lexer{
		desc:      &SessionDescription{},
		baseLexer: baseLexer{value: "audio 9 RTP/AVP"},
	}

	st, err := unmarshalMediaDescription(l)
	assert.Nil(t, st)
	assert.ErrorIs(t, err, io.EOF)
}

func TestUnmarshalMediaDescription_SetsPortRange(t *testing.T) {
	// valid port with a range
	lex := &lexer{
		desc:      &SessionDescription{},
		cache:     &unmarshalCache{},
		baseLexer: baseLexer{value: "video 1234/7 RTP/AVP 99\r\n"},
	}

	st, err := unmarshalMediaDescription(lex)
	assert.NoError(t, err)
	assert.NotNil(t, st)

	if assert.Len(t, lex.desc.MediaDescriptions, 1) {
		md := lex.desc.MediaDescriptions[0]

		if assert.NotNil(t, md.MediaName.Port.Range, "Range should be set when <port>/<range> is provided") {
			assert.Equal(t, 7, *md.MediaName.Port.Range)
		}
	}
}

func TestUnmarshalMediaTitle_Error_ReadLine(t *testing.T) {
	// empty input
	l := &lexer{
		desc:      &SessionDescription{},
		baseLexer: baseLexer{value: ""},
	}

	st, err := unmarshalMediaTitle(l)
	assert.Nil(t, st)
	assert.ErrorIs(t, err, io.EOF)
}

func TestUnmarshalMediaConnectionInformation_Error_FromInner(t *testing.T) {
	// empty input
	l := &lexer{
		desc: &SessionDescription{
			MediaDescriptions: []*MediaDescription{{}},
		},
		baseLexer: baseLexer{value: ""},
	}

	st, err := unmarshalMediaConnectionInformation(l)
	assert.Nil(t, st)
	assert.ErrorIs(t, err, io.EOF)
}

func TestUnmarshalMediaBandwidth_Error_ReadLine(t *testing.T) {
	l := &lexer{
		desc: &SessionDescription{
			MediaDescriptions: []*MediaDescription{{}},
		},
		baseLexer: baseLexer{value: ""},
	}

	st, err := unmarshalMediaBandwidth(l)
	assert.Nil(t, st)
	assert.ErrorIs(t, err, io.EOF)
}

func TestUnmarshalMediaBandwidth_Error_InvalidBandwidth(t *testing.T) {
	l := &lexer{
		desc: &SessionDescription{
			MediaDescriptions: []*MediaDescription{{}},
		},
		baseLexer: baseLexer{value: "bad\r\n"},
	}

	st, err := unmarshalMediaBandwidth(l)
	assert.Nil(t, st)
	assert.ErrorIs(t, err, errSDPInvalidSyntax)
}

func TestUnmarshalMediaEncryptionKey_Error_ReadLine(t *testing.T) {
	// empty input
	l := &lexer{
		desc:      &SessionDescription{},
		baseLexer: baseLexer{value: ""},
	}

	st, err := unmarshalMediaEncryptionKey(l)
	assert.Nil(t, st)
	assert.ErrorIs(t, err, io.EOF)
}

func TestUnmarshalMediaAttribute_Error_ReadLine(t *testing.T) {
	// empty input
	l := &lexer{
		desc:      &SessionDescription{},
		baseLexer: baseLexer{value: ""},
	}

	st, err := unmarshalMediaAttribute(l)
	assert.Nil(t, st)
	assert.ErrorIs(t, err, io.EOF)
}

func TestTimeShorthand_MinutesAndSeconds(t *testing.T) {
	t.Run("minutes (m)", func(t *testing.T) {
		assert.Equal(t, int64(60), timeShorthand('m'))
	})

	t.Run("seconds (s)", func(t *testing.T) {
		assert.Equal(t, int64(1), timeShorthand('s'))
	})
}
