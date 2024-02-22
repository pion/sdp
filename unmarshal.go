// SPDX-FileCopyrightText: 2023 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

package sdp

import (
	"errors"
	"fmt"
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
//
//	https://tools.ietf.org/html/rfc4566#section-5
//	https://tools.ietf.org/html/draft-ietf-rtcweb-jsep-24
//
// https://tools.ietf.org/html/rfc4566#section-5
// Session description
//
//	v=  (protocol version)
//	o=  (originator and session identifier)
//	s=  (session name)
//	i=* (session information)
//	u=* (URI of description)
//	e=* (email address)
//	p=* (phone number)
//	c=* (connection information -- not required if included in
//	     all media)
//	b=* (zero or more bandwidth information lines)
//	One or more time descriptions ("t=" and "r=" lines; see below)
//	z=* (time zone adjustments)
//	k=* (encryption key)
//	a=* (zero or more session attribute lines)
//	Zero or more media descriptions
//
// Time description
//
//	t=  (time the session is active)
//	r=* (zero or more repeat times)
//
// Media description, if present
//
//	m=  (media name and transport address)
//	i=* (media title)
//	c=* (connection information -- optional if included at
//	     session level)
//	b=* (zero or more bandwidth information lines)
//	k=* (encryption key)
//	a=* (zero or more media attribute lines)
//
// In order to generate the following state table and draw subsequent
// deterministic finite-state automota ("DFA") the following regex was used to
// derive the DFA:
//
//	vosi?u?e?p?c?b*(tr*)+z?k?a*(mi?c?b*k?a*)*
//
// possible place and state to exit:
//
//	**   * * *  ** * * * *
//	99   1 1 1  11 1 1 1 1
//	     3 1 1  26 5 5 4 4
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
func (s *SessionDescription) Unmarshal(value string) error {
	l := new(lexer)
	l.desc = s
	l.value = value

	// stats := struct{ b, t, a, m int }{}
	// for {
	// 	name, err := l.readFieldName()
	// 	if err != nil {
	// 		break
	// 	}
	// 	switch name {
	// 	case 'b':
	// 		stats.b++
	// 	case 't':
	// 		stats.t++
	// 	case 'a':
	// 		stats.a++
	// 	case 'm':
	// 		stats.m++
	// 	}
	// 	if _, err := l.readLine(); err != nil {
	// 		break
	// 	}
	// }

	// l.reset()

	// s.Bandwidth = make([]Bandwidth, 0, stats.b)
	// s.TimeDescriptions = make([]TimeDescription, 0, stats.t)
	// l.attrs = make([]Attribute, 0, stats.a)
	// // s.Attributes = make([]Attribute, 0, stats.a)
	// s.MediaDescriptions = make([]MediaDescription, 0, stats.m)

	for state := s1; state != nil; {
		var err error
		state, err = state(l)
		if err != nil {
			return err
		}
	}
	return nil
}

func s1(l *lexer) (stateFn, error) {
	return l.handleType(func(name attrName) stateFn {
		if name == 'v' {
			return unmarshalProtocolVersion
		}
		return nil
	})
}

func s2(l *lexer) (stateFn, error) {
	return l.handleType(func(name attrName) stateFn {
		if name == 'o' {
			return unmarshalOrigin
		}
		return nil
	})
}

func s3(l *lexer) (stateFn, error) {
	return l.handleType(func(name attrName) stateFn {
		if name == 's' {
			return unmarshalSessionName
		}
		return nil
	})
}

func s4(l *lexer) (stateFn, error) {
	return l.handleType(func(name attrName) stateFn {
		switch name {
		case 'i':
			return unmarshalSessionInformation
		case 'u':
			return unmarshalURI
		case 'e':
			return unmarshalEmail
		case 'p':
			return unmarshalPhone
		case 'c':
			return unmarshalSessionConnectionInformation
		case 'b':
			return unmarshalSessionBandwidth
		case 't':
			return unmarshalTiming
		}
		return nil
	})
}

func s5(l *lexer) (stateFn, error) {
	return l.handleType(func(name attrName) stateFn {
		switch name {
		case 'b':
			return unmarshalSessionBandwidth
		case 't':
			return unmarshalTiming
		}
		return nil
	})
}

func s6(l *lexer) (stateFn, error) {
	return l.handleType(func(name attrName) stateFn {
		switch name {
		case 'p':
			return unmarshalPhone
		case 'c':
			return unmarshalSessionConnectionInformation
		case 'b':
			return unmarshalSessionBandwidth
		case 't':
			return unmarshalTiming
		}
		return nil
	})
}

func s7(l *lexer) (stateFn, error) {
	return l.handleType(func(name attrName) stateFn {
		switch name {
		case 'u':
			return unmarshalURI
		case 'e':
			return unmarshalEmail
		case 'p':
			return unmarshalPhone
		case 'c':
			return unmarshalSessionConnectionInformation
		case 'b':
			return unmarshalSessionBandwidth
		case 't':
			return unmarshalTiming
		}
		return nil
	})
}

func s8(l *lexer) (stateFn, error) {
	return l.handleType(func(name attrName) stateFn {
		switch name {
		case 'c':
			return unmarshalSessionConnectionInformation
		case 'b':
			return unmarshalSessionBandwidth
		case 't':
			return unmarshalTiming
		}
		return nil
	})
}

func s9(l *lexer) (stateFn, error) {
	return l.handleType(func(name attrName) stateFn {
		switch name {
		case 'z':
			return unmarshalTimeZones
		case 'k':
			return unmarshalSessionEncryptionKey
		case 'a':
			return unmarshalSessionAttribute
		case 'r':
			return unmarshalRepeatTimes
		case 't':
			return unmarshalTiming
		case 'm':
			return unmarshalMediaDescription
		}
		return nil
	})
}

func s10(l *lexer) (stateFn, error) {
	return l.handleType(func(name attrName) stateFn {
		switch name {
		case 'e':
			return unmarshalEmail
		case 'p':
			return unmarshalPhone
		case 'c':
			return unmarshalSessionConnectionInformation
		case 'b':
			return unmarshalSessionBandwidth
		case 't':
			return unmarshalTiming
		}
		return nil
	})
}

func s11(l *lexer) (stateFn, error) {
	return l.handleType(func(name attrName) stateFn {
		switch name {
		case 'a':
			return unmarshalSessionAttribute
		case 'm':
			return unmarshalMediaDescription
		}
		return nil
	})
}

func s12(l *lexer) (stateFn, error) {
	return l.handleType(func(name attrName) stateFn {
		switch name {
		case 'a':
			return unmarshalMediaAttribute
		case 'k':
			return unmarshalMediaEncryptionKey
		case 'b':
			return unmarshalMediaBandwidth
		case 'c':
			return unmarshalMediaConnectionInformation
		case 'i':
			return unmarshalMediaTitle
		case 'm':
			return unmarshalMediaDescription
		}
		return nil
	})
}

func s13(l *lexer) (stateFn, error) {
	return l.handleType(func(name attrName) stateFn {
		switch name {
		case 'a':
			return unmarshalSessionAttribute
		case 'k':
			return unmarshalSessionEncryptionKey
		case 'm':
			return unmarshalMediaDescription
		}
		return nil
	})
}

func s14(l *lexer) (stateFn, error) {
	return l.handleType(func(name attrName) stateFn {
		switch name {
		case 'a':
			return unmarshalMediaAttribute
		case 'k':
			// Non-spec ordering
			return unmarshalMediaEncryptionKey
		case 'b':
			// Non-spec ordering
			return unmarshalMediaBandwidth
		case 'c':
			// Non-spec ordering
			return unmarshalMediaConnectionInformation
		case 'i':
			// Non-spec ordering
			return unmarshalMediaTitle
		case 'm':
			return unmarshalMediaDescription
		}
		return nil
	})
}

func s15(l *lexer) (stateFn, error) {
	return l.handleType(func(name attrName) stateFn {
		switch name {
		case 'a':
			return unmarshalMediaAttribute
		case 'k':
			return unmarshalMediaEncryptionKey
		case 'b':
			return unmarshalMediaBandwidth
		case 'c':
			return unmarshalMediaConnectionInformation
		case 'i':
			// Non-spec ordering
			return unmarshalMediaTitle
		case 'm':
			return unmarshalMediaDescription
		}
		return nil
	})
}

func s16(l *lexer) (stateFn, error) {
	return l.handleType(func(name attrName) stateFn {
		switch name {
		case 'a':
			return unmarshalMediaAttribute
		case 'k':
			return unmarshalMediaEncryptionKey
		case 'c':
			return unmarshalMediaConnectionInformation
		case 'b':
			return unmarshalMediaBandwidth
		case 'i':
			// Non-spec ordering
			return unmarshalMediaTitle
		case 'm':
			return unmarshalMediaDescription
		}
		return nil
	})
}

func unmarshalProtocolVersion(l *lexer) (stateFn, error) {
	version, err := l.readUint64Field()
	if err != nil {
		return nil, err
	}

	// As off the latest draft of the rfc this value is required to be 0.
	// https://tools.ietf.org/html/draft-ietf-rtcweb-jsep-24#section-5.8.1
	if version != 0 {
		return nil, fmt.Errorf("%w `%d`", errSDPInvalidValue, version)
	}

	if err := l.nextLine(); err != nil {
		return nil, err
	}

	return s2, nil
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
	if !anyOf(l.desc.Origin.NetworkType, kIn) {
		return nil, fmt.Errorf("%w `%s`", errSDPInvalidValue, l.desc.Origin.NetworkType)
	}

	l.desc.Origin.AddressType, err = l.readField()
	if err != nil {
		return nil, err
	}

	// Set according to currently registered with IANA
	// https://tools.ietf.org/html/rfc4566#section-8.2.7
	if !anyOf(l.desc.Origin.AddressType, kIp4, kIp6) {
		return nil, fmt.Errorf("%w `%s`", errSDPInvalidValue, l.desc.Origin.AddressType)
	}

	l.desc.Origin.UnicastAddress, err = l.readField()
	if err != nil {
		return nil, err
	}

	if err := l.nextLine(); err != nil {
		return nil, err
	}

	return s3, nil
}

func unmarshalSessionName(l *lexer) (stateFn, error) {
	value, err := l.readLine()
	if err != nil {
		return nil, err
	}

	l.desc.SessionName = SessionName(value)
	return s4, nil
}

func unmarshalSessionInformation(l *lexer) (stateFn, error) {
	value, err := l.readLine()
	if err != nil {
		return nil, err
	}

	l.desc.SessionInformation = Information(value)
	return s7, nil
}

func unmarshalURI(l *lexer) (stateFn, error) {
	value, err := l.readLine()
	if err != nil {
		return nil, err
	}

	l.desc.URI = URI(value)

	return s10, nil
}

func unmarshalEmail(l *lexer) (stateFn, error) {
	value, err := l.readLine()
	if err != nil {
		return nil, err
	}

	l.desc.EmailAddress = EmailAddress(value)
	return s6, nil
}

func unmarshalPhone(l *lexer) (stateFn, error) {
	value, err := l.readLine()
	if err != nil {
		return nil, err
	}

	l.desc.PhoneNumber = PhoneNumber(value)
	return s8, nil
}

func unmarshalSessionConnectionInformation(l *lexer) (stateFn, error) {
	var err error
	l.desc.ConnectionInformation, err = l.unmarshalConnectionInformation()
	if err != nil {
		return nil, err
	}
	return s5, nil
}

func (l *lexer) unmarshalConnectionInformation() (c ConnectionInformation, err error) {
	c.NetworkType, err = l.readField()
	if err != nil {
		return c, err
	}

	// Set according to currently registered with IANA
	// https://tools.ietf.org/html/rfc4566#section-8.2.6
	if !anyOf(c.NetworkType, kIn) {
		return c, fmt.Errorf("%w `%s`", errSDPInvalidValue, c.NetworkType)
	}

	c.AddressType, err = l.readField()
	if err != nil {
		return c, err
	}

	// Set according to currently registered with IANA
	// https://tools.ietf.org/html/rfc4566#section-8.2.7
	if !anyOf(c.AddressType, kIp4, kIp6) {
		return c, fmt.Errorf("%w `%s`", errSDPInvalidValue, c.AddressType)
	}

	address, err := l.readField()
	if err != nil {
		return c, err
	}

	c.Address.Address = address

	if err := l.nextLine(); err != nil {
		return c, err
	}

	return c, nil
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

	return s5, nil
}

func unmarshalBandwidth(value string) (*Bandwidth, error) {
	i := strings.IndexRune(value, ':')
	if i == -1 {
		return nil, fmt.Errorf("%w `b=%v`", errSDPInvalidValue, value)
	}

	typ := value[:i]
	experimental := strings.HasPrefix(typ, kExperimental)
	if experimental {
		typ = strings.TrimPrefix(typ, kExperimental)
	} else if !anyOf(typ, kCt, kAs, kTias, kRs, kRr) {
		// Set according to currently registered with IANA
		// https://tools.ietf.org/html/rfc4566#section-5.8
		// https://tools.ietf.org/html/rfc3890#section-6.2
		// https://tools.ietf.org/html/rfc3556#section-2
		return nil, fmt.Errorf("%w `%s`", errSDPInvalidValue, typ)
	}

	bandwidth, ok := parseUint(value[i+1:], 64)
	if !ok {
		return nil, fmt.Errorf("%w `%s`", errSDPInvalidNumericValue, value[i+1:])
	}

	return &Bandwidth{
		Experimental: experimental,
		Type:         typ,
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
	return s9, nil
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
		return nil, fmt.Errorf("%w `%s`", errSDPInvalidValue, field)
	}

	field, err = l.readField()
	if err != nil {
		return nil, err
	}

	newRepeatTime.Duration, err = parseTimeUnits(field)
	if err != nil {
		return nil, fmt.Errorf("%w `%s`", errSDPInvalidValue, field)
	}

	for {
		field, err := l.readField()
		if err != nil {
			return nil, err
		}
		if len(field) == 0 {
			break
		}
		offset, err := parseTimeUnits(field)
		if err != nil {
			return nil, fmt.Errorf("%w `%s`", errSDPInvalidValue, field)
		}
		newRepeatTime.Offsets = append(newRepeatTime.Offsets, offset)
	}

	if err := l.nextLine(); err != nil {
		return nil, err
	}

	latestTimeDesc.RepeatTimes = append(latestTimeDesc.RepeatTimes, newRepeatTime)
	return s9, nil
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

		if len(offset) == 0 {
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

	return s13, nil
}

func unmarshalSessionEncryptionKey(l *lexer) (stateFn, error) {
	value, err := l.readLine()
	if err != nil {
		return nil, err
	}

	l.desc.EncryptionKey = EncryptionKey(value)
	return s11, nil
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
	return s11, nil
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
	if !anyOf(field, kAudio, kVideo, kText, kApplication, kMessage) {
		return nil, fmt.Errorf("%w `%s`", errSDPInvalidValue, field)
	}
	newMediaDesc.MediaName.Media = field

	// <port>
	field, err = l.readField()
	if err != nil {
		return nil, err
	}

	i := strings.IndexRune(field, '/')
	if i == -1 {
		i = len(field)
	} else {
		portRange, ok := parseUint(field[i+1:], 16)
		if !ok {
			return nil, fmt.Errorf("%w `%s`", errSDPInvalidValue, field[i+1:])
		}
		newMediaDesc.MediaName.Port.Range = uint16(portRange)
	}

	port, ok := parseUint(field[:i], 16)
	if !ok {
		return nil, fmt.Errorf("%w `%s`", errSDPInvalidPortValue, field[:i])
	}
	newMediaDesc.MediaName.Port.Value = uint16(port)

	// <proto>
	field, err = l.readField()
	if err != nil {
		return nil, err
	}

	// Set according to currently registered with IANA
	// https://tools.ietf.org/html/rfc4566#section-5.14
	// https://tools.ietf.org/html/rfc4975#section-8.1
	newMediaDesc.MediaName.Protos = make([]string, 0, countSegments(field, '/'))
	for pos := 0; pos < len(field); pos++ {
		field = field[pos:]
		if pos = strings.IndexRune(field, '/'); pos == -1 {
			pos = len(field)
		}
		proto := field[:pos]
		if !anyOf(proto, kUdp, kRtp, kAvp, kSavp, kSavpf, kTls, kDtls, kSctp, kAvpf, kTcp, kMsrp) {
			return nil, fmt.Errorf("%w `%s`", errSDPInvalidValue, proto)
		}
		newMediaDesc.MediaName.Protos = append(newMediaDesc.MediaName.Protos, proto)
	}

	// <fmt>...
	for {
		field, err = l.readField()
		if err != nil {
			return nil, err
		}
		if len(field) == 0 {
			break
		}
		newMediaDesc.MediaName.Formats = append(newMediaDesc.MediaName.Formats, field)
	}

	if err := l.nextLine(); err != nil {
		return nil, err
	}

	l.desc.MediaDescriptions = append(l.desc.MediaDescriptions, newMediaDesc)
	return s12, nil
}

func unmarshalMediaTitle(l *lexer) (stateFn, error) {
	value, err := l.readLine()
	if err != nil {
		return nil, err
	}

	latestMediaDesc := l.desc.MediaDescriptions[len(l.desc.MediaDescriptions)-1]
	latestMediaDesc.MediaTitle = Information(value)
	return s16, nil
}

func unmarshalMediaConnectionInformation(l *lexer) (stateFn, error) {
	var err error
	latestMediaDesc := l.desc.MediaDescriptions[len(l.desc.MediaDescriptions)-1]
	latestMediaDesc.ConnectionInformation, err = l.unmarshalConnectionInformation()
	if err != nil {
		return nil, err
	}
	return s15, nil
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
	return s15, nil
}

func unmarshalMediaEncryptionKey(l *lexer) (stateFn, error) {
	value, err := l.readLine()
	if err != nil {
		return nil, err
	}

	latestMediaDesc := l.desc.MediaDescriptions[len(l.desc.MediaDescriptions)-1]
	latestMediaDesc.EncryptionKey = EncryptionKey(value)
	return s14, nil
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
	return s14, nil
}

func parseTimeUnits(value string) (num int64, err error) {
	if len(value) == 0 {
		return 0, fmt.Errorf("%w `%s`", errSDPInvalidValue, value)
	}
	k, ok := timeShorthand(value[len(value)-1])
	if ok {
		value = value[:len(value)-1]
	}
	num, ok = parseInt(value)
	if !ok {
		return 0, fmt.Errorf("%w `%s`", errSDPInvalidValue, value)
	}
	return num * k, nil
}

func timeShorthand(b byte) (int64, bool) {
	// Some time offsets in the protocol can be provided with a shorthand
	// notation. This code ensures to convert it to NTP timestamp format.
	switch b {
	case 'd': // days
		return 86400, true
	case 'h': // hours
		return 3600, true
	case 'm': // minutes
		return 60, true
	case 's': // seconds (allowed for completeness)
		return 1, true
	default:
		return 1, false
	}
}

func parseInt(value string) (int64, bool) {
	sign := int64(1)
	if len(value) != 0 && value[0] == '-' {
		sign = -1
		value = value[1:]
	}
	n, ok := parseUint(value, 64)
	if !ok {
		return 0, false
	}
	return sign * int64(n), true
}

func parseUint(value string, bits int) (uint64, bool) {
	var n uint64
	for _, ch := range value {
		if ch < '0' || ch > '9' {
			return 0, false
		}

		n = n*10 + uint64(ch-'0')
	}
	return n, n <= uint64(1<<bits)-1
}

func countSegments(value string, r rune) int {
	n := 1
	for pos := 0; pos < len(value); {
		i := strings.IndexRune(value[pos:], r)
		if i == -1 {
			break
		}
		pos += i + 1
		n++
	}
	return n
}
