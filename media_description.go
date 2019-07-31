package sdp

import (
	"errors"
	"strconv"
	"strings"
)

// MediaDescription represents a media type.
// https://tools.ietf.org/html/rfc4566#section-5.14
type MediaDescription struct {
	// m=<media> <port>/<number of ports> <proto> <fmt> ...
	// https://tools.ietf.org/html/rfc4566#section-5.14
	MediaName MediaName

	// i=<session description>
	// https://tools.ietf.org/html/rfc4566#section-5.4
	MediaTitle *Information

	// c=<nettype> <addrtype> <connection-address>
	// https://tools.ietf.org/html/rfc4566#section-5.7
	ConnectionInformation *ConnectionInformation

	// b=<bwtype>:<bandwidth>
	// https://tools.ietf.org/html/rfc4566#section-5.8
	Bandwidth []Bandwidth

	// k=<method>
	// k=<method>:<encryption key>
	// https://tools.ietf.org/html/rfc4566#section-5.12
	EncryptionKey *EncryptionKey

	// a=<attribute>
	// a=<attribute>:<value>
	// https://tools.ietf.org/html/rfc4566#section-5.13
	Attributes []Attribute
}

// Attribute returns the value of an attribute and if it exists
func (s *MediaDescription) Attribute(key string) (string, bool) {
	for _, a := range s.Attributes {
		if a.Key == key {
			return a.Value, true
		}
	}
	return "", false
}

// RangedPort supports special format for the media field "m=" port value. If
// it may be necessary to specify multiple transport ports, the protocol allows
// to write it as: <port>/<number of ports> where number of ports is a an
// offsetting range.
type RangedPort struct {
	Value int
	Range *int
}

func (p *RangedPort) String() string {
	output := strconv.Itoa(p.Value)
	if p.Range != nil {
		output += "/" + strconv.Itoa(*p.Range)
	}
	return output
}

// MediaName describes the "m=" field storage structure.
type MediaName struct {
	Media   string
	Port    RangedPort
	Protos  []string
	Formats []string
}

func (m *MediaName) String() *string {
	output := strings.Join([]string{
		m.Media,
		m.Port.String(),
		strings.Join(m.Protos, "/"),
		strings.Join(m.Formats, " "),
	}, " ")
	return &output
}

type RTCPFeedback struct {
	Type string
	Parameter string
}
// MediaFormat contains information corresponding to one m= line in SDP.
// Most of these are source attributes as defined in RFC 5576
type MediaFormat struct {
	MediaType   string
	PayloadType int
	// a=rtpmap <payload type> <EncodingName>
	// such as a=rtpmap 96 H264/90000
	EncodingName string
	// a=fmtp:<payload type> <format parameters string>
	// a=fmtp:96 level-asymmetry-allowed=1;packetization-mode=1;profile-level-id=640c1f
	Parameters string
	// Array of rtcpfeeback attributes
	// a=rtcp-fb:<payload type> <atttribute>
	RTCPFeedback []RTCPFeedback
}

// MediFormats returns an array of MediaFormat structs, one
// for each m= line of the MediaDescription
func (md *MediaDescription) MediaFormats() ([]*MediaFormat, error) {
	ret := make([]*MediaFormat, 0)              // the returned array
	var formats = make(map[string]*MediaFormat) // k: payload type as string, v: *MediaFormats under construction
	for _, fmt := range md.MediaName.Formats {
		payloadType, err := strconv.Atoi(fmt)
		if err != nil {
			return nil, errors.New("format parse error")
		}
		mf := &MediaFormat{MediaType: md.MediaName.Media, PayloadType: payloadType}
		formats[fmt] = mf
		ret = append(ret, mf)
	}
	for _, a := range md.Attributes {
		switch a.Key {
		case "rtpmap":
			splits := strings.Split(a.Value, " ")
			if len(splits) != 2 {
				return nil, errors.New("error parsing rtpmap line")
			}
			mf, ok := formats[splits[0]]
			if !ok {
				return nil, errors.New("unexpected payload type in rtpmap")
			}
			mf.EncodingName = splits[1]
		case "fmtp":
			splits := strings.Split(a.Value, " ")
			if len(splits) != 2 {
				return nil, errors.New("error parsing fmtp line")
			}
			mf, ok := formats[splits[0]]
			if !ok {
				return nil, errors.New("unexpected payload type in fmtp")
			}
			mf.Parameters = splits[1]
		case "rtcp-fb":
			splits := strings.Split(a.Value, " ")
			if len(splits) < 2 {
				return nil, errors.New("error parsing fmtp line")
			}
			mf, ok := formats[splits[0]]
			if !ok {
				return nil, errors.New("unexpected payload type in rtc-fp")
			}
			fb := RTCPFeedback{Type:splits[1]}
			if len(splits) > 2 {
				fb.Parameter = strings.Join(splits[2:], " ")
			}
			mf.RTCPFeedback = append(mf.RTCPFeedback, fb)
		}
	}
	return ret, nil
}

// SameFormat compares one MediaFormat with another and returns
// true if they match (excluding rtcp-fb lines)
func (f *MediaFormat) SameFormat(format *MediaFormat) bool {
	if format.PayloadType != f.PayloadType {
		return false
	}
	if format.MediaType != f.MediaType {
		return false
	}
	if format.Parameters != f.Parameters {
		return false
	}
	if format.EncodingName != f.EncodingName {
		return false
	}
	return true
}
