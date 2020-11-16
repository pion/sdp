package sdp

import (
	"errors"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"strings"
)

var (
	errSDPInvalidSyntax       = errors.New("sdp: invalid syntax")
	errSDPInvalidNumericValue = errors.New("sdp: invalid numeric value")
	errSDPInvalidValue        = errors.New("sdp: invalid value")
	errSDPInvalidPortValue    = errors.New("sdp: invalid port value")
)

// Unmarshal is the primary function that deserializes the session description
// message and stores it inside of a structured SessionDescription object.
//
// The States Transition Table describes the computation flow between functions
// (namely s1, s2, s3, ...) for a parsing procedure that complies with the
// specifications laid out by the rfc4566#section-5 as well as by JavaScript
// Session Establishment Protocol draft. Links:
// 		https://tools.ietf.org/html/rfc4566#section-5
// 		https://tools.ietf.org/html/draft-ietf-rtcweb-jsep-24
//
// https://tools.ietf.org/html/rfc4566#section-5
// Session description
//    v=  (protocol version)
//    o=  (originator and session identifier)
//    s=  (session name)
//    i=* (session information)
//    u=* (URI of description)
//    e=* (email address)
//    p=* (phone number)
//    c=* (connection information -- not required if included in
//         all media)
//    b=* (zero or more bandwidth information lines)
//    One or more time descriptions ("t=" and "r=" lines; see below)
//    z=* (time zone adjustments)
//    k=* (encryption key)
//    a=* (zero or more session attribute lines)
//    Zero or more media descriptions
//
// Time description
//    t=  (time the session is active)
//    r=* (zero or more repeat times)
//
// Media description, if present
//    m=  (media name and transport address)
//    i=* (media title)
//    c=* (connection information -- optional if included at
//         session level)
//    b=* (zero or more bandwidth information lines)
//    k=* (encryption key)
//    a=* (zero or more media attribute lines)
//
// In order to generate the following state table and draw subsequent
// deterministic finite-state automota ("DFA") the following regex was used to
// derive the DFA:
//    vosi?u?e?p?c?b*(tr*)+z?k?a*(mi?c?b*k?a*)*
// possible place and state to exit:
//                    **   * * *  ** * * * *
//                    99   1 1 1  11 1 1 1 1
//                         3 1 1  26 5 5 4 4
//
// Please pay close attention to the `k`, and `a` parsing states. In the table
// below in order to distinguish between the states belonging to the media
// description as opposed to the session description, the states are marked
// with an asterisk ("a*", "k*").
// +--------+----+-------+----+-----+----+-----+---+----+----+---+---+-----+---+---+----+---+----+
// | STATES | a* | a*,k* | a  | a,k | b  | b,c | e | i  | m  | o | p | r,t | s | t | u  | v | z  |
// +--------+----+-------+----+-----+----+-----+---+----+----+---+---+-----+---+---+----+---+----+
// |   s1   |    |       |    |     |    |     |   |    |    |   |   |     |   |   |    | 2 |    |
// |   s2   |    |       |    |     |    |     |   |    |    | 3 |   |     |   |   |    |   |    |
// |   s3   |    |       |    |     |    |     |   |    |    |   |   |     | 4 |   |    |   |    |
// |   s4   |    |       |    |     |    |   5 | 6 |  7 |    |   | 8 |     |   | 9 | 10 |   |    |
// |   s5   |    |       |    |     |  5 |     |   |    |    |   |   |     |   | 9 |    |   |    |
// |   s6   |    |       |    |     |    |   5 |   |    |    |   | 8 |     |   | 9 |    |   |    |
// |   s7   |    |       |    |     |    |   5 | 6 |    |    |   | 8 |     |   | 9 | 10 |   |    |
// |   s8   |    |       |    |     |    |   5 |   |    |    |   |   |     |   | 9 |    |   |    |
// |   s9   |    |       |    |  11 |    |     |   |    | 12 |   |   |   9 |   |   |    |   | 13 |
// |   s10  |    |       |    |     |    |   5 | 6 |    |    |   | 8 |     |   | 9 |    |   |    |
// |   s11  |    |       | 11 |     |    |     |   |    | 12 |   |   |     |   |   |    |   |    |
// |   s12  |    |    14 |    |     |    |  15 |   | 16 | 12 |   |   |     |   |   |    |   |    |
// |   s13  |    |       |    |  11 |    |     |   |    | 12 |   |   |     |   |   |    |   |    |
// |   s14  | 14 |       |    |     |    |     |   |    | 12 |   |   |     |   |   |    |   |    |
// |   s15  |    |    14 |    |     | 15 |     |   |    | 12 |   |   |     |   |   |    |   |    |
// |   s16  |    |    14 |    |     |    |  15 |   |    | 12 |   |   |     |   |   |    |   |    |
// +--------+----+-------+----+-----+----+-----+---+----+----+---+---+-----+---+---+----+---+----+
func (s *SessionDescription) Unmarshal(value []byte) (err error) {
	l := new(lexer)
	l.desc = s
	l.data = value
	for state := sFn[1]; state != nil; {
		state, err = state(l)
		if err != nil {
			return
		}
	}
	return nil
}

var sFn [16 + 1]func(l *lexer) (stateFn, error)

func init() {
	for i, m := range []map[string]stateFn{
		0: {
			// skip zero index
		},
		1: {
			"v=": unmarshalProtocolVersion,
		},
		2: {
			"o=": unmarshalOrigin,
		},
		3: {
			"s=": unmarshalSessionName,
		},
		4: {
			"b=": unmarshalSessionBandwidth,
			"c=": unmarshalSessionConnectionInformation,
			"e=": unmarshalEmail,
			"i=": unmarshalSessionInformation,
			"p=": unmarshalPhone,
			"t=": unmarshalTiming,
			"u=": unmarshalURI,
		},
		5: {
			"b=": unmarshalSessionBandwidth,
			"t=": unmarshalTiming,
		},
		6: {
			"b=": unmarshalSessionBandwidth,
			"c=": unmarshalSessionConnectionInformation,
			"p=": unmarshalPhone,
			"t=": unmarshalTiming,
		},
		7: {
			"b=": unmarshalSessionBandwidth,
			"c=": unmarshalSessionConnectionInformation,
			"e=": unmarshalEmail,
			"p=": unmarshalPhone,
			"t=": unmarshalTiming,
			"u=": unmarshalURI,
		},
		8: {
			"b=": unmarshalSessionBandwidth,
			"c=": unmarshalSessionConnectionInformation,
			"t=": unmarshalTiming,
		},
		9: {
			"a=": unmarshalSessionAttribute,
			"k=": unmarshalSessionEncryptionKey,
			"m=": unmarshalMediaDescription,
			"r=": unmarshalRepeatTimes,
			"t=": unmarshalTiming,
			"z=": unmarshalTimeZones,
		},
		10: {
			"b=": unmarshalSessionBandwidth,
			"c=": unmarshalSessionConnectionInformation,
			"e=": unmarshalEmail,
			"p=": unmarshalPhone,
			"t=": unmarshalTiming,
		},
		11: {
			"a=": unmarshalSessionAttribute,
			"m=": unmarshalMediaDescription,
		},
		12: {
			"a=": unmarshalMediaAttribute,
			"b=": unmarshalMediaBandwidth,
			"c=": unmarshalMediaConnectionInformation,
			"i=": unmarshalMediaTitle,
			"k=": unmarshalMediaEncryptionKey,
			"m=": unmarshalMediaDescription,
		},
		13: {
			"a=": unmarshalSessionAttribute,
			"k=": unmarshalSessionEncryptionKey,
			"m=": unmarshalMediaDescription,
		},
		14: {
			"a=": unmarshalMediaAttribute,
			"b=": unmarshalMediaBandwidth,             // Non-spec ordering
			"c=": unmarshalMediaConnectionInformation, // Non-spec ordering
			"i=": unmarshalMediaTitle,                 // Non-spec ordering
			"k=": unmarshalMediaEncryptionKey,         // Non-spec ordering
			"m=": unmarshalMediaDescription,
		},
		15: {
			"a=": unmarshalMediaAttribute,
			"b=": unmarshalMediaBandwidth,
			"c=": unmarshalMediaConnectionInformation,
			"i=": unmarshalMediaTitle, // Non-spec ordering
			"k=": unmarshalMediaEncryptionKey,
			"m=": unmarshalMediaDescription,
		},
		16: {
			"a=": unmarshalMediaAttribute,
			"b=": unmarshalMediaBandwidth,
			"c=": unmarshalMediaConnectionInformation,
			"i=": unmarshalMediaTitle, // Non-spec ordering
			"k=": unmarshalMediaEncryptionKey,
			"m=": unmarshalMediaDescription,
		},
	} {
		sFn[i] = newStateFn(m)
	}
}

func newStateFn(m map[string]stateFn) func(l *lexer) (stateFn, error) {
	return func(l *lexer) (stateFn, error) {
		typ, err := l.readType()
		if err == io.EOF && typ == "" {
			return nil, nil
		} else if err != nil {
			return nil, err
		}

		if fn, ok := m[typ]; ok {
			return fn, nil
		}

		return nil, fmt.Errorf("%w `%v`", errSDPInvalidSyntax, typ)
	}
}

func unmarshalProtocolVersion(l *lexer) (stateFn, error) {
	version, err := l.readUint64Field()
	if err != nil {
		return nil, err
	}

	// As off the latest draft of the rfc this value is required to be 0.
	// https://tools.ietf.org/html/draft-ietf-rtcweb-jsep-24#section-5.8.1
	if version != 0 {
		return nil, fmt.Errorf("%w `%v`", errSDPInvalidValue, version)
	}

	if err := l.nextLine(); err != nil {
		return nil, err
	}

	return sFn[2], nil
}

func unmarshalOrigin(l *lexer) (stateFn, error) {
	var err error

	l.desc.Origin.Username, err = l.readField()
	if err != nil {
		return nil, err
	}

	l.desc.Origin.SessionID, err = l.readUint64Field()
	if err != nil {
		return nil, err
	}

	l.desc.Origin.SessionVersion, err = l.readUint64Field()
	if err != nil {
		return nil, err
	}

	l.desc.Origin.NetworkType, err = l.readField()
	if err != nil {
		return nil, err
	}

	// Set according to currently registered with IANA
	// https://tools.ietf.org/html/rfc4566#section-8.2.6
	if !anyOf(l.desc.Origin.NetworkType, "IN") {
		return nil, fmt.Errorf("%w `%v`", errSDPInvalidValue, l.desc.Origin.NetworkType)
	}

	l.desc.Origin.AddressType, err = l.readField()
	if err != nil {
		return nil, err
	}

	// Set according to currently registered with IANA
	// https://tools.ietf.org/html/rfc4566#section-8.2.7
	if !anyOf(l.desc.Origin.AddressType, "IP4", "IP6") {
		return nil, fmt.Errorf("%w `%v`", errSDPInvalidValue, l.desc.Origin.AddressType)
	}

	l.desc.Origin.UnicastAddress, err = l.readField()
	if err != nil {
		return nil, err
	}

	if err := l.nextLine(); err != nil {
		return nil, err
	}

	return sFn[3], nil
}

func unmarshalSessionName(l *lexer) (stateFn, error) {
	value, err := l.readLine()
	if err != nil {
		return nil, err
	}

	l.desc.SessionName = SessionName(value)
	return sFn[4], nil
}

func unmarshalSessionInformation(l *lexer) (stateFn, error) {
	value, err := l.readLine()
	if err != nil {
		return nil, err
	}

	sessionInformation := Information(value)
	l.desc.SessionInformation = &sessionInformation
	return sFn[7], nil
}

func unmarshalURI(l *lexer) (stateFn, error) {
	value, err := l.readLine()
	if err != nil {
		return nil, err
	}

	l.desc.URI, err = url.Parse(value)
	if err != nil {
		return nil, err
	}

	return sFn[10], nil
}

func unmarshalEmail(l *lexer) (stateFn, error) {
	value, err := l.readLine()
	if err != nil {
		return nil, err
	}

	emailAddress := EmailAddress(value)
	l.desc.EmailAddress = &emailAddress
	return sFn[6], nil
}

func unmarshalPhone(l *lexer) (stateFn, error) {
	value, err := l.readLine()
	if err != nil {
		return nil, err
	}

	phoneNumber := PhoneNumber(value)
	l.desc.PhoneNumber = &phoneNumber
	return sFn[8], nil
}

func unmarshalSessionConnectionInformation(l *lexer) (stateFn, error) {
	var err error
	l.desc.ConnectionInformation, err = l.unmarshalConnectionInformation()
	if err != nil {
		return nil, err
	}
	return sFn[5], nil
}

func (l *lexer) unmarshalConnectionInformation() (*ConnectionInformation, error) {
	var err error
	var c ConnectionInformation

	c.NetworkType, err = l.readField()
	if err != nil {
		return nil, err
	}

	// Set according to currently registered with IANA
	// https://tools.ietf.org/html/rfc4566#section-8.2.6
	if !anyOf(c.NetworkType, "IN") {
		return nil, fmt.Errorf("%w `%v`", errSDPInvalidValue, c.NetworkType)
	}

	c.AddressType, err = l.readField()
	if err != nil {
		return nil, err
	}

	// Set according to currently registered with IANA
	// https://tools.ietf.org/html/rfc4566#section-8.2.7
	if !anyOf(c.AddressType, "IP4", "IP6") {
		return nil, fmt.Errorf("%w `%v`", errSDPInvalidValue, c.AddressType)
	}

	address, err := l.readField()
	if err != nil {
		return nil, err
	}

	if address != "" {
		c.Address = new(Address)
		c.Address.Address = address
	}

	if err := l.nextLine(); err != nil {
		return nil, err
	}

	return &c, nil
}

func unmarshalSessionBandwidth(l *lexer) (stateFn, error) {
	value, err := l.readLine()
	if err != nil {
		return nil, err
	}

	bandwidth, err := unmarshalBandwidth(value)
	if err != nil {
		return nil, fmt.Errorf("%w `b=%v`", errSDPInvalidValue, value)
	}
	l.desc.Bandwidth = append(l.desc.Bandwidth, *bandwidth)

	return sFn[5], nil
}

func unmarshalBandwidth(value string) (*Bandwidth, error) {
	parts := strings.Split(value, ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("%w `b=%v`", errSDPInvalidValue, parts)
	}

	experimental := strings.HasPrefix(parts[0], "X-")
	if experimental {
		parts[0] = strings.TrimPrefix(parts[0], "X-")
	} else if !anyOf(parts[0], "CT", "AS") {
		// Set according to currently registered with IANA
		// https://tools.ietf.org/html/rfc4566#section-5.8
		return nil, fmt.Errorf("%w `%v`", errSDPInvalidValue, parts[0])
	}

	bandwidth, err := strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("%w `%v`", errSDPInvalidNumericValue, parts[1])
	}

	return &Bandwidth{
		Experimental: experimental,
		Type:         parts[0],
		Bandwidth:    bandwidth,
	}, nil
}

func unmarshalTiming(l *lexer) (stateFn, error) {
	var err error
	var td TimeDescription

	td.Timing.StartTime, err = l.readUint64Field()
	if err != nil {
		return nil, err
	}

	td.Timing.StopTime, err = l.readUint64Field()
	if err != nil {
		return nil, err
	}

	if err := l.nextLine(); err != nil {
		return nil, err
	}

	l.desc.TimeDescriptions = append(l.desc.TimeDescriptions, td)
	return sFn[9], nil
}

func unmarshalRepeatTimes(l *lexer) (stateFn, error) {
	var err error
	var newRepeatTime RepeatTime

	latestTimeDesc := &l.desc.TimeDescriptions[len(l.desc.TimeDescriptions)-1]

	field, err := l.readField()
	if err != nil {
		return nil, err
	}

	newRepeatTime.Interval, err = parseTimeUnits(field)
	if err != nil {
		return nil, fmt.Errorf("%w `%v`", errSDPInvalidValue, field)
	}

	field, err = l.readField()
	if err != nil {
		return nil, err
	}

	newRepeatTime.Duration, err = parseTimeUnits(field)
	if err != nil {
		return nil, fmt.Errorf("%w `%v`", errSDPInvalidValue, field)
	}

	for {
		field, err := l.readField()
		if err != nil {
			return nil, err
		}
		if field == "" {
			break
		}
		offset, err := parseTimeUnits(field)
		if err != nil {
			return nil, fmt.Errorf("%w `%v`", errSDPInvalidValue, field)
		}
		newRepeatTime.Offsets = append(newRepeatTime.Offsets, offset)
	}

	if err := l.nextLine(); err != nil {
		return nil, err
	}

	latestTimeDesc.RepeatTimes = append(latestTimeDesc.RepeatTimes, newRepeatTime)
	return sFn[9], nil
}

func unmarshalTimeZones(l *lexer) (stateFn, error) {
	// These fields are transimitted in pairs
	// z=<adjustment time> <offset> <adjustment time> <offset> ....
	// so we are making sure that there are actually multiple of 2 total.
	for {
		var err error
		var timeZone TimeZone

		timeZone.AdjustmentTime, err = l.readUint64Field()
		if err != nil {
			return nil, err
		}

		offset, err := l.readField()
		if err != nil {
			return nil, err
		}

		if offset == "" {
			break
		}

		timeZone.Offset, err = parseTimeUnits(offset)
		if err != nil {
			return nil, err
		}

		l.desc.TimeZones = append(l.desc.TimeZones, timeZone)
	}

	if err := l.nextLine(); err != nil {
		return nil, err
	}

	return sFn[13], nil
}

func unmarshalSessionEncryptionKey(l *lexer) (stateFn, error) {
	value, err := l.readLine()
	if err != nil {
		return nil, err
	}

	encryptionKey := EncryptionKey(value)
	l.desc.EncryptionKey = &encryptionKey
	return sFn[11], nil
}

func unmarshalSessionAttribute(l *lexer) (stateFn, error) {
	value, err := l.readLine()
	if err != nil {
		return nil, err
	}

	i := strings.IndexRune(value, ':')
	var a Attribute
	if i > 0 {
		a = NewAttribute(value[:i], value[i+1:])
	} else {
		a = NewPropertyAttribute(value)
	}

	l.desc.Attributes = append(l.desc.Attributes, a)
	return sFn[11], nil
}

func unmarshalMediaDescription(l *lexer) (stateFn, error) {
	var newMediaDesc MediaDescription

	// <media>
	field, err := l.readField()
	if err != nil {
		return nil, err
	}

	// Set according to currently registered with IANA
	// https://tools.ietf.org/html/rfc4566#section-5.14
	if !anyOf(field, "audio", "video", "text", "application", "message") {
		return nil, fmt.Errorf("%w `%v`", errSDPInvalidValue, field)
	}
	newMediaDesc.MediaName.Media = field

	// <port>
	field, err = l.readField()
	if err != nil {
		return nil, err
	}
	parts := strings.Split(field, "/")
	newMediaDesc.MediaName.Port.Value, err = parsePort(parts[0])
	if err != nil {
		return nil, fmt.Errorf("%w `%v`", errSDPInvalidPortValue, parts[0])
	}

	if len(parts) > 1 {
		portRange, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("%w `%v`", errSDPInvalidValue, parts)
		}
		newMediaDesc.MediaName.Port.Range = &portRange
	}

	// <proto>
	field, err = l.readField()
	if err != nil {
		return nil, err
	}

	// Set according to currently registered with IANA
	// https://tools.ietf.org/html/rfc4566#section-5.14
	for _, proto := range strings.Split(field, "/") {
		if !anyOf(proto, "UDP", "RTP", "AVP", "SAVP", "SAVPF", "TLS", "DTLS", "SCTP", "AVPF") {
			return nil, fmt.Errorf("%w `%v`", errSDPInvalidNumericValue, field)
		}
		newMediaDesc.MediaName.Protos = append(newMediaDesc.MediaName.Protos, proto)
	}

	// <fmt>...
	for {
		field, err = l.readField()
		if err != nil {
			return nil, err
		}
		if field == "" {
			break
		}
		newMediaDesc.MediaName.Formats = append(newMediaDesc.MediaName.Formats, field)
	}

	if err := l.nextLine(); err != nil {
		return nil, err
	}

	l.desc.MediaDescriptions = append(l.desc.MediaDescriptions, &newMediaDesc)
	return sFn[12], nil
}

func unmarshalMediaTitle(l *lexer) (stateFn, error) {
	value, err := l.readLine()
	if err != nil {
		return nil, err
	}

	latestMediaDesc := l.desc.MediaDescriptions[len(l.desc.MediaDescriptions)-1]
	mediaTitle := Information(value)
	latestMediaDesc.MediaTitle = &mediaTitle
	return sFn[16], nil
}

func unmarshalMediaConnectionInformation(l *lexer) (stateFn, error) {
	var err error
	latestMediaDesc := l.desc.MediaDescriptions[len(l.desc.MediaDescriptions)-1]
	latestMediaDesc.ConnectionInformation, err = l.unmarshalConnectionInformation()
	if err != nil {
		return nil, err
	}
	return sFn[15], nil
}

func unmarshalMediaBandwidth(l *lexer) (stateFn, error) {
	value, err := l.readLine()
	if err != nil {
		return nil, err
	}

	latestMediaDesc := l.desc.MediaDescriptions[len(l.desc.MediaDescriptions)-1]
	bandwidth, err := unmarshalBandwidth(value)
	if err != nil {
		return nil, fmt.Errorf("%w `b=%v`", errSDPInvalidSyntax, value)
	}
	latestMediaDesc.Bandwidth = append(latestMediaDesc.Bandwidth, *bandwidth)
	return sFn[15], nil
}

func unmarshalMediaEncryptionKey(l *lexer) (stateFn, error) {
	value, err := l.readLine()
	if err != nil {
		return nil, err
	}

	latestMediaDesc := l.desc.MediaDescriptions[len(l.desc.MediaDescriptions)-1]
	encryptionKey := EncryptionKey(value)
	latestMediaDesc.EncryptionKey = &encryptionKey
	return sFn[14], nil
}

func unmarshalMediaAttribute(l *lexer) (stateFn, error) {
	value, err := l.readLine()
	if err != nil {
		return nil, err
	}

	i := strings.IndexRune(value, ':')
	var a Attribute
	if i > 0 {
		a = NewAttribute(value[:i], value[i+1:])
	} else {
		a = NewPropertyAttribute(value)
	}

	latestMediaDesc := l.desc.MediaDescriptions[len(l.desc.MediaDescriptions)-1]
	latestMediaDesc.Attributes = append(latestMediaDesc.Attributes, a)
	return sFn[14], nil
}

func parseTimeUnits(value string) (int64, error) {
	// Some time offsets in the protocol can be provided with a shorthand
	// notation. This code ensures to convert it to NTP timestamp format.
	//      d - days (86400 seconds)
	//      h - hours (3600 seconds)
	//      m - minutes (60 seconds)
	//      s - seconds (allowed for completeness)
	switch value[len(value)-1:] {
	case "d":
		num, err := strconv.ParseInt(value[:len(value)-1], 10, 64)
		if err != nil {
			return 0, fmt.Errorf("%w `%v`", errSDPInvalidValue, value)
		}
		return num * 86400, nil
	case "h":
		num, err := strconv.ParseInt(value[:len(value)-1], 10, 64)
		if err != nil {
			return 0, fmt.Errorf("%w `%v`", errSDPInvalidValue, value)
		}
		return num * 3600, nil
	case "m":
		num, err := strconv.ParseInt(value[:len(value)-1], 10, 64)
		if err != nil {
			return 0, fmt.Errorf("%w `%v`", errSDPInvalidValue, value)
		}
		return num * 60, nil
	}

	num, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%w `%v`", errSDPInvalidValue, value)
	}

	return num, nil
}

func parsePort(value string) (int, error) {
	port, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("%w `%v`", errSDPInvalidPortValue, port)
	}

	if port < 0 || port > 65536 {
		return 0, fmt.Errorf("%w -- out of range `%v`", errSDPInvalidPortValue, port)
	}

	return port, nil
}
