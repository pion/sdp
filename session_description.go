// SPDX-FileCopyrightText: 2023 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

package sdp

import (
	"bytes"
	"fmt"
	"strconv"
)

// SessionDescription is a a well-defined format for conveying sufficient
// information to discover and participate in a multimedia session.
type SessionDescription struct {
	// v=0
	// https://tools.ietf.org/html/rfc4566#section-5.1
	Version Version

	// o=<username> <sess-id> <sess-version> <nettype> <addrtype> <unicast-address>
	// https://tools.ietf.org/html/rfc4566#section-5.2
	Origin Origin

	// s=<session name>
	// https://tools.ietf.org/html/rfc4566#section-5.3
	SessionName SessionName

	// i=<session description>
	// https://tools.ietf.org/html/rfc4566#section-5.4
	SessionInformation Information

	// u=<uri>
	// https://tools.ietf.org/html/rfc4566#section-5.5
	URI URI

	// e=<email-address>
	// https://tools.ietf.org/html/rfc4566#section-5.6
	EmailAddress EmailAddress

	// p=<phone-number>
	// https://tools.ietf.org/html/rfc4566#section-5.6
	PhoneNumber PhoneNumber

	// c=<nettype> <addrtype> <connection-address>
	// https://tools.ietf.org/html/rfc4566#section-5.7
	ConnectionInformation ConnectionInformation

	// b=<bwtype>:<bandwidth>
	// https://tools.ietf.org/html/rfc4566#section-5.8
	Bandwidth []Bandwidth

	// https://tools.ietf.org/html/rfc4566#section-5.9
	// https://tools.ietf.org/html/rfc4566#section-5.10
	TimeDescriptions []TimeDescription

	// z=<adjustment time> <offset> <adjustment time> <offset> ...
	// https://tools.ietf.org/html/rfc4566#section-5.11
	TimeZones TimeZoneSet

	// k=<method>
	// k=<method>:<encryption key>
	// https://tools.ietf.org/html/rfc4566#section-5.12
	EncryptionKey EncryptionKey

	// a=<attribute>
	// a=<attribute>:<value>
	// https://tools.ietf.org/html/rfc4566#section-5.13
	Attributes []Attribute

	// https://tools.ietf.org/html/rfc4566#section-5.14
	MediaDescriptions []MediaDescription
}

// Attribute returns the value of an attribute and if it exists
func (s *SessionDescription) Attribute(key []byte) ([]byte, bool) {
	for _, a := range s.Attributes {
		if bytes.Equal(a.Key, key) {
			return a.Value, true
		}
	}
	return nil, false
}

// Version describes the value provided by the "v=" field which gives
// the version of the Session Description Protocol.
type Version int

func (v Version) Len() int {
	return uintLen(uint64(v))
}

func (v Version) AppendTo(b []byte) []byte {
	return strconv.AppendUint(b, uint64(v), 10)
}

func (v Version) String() string {
	return strconv.Itoa(int(v))
}

// Origin defines the structure for the "o=" field which provides the
// originator of the session plus a session identifier and version number.
type Origin struct {
	Username       []byte
	SessionID      uint64
	SessionVersion uint64
	NetworkType    []byte
	AddressType    []byte
	UnicastAddress []byte
}

func (o Origin) Len() int {
	n := len(o.Username)
	n += uintLen(o.SessionID) + 1
	n += uintLen(o.SessionVersion) + 1
	n += len(o.NetworkType) + 1
	n += len(o.AddressType) + 1
	n += len(o.UnicastAddress) + 1
	return n
}

func (o Origin) AppendTo(b []byte) []byte {
	b = growByteSlice(b, o.Len())
	b = append(b, o.Username...)
	b = append(b, ' ')
	b = strconv.AppendUint(b, uint64(o.SessionID), 10)
	b = append(b, ' ')
	b = strconv.AppendUint(b, uint64(o.SessionVersion), 10)
	b = append(b, ' ')
	b = append(b, o.NetworkType...)
	b = append(b, ' ')
	b = append(b, o.AddressType...)
	b = append(b, ' ')
	b = append(b, o.UnicastAddress...)
	return b
}

func (o Origin) String() string {
	return fmt.Sprintf(
		"%v %d %d %v %v %v",
		o.Username,
		o.SessionID,
		o.SessionVersion,
		o.NetworkType,
		o.AddressType,
		o.UnicastAddress,
	)
}

// SessionName describes a structured representations for the "s=" field
// and is the textual session name.
type SessionName []byte

func (s SessionName) Defined() bool {
	return len(s) != 0
}

func (s SessionName) Len() int {
	return len(s)
}

func (s SessionName) AppendTo(b []byte) []byte {
	return append(b, s...)
}

func (s SessionName) String() string {
	return string(s)
}

// EmailAddress describes a structured representations for the "e=" line
// which specifies email contact information for the person responsible for
// the conference.
type EmailAddress []byte

func (e EmailAddress) Defined() bool {
	return len(e) != 0
}

func (e EmailAddress) Len() int {
	return len(e)
}

func (e EmailAddress) AppendTo(b []byte) []byte {
	return append(b, e...)
}

func (e EmailAddress) String() string {
	return string(e)
}

// PhoneNumber describes a structured representations for the "p=" line
// specify phone contact information for the person responsible for the
// conference.
type PhoneNumber []byte

func (p PhoneNumber) Defined() bool {
	return len(p) != 0
}

func (p PhoneNumber) Len() int {
	return len(p)
}

func (p PhoneNumber) AppendTo(b []byte) []byte {
	return append(b, p...)
}

func (p PhoneNumber) String() string {
	return string(p)
}

type TimeZoneSet []TimeZone

func (s TimeZoneSet) Defined() bool {
	return len(s) != 0
}

func (s TimeZoneSet) Len() (n int) {
	for i, z := range s {
		if i > 0 {
			n++
		}
		n += z.Len()
	}
	return n
}

func (s TimeZoneSet) AppendTo(b []byte) []byte {
	b = growByteSlice(b, s.Len())
	for i, z := range s {
		if i > 0 {
			b = append(b, ' ')
		}
		b = z.AppendTo(b)
	}
	return b
}

// TimeZone defines the structured object for "z=" line which describes
// repeated sessions scheduling.
type TimeZone struct {
	AdjustmentTime uint64
	Offset         int64
}

func (z TimeZone) Len() int {
	n := uintLen(z.AdjustmentTime) + 1
	if z.Offset < 0 {
		n += uintLen(uint64(-z.Offset)) + 1
	} else {
		n += uintLen(uint64(z.Offset))
	}
	return n
}

func (z TimeZone) AppendTo(b []byte) []byte {
	b = growByteSlice(b, z.Len())
	b = strconv.AppendUint(b, z.AdjustmentTime, 10)
	b = append(b, ' ')
	b = strconv.AppendInt(b, z.Offset, 10)
	return b
}

func (z TimeZone) String() string {
	return strconv.FormatUint(z.AdjustmentTime, 10) + " " + strconv.FormatInt(z.Offset, 10)
}
