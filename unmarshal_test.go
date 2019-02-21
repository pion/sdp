package sdp

import (
	"testing"
)

const (
	BaseSDP = "v=0\r\n" +
		"o=jdoe 2890844526 2890842807 IN IP4 10.47.16.5\r\n" +
		"s=SDP Seminar\r\n"

	SessionInformationSDP = BaseSDP +
		"i=A Seminar on the session description protocol\r\n" +
		"t=3034423619 3042462419\r\n"

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

	// The expected value looks a bit different for the same reason as mentioned
	// above regarding RepeatTimes.
	TimeZonesSDP = TimingSDP +
		"r=2882844526 -1h 2898848070 0\r\n"

	TimeZonesSDPExpected = TimingSDP +
		"r=2882844526 -3600 2898848070 0\r\n"

	SessionEncryptionKeySDP = TimingSDP +
		"k=prompt\r\n"

	SessionAttributesSDP = TimingSDP +
		"a=rtpmap:96 opus/48000\r\n"

	MediaNameSDP = TimingSDP +
		"m=video 51372 RTP/AVP 99\r\n" +
		"m=audio 54400 RTP/SAVPF 0 96\r\n"

	MediaTitleSDP = MediaNameSDP +
		"i=Vivamus a posuere nisl\r\n"

	MediaConnectionInformationSDP = MediaNameSDP +
		"c=IN IP4 203.0.113.1\r\n"

	MediaBandwidthSDP = MediaNameSDP +
		"b=X-YZ:128\r\n" +
		"b=AS:12345\r\n"

	MediaEncryptionKeySDP = MediaNameSDP +
		"k=prompt\r\n"

	MediaAttributesSDP = MediaNameSDP +
		"a=rtpmap:99 h263-1998/90000\r\n" +
		"a=candidate:0 1 UDP 2113667327 203.0.113.1 54400 typ host\r\n"

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
		Name string
		SDP  string
	}{
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
			SDP:  string(EmailAddressSDP),
		},
		{
			Name: "PhoneNumber",
			SDP:  PhoneNumberSDP,
		},
		{
			Name: "SessionConnectionInformation",
			SDP:  SessionConnectionInformationSDP,
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
			Name: "SessionAttributes",
			SDP:  SessionAttributesSDP,
		},
		{
			Name: "MediaName",
			SDP:  MediaNameSDP,
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
			Name: "MediaConnectionInformation",
			SDP:  MediaConnectionInformationSDP,
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
			Name: "MediaAttributes",
			SDP:  MediaAttributesSDP,
		},
		{
			Name: "CanonicalUnmarshal",
			SDP:  CanonicalUnmarshalSDP,
		},
	} {
		sd := &SessionDescription{}

		err := sd.Unmarshal([]byte(test.SDP))
		if got, want := err, error(nil); got != want {
			t.Fatalf("Unmarshal(%s): err=%v, want %v", test.Name, got, want)
		}

		actual, err := sd.Marshal()
		if got, want := err, error(nil); got != want {
			t.Fatalf("Marshal(): err=%v, want %v", got, want)
		}
		if got, want := string(actual), test.SDP; got != want {
			t.Fatalf("Marshal(%s) = %q, want %q", test.Name, got, want)
		}
	}
}

func TestUnmarshalRepeatTimes(t *testing.T) {
	sd := &SessionDescription{}
	if err := sd.Unmarshal([]byte(RepeatTimesSDP)); err != nil {
		t.Errorf("error: %v", err)
	}

	actual, err := sd.Marshal()
	if got, want := err, error(nil); got != want {
		t.Fatalf("Marshal(): err=%v, want %v", got, want)
	}
	if string(actual) != RepeatTimesSDPExpected {
		t.Errorf("error:\n\nEXPECTED:\n%v\nACTUAL:\n%v", RepeatTimesSDPExpected, actual)
	}
}

func TestUnmarshalTimeZones(t *testing.T) {
	sd := &SessionDescription{}
	if err := sd.Unmarshal([]byte(TimeZonesSDP)); err != nil {
		t.Errorf("error: %v", err)
	}

	actual, err := sd.Marshal()
	if got, want := err, error(nil); got != want {
		t.Fatalf("Marshal(): err=%v, want %v", got, want)
	}
	if string(actual) != TimeZonesSDPExpected {
		t.Errorf("error:\n\nEXPECTED:\n%v\nACTUAL:\n%v", TimeZonesSDPExpected, actual)
	}
}
