package sdp

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	exampleVersion                  = "0"
	exampleVersionLine              = versionKey + exampleVersion + endline
	exampleOrigin                   = "jdoe 2890844526 2890842807 IN IP4 10.47.16.5"
	exampleOriginLine               = originKey + exampleOrigin + endline
	exampleSession                  = "SDP Seminar"
	exampleSessionLine              = sessionKey + exampleSession + endline
	exampleInfo1                    = "A Seminar on the session description protocol"
	exampleInfo1Line                = infoKey + exampleInfo1 + endline
	exampleInfo2                    = "Vivamus a posuere nisl"
	exampleInfo2Line                = infoKey + exampleInfo2 + endline
	exampleURI                      = "http://www.example.com/seminars/sdp.pdf"
	exampleURILine                  = uriKey + exampleURI + endline
	exampleEmail                    = "j.doe@example.com (Jane Doe)"
	exampleEmailLine                = emailKey + exampleEmail + endline
	examplePhone                    = "+1 617 555-6011"
	examplePhoneLine                = phoneKey + examplePhone + endline
	exampleConnection               = "IN IP4 224.2.17.12/127"
	exampleConnectionLine           = connectionKey + exampleConnection + endline
	exampleBandwidth1               = "X-YZ:128"
	exampleBandwidth1Line           = bandwidthKey + exampleBandwidth1 + endline
	exampleBandwidth2               = "AS:12345"
	exampleBandwidth2Line           = bandwidthKey + exampleBandwidth2 + endline
	exampleTiming1                  = "3034423619 3042462419"
	exampleTiming1Line              = timingKey + exampleTiming1 + endline
	exampleTiming2                  = "2873397496 2873404696"
	exampleTiming2Line              = timingKey + exampleTiming2 + endline
	exampleRepeatTime1              = "604800 3600 0 90000"
	exampleRepeatTime1Line          = repeatTimeKey + exampleRepeatTime1 + endline
	exampleRepeatTime2              = "3d 2h 0 21h"
	exampleRepeatTime2Line          = repeatTimeKey + exampleRepeatTime2 + endline
	exampleRepeatTime3              = "259200 7200 0 75600"
	exampleRepeatTime3Line          = repeatTimeKey + exampleRepeatTime3 + endline
	exampleTimeZone1                = "2882844526 -1h"
	exampleTimeZone2                = "2898848070 0"
	exampleTimeZone3                = "2882844526 -3600"
	exampleTimeZones1               = exampleTimeZone1 + " " + exampleTimeZone2
	exampleTimeZones1Line           = timeZonesKey + exampleTimeZones1 + endline
	exampleTimeZones2               = exampleTimeZone3 + " " + exampleTimeZone2
	exampleTimeZones2Line           = timeZonesKey + exampleTimeZones2 + endline
	exampleEncryptionKey            = "prompt"
	exampleEncryptionKeyLine        = encryptionKey + exampleEncryptionKey + endline
	exampleMedia1                   = "video 51372 RTP/AVP 99"
	exampleMedia1Line               = mediaKey + exampleMedia1 + endline
	exampleMedia2                   = "audio 54400 RTP/SAVPF 0 96"
	exampleMedia2Line               = mediaKey + exampleMedia2 + endline
	exampleAttrCategory             = "cat:a.b.c.d"
	exampleAttrCategoryLine         = attributeKey + exampleAttrCategory + endline
	exampleAttrKeywords             = "keywds:a.b.c.d"
	exampleAttrKeywordsLine         = attributeKey + exampleAttrKeywords + endline
	exampleAttrTool                 = "tool:unittest-v1.0"
	exampleAttrToolLine             = attributeKey + exampleAttrTool + endline
	exampleAttrPtime                = "ptime:1538271111099"
	exampleAttrPtimeLine            = attributeKey + exampleAttrPtime + endline
	exampleAttrMaxptime             = "maxptime:1538271111099"
	exampleAttrMaxptimeLine         = attributeKey + exampleAttrMaxptime + endline
	exampleAttrRtpmap1              = "rtpmap:96 opus/48000"
	exampleAttrRtpmap1Line          = attributeKey + exampleAttrRtpmap1 + endline
	exampleAttrRtpmap2              = "rtpmap:99 h263-1998/90000"
	exampleAttrRtpmap2Line          = attributeKey + exampleAttrRtpmap2 + endline
	exampleAttrRtpmap3              = "rtpmap:98 L16/11025/2"
	exampleAttrRtpmap3Line          = attributeKey + exampleAttrRtpmap3 + endline
	exampleAttrRecvonly             = "recvonly"
	exampleAttrRecvonlyLine         = attributeKey + exampleAttrRecvonly + endline
	exampleAttrSendrecv             = "sendrecv"
	exampleAttrSendrecvLine         = attributeKey + exampleAttrSendrecv + endline
	exampleAttrSendonly             = "sendonly"
	exampleAttrSendonlyLine         = attributeKey + exampleAttrSendonly + endline
	exampleAttrInactive             = "inactive"
	exampleAttrInactiveLine         = attributeKey + exampleAttrInactive + endline
	exampleAttrOrient1              = "orient:portrait"
	exampleAttrOrient1Line          = attributeKey + exampleAttrOrient1 + endline
	exampleAttrOrient2              = "orient:landscape"
	exampleAttrOrient2Line          = attributeKey + exampleAttrOrient2 + endline
	exampleAttrOrient3              = "orient:seascape"
	exampleAttrOrient3Line          = attributeKey + exampleAttrOrient3 + endline
	exampleAttrType                 = "type:H332"
	exampleAttrTypeLine             = attributeKey + exampleAttrType + endline
	exampleAttrCharset              = "charset:ISO-8859-1"
	exampleAttrCharsetLine          = attributeKey + exampleAttrCharset + endline
	exampleAttrSdplang              = "sdplang:en-US"
	exampleAttrSdplangLine          = attributeKey + exampleAttrSdplang + endline
	exampleAttrLang                 = "lang:en-US"
	exampleAttrLangLine             = attributeKey + exampleAttrLang + endline
	exampleAttrFramerate            = "framerate:30.4"
	exampleAttrFramerateLine        = attributeKey + exampleAttrFramerate + endline
	exampleAttrQuality              = "quality:7"
	exampleAttrQualityLine          = attributeKey + exampleAttrQuality + endline
	exampleAttrFmtp1                = "fmtp:18 annexb=yes;annexc=no"
	exampleAttrFmtp1Line            = attributeKey + exampleAttrFmtp1 + endline
	exampleAttrFmtp2                = "fmtp:18 annexb=yes"
	exampleAttrFmtp2Line            = attributeKey + exampleAttrFmtp2 + endline
	exampleAttrCandidate1           = "candidate:0 1 UDP 2113667327 203.0.113.1 54400 typ host"
	exampleAttrCandidate1Line       = attributeKey + exampleAttrCandidate1 + endline
	exampleAttrCandidate2           = "candidate:2 1 UDP 1694498815 192.0.2.3 45664 typ srflx raddr 10.0.1.1 rport 8998"
	exampleAttrCandidate2Line       = attributeKey + exampleAttrCandidate2 + endline
	exampleAttrCandidate3           = "candidate:2 1 TCP 2124414975 10.0.1.1 8998 typ host tcptype passive"
	exampleAttrCandidate3Line       = attributeKey + exampleAttrCandidate3 + endline
	exampleAttrRemoteCandidates     = "remote-candidates:1 1.160.10.240 5000 2 3ffe:1900:4545:3:200:f8ff:fe21:67cf 5000"
	exampleAttrRemoteCandidatesLine = attributeKey + exampleAttrRemoteCandidates + endline
	exampleAttrEndOfCandidates      = "end-of-candidates"
	exampleAttrEndOfCandidatesLine  = attributeKey + exampleAttrEndOfCandidates + endline
	exampleAttrIceLite              = "ice-lite"
	exampleAttrIceLiteLine          = attributeKey + exampleAttrIceLite + endline
	exampleAttrIceMismatch          = "ice-mismatch"
	exampleAttrIceMismatchLine      = attributeKey + exampleAttrIceMismatch + endline
	exampleAttrIceUfrag             = "ice-ufrag:8hhY"
	exampleAttrIceUfragLine         = attributeKey + exampleAttrIceUfrag + endline
	exampleAttrIcePwd               = "ice-pwd:asd88fgpdd777uzjYhagZg"
	exampleAttrIcePwdLine           = attributeKey + exampleAttrIcePwd + endline
	exampleAttrIceOptions           = "ice-options:rtp+ecn"
	exampleAttrIceOptionsLine       = attributeKey + exampleAttrIceOptions + endline
	exampleAttrGroup1               = "group:LS 1 2"
	exampleAttrGroup1Line           = attributeKey + exampleAttrGroup1 + endline
	exampleAttrGroup2               = "group:LS"
	exampleAttrGroup2Line           = attributeKey + exampleAttrGroup2 + endline
	exampleAttrMid                  = "mid:1"
	exampleAttrMidLine              = attributeKey + exampleAttrMid + endline
	exampleAttrFingerprint          = "fingerprint:sha-1 4A:AD:B9:B1:3F:82:18:3B:54:02:12:DF:3E:5D:49:6B:19:E5:7C:AB"
	exampleAttrFingerprintLine      = attributeKey + exampleAttrFingerprint + endline
	exampleAttrIdentity1            = "identity:eyJpZHAiOnsiZG9tYWluIjoiZXhhbXBsZS5vcmciLCJwcm90b2NvbCI6ImJvZ3VzIn0sImFzc2VydGlvbiI6IntcImlkZW50aXR5XCI6XCJib2JAZXhhbXBsZS5vcmdcIixcImNvbnRlbnRzXCI6XCJhYmNkZWZnaGlqa2xtbm9wcXJzdHV2d3l6XCIsXCJzaWduYXR1cmVcIjpcIjAxMDIwMzA0MDUwNlwifSJ9 annexb=yes"
	exampleAttrIdentity1Line        = attributeKey + exampleAttrIdentity1 + endline
	exampleAttrIdentity2            = "identity:eyJpZHAiOnsiZG9tYWluIjoiZXhhbXBsZS5vcmciLCJwcm90b2NvbCI6ImJvZ3VzIn0sImFzc2VydGlvbiI6IntcImlkZW50aXR5XCI6XCJib2JAZXhhbXBsZS5vcmdcIixcImNvbnRlbnRzXCI6XCJhYmNkZWZnaGlqa2xtbm9wcXJzdHV2d3l6XCIsXCJzaWduYXR1cmVcIjpcIjAxMDIwMzA0MDUwNlwifSJ9 annexb=yes;annexc=no"
	exampleAttrIdentity2Line        = attributeKey + exampleAttrIdentity2 + endline
	exampleAttrBundleOnly           = "bundle-only"
	exampleAttrBundleOnlyLine       = attributeKey + exampleAttrBundleOnly + endline
	exampleAttrSctpPort             = "sctp-port:5000"
	exampleAttrSctpPortLine         = attributeKey + exampleAttrSctpPort + endline
	exampleAttrMaxMessageSize       = "max-message-size:1024"
	exampleAttrMaxMessageSizeLine   = attributeKey + exampleAttrMaxMessageSize + endline
	exampleAttrSetup                = "setup:actpass"
	exampleAttrSetupLine            = attributeKey + exampleAttrSetup + endline
	exampleAttrTlsId                = "tls-id:abc3de65cddef001be82"
	exampleAttrTlsIdLine            = attributeKey + exampleAttrTlsId + endline
	exampleAttrExtmap1              = "extmap:1 http://example.com/082005/ext.htm#ttime"
	exampleAttrExtmap1Line          = attributeKey + exampleAttrExtmap1 + endline
	exampleAttrExtmap2              = "extmap:2/sendrecv http://example.com/082005/ext.htm#xmeta short"
	exampleAttrExtmap2Line          = attributeKey + exampleAttrExtmap2 + endline
	failingAttrExtmap1              = "extmap:257/sendrecv http://example.com/082005/ext.htm#xmeta short"
	failingAttrExtmap1Line          = attributeKey + failingAttrExtmap1 + endline
	failingAttrExtmap2              = "extmap:2/blorg http://example.com/082005/ext.htm#xmeta short"
	failingAttrExtmap2Line          = attributeKey + failingAttrExtmap2 + endline
	exampleAttrSsrc1                = "ssrc:12345 previous-ssrc:54321 3241"
	exampleAttrSsrc1Line            = attributeKey + exampleAttrSsrc1 + endline
	exampleAttrSsrc2                = "ssrc:12345 cname:another-user@example.com"
	exampleAttrSsrc2Line            = attributeKey + exampleAttrSsrc2 + endline
	exampleAttrSsrc3                = "ssrc:12345 fmtp:18 annexb=yes;annexc=no"
	exampleAttrSsrc3Line            = attributeKey + exampleAttrSsrc3 + endline
	exampleAttrSsrc4                = "ssrc:12345 fmtp:18 annexb=yes"
	exampleAttrSsrc4Line            = attributeKey + exampleAttrSsrc4 + endline
	exampleAttrSsrcGroup1           = "ssrc-group:FEC 1 2 3"
	exampleAttrSsrcGroup1Line       = attributeKey + exampleAttrSsrcGroup1 + endline
	exampleAttrSsrcGroup2           = "ssrc-group:FEC"
	exampleAttrSsrcGroup2Line       = attributeKey + exampleAttrSsrcGroup2 + endline
	exampleAttrRtcp                 = "rtcp:53020 IN IP6 2001:2345:6789:ABCD:EF01:2345:6789:ABCD"
	exampleAttrRtcpLine             = attributeKey + exampleAttrRtcp + endline
	exampleAttrRtcpMuxOnly          = "rtcp-mux-only"
	exampleAttrRtcpMuxOnlyLine      = attributeKey + exampleAttrRtcpMuxOnly + endline
	exampleAttrRtcpFb1              = "rtcp-fb:96 nack"
	exampleAttrRtcpFb1Line          = attributeKey + exampleAttrRtcpFb1 + endline
	exampleAttrRtcpFb2              = "rtcp-fb:* nack"
	exampleAttrRtcpFb2Line          = attributeKey + exampleAttrRtcpFb2 + endline
	exampleAttrRtcpFb3              = "rtcp-fb:98 nack rpsi"
	exampleAttrRtcpFb3Line          = attributeKey + exampleAttrRtcpFb3 + endline
	exampleAttrRtcpMux              = "rtcp-mux"
	exampleAttrRtcpMuxLine          = attributeKey + exampleAttrRtcpMux + endline
	exampleAttrRtcpRsize            = "rtcp-rsize"
	exampleAttrRtcpRsizeLine        = attributeKey + exampleAttrRtcpRsize + endline
	exampleAttrMsid                 = "msid:61317484-2ed4-49d7-9eb7-1414322a7aae f30bdb4a-1497-49b5-3198-e0c9a23172e0"
	exampleAttrMsidLine             = attributeKey + exampleAttrMsid + endline
	exampleAttrImageattr            = "imageattr:97 send [x=720,y=576,sar=[0.91,1.0,1.09,1.45]]"
	exampleAttrImageattrLine        = attributeKey + exampleAttrImageattr + endline
	exampleAttrRid                  = "rid:1 send max-width=1280;max-height=720;max-fps=30"
	exampleAttrRidLine              = attributeKey + exampleAttrRid + endline
	exampleAttrSimulcast            = "simulcast:send 1;2,3 recv 4"
	exampleAttrSimulcastLine        = attributeKey + exampleAttrSimulcast + endline
)

const (
	baseSDP              = exampleVersionLine + exampleOriginLine + exampleSessionLine
	sessionInfoSDP       = baseSDP + exampleInfo1Line + exampleTiming1Line
	uriSDP               = baseSDP + exampleURILine + exampleTiming1Line
	emailSDP             = baseSDP + exampleEmailLine + exampleTiming1Line
	phoneSDP             = baseSDP + examplePhoneLine + exampleTiming1Line
	sessionConnectionSDP = baseSDP + exampleConnectionLine + exampleTiming1Line
	sessionBandwidthSDP  = baseSDP + exampleBandwidth1Line + exampleBandwidth2Line + exampleTiming1Line
	timingSDP            = baseSDP + exampleTiming1Line

	// Short hand time notation is converted into NTP timestamp format in
	// seconds. Because of that unittest comparisons will fail as the same time
	// will be expressed in different units.
	repeatTimesSDP         = timingSDP + exampleRepeatTime1Line + exampleRepeatTime2Line
	repeatTimesSDPExpected = timingSDP + exampleRepeatTime1Line + exampleRepeatTime3Line

	// The expected value looks a bit different for the same reason as mentioned
	// above regarding RepeatTimes.
	timeZonesSDP         = timingSDP + exampleTimeZones1Line
	timeZonesSDPExpected = timingSDP + exampleTimeZones2Line

	sessionEncryptionKeySDP = timingSDP + exampleEncryptionKeyLine
	sessionAttributesSDP    = timingSDP + exampleAttrSdplangLine
	mediaSDP                = timingSDP + exampleMedia1Line + exampleMedia2Line
	mediaInfoSDP            = mediaSDP + exampleInfo2Line
	mediaConnectionSDP      = mediaSDP + exampleConnectionLine
	mediaBandwidthSDP       = mediaSDP + exampleBandwidth1Line + exampleBandwidth2Line
	mediaEncryptionKeySDP   = mediaSDP + exampleEncryptionKeyLine
	mediaAttributesSDP      = mediaSDP + exampleAttrRtpmap2Line + exampleAttrCandidate1Line
	canonicalUnmarshalSDP   = exampleVersionLine +
		exampleOriginLine +
		exampleSessionLine +
		exampleInfo1Line +
		exampleURILine +
		exampleEmailLine +
		examplePhoneLine +
		exampleConnectionLine +
		exampleBandwidth1Line +
		exampleBandwidth2Line +
		exampleTiming1Line +
		exampleTiming2Line +
		exampleRepeatTime1Line +
		exampleTimeZones2Line +
		exampleEncryptionKeyLine +
		exampleAttrRecvonlyLine +
		exampleAttrIceLiteLine +
		exampleAttrIceUfragLine +
		exampleAttrIcePwdLine +
		exampleAttrGroup1Line +
		exampleAttrSetupLine +
		exampleAttrExtmap1Line +
		exampleAttrIdentity2Line +
		exampleAttrFingerprintLine +

		exampleMedia1Line +
		exampleInfo2Line +
		exampleConnectionLine +
		exampleBandwidth1Line +
		exampleEncryptionKeyLine +
		exampleAttrSendrecvLine +
		exampleAttrTlsIdLine +
		exampleAttrIceMismatchLine +
		exampleAttrSsrc1Line +
		exampleAttrSsrc2Line +
		exampleAttrSsrc3Line +
		exampleAttrSsrc4Line +
		exampleAttrRtcpLine +
		exampleAttrRtcpMuxLine +
		exampleAttrRtcpRsizeLine +
		exampleAttrRtcpMuxOnlyLine +
		exampleAttrFingerprintLine +

		exampleMedia2Line +
		exampleAttrSctpPortLine +
		exampleAttrMaxMessageSizeLine +
		exampleAttrRtpmap2Line +
		exampleAttrMidLine +
		exampleAttrSetupLine +
		exampleAttrBundleOnlyLine +
		exampleAttrCandidate1Line +
		exampleAttrEndOfCandidatesLine +
		exampleAttrRemoteCandidatesLine +
		exampleAttrSsrcGroup1Line +
		exampleAttrSsrcGroup2Line +
		exampleAttrExtmap2Line

	canonicalMarshalSDP = "v=0\r\n" +
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
		"a=recvonly\r\n" +
		"m=audio 49170 RTP/AVP 0\r\n" +
		"i=Vivamus a posuere nisl\r\n" +
		"c=IN IP4 203.0.113.1\r\n" +
		"b=X-YZ:128\r\n" +
		"k=prompt\r\n" +
		"a=sendrecv\r\n" +
		"m=video 51372 RTP/AVP 99\r\n" +
		"a=rtpmap:99 h263-1998/90000\r\n" +
		"a=candidate:0 1 UDP 2113667327 203.0.113.1 54400 typ host\r\n"
)

func TestSessionDescription_Unmarshal(t *testing.T) {
	tests := []struct {
		sdp      string
		expected string
	}{
		{sessionInfoSDP, sessionInfoSDP},
		{uriSDP, uriSDP},
		{emailSDP, emailSDP},
		{phoneSDP, phoneSDP},
		{sessionConnectionSDP, sessionConnectionSDP},
		{sessionBandwidthSDP, sessionBandwidthSDP},
		{repeatTimesSDP, repeatTimesSDPExpected},
		{timeZonesSDP, timeZonesSDPExpected},
		{sessionEncryptionKeySDP, sessionEncryptionKeySDP},
		{sessionAttributesSDP, sessionAttributesSDP},
		{mediaSDP, mediaSDP},
		{mediaInfoSDP, mediaInfoSDP},
		{mediaConnectionSDP, mediaConnectionSDP},
		{mediaBandwidthSDP, mediaBandwidthSDP},
		{mediaEncryptionKeySDP, mediaEncryptionKeySDP},
		{mediaAttributesSDP, mediaAttributesSDP},
		{canonicalUnmarshalSDP, canonicalUnmarshalSDP},
	}

	for i, u := range tests {
		sd := &SessionDescription{}
		err := sd.Unmarshal(u.sdp)
		assert.Nil(t, err, "%d: %+v", i, err)
		assert.Equal(t, u.expected, sd.Marshal(), "%d", i)
	}
}

func TestSessionDescription_Marshal(t *testing.T) {
	sd := &SessionDescription{
		Version: Version{Value: 0},
		Origin: Origin{
			Username:       "jdoe",
			SessionID:      uint64(2890844526),
			SessionVersion: uint64(2890842807),
			NetworkType:    "IN",
			AddressType:    "IP4",
			UnicastAddress: "10.47.16.5",
		},
		Session:     Session{Value: "SDP Seminar"},
		Information: &Information{Value: "A Seminar on the session description protocol"},
		URI: func() *URL {
			uri := &URL{}
			if err := uri.Unmarshal("http://www.example.com/seminars/sdp.pdf"); err != nil {
				return nil
			}
			return uri
		}(),
		EmailAddress: &EmailAddress{Value: "j.doe@example.com (Jane Doe)"},
		PhoneNumber:  &PhoneNumber{Value: "+1 617 555-6011"},
		Connection: &Connection{
			NetworkType: "IN",
			AddressType: "IP4",
			Address: &Address{
				IP:  net.ParseIP("224.2.17.12"),
				TTL: &(&struct{ x int }{127}).x,
			},
		},
		Bandwidths: Bandwidths{
			{
				Experimental: true,
				Type:         "YZ",
				Bandwidth:    128,
			},
			{
				Type:      "AS",
				Bandwidth: 12345,
			},
		},
		TimeDescriptions: TimeDescriptions{
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
		EncryptionKey: &EncryptionKey{Value: "prompt"},
		Attributes: []Attribute{
			&RecvOnly{},
		},
		MediaDescriptions: []MediaDescription{
			{
				Media: Media{
					Type: MediaTypeAudio,
					Port: RangedPort{
						Value: 49170,
					},
					Protos:  []string{"RTP", "AVP"},
					Formats: []string{"0"},
				},
				Information: &Information{Value: "Vivamus a posuere nisl"},
				Connection: &Connection{
					NetworkType: "IN",
					AddressType: "IP4",
					Address: &Address{
						IP: net.ParseIP("203.0.113.1"),
					},
				},
				Bandwidths: Bandwidths{
					{
						Experimental: true,
						Type:         "YZ",
						Bandwidth:    128,
					},
				},
				EncryptionKey: &EncryptionKey{Value: "prompt"},
				Attributes: []Attribute{
					&SendRecv{},
				},
			},
			{
				Media: Media{
					Type: MediaTypeVideo,
					Port: RangedPort{
						Value: 51372,
					},
					Protos:  []string{"RTP", "AVP"},
					Formats: []string{"99"},
				},
				Attributes: []Attribute{
					&RtpMap{
						Payload:   99,
						Encoding:  "h263-1998",
						ClockRate: 90000,
					},
					&Candidate{
						Foundation: "0",
						Component:  1,
						Protocol:   ProtocolUDP,
						Priority:   2113667327,
						Addr:       "203.0.113.1",
						Port:       54400,
						Type:       CandidateTypeHost,
					},
				},
			},
		},
	}

	actual := sd.Marshal()
	assert.Equal(t, canonicalMarshalSDP, actual)
}
