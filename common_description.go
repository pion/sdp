// SPDX-FileCopyrightText: 2023 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

package sdp

import (
	"strconv"
	"strings"
)

// Information describes the "i=" field which provides textual information
// about the session.
type Information string

func (i Information) String() string {
	return string(i)
}

func (i Information) marshalSize() (size int) {
	return len(i)
}

// ConnectionInformation defines the representation for the "c=" field
// containing connection data.
type ConnectionInformation struct {
	NetworkType string
	AddressType string
	Address     *Address
}

func (c ConnectionInformation) String() string {
	parts := []string{c.NetworkType, c.AddressType}
	if c.Address != nil {
		parts = append(parts, c.Address.String())
	}
	return strings.Join(parts, " ")
}

func (c ConnectionInformation) marshalSize() (size int) {
	size = len(c.NetworkType)
	size += 1 + len(c.AddressType)
	if c.Address != nil {
		size += 1 + c.Address.marshalSize()
	}

	return
}

// Address desribes a structured address token from within the "c=" field.
type Address struct {
	Address string
	TTL     *int
	Range   *int
}

func (c *Address) String() string {
	var parts []string
	parts = append(parts, c.Address)
	if c.TTL != nil {
		parts = append(parts, strconv.Itoa(*c.TTL))
	}

	if c.Range != nil {
		parts = append(parts, strconv.Itoa(*c.Range))
	}

	return strings.Join(parts, "/")
}

func (c Address) marshalSize() (size int) {
	size = len(c.Address)
	if c.TTL != nil {
		size += 1 + lenUint(uint64(*c.TTL))
	}
	if c.Range != nil {
		size += 1 + lenUint(uint64(*c.Range))
	}

	return
}

// Bandwidth describes an optional field which denotes the proposed bandwidth
// to be used by the session or media.
type Bandwidth struct {
	Experimental bool
	Type         string
	Bandwidth    uint64
}

func (b Bandwidth) String() string {
	var output string
	if b.Experimental {
		output += "X-"
	}
	output += b.Type + ":" + strconv.FormatUint(b.Bandwidth, 10)
	return output
}

func (b Bandwidth) marshalSize() (size int) {
	if b.Experimental {
		size += 2
	}

	size += len(b.Type) + 1 + lenUint(b.Bandwidth)
	return
}

// EncryptionKey describes the "k=" which conveys encryption key information.
type EncryptionKey string

func (e EncryptionKey) String() string {
	return string(e)
}

func (e EncryptionKey) marshalSize() (size int) {
	return len(e.String())
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

func (a Attribute) String() string {
	output := a.Key
	if len(a.Value) > 0 {
		output += ":" + a.Value
	}
	return output
}

func (a Attribute) marshalSize() (size int) {
	size = len(a.Key)
	if len(a.Value) > 0 {
		size += 1 + len(a.Value)
	}

	return size
}

// IsICECandidate returns true if the attribute key equals "candidate".
func (a Attribute) IsICECandidate() bool {
	return a.Key == "candidate"
}
