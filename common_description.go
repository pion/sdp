// SPDX-FileCopyrightText: 2023 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

package sdp

import (
	"bytes"
	"strconv"
)

// Information describes the "i=" field which provides textual information
// about the session.
type Information []byte

func (t Information) Defined() bool {
	return len(t) != 0
}

func (t Information) Len() int {
	return len(t)
}

func (t Information) AppendTo(b []byte) []byte {
	return append(b, t...)
}

// URI describes the "u=" field which provides the uri.
type URI []byte

func (t URI) Defined() bool {
	return len(t) != 0
}

func (t URI) Len() int {
	return len(t)
}

func (t URI) AppendTo(b []byte) []byte {
	return append(b, t...)
}

// ConnectionInformation defines the representation for the "c=" field
// containing connection data.
type ConnectionInformation struct {
	NetworkType []byte
	AddressType []byte
	Address     Address
}

func (t ConnectionInformation) Defined() bool {
	return len(t.NetworkType) > 0
}

func (t ConnectionInformation) Len() int {
	n := t.Address.Len()
	if n > 0 {
		n++
	}
	n += len(t.NetworkType) + len(t.AddressType) + 1
	return n
}

func (t ConnectionInformation) AppendTo(b []byte) []byte {
	b = growByteSlice(b, t.Len())
	b = append(b, t.NetworkType...)
	b = append(b, ' ')
	b = append(b, t.AddressType...)
	if t.Address.Len() != 0 {
		b = append(b, ' ')
		b = t.Address.AppendTo(b)
	}
	return b
}

// Address desribes a structured address token from within the "c=" field.
type Address struct {
	Address []byte
	TTL     uint64
	Range   uint64
}

func (t Address) Len() int {
	n := len(t.Address)
	if t.TTL != 0 {
		n += uintLen(t.TTL) + 1
	}
	if t.Range != 0 {
		n += uintLen(t.Range) + 1
	}
	return n
}

func (t Address) AppendTo(b []byte) []byte {
	b = growByteSlice(b, t.Len())
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
	Experimental bool
	Type         []byte
	Bandwidth    uint64
}

func (t Bandwidth) Len() int {
	n := len(t.Type) + uintLen(t.Bandwidth) + 1
	if t.Experimental {
		n += 2
	}
	return n
}

func (t Bandwidth) AppendTo(b []byte) []byte {
	b = growByteSlice(b, t.Len())
	if t.Experimental {
		b = append(b, "X-"...)
	}
	b = append(b, t.Type...)
	b = append(b, ':')
	b = strconv.AppendUint(b, t.Bandwidth, 10)
	return b
}

// EncryptionKey describes the "k=" which conveys encryption key information.
type EncryptionKey []byte

func (t EncryptionKey) Defined() bool {
	return len(t) != 0
}

func (t EncryptionKey) Len() int {
	return len(t)
}

func (t EncryptionKey) AppendTo(b []byte) []byte {
	return append(b, t...)
}

// Attribute describes the "a=" field which represents the primary means for
// extending SDP.
type Attribute struct {
	Key   []byte
	Value []byte
}

// NewPropertyAttribute constructs a new attribute
func NewPropertyAttribute(key []byte) Attribute {
	return Attribute{
		Key: key,
	}
}

// NewAttribute constructs a new attribute
func NewAttribute(key, value []byte) Attribute {
	return Attribute{
		Key:   key,
		Value: value,
	}
}

func (t Attribute) Len() int {
	n := len(t.Key)
	if t.Value != nil {
		n += len(t.Value) + 1
	}
	return n
}

func (t Attribute) AppendTo(b []byte) []byte {
	b = growByteSlice(b, t.Len())
	b = append(b, t.Key...)
	if len(t.Value) > 0 {
		b = append(b, ':')
		b = append(b, t.Value...)
	}
	return b
}

// IsICECandidate returns true if the attribute key equals "candidate".
func (a Attribute) IsICECandidate() bool {
	return bytes.Equal(a.Key, []byte("candidate"))
}
