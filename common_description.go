// SPDX-FileCopyrightText: 2023 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

package sdp

import (
	"net/url"
	"strconv"
)

// Information describes the "i=" field which provides textual information
// about the session.
type Information string

func (t Information) Defined() bool {
	return len(t) != 0
}

func (t Information) ByteLen() int {
	return len(t)
}

func (t Information) MarshalAppend(b []byte) []byte {
	return append(b, t...)
}

// URI describes the "u=" field which provides the uri.
type URI string

func (t URI) Defined() bool {
	return len(t) != 0
}

func (t URI) ByteLen() int {
	return len(t)
}

func (t URI) MarshalAppend(b []byte) []byte {
	return append(b, t...)
}

func (t URI) URL() (*url.URL, error) {
	return url.Parse(string(t))
}

// ConnectionInformation defines the representation for the "c=" field
// containing connection data.
type ConnectionInformation struct {
	NetworkType string
	AddressType string
	Address     Address
}

func (t ConnectionInformation) Defined() bool {
	return len(t.NetworkType) > 0
}

func (t ConnectionInformation) ByteLen() int {
	n := len(t.NetworkType) + len(t.AddressType) + 1
	if t.Address.Defined() {
		n += t.Address.ByteLen() + 1
	}
	return n
}

func (t ConnectionInformation) MarshalAppend(b []byte) []byte {
	b = growByteSlice(b, t.ByteLen())
	b = append(b, t.NetworkType...)
	b = append(b, ' ')
	b = append(b, t.AddressType...)
	if t.Address.Defined() {
		b = append(b, ' ')
		b = t.Address.MarshalAppend(b)
	}
	return b
}

// Address desribes a structured address token from within the "c=" field.
type Address struct {
	Address string
	TTL     uint64
	Range   uint64
}

func (t Address) Defined() bool {
	return len(t.Address) != 0
}

func (t Address) ByteLen() int {
	n := len(t.Address)
	if t.TTL != 0 {
		n += uintLen(t.TTL) + 1
	}
	if t.Range != 0 {
		n += uintLen(t.Range) + 1
	}
	return n
}

func (t Address) MarshalAppend(b []byte) []byte {
	b = growByteSlice(b, t.ByteLen())
	b = append(b, t.Address...)
	if t.TTL != 0 {
		b = append(b, '/')
		b = strconv.AppendUint(b, t.TTL, 10)
	}
	if t.Range != 0 {
		b = append(b, '/')
		b = strconv.AppendUint(b, t.Range, 10)
	}
	return b
}

// Bandwidth describes an optional field which denotes the proposed bandwidth
// to be used by the session or media.
type Bandwidth struct {
	Type      string
	Bandwidth uint64
}

func (t Bandwidth) ByteLen() int {
	return len(t.Type) + uintLen(t.Bandwidth) + 1
}

func (t Bandwidth) MarshalAppend(b []byte) []byte {
	b = growByteSlice(b, t.ByteLen())
	b = append(b, t.Type...)
	b = append(b, ':')
	b = strconv.AppendUint(b, t.Bandwidth, 10)
	return b
}

// EncryptionKey describes the "k=" which conveys encryption key information.
type EncryptionKey string

func (t EncryptionKey) Defined() bool {
	return len(t) != 0
}

func (t EncryptionKey) ByteLen() int {
	return len(t)
}

func (t EncryptionKey) MarshalAppend(b []byte) []byte {
	return append(b, t...)
}

// Attribute describes the "a=" field which represents the primary means for
// extending SDP.
type Attribute struct {
	Key   string
	Value string
}

// NewPropertyAttribute constructs a new attribute
func NewPropertyAttribute(key string) Attribute {
	return Attribute{
		Key: key,
	}
}

// NewAttribute constructs a new attribute
func NewAttribute(key, value string) Attribute {
	return Attribute{
		Key:   key,
		Value: value,
	}
}

func (t Attribute) ByteLen() int {
	n := len(t.Key)
	if t.Value != "" {
		n += len(t.Value) + 1
	}
	return n
}

func (t Attribute) MarshalAppend(b []byte) []byte {
	b = growByteSlice(b, t.ByteLen())
	b = append(b, t.Key...)
	if len(t.Value) > 0 {
		b = append(b, ':')
		b = append(b, t.Value...)
	}
	return b
}

// IsICECandidate returns true if the attribute key equals "candidate".
func (a Attribute) IsICECandidate() bool {
	return a.Key == "candidate"
}
