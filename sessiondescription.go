// Package sdp implements the Session Description Protocol as defined in RFC 4566.
//
// The States Sransition Table below describes the computation flow between
// functions (namely s1, s2, s3, ...) for a decoding procedures defined in
// [RFC4566]. Additionally, the session description and its fields defined below
// provide additional functionalities outside the scope of [RFC4566] in order to
// provide support for specifications described in [JSEP], [ICE], [RTP], [SCTP],
// [BUNDLE], and more.
//
// Visualization
//
// Below is a generalized textual representation of the session discription as
// described in [RFC4566]. It is presented here for guidance of naming and
// structure of the general SDP. Some implementations may chose not to use all
// available fields for necessary functionalities.
//  Session description
//     v=  (protocol version)
//     o=  (originator and session identifier)
//     s=  (session name)
//     i=* (session information)
//     u=* (URI of description)
//     e=* (email address)
//     p=* (phone number)
//     c=* (connection information -- not required if included in all media)
//     b=* (zero or more bandwidth information lines)
//     One or more time descriptions ("t=" and "r=" lines; see below)
//     z=* (time zone adjustments)
//     k=* (encryption key)
//     a=* (zero or more session attribute lines)
//     Zero or more media descriptions
//
//  Time description
//     t=  (time the session is active)
//     r=* (zero or more repeat times)
//
//  Type description, if present
//     m=  (media name and transport address)
//     i=* (media title)
//     c=* (connection information -- optional if included at session level)
//     b=* (zero or more bandwidth information lines)
//     k=* (encryption key)
//     a=* (zero or more media attribute lines)
//
// Example
//
// A general example SDP description with two media sections.
//  v=0
//  o=jdoe 2890844526 2890842807 IN IP4 10.47.16.5
//  s=SDP Seminar
//  i=A Seminar on the session description protocol
//  u=http://www.example.com/seminars/sdp.pdf
//  e=j.doe@example.com (Jane Doe)
//  c=IN IP4 224.2.17.12/127
//  t=2873397496 2873404696
//  a=recvonly
//  m=audio 49170 RTP/AVP 0
//  m=video 51372 RTP/AVP 99
//  a=rtpmap:99 h263-1998/90000
//
// Implementation Details
//
// In order to generate the following state transition table and draw subsequent
// deterministic finite-state automota ("DFA") the following regex was used to
// derive the DFA:
// 	vosi?u?e?p?c?b*(tr*)+z?k?a*(mi?c?b*k?a*)*
//
// Each of the letters in the regex above represents the key of a particular sdp
// line. (e.g. v for v=, o for o=) As the decoding process scans one line at a
// time, the decoder parses said line and determines what line to expect next
// based on the state to which it transitions.
//
// Note: In the table below in order to distinguish between the states belonging
// to the media description as opposed to the session description, with regards
// to the "k" and "a" decoding states, the states are marked with an asterisk
// "a*", "k*". These asterisk marked states represent that a= and k= lines of
// the media description sections.
// 	+--------+----+-------+----+-----+----+-----+---+----+----+---+---+-----+---+---+----+---+----+
// 	| STATES | a* | a*,k* | a  | a,k | b  | b,c | e | i  | m  | o | p | r,t | s | t | u  | v | z  |
// 	+--------+----+-------+----+-----+----+-----+---+----+----+---+---+-----+---+---+----+---+----+
// 	|   s1   |    |       |    |     |    |     |   |    |    |   |   |     |   |   |    | 2 |    |
// 	|   s2   |    |       |    |     |    |     |   |    |    | 3 |   |     |   |   |    |   |    |
// 	|   s3   |    |       |    |     |    |     |   |    |    |   |   |     | 4 |   |    |   |    |
// 	|   s4   |    |       |    |     |    |   5 | 6 |  7 |    |   | 8 |     |   | 9 | 10 |   |    |
// 	|   s5   |    |       |    |     |  5 |     |   |    |    |   |   |     |   | 9 |    |   |    |
// 	|   s6   |    |       |    |     |    |   5 |   |    |    |   | 8 |     |   | 9 |    |   |    |
// 	|   s7   |    |       |    |     |    |   5 | 6 |    |    |   | 8 |     |   | 9 | 10 |   |    |
// 	|   s8   |    |       |    |     |    |   5 |   |    |    |   |   |     |   | 9 |    |   |    |
// 	|   s9   |    |       |    |  11 |    |     |   |    | 12 |   |   |   9 |   |   |    |   | 13 |
// 	|   s10  |    |       |    |     |    |   5 | 6 |    |    |   | 8 |     |   | 9 |    |   |    |
// 	|   s11  |    |       | 11 |     |    |     |   |    | 12 |   |   |     |   |   |    |   |    |
// 	|   s12  |    |    14 |    |     |    |  15 |   | 16 | 12 |   |   |     |   |   |    |   |    |
// 	|   s13  |    |       |    |  11 |    |     |   |    | 12 |   |   |     |   |   |    |   |    |
// 	|   s14  | 14 |       |    |     |    |     |   |    | 12 |   |   |     |   |   |    |   |    |
// 	|   s15  |    |    14 |    |     | 15 |     |   |    | 12 |   |   |     |   |   |    |   |    |
// 	|   s16  |    |    14 |    |     |    |  15 |   |    | 12 |   |   |     |   |   |    |   |    |
// 	+--------+----+-------+----+-----+----+-----+---+----+----+---+---+-----+---+---+----+---+----+
//
// References
//
// [RFC4566] https://tools.ietf.org/html/rfc4566
//
// [JSEP] https://tools.ietf.org/html/draft-ietf-rtcweb-jsep-24
//
// [ICE] https://tools.ietf.org/html/rfc5245
//
// [RTP] https://tools.ietf.org/html/rfc3550
//
// [SCTP] https://tools.ietf.org/html/rfc4960
//
// [BUNDLE] https://tools.ietf.org/html/draft-ietf-mmusic-sdp-bundle-negotiation-53
package sdp

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/pions/webrtc/pkg/rtcerr"

	"github.com/pkg/errors"
)

// SessionDescription is a a well-defined format for conveying sufficient
// information to discover and participate in a multimedia session.
type SessionDescription struct {
	// v=0
	Version Version

	// o=<username> <sess-id> <sess-version> <nettype> <addrtype> <unicast-address>
	Origin Origin

	// s=<session name>
	Session Session

	// i=<session description>
	Information *Information

	// u=<uri>
	URI *URL

	// e=<email-address>
	EmailAddress *EmailAddress

	// p=<phone-number>
	PhoneNumber *PhoneNumber

	// c=<nettype> <addrtype> <connection-address>
	Connection *Connection

	// b=<bwtype>:<bandwidth>
	Bandwidths Bandwidths

	TimeDescriptions TimeDescriptions

	// z=<adjustment time> <offset> <adjustment time> <offset> ...
	TimeZones TimeZones

	// k=<method>
	// k=<method>:<encryption key>
	EncryptionKey *EncryptionKey

	// a=<attribute>
	// a=<attribute>:<value>
	Attributes SessionAttributes

	MediaDescriptions MediaDescriptions
}

func (s *SessionDescription) Clone() *SessionDescription {
	desc := &SessionDescription{}
	desc.Version = *s.Version.Clone()
	desc.Origin = *s.Origin.Clone()
	desc.Session = *s.Session.Clone()

	if s.Information != nil {
		desc.Information = s.Information.Clone()
	}

	if s.URI != nil {
		desc.URI = s.URI.Clone()
	}

	if s.EmailAddress != nil {
		desc.EmailAddress = s.EmailAddress.Clone()
	}

	if s.PhoneNumber != nil {
		desc.PhoneNumber = s.PhoneNumber.Clone()
	}

	if s.Connection != nil {
		desc.Connection = s.Connection.Clone()
	}

	if len(s.Bandwidths) > 0 {
		desc.Bandwidths = *s.Bandwidths.Clone()
	}

	if len(s.TimeDescriptions) > 0 {
		desc.TimeDescriptions = *s.TimeDescriptions.Clone()
	}

	if len(s.TimeZones) > 0 {
		desc.TimeZones = *s.TimeZones.Clone()
	}

	if s.EncryptionKey != nil {
		desc.EncryptionKey = s.EncryptionKey.Clone()
	}

	if len(s.Attributes) > 0 {
		desc.Attributes = *s.Attributes.Clone()
	}

	if len(s.MediaDescriptions) > 0 {
		desc.MediaDescriptions = *s.MediaDescriptions.Clone()
	}

	return desc
}

// Marshal converts an SDP struct to text.
func (s *SessionDescription) Marshal() (raw string) {
	raw += s.Version.Marshal()
	raw += s.Origin.Marshal()
	raw += s.Session.Marshal()

	if s.Information != nil {
		raw += s.Information.Marshal()
	}

	if s.URI != nil {
		raw += s.URI.Marshal()
	}

	if s.EmailAddress != nil {
		raw += s.EmailAddress.Marshal()
	}

	if s.PhoneNumber != nil {
		raw += s.PhoneNumber.Marshal()
	}

	if s.Connection != nil {
		raw += s.Connection.Marshal()
	}

	for _, b := range s.Bandwidths {
		raw += b.Marshal()
	}

	for _, td := range s.TimeDescriptions {
		raw += td.Timing.Marshal()
		for _, r := range td.RepeatTimes {
			raw += r.Marshal()
		}
	}

	if len(s.TimeZones) > 0 {
		raw += s.TimeZones.Marshal()
	}

	if s.EncryptionKey != nil {
		raw += s.EncryptionKey.Marshal()
	}

	for _, a := range s.Attributes {
		raw += a.Marshal()
	}

	for _, md := range s.MediaDescriptions {
		raw += md.Media.Marshal()

		if md.Information != nil {
			raw += md.Information.Marshal()
		}

		if md.Connection != nil {
			raw += md.Connection.Marshal()
		}

		for _, b := range md.Bandwidths {
			raw += b.Marshal()
		}

		if md.EncryptionKey != nil {
			raw += md.EncryptionKey.Marshal()
		}

		for _, a := range md.Attributes {
			raw += a.Marshal()
		}
	}
	return raw
}

// Unmarshal is the primary function that deserializes the session description
// message and stores it inside of a structured SessionDescription object.
func (s *SessionDescription) Unmarshal(value string) error {
	l := &lexer{
		desc:  s,
		input: bufio.NewReader(strings.NewReader(value)),
	}
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
	key, err := readType(l.input)
	if err != nil {
		return nil, err
	}

	if key == versionKey {
		return unmarshalProtocolVersion, nil
	}

	return nil, errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", key)}, pkgName)
}

func s2(l *lexer) (stateFn, error) {
	key, err := readType(l.input)
	if err != nil {
		return nil, err
	}

	if key == originKey {
		return unmarshalOrigin, nil
	}

	return nil, errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", key)}, pkgName)
}

func s3(l *lexer) (stateFn, error) {
	key, err := readType(l.input)
	if err != nil {
		return nil, err
	}

	if key == sessionKey {
		return unmarshalSession, nil
	}

	return nil, errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", key)}, pkgName)
}

func s4(l *lexer) (stateFn, error) {
	key, err := readType(l.input)
	if err != nil {
		return nil, err
	}

	switch key {
	case infoKey:
		return unmarshalSessionInformation, nil
	case uriKey:
		return unmarshalURI, nil
	case emailKey:
		return unmarshalEmail, nil
	case phoneKey:
		return unmarshalPhone, nil
	case connectionKey:
		return unmarshalSessionConnection, nil
	case bandwidthKey:
		return unmarshalSessionBandwidth, nil
	case timingKey:
		return unmarshalTiming, nil
	}

	return nil, errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", key)}, pkgName)
}

func s5(l *lexer) (stateFn, error) {
	key, err := readType(l.input)
	if err != nil {
		return nil, err
	}

	switch key {
	case bandwidthKey:
		return unmarshalSessionBandwidth, nil
	case timingKey:
		return unmarshalTiming, nil
	}

	return nil, errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", key)}, pkgName)
}

func s6(l *lexer) (stateFn, error) {
	key, err := readType(l.input)
	if err != nil {
		return nil, err
	}

	switch key {
	case phoneKey:
		return unmarshalPhone, nil
	case connectionKey:
		return unmarshalSessionConnection, nil
	case bandwidthKey:
		return unmarshalSessionBandwidth, nil
	case timingKey:
		return unmarshalTiming, nil
	}

	return nil, errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", key)}, pkgName)
}

func s7(l *lexer) (stateFn, error) {
	key, err := readType(l.input)
	if err != nil {
		return nil, err
	}

	switch key {
	case uriKey:
		return unmarshalURI, nil
	case emailKey:
		return unmarshalEmail, nil
	case phoneKey:
		return unmarshalPhone, nil
	case connectionKey:
		return unmarshalSessionConnection, nil
	case bandwidthKey:
		return unmarshalSessionBandwidth, nil
	case timingKey:
		return unmarshalTiming, nil
	}

	return nil, errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", key)}, pkgName)
}

func s8(l *lexer) (stateFn, error) {
	key, err := readType(l.input)
	if err != nil {
		return nil, err
	}

	switch key {
	case connectionKey:
		return unmarshalSessionConnection, nil
	case bandwidthKey:
		return unmarshalSessionBandwidth, nil
	case timingKey:
		return unmarshalTiming, nil
	}

	return nil, errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", key)}, pkgName)
}

func s9(l *lexer) (stateFn, error) {
	key, err := readType(l.input)
	if err != nil {
		if errors.Cause(err) == io.EOF && key == "" {
			return nil, nil
		}
		return nil, err
	}

	switch key {
	case timeZonesKey:
		return unmarshalTimeZones, nil
	case encryptionKey:
		return unmarshalSessionEncryptionKey, nil
	case attributeKey:
		return unmarshalSessionAttribute, nil
	case repeatTimeKey:
		return unmarshalRepeatTimes, nil
	case timingKey:
		return unmarshalTiming, nil
	case mediaKey:
		return unmarshalMediaDescription, nil
	}

	return nil, errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", key)}, pkgName)
}

func s10(l *lexer) (stateFn, error) {
	key, err := readType(l.input)
	if err != nil {
		return nil, err
	}

	switch key {
	case emailKey:
		return unmarshalEmail, nil
	case phoneKey:
		return unmarshalPhone, nil
	case connectionKey:
		return unmarshalSessionConnection, nil
	case bandwidthKey:
		return unmarshalSessionBandwidth, nil
	case timingKey:
		return unmarshalTiming, nil
	}

	return nil, errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", key)}, pkgName)
}

func s11(l *lexer) (stateFn, error) {
	key, err := readType(l.input)
	if err != nil {
		if errors.Cause(err) == io.EOF && key == "" {
			return nil, nil
		}
		return nil, err
	}

	switch key {
	case attributeKey:
		return unmarshalSessionAttribute, nil
	case mediaKey:
		return unmarshalMediaDescription, nil
	}

	return nil, errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", key)}, pkgName)
}

func s12(l *lexer) (stateFn, error) {
	key, err := readType(l.input)
	if err != nil {
		if errors.Cause(err) == io.EOF && key == "" {
			return nil, nil
		}
		return nil, err
	}

	switch key {
	case attributeKey:
		return unmarshalMediaAttribute, nil
	case encryptionKey:
		return unmarshalMediaEncryptionKey, nil
	case bandwidthKey:
		return unmarshalMediaBandwidth, nil
	case connectionKey:
		return unmarshalMediaConnection, nil
	case infoKey:
		return unmarshalMediaInformation, nil
	case mediaKey:
		return unmarshalMediaDescription, nil
	}

	return nil, errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", key)}, pkgName)
}

func s13(l *lexer) (stateFn, error) {
	key, err := readType(l.input)
	if err != nil {
		if errors.Cause(err) == io.EOF && key == "" {
			return nil, nil
		}
		return nil, err
	}

	switch key {
	case attributeKey:
		return unmarshalSessionAttribute, nil
	case encryptionKey:
		return unmarshalSessionEncryptionKey, nil
	case mediaKey:
		return unmarshalMediaDescription, nil
	}

	return nil, errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", key)}, pkgName)
}

func s14(l *lexer) (stateFn, error) {
	key, err := readType(l.input)
	if err != nil {
		if errors.Cause(err) == io.EOF && key == "" {
			return nil, nil
		}
		return nil, err
	}

	switch key {
	case attributeKey:
		return unmarshalMediaAttribute, nil
	case mediaKey:
		return unmarshalMediaDescription, nil
	}

	return nil, errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", key)}, pkgName)
}

func s15(l *lexer) (stateFn, error) {
	key, err := readType(l.input)
	if err != nil {
		if errors.Cause(err) == io.EOF && key == "" {
			return nil, nil
		}
		return nil, err
	}

	switch key {
	case attributeKey:
		return unmarshalMediaAttribute, nil
	case encryptionKey:
		return unmarshalMediaEncryptionKey, nil
	case bandwidthKey:
		return unmarshalMediaBandwidth, nil
	case mediaKey:
		return unmarshalMediaDescription, nil
	}

	return nil, errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", key)}, pkgName)
}

func s16(l *lexer) (stateFn, error) {
	key, err := readType(l.input)
	if err != nil {
		if errors.Cause(err) == io.EOF && key == "" {
			return nil, nil
		}
		return nil, err
	}

	switch key {
	case attributeKey:
		return unmarshalMediaAttribute, nil
	case encryptionKey:
		return unmarshalMediaEncryptionKey, nil
	case connectionKey:
		return unmarshalMediaConnection, nil
	case bandwidthKey:
		return unmarshalMediaBandwidth, nil
	case mediaKey:
		return unmarshalMediaDescription, nil
	}

	return nil, errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", key)}, pkgName)
}

func unmarshalProtocolVersion(l *lexer) (stateFn, error) {
	value, err := readValue(l.input)
	if err != nil {
		return nil, err
	}

	version := Version{}
	if err := version.Unmarshal(value); err != nil {
		return nil, err
	}

	l.desc.Version = version
	return s2, nil
}

func unmarshalOrigin(l *lexer) (stateFn, error) {
	value, err := readValue(l.input)
	if err != nil {
		return nil, err
	}

	origin := Origin{}
	if err := origin.Unmarshal(value); err != nil {
		return nil, err
	}

	l.desc.Origin = origin
	return s3, nil
}

func unmarshalSession(l *lexer) (stateFn, error) {
	value, err := readValue(l.input)
	if err != nil {
		return nil, err
	}

	session := Session{}
	if err := session.Unmarshal(value); err != nil {
		return nil, err
	}

	l.desc.Session = session
	return s4, nil
}

func unmarshalSessionInformation(l *lexer) (stateFn, error) {
	value, err := readValue(l.input)
	if err != nil {
		return nil, err
	}

	info := Information{}
	if err := info.Unmarshal(value); err != nil {
		return nil, err
	}

	l.desc.Information = &info
	return s7, nil
}

func unmarshalURI(l *lexer) (stateFn, error) {
	value, err := readValue(l.input)
	if err != nil {
		return nil, err
	}

	uri := URL{}
	if err := uri.Unmarshal(value); err != nil {
		return nil, err
	}

	l.desc.URI = &uri
	return s10, nil
}

func unmarshalEmail(l *lexer) (stateFn, error) {
	value, err := readValue(l.input)
	if err != nil {
		return nil, err
	}

	email := EmailAddress{}
	if err := email.Unmarshal(value); err != nil {
		return nil, err
	}

	l.desc.EmailAddress = &email
	return s6, nil
}

func unmarshalPhone(l *lexer) (stateFn, error) {
	value, err := readValue(l.input)
	if err != nil {
		return nil, err
	}

	phone := PhoneNumber{}
	if err := phone.Unmarshal(value); err != nil {
		return nil, err
	}

	l.desc.PhoneNumber = &phone
	return s8, nil
}

func unmarshalSessionConnection(l *lexer) (stateFn, error) {
	value, err := readValue(l.input)
	if err != nil {
		return nil, err
	}

	connection := Connection{}
	if err := connection.Unmarshal(value); err != nil {
		return nil, err
	}

	l.desc.Connection = &connection
	return s5, nil
}

func unmarshalSessionBandwidth(l *lexer) (stateFn, error) {
	value, err := readValue(l.input)
	if err != nil {
		return nil, err
	}

	bandwidth := Bandwidth{}
	if err := bandwidth.Unmarshal(value); err != nil {
		return nil, err
	}

	l.desc.Bandwidths = append(l.desc.Bandwidths, bandwidth)
	return s5, nil
}

func unmarshalTiming(l *lexer) (stateFn, error) {
	value, err := readValue(l.input)
	if err != nil {
		return nil, err
	}

	td := TimeDescription{}
	if err := td.Timing.Unmarshal(value); err != nil {
		return nil, err
	}

	l.desc.TimeDescriptions = append(l.desc.TimeDescriptions, td)
	return s9, nil
}

func unmarshalRepeatTimes(l *lexer) (stateFn, error) {
	value, err := readValue(l.input)
	if err != nil {
		return nil, err
	}

	repeatTime := RepeatTime{}
	if err := repeatTime.Unmarshal(value); err != nil {
		return nil, err
	}

	latestTimeDesc := &l.desc.TimeDescriptions[len(l.desc.TimeDescriptions)-1]
	latestTimeDesc.RepeatTimes = append(latestTimeDesc.RepeatTimes, repeatTime)
	return s9, nil
}

func unmarshalTimeZones(l *lexer) (stateFn, error) {
	value, err := readValue(l.input)
	if err != nil {
		return nil, err
	}

	timeZones := TimeZones{}
	if err := timeZones.Unmarshal(value); err != nil {
		return nil, err
	}

	l.desc.TimeZones = append(l.desc.TimeZones, timeZones...)
	return s13, nil
}

func unmarshalSessionEncryptionKey(l *lexer) (stateFn, error) {
	value, err := readValue(l.input)
	if err != nil {
		return nil, err
	}

	encryptionKey := EncryptionKey{}
	if err := encryptionKey.Unmarshal(value); err != nil {
		return nil, err
	}

	l.desc.EncryptionKey = &encryptionKey
	return s11, nil
}

func unmarshalSessionAttribute(l *lexer) (stateFn, error) {
	value, err := readValue(l.input)
	if err != nil {
		return nil, err
	}

	if err := l.desc.Attributes.Unmarshal(value); err != nil {
		return nil, err
	}
	return s11, nil
}

func unmarshalMediaDescription(l *lexer) (stateFn, error) {
	value, err := readValue(l.input)
	if err != nil {
		return nil, err
	}

	mediaDesc := MediaDescription{}
	if err := mediaDesc.Media.Unmarshal(value); err != nil {
		return nil, err
	}

	l.desc.MediaDescriptions = append(l.desc.MediaDescriptions, mediaDesc)
	return s12, nil
}

func unmarshalMediaInformation(l *lexer) (stateFn, error) {
	value, err := readValue(l.input)
	if err != nil {
		return nil, err
	}

	info := Information{}
	if err := info.Unmarshal(value); err != nil {
		return nil, err
	}

	lastMediaDesc := &l.desc.MediaDescriptions[len(l.desc.MediaDescriptions)-1]
	lastMediaDesc.Information = &info
	return s16, nil
}

func unmarshalMediaConnection(l *lexer) (stateFn, error) {
	value, err := readValue(l.input)
	if err != nil {
		return nil, err
	}

	connection := Connection{}
	if err := connection.Unmarshal(value); err != nil {
		return nil, err
	}

	lastMediaDesc := &l.desc.MediaDescriptions[len(l.desc.MediaDescriptions)-1]
	lastMediaDesc.Connection = &connection
	return s15, nil
}

func unmarshalMediaBandwidth(l *lexer) (stateFn, error) {
	value, err := readValue(l.input)
	if err != nil {
		return nil, err
	}

	bandwidth := Bandwidth{}
	if err := bandwidth.Unmarshal(value); err != nil {
		return nil, err
	}

	lastMediaDesc := &l.desc.MediaDescriptions[len(l.desc.MediaDescriptions)-1]
	lastMediaDesc.Bandwidths = append(lastMediaDesc.Bandwidths, bandwidth)
	return s15, nil
}

func unmarshalMediaEncryptionKey(l *lexer) (stateFn, error) {
	value, err := readValue(l.input)
	if err != nil {
		return nil, err
	}

	encryptionKey := EncryptionKey{}
	if err := encryptionKey.Unmarshal(value); err != nil {
		return nil, err
	}

	lastMediaDesc := &l.desc.MediaDescriptions[len(l.desc.MediaDescriptions)-1]
	lastMediaDesc.EncryptionKey = &encryptionKey
	return s14, nil
}

func unmarshalMediaAttribute(l *lexer) (stateFn, error) {
	value, err := readValue(l.input)
	if err != nil {
		return nil, err
	}

	lastMediaDesc := &l.desc.MediaDescriptions[len(l.desc.MediaDescriptions)-1]
	if err := lastMediaDesc.Attributes.Unmarshal(value); err != nil {
		return nil, err
	}
	return s14, nil
}
