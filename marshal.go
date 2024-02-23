// SPDX-FileCopyrightText: 2023 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

package sdp

func (s *SessionDescription) ByteLen() int {
	n := s.Version.ByteLen() + 4
	n += s.Origin.ByteLen() + 4
	n += s.SessionName.ByteLen() + 4
	if s.SessionInformation.Defined() {
		n += s.SessionInformation.ByteLen() + 4
	}
	if s.URI.Defined() {
		n += s.URI.ByteLen() + 4
	}
	if s.EmailAddress.Defined() {
		n += s.EmailAddress.ByteLen() + 4
	}
	if s.PhoneNumber.Defined() {
		n += s.PhoneNumber.ByteLen() + 4
	}
	if s.ConnectionInformation.Defined() {
		n += s.ConnectionInformation.ByteLen() + 4
	}
	for _, bw := range s.Bandwidth {
		n += bw.ByteLen() + 4
	}
	for _, td := range s.TimeDescriptions {
		n += td.Timing.ByteLen() + 4
		for _, r := range td.RepeatTimes {
			n += r.ByteLen() + 4
		}
	}
	if s.TimeZones.Defined() {
		n += s.TimeZones.ByteLen() + 4
	}
	if s.EncryptionKey.Defined() {
		n += s.EncryptionKey.ByteLen() + 4
	}
	for _, a := range s.Attributes {
		n += a.ByteLen() + 4
	}
	for _, md := range s.MediaDescriptions {
		n += md.MediaName.ByteLen() + 4
		if md.MediaTitle.Defined() {
			n += md.MediaTitle.ByteLen() + 4
		}
		if md.ConnectionInformation.Defined() {
			n += md.ConnectionInformation.ByteLen() + 4
		}
		for _, bw := range md.Bandwidth {
			n += bw.ByteLen() + 4
		}
		if md.EncryptionKey.Defined() {
			n += md.EncryptionKey.ByteLen() + 4
		}
		for _, a := range md.Attributes {
			n += a.ByteLen() + 4
		}
	}
	return n
}

// Marshal takes a SDP struct to text
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
func (s *SessionDescription) Marshal() ([]byte, error) {
	return s.MarshalAppend(nil), nil
}

func (s *SessionDescription) MarshalAppend(b []byte) []byte {
	b = growByteSlice(b, s.ByteLen())
	b = append(b, "v="...)
	b = s.Version.MarshalAppend(b)
	b = append(b, "\r\n"...)
	b = append(b, "o="...)
	b = s.Origin.MarshalAppend(b)
	b = append(b, "\r\n"...)
	b = append(b, "s="...)
	b = s.SessionName.MarshalAppend(b)
	b = append(b, "\r\n"...)
	if s.SessionInformation.Defined() {
		b = append(b, "i="...)
		b = s.SessionInformation.MarshalAppend(b)
		b = append(b, "\r\n"...)
	}
	if s.URI.Defined() {
		b = append(b, "u="...)
		b = s.URI.MarshalAppend(b)
		b = append(b, "\r\n"...)
	}
	if s.EmailAddress.Defined() {
		b = append(b, "e="...)
		b = s.EmailAddress.MarshalAppend(b)
		b = append(b, "\r\n"...)
	}
	if s.PhoneNumber.Defined() {
		b = append(b, "p="...)
		b = s.PhoneNumber.MarshalAppend(b)
		b = append(b, "\r\n"...)
	}
	if s.ConnectionInformation.Defined() {
		b = append(b, "c="...)
		b = s.ConnectionInformation.MarshalAppend(b)
		b = append(b, "\r\n"...)
	}
	for _, bw := range s.Bandwidth {
		b = append(b, "b="...)
		b = bw.MarshalAppend(b)
		b = append(b, "\r\n"...)
	}
	for _, td := range s.TimeDescriptions {
		b = append(b, "t="...)
		b = td.Timing.MarshalAppend(b)
		b = append(b, "\r\n"...)
		for _, r := range td.RepeatTimes {
			b = append(b, "r="...)
			b = r.MarshalAppend(b)
			b = append(b, "\r\n"...)
		}
	}
	if s.TimeZones.Defined() {
		b = append(b, "z="...)
		b = s.TimeZones.MarshalAppend(b)
		b = append(b, "\r\n"...)
	}
	if s.EncryptionKey.Defined() {
		b = append(b, "k="...)
		b = s.EncryptionKey.MarshalAppend(b)
		b = append(b, "\r\n"...)
	}
	for _, a := range s.Attributes {
		b = append(b, "a="...)
		b = a.MarshalAppend(b)
		b = append(b, "\r\n"...)
	}
	for _, md := range s.MediaDescriptions {
		b = append(b, "m="...)
		b = md.MediaName.MarshalAppend(b)
		b = append(b, "\r\n"...)
		if md.MediaTitle.Defined() {
			b = append(b, "i="...)
			b = md.MediaTitle.MarshalAppend(b)
			b = append(b, "\r\n"...)
		}
		if md.ConnectionInformation.Defined() {
			b = append(b, "c="...)
			b = md.ConnectionInformation.MarshalAppend(b)
			b = append(b, "\r\n"...)
		}
		for _, bw := range md.Bandwidth {
			b = append(b, "b="...)
			b = bw.MarshalAppend(b)
			b = append(b, "\r\n"...)
		}
		if md.EncryptionKey.Defined() {
			b = append(b, "k="...)
			b = md.EncryptionKey.MarshalAppend(b)
			b = append(b, "\r\n"...)
		}
		for _, a := range md.Attributes {
			b = append(b, "a="...)
			b = a.MarshalAppend(b)
			b = append(b, "\r\n"...)
		}
	}
	return b
}

func intLen(n int64) int {
	if n < 0 {
		return uintLen(uint64(-n)) + 1
	}
	return uintLen(uint64(n))
}

func uintLen(n uint64) int {
	if n == 0 {
		return 1
	}
	return log10(n)
}

func log10(n uint64) int {
	switch {
	case n == 0:
		return 0
	case n < 1e1:
		return 1
	case n < 1e2:
		return 2
	case n < 1e3:
		return 3
	case n < 1e4:
		return 4
	case n < 1e5:
		return 5
	case n < 1e6:
		return 6
	case n < 1e7:
		return 7
	case n < 1e8:
		return 8
	case n < 1e9:
		return 9
	case n < 1e10:
		return 10
	case n < 1e11:
		return 11
	case n < 1e12:
		return 12
	case n < 1e13:
		return 13
	case n < 1e14:
		return 14
	case n < 1e15:
		return 15
	case n < 1e16:
		return 16
	case n < 1e17:
		return 17
	case n < 1e18:
		return 18
	case n < 1e19:
		return 19
	default:
		return 20
	}
}

// increase capacity of byte slice to accommodate at least n bytes
func growByteSlice(b []byte, n int) []byte {
	if cap(b)-len(b) >= n {
		return b
	}
	bc := make([]byte, len(b), len(b)+n)
	copy(bc, b)
	return bc
}
