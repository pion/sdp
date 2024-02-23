// SPDX-FileCopyrightText: 2023 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

package sdp

import (
	"fmt"
	"sync"
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
func (s *SessionDescription) Unmarshal(value string) (err error) {
	c := cachePool.Get().(*cache)
	*s, err = parseDescription(value, c)
	cachePool.Put(c)
	return
}

var cachePool = sync.Pool{
	New: func() interface{} {
		return &cache{}
	},
}

type cache struct {
	ss []string
	b  []Bandwidth
	t  []TimeDescription
	rt []RepeatTime
	a  []Attribute
	m  []MediaDescription
}

func (c *cache) reset() {
	c.ss = c.ss[:0]
	c.b = c.b[:0]
	c.t = c.t[:0]
	c.rt = c.rt[:0]
	c.a = c.a[:0]
	c.m = c.m[:0]
}

func (c *cache) cloneBandwidth() []Bandwidth {
	if len(c.b) == 0 {
		return nil
	}
	s := make([]Bandwidth, len(c.b))
	copy(s, c.b)
	c.b = c.b[:0]
	return s
}

func (c *cache) getBandwidth() *Bandwidth {
	c.b = append(c.b, Bandwidth{})
	return &c.b[len(c.b)-1]
}

func (c *cache) cloneTimeDescription() []TimeDescription {
	if len(c.t) == 0 {
		return nil
	}
	s := make([]TimeDescription, len(c.t))
	copy(s, c.t)
	c.t = c.t[:0]
	return s
}

func (c *cache) getTimeDescription() *TimeDescription {
	c.t = append(c.t, TimeDescription{})
	return &c.t[len(c.t)-1]
}

func (c *cache) cloneRepeatTime() []RepeatTime {
	if len(c.rt) == 0 {
		return nil
	}
	s := make([]RepeatTime, len(c.rt))
	copy(s, c.rt)
	c.rt = c.rt[:0]
	return s
}

func (c *cache) getRepeatTime() *RepeatTime {
	c.rt = append(c.rt, RepeatTime{})
	return &c.rt[len(c.rt)-1]
}

func (c *cache) cloneAttribute() []Attribute {
	if len(c.a) == 0 {
		return nil
	}
	s := make([]Attribute, len(c.a))
	copy(s, c.a)
	c.a = c.a[:0]
	return s
}

func (c *cache) getAttribute() *Attribute {
	c.a = append(c.a, Attribute{})
	return &c.a[len(c.a)-1]
}

func (c *cache) cloneMediaDescription() []MediaDescription {
	if len(c.m) == 0 {
		return nil
	}
	s := make([]MediaDescription, len(c.m))
	copy(s, c.m)
	c.m = c.m[:0]
	return s
}

func (c *cache) getMediaDescription() *MediaDescription {
	c.m = append(c.m, MediaDescription{})
	return &c.m[len(c.m)-1]
}

func parseDescription(s string, c *cache) (d SessionDescription, err error) {
	c.reset()

	if len(s) == 0 || s[0] != 'v' {
		return d, fmt.Errorf("expected version found: %q", s)
	}
	s, err = parseVersion(&d.Version, s[1:])
	if err != nil {
		return d, err
	}

	if len(s) == 0 || s[0] != 'o' {
		return d, fmt.Errorf("expected origin found: %q", s)
	}
	s, err = parseOrigin(&d.Origin, s[1:])
	if err != nil {
		return d, err
	}

	if len(s) == 0 || s[0] != 's' {
		return d, fmt.Errorf("expected origin found: %q", s)
	}
	s, err = parseStringField((*string)(&d.SessionName), s[1:])
	if err != nil {
		return d, err
	}

	if len(s) != 0 && s[0] == 'i' {
		s, err = parseStringField((*string)(&d.SessionInformation), s[1:])
		if err != nil {
			return d, fmt.Errorf("cannot parse session information: %s", err)
		}
	}
	if len(s) != 0 && s[0] == 'u' {
		s, err = parseStringField((*string)(&d.URI), s[1:])
		if err != nil {
			return d, fmt.Errorf("cannot parse uri: %s", err)
		}
	}
	if len(s) != 0 && s[0] == 'e' {
		s, err = parseStringField((*string)(&d.EmailAddress), s[1:])
		if err != nil {
			return d, fmt.Errorf("cannot parse email address: %s", err)
		}
	}
	if len(s) != 0 && s[0] == 'p' {
		s, err = parseStringField((*string)(&d.PhoneNumber), s[1:])
		if err != nil {
			return d, fmt.Errorf("cannot parse phone number: %s", err)
		}
	}
	if len(s) != 0 && s[0] == 'c' {
		s, err = parseConnectionInformation(&d.ConnectionInformation, s[1:], c)
		if err != nil {
			return d, fmt.Errorf("cannot parse connection information: %s", err)
		}
	}

	for len(s) != 0 && s[0] == 'b' {
		if s, err = parseBandwidth(c.getBandwidth(), s[1:], c); err != nil {
			return d, fmt.Errorf("cannot parse bandwidth: %s", err)
		}
	}
	d.Bandwidth = c.cloneBandwidth()

	for len(s) != 0 && s[0] == 't' {
		t := c.getTimeDescription()
		if s, err = parseTiming(&t.Timing, s[1:]); err != nil {
			return d, fmt.Errorf("cannot parse time description: %s", err)
		}
		for len(s) != 0 && s[0] == 'r' {
			s, err = parseRepeatTime(c.getRepeatTime(), s[1:], c)
			if err != nil {
				return d, fmt.Errorf("cannot parse repeat times: %s", err)
			}
		}
		t.RepeatTimes = c.cloneRepeatTime()
	}
	d.TimeDescriptions = c.cloneTimeDescription()

	if len(s) != 0 && s[0] == 'z' {
		d.TimeZones, s, err = parseTimeZone(s[1:], c)
		if err != nil {
			return d, fmt.Errorf("cannot parse time zone adjustments: %s", err)
		}
	}

	if len(s) != 0 && s[0] == 'k' {
		s, err = parseStringField((*string)(&d.EncryptionKey), s[1:])
		if err != nil {
			return d, fmt.Errorf("cannot parse encryption key: %s", err)
		}
	}

	for len(s) != 0 && s[0] == 'a' {
		if s, err = parseAttribute(c.getAttribute(), s[1:], c); err != nil {
			return d, fmt.Errorf("cannot parse session attribute: %s", err)
		}
	}
	d.Attributes = c.cloneAttribute()

	for len(s) != 0 && s[0] == 'm' {
		m := c.getMediaDescription()
		if s, err = parseMediaName(&m.MediaName, s[1:], c); err != nil {
			return d, fmt.Errorf("cannot parse media description: %s", err)
		}

		for ok := true; ok && len(s) != 0 && s[0] != 'm'; {
			ok = false
			if len(s) != 0 && s[0] == 'i' {
				s, err = parseStringField((*string)(&m.MediaTitle), s[1:])
				if err != nil {
					return d, fmt.Errorf("cannot parse session information: %s", err)
				}
				ok = true
			}

			if len(s) != 0 && s[0] == 'c' {
				s, err = parseConnectionInformation(&m.ConnectionInformation, s[1:], c)
				if err != nil {
					return d, fmt.Errorf("cannot parse connection information: %s", err)
				}
				ok = true
			}

			for len(s) != 0 && s[0] == 'b' {
				if s, err = parseBandwidth(c.getBandwidth(), s[1:], c); err != nil {
					return d, fmt.Errorf("cannot parse bandwidth: %s", err)
				}
				ok = true
			}

			if len(s) != 0 && s[0] == 'k' {
				s, err = parseStringField((*string)(&m.EncryptionKey), s[1:])
				if err != nil {
					return d, fmt.Errorf("cannot parse encryption key: %s", err)
				}
				ok = true
			}

			for len(s) != 0 && s[0] == 'a' {
				if s, err = parseAttribute(c.getAttribute(), s[1:], c); err != nil {
					return d, fmt.Errorf("cannot parse attribute: %s", err)
				}
				ok = true
			}
		}
		m.Bandwidth = c.cloneBandwidth()
		m.Attributes = c.cloneAttribute()
	}
	d.MediaDescriptions = c.cloneMediaDescription()

	return d, nil
}

func parseVersion(f *Version, s string) (string, error) {
	s, err := skipEq(s)
	if err != nil {
		return s, err
	}
	n, s, err := parseUint8(s)
	if err != nil {
		return s, err
	}
	*f = Version(n)
	return skipNewLine(s), nil
}

func parseOrigin(f *Origin, s string) (string, error) {
	s, err := skipEq(s)
	if err != nil {
		return s, err
	}

	f.Username, s, err = parseNonWSString(s)
	if err != nil {
		return s, err
	}

	f.SessionID, s, err = parseUint64(s)
	if err != nil {
		return s, err
	}
	s = skipWS(s)

	f.SessionVersion, s, err = parseUint64(s)
	if err != nil {
		return s, err
	}
	s = skipWS(s)

	f.NetworkType, s, err = parseNonWSString(s)
	if err != nil {
		return s, err
	}

	f.AddressType, s, err = parseNonWSString(s)
	if err != nil {
		return s, err
	}

	f.UnicastAddress, s, err = parseNonWSString(s)
	if err != nil {
		return s, err
	}

	return skipNewLine(s), nil
}

func parseStringField(f *string, s string) (string, error) {
	s, err := skipEq(s)
	if err != nil {
		return s, err
	}

	*f, s, err = parseLine(s)
	if err != nil {
		return s, err
	}

	return skipNewLine(s), nil
}

func parseConnectionInformation(f *ConnectionInformation, s string, c *cache) (string, error) {
	s, err := skipEq(s)
	if err != nil {
		return s, err
	}

	f.NetworkType, s, err = parseNonWSString(s)
	if err != nil {
		return s, err
	}

	f.AddressType, s, err = parseNonWSString(s)
	if err != nil {
		return s, err
	}

	addr, s, err := parseDelimitedStringSlice(s, '/', c)
	if err != nil {
		return s, err
	}

	switch len(addr) {
	case 0:
		return s, fmt.Errorf("expected connection address found %q", s)
	case 1:
		f.Address.Address = addr[0]
	case 2:
		f.Address.Address = addr[0]
		f.Address.Range, _, err = parseUint64(addr[1])
		if err != nil {
			return s, fmt.Errorf("cannot parse connection address range: %s", err)
		}
	case 3:
		f.Address.Address = addr[0]
		f.Address.TTL, _, err = parseUint64(addr[1])
		if err != nil {
			return s, fmt.Errorf("cannot parse connection address ttl: %s", err)
		}
		f.Address.Range, _, err = parseUint64(addr[2])
		if err != nil {
			return s, fmt.Errorf("cannot parse connection address range: %s", err)
		}
	}

	return skipNewLine(s), nil
}

func parseBandwidth(f *Bandwidth, s string, c *cache) (string, error) {
	s, err := skipEq(s)
	if err != nil {
		return s, err
	}

	f.Type, s, err = parseDelimitedString(s, ':')
	if err != nil {
		return s, fmt.Errorf("cannot parse bandwidth: %s", err)
	}

	f.Bandwidth, s, err = parseUint64(s)
	if err != nil {
		return s, fmt.Errorf("cannot parse bandwidth: %s", err)
	}

	return skipNewLine(s), nil
}

func parseTiming(f *Timing, s string) (string, error) {
	s, err := skipEq(s)
	if err != nil {
		return s, err
	}

	f.StartTime, s, err = parseUint64(s)
	if err != nil {
		return s, fmt.Errorf("cannot parse timing: %s", err)
	}
	s = skipWS(s)

	f.StopTime, s, err = parseUint64(s)
	if err != nil {
		return s, fmt.Errorf("cannot parse timing: %s", err)
	}

	return skipNewLine(s), nil
}

func parseRepeatTime(f *RepeatTime, s string, c *cache) (string, error) {
	s, err := skipEq(s)
	if err != nil {
		return s, err
	}

	f.Interval, s, err = parseTime(s)
	if err != nil {
		return s, fmt.Errorf("cannot parse repeat time interval: %s", err)
	}

	f.Duration, s, err = parseTime(s)
	if err != nil {
		return s, fmt.Errorf("cannot parse repeat time duration: %s", err)
	}

	offsets, s, err := parseDelimitedStringSlice(s, ' ', c)
	if err != nil {
		return s, fmt.Errorf("cannot parse repeat time offsets: %s", err)
	}
	f.Offsets = make([]int64, len(offsets))
	for i := 0; i < len(offsets); i++ {
		f.Offsets[i], _, err = parseTime(offsets[i])
		if err != nil {
			return s, fmt.Errorf("cannot parse repeat time offsets: %s", err)
		}
	}

	return skipNewLine(s), nil
}

func parseTimeZone(s string, c *cache) ([]TimeZone, string, error) {
	s, err := skipEq(s)
	if err != nil {
		return nil, s, err
	}

	parts, s, err := parseDelimitedStringSlice(s, ' ', c)
	if err != nil {
		return nil, s, fmt.Errorf("cannot parse time zones: %s", err)
	}
	if len(parts)%2 != 0 {
		return nil, s, fmt.Errorf("cannot parse time zones: no adjustment time for offset")
	}

	tzs := make([]TimeZone, len(parts)/2)
	for i := 0; i < len(tzs); i++ {
		tzs[i].AdjustmentTime, _, err = parseUint64(parts[i*2])
		if err != nil {
			return nil, s, fmt.Errorf("cannot parse time zones: %s", err)
		}

		tzs[i].Offset, _, err = parseTime(parts[i*2+1])
		if err != nil {
			return nil, s, fmt.Errorf("cannot parse time zones: %s", err)
		}
	}

	return tzs, skipNewLine(s), nil
}

func parseAttribute(f *Attribute, s string, c *cache) (string, error) {
	s, err := skipEq(s)
	if err != nil {
		return s, err
	}

	for i := 0; i < len(s); i++ {
		if s[i] == ':' {
			f.Key = s[:i]
			f.Value, s, err = parseLine(s[i+1:])
			if err != nil {
				return s, fmt.Errorf("cannot parse attribute: %s", err)
			}
			return s, nil
		}
		if s[i] == '\r' && len(s) >= i+1 && s[i+1] == '\n' {
			f.Key = s[:i]
			return s[i+2:], nil
		}
		if s[i] == '\n' {
			f.Key = s[:i]
			return s[i+1:], nil
		}
	}

	return s, fmt.Errorf("cannot parse attribute: unexpected eof")
}

func parseMediaName(f *MediaName, s string, c *cache) (string, error) {
	s, err := skipEq(s)
	if err != nil {
		return s, err
	}

	f.Media, s, err = parseNonWSString(s)
	if err != nil {
		return s, fmt.Errorf("cannot parse media name: %s", err)
	}

	port, s, err := parseDelimitedStringSlice(s, '/', c)
	if err != nil {
		return s, fmt.Errorf("cannot parse media name port: %s", err)
	}

	if len(port) != 0 {
		f.Port.Value, _, err = parseUint16(port[0])
		if err != nil {
			return s, fmt.Errorf("cannot parse media name port: %s", err)
		}
	}
	if len(port) == 2 {
		f.Port.Range, _, err = parseUint16(port[1])
		if err != nil {
			return s, fmt.Errorf("cannot parse media name port: %s", err)
		}
	}

	protos, s, err := parseDelimitedStringSlice(s, '/', c)
	if err != nil {
		return s, fmt.Errorf("cannot parse media name proto: %s", err)
	}
	f.Protos = make([]string, len(protos))
	copy(f.Protos, protos)

	formats, s, err := parseDelimitedStringSlice(s, ' ', c)
	if err != nil {
		return s, fmt.Errorf("cannot parse media name formats: %s", err)
	}
	f.Formats = make([]string, len(formats))
	copy(f.Formats, formats)

	return skipNewLine(s), nil
}

func skipEq(s string) (string, error) {
	if len(s) == 0 || s[0] != '=' {
		return s, fmt.Errorf("expected '=' found %q", s)
	}
	return s[1:], nil
}

func skipWS(s string) string {
	if len(s) == 0 || s[0] != ' ' {
		return s
	}
	return skipWSSlow(s)
}

func skipWSSlow(s string) string {
	for i := 0; i < len(s); i++ {
		if s[i] != ' ' {
			return s[i:]
		}
	}
	return ""
}

func skipNewLine(s string) string {
	for i := 0; i < len(s); i++ {
		if s[i] != '\r' && s[i] != '\n' {
			return s[i:]
		}
	}
	return ""
}

func parseLine(s string) (string, string, error) {
	for i := 0; i < len(s); i++ {
		if s[i] == '\r' && len(s) >= i+1 && s[i+1] == '\n' {
			return s[:i], s[i+2:], nil
		}
		if s[i] == '\n' {
			return s[:i], s[i+1:], nil
		}
	}
	return "", s, fmt.Errorf("unexpected eof")
}

func parseNonWSString(s string) (string, string, error) {
	for i := 0; i < len(s); i++ {
		if s[i] == ' ' {
			return s[:i], s[i+1:], nil
		}
		if s[i] == '\r' && len(s) >= i+1 && s[i+1] == '\n' {
			return s[:i], s[i:], nil
		}
		if s[i] == '\n' {
			return s[:i], s[i:], nil
		}
	}
	return "", s, fmt.Errorf("unexpected eof")
}

func parseDelimitedString(s string, ch byte) (string, string, error) {
	for i := 0; i < len(s); i++ {
		if s[i] == ch {
			return s[:i], s[i+1:], nil
		}
	}
	return "", s, fmt.Errorf("expected %q found %q", ch, s)
}

func parseDelimitedStringSlice(s string, ch byte, c *cache) ([]string, string, error) {
	c.ss = c.ss[:0]
	for i := 0; i < len(s); i++ {
		if s[i] == ch {
			c.ss = append(c.ss, s[:i])
			s = s[i+1:]
			i = 0
		}
		if s[i] == ' ' {
			c.ss = append(c.ss, s[:i])
			return c.ss, s[i+1:], nil
		}
		if s[i] == '\r' && len(s) >= i+1 && s[i+1] == '\n' {
			c.ss = append(c.ss, s[:i])
			return c.ss, s[i+2:], nil
		}
		if s[i] == '\n' {
			c.ss = append(c.ss, s[:i])
			return c.ss, s[i+1:], nil
		}
	}
	return nil, s, fmt.Errorf("expected %q found %q", ch, s)
}

func parseTime(s string) (int64, string, error) {
	n, s, err := parseInt64(s)
	if err != nil {
		return 0, s, err
	}

	if len(s) == 0 {
		return n, s, nil
	}
	switch s[0] {
	case 'd':
		n *= 86400
		s = s[1:]
	case 'h':
		n *= 3600
		s = s[1:]
	case 'm':
		n *= 60
		s = s[1:]
	case 's':
		s = s[1:]
	}

	return n, skipWS(s), nil
}

func parseInt64(s string) (int64, string, error) {
	sign := int64(1)
	if len(s) != 0 && s[0] == '-' {
		sign = -1
		s = s[1:]
	}
	n, s, err := parseUintN(s, 63)
	if err != nil {
		return 0, s, err
	}
	return sign * int64(n), s, nil
}

func parseUint8(s string) (uint8, string, error) {
	n, s, err := parseUintN(s, 8)
	return uint8(n), s, err
}

func parseUint16(s string) (uint16, string, error) {
	n, s, err := parseUintN(s, 16)
	return uint16(n), s, err
}

func parseUint32(s string) (uint32, string, error) {
	n, s, err := parseUintN(s, 32)
	return uint32(n), s, err
}

func parseUintN(s string, bits int) (uint64, string, error) {
	n, s, err := parseUint64(s)
	if err != nil {
		return 0, s, err
	}
	if n > uint64(1<<bits)-1 {
		return 0, s, fmt.Errorf("value out of range: %d", n)
	}
	return n, s, nil
}

func parseUint64(s string) (uint64, string, error) {
	var n uint64
	var i int
	for ; i < len(s) && s[i] >= '0' && s[i] <= '9'; i++ {
		n = n*10 + uint64(s[i]-'0')
	}
	if i == 0 {
		return 0, s, fmt.Errorf("expected number found %q", s)
	}
	return n, s[i:], nil
}
