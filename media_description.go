// SPDX-FileCopyrightText: 2023 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

package sdp

import (
	"bytes"
	"strconv"
)

// MediaDescription represents a media type.
// https://tools.ietf.org/html/rfc4566#section-5.14
type MediaDescription struct {
	// m=<media> <port>/<number of ports> <proto> <fmt> ...
	// https://tools.ietf.org/html/rfc4566#section-5.14
	MediaName MediaName

	// i=<session description>
	// https://tools.ietf.org/html/rfc4566#section-5.4
	MediaTitle Information

	// c=<nettype> <addrtype> <connection-address>
	// https://tools.ietf.org/html/rfc4566#section-5.7
	ConnectionInformation ConnectionInformation

	// b=<bwtype>:<bandwidth>
	// https://tools.ietf.org/html/rfc4566#section-5.8
	Bandwidth []Bandwidth

	// k=<method>
	// k=<method>:<encryption key>
	// https://tools.ietf.org/html/rfc4566#section-5.12
	EncryptionKey EncryptionKey

	// a=<attribute>
	// a=<attribute>:<value>
	// https://tools.ietf.org/html/rfc4566#section-5.13
	Attributes []Attribute
}

// Attribute returns the value of an attribute and if it exists
func (d *MediaDescription) Attribute(key []byte) ([]byte, bool) {
	for _, a := range d.Attributes {
		if bytes.Equal(a.Key, key) {
			return a.Value, true
		}
	}
	return nil, false
}

// RangedPort supports special format for the media field "m=" port value. If
// it may be necessary to specify multiple transport ports, the protocol allows
// to write it as: <port>/<number of ports> where number of ports is a an
// offsetting range.
type RangedPort struct {
	Value uint16
	Range uint16
}

func (p RangedPort) Len() int {
	n := uintLen(uint64(p.Value))
	if p.Range != 0 {
		n += uintLen(uint64(p.Range)) + 1
	}
	return n
}

func (p RangedPort) AppendTo(b []byte) []byte {
	b = growByteSlice(b, p.Len())
	b = strconv.AppendUint(b, uint64(p.Value), 10)
	if p.Range != 0 {
		b = append(b, '/')
		b = strconv.AppendUint(b, uint64(p.Range), 10)
	}
	return b
}

// MediaName describes the "m=" field storage structure.
type MediaName struct {
	Media   []byte
	Port    RangedPort
	Protos  [][]byte
	Formats [][]byte
}

func (m MediaName) Len() int {
	n := len(m.Media) + m.Port.Len() + 1
	for i := range m.Protos {
		n += len(m.Protos[i]) + 1
	}
	for i := range m.Formats {
		n += len(m.Formats[i]) + 1
	}
	return n
}

func (m MediaName) AppendTo(b []byte) []byte {
	b = growByteSlice(b, m.Len())
	b = append(b, m.Media...)
	b = append(b, ' ')
	b = m.Port.AppendTo(b)
	for i := range m.Protos {
		if i == 0 {
			b = append(b, ' ')
		} else {
			b = append(b, '/')
		}
		b = append(b, m.Protos[i]...)
	}
	for i := range m.Formats {
		b = append(b, ' ')
		b = append(b, m.Formats[i]...)
	}
	return b
}
