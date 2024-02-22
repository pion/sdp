// SPDX-FileCopyrightText: 2023 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

package sdp

func (s *SessionDescription) Len() int {
	n := s.Version.Len() + 4
	n += s.Origin.Len() + 4
	n += s.SessionName.Len() + 4
	if s.SessionInformation.Defined() {
		n += s.SessionInformation.Len() + 4
	}
	if s.URI.Defined() {
		n += s.URI.Len() + 4
	}
	if s.EmailAddress.Defined() {
		n += s.EmailAddress.Len() + 4
	}
	if s.PhoneNumber.Defined() {
		n += s.PhoneNumber.Len() + 4
	}
	if s.ConnectionInformation.Defined() {
		n += s.ConnectionInformation.Len() + 4
	}
	for _, bw := range s.Bandwidth {
		n += bw.Len() + 4
	}
	for _, td := range s.TimeDescriptions {
		n += td.Timing.Len() + 4
		for _, r := range td.RepeatTimes {
			n += r.Len() + 4
		}
	}
	if s.TimeZones.Defined() {
		n += s.TimeZones.Len() + 4
	}
	if s.EncryptionKey.Defined() {
		n += s.EncryptionKey.Len() + 4
	}
	for _, a := range s.Attributes {
		n += a.Len() + 4
	}
	for _, md := range s.MediaDescriptions {
		n += md.MediaName.Len() + 4
		if md.MediaTitle.Defined() {
			n += md.MediaTitle.Len() + 4
		}
		if md.ConnectionInformation.Defined() {
			n += md.ConnectionInformation.Len() + 4
		}
		for _, bw := range md.Bandwidth {
			n += bw.Len() + 4
		}
		if md.EncryptionKey.Defined() {
			n += md.EncryptionKey.Len() + 4
		}
		for _, a := range md.Attributes {
			n += a.Len() + 4
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
	return s.AppendTo(nil), nil
}

func (s *SessionDescription) AppendTo(b []byte) []byte {
	b = growByteSlice(b, s.Len())
	b = appendAttribute(b, "v=", s.Version)
	b = appendAttribute(b, "o=", s.Origin)
	b = appendAttribute(b, "s=", s.SessionName)
	if s.SessionInformation.Defined() {
		b = appendAttribute(b, "i=", s.SessionInformation)
	}
	if s.URI.Defined() {
		b = appendAttribute(b, "u=", s.URI)
	}
	if s.EmailAddress.Defined() {
		b = appendAttribute(b, "e=", s.EmailAddress)
	}
	if s.PhoneNumber.Defined() {
		b = appendAttribute(b, "p=", s.PhoneNumber)
	}
	if s.ConnectionInformation.Defined() {
		b = appendAttribute(b, "c=", s.ConnectionInformation)
	}
	for _, bw := range s.Bandwidth {
		b = appendAttribute(b, "b=", bw)
	}
	for _, td := range s.TimeDescriptions {
		b = appendAttribute(b, "t=", td.Timing)
		for _, r := range td.RepeatTimes {
			b = appendAttribute(b, "r=", r)
		}
	}
	if s.TimeZones.Defined() {
		b = appendAttribute(b, "z=", s.TimeZones)
	}
	if s.EncryptionKey.Defined() {
		b = appendAttribute(b, "k=", s.EncryptionKey)
	}
	for _, a := range s.Attributes {
		b = appendAttribute(b, "a=", a)
	}
	for _, md := range s.MediaDescriptions {
		b = appendAttribute(b, "m=", md.MediaName)
		if md.MediaTitle.Defined() {
			b = appendAttribute(b, "i=", md.MediaTitle)
		}
		if md.ConnectionInformation.Defined() {
			b = appendAttribute(b, "c=", md.ConnectionInformation)
		}
		for _, bw := range md.Bandwidth {
			b = appendAttribute(b, "b=", bw)
		}
		if md.EncryptionKey.Defined() {
			b = appendAttribute(b, "k=", md.EncryptionKey)
		}
		for _, a := range md.Attributes {
			b = appendAttribute(b, "a=", a)
		}
	}
	return b
}

func appendAttribute(b []byte, name string, a interface{ AppendTo([]byte) []byte }) []byte {
	b = append(b, name...)
	b = a.AppendTo(b)
	b = append(b, "\r\n"...)
	return b
}
