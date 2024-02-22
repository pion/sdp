// SPDX-FileCopyrightText: 2023 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

package sdp

import (
	"bytes"
	"fmt"
	"strconv"
)

// Default ext values
const (
	DefExtMapValueABSSendTime     = 1
	DefExtMapValueTransportCC     = 2
	DefExtMapValueSDESMid         = 3
	DefExtMapValueSDESRTPStreamID = 4

	ABSSendTimeURI     = "http://www.webrtc.org/experiments/rtp-hdrext/abs-send-time"
	TransportCCURI     = "http://www.ietf.org/id/draft-holmer-rmcat-transport-wide-cc-extensions-01"
	SDESMidURI         = "urn:ietf:params:rtp-hdrext:sdes:mid"
	SDESRTPStreamIDURI = "urn:ietf:params:rtp-hdrext:sdes:rtp-stream-id"
	AudioLevelURI      = "urn:ietf:params:rtp-hdrext:ssrc-audio-level"
)

// ExtMap represents the activation of a single RTP header extension
type ExtMap struct {
	Value     int
	Direction []byte
	URI       []byte
	ExtAttr   []byte
}

// Clone converts this object to an Attribute
func (e ExtMap) Clone() Attribute {
	return Attribute{Key: kExtmap, Value: e.AppendTo(nil)}
}

// Unmarshal creates an Extmap from a string
func (e ExtMap) Unmarshal(raw []byte) error {
	parts := bytes.SplitN(raw, kColon, 2)
	if len(parts) != 2 {
		return fmt.Errorf("%w: %v", errSyntaxError, raw)
	}

	fields := bytes.Fields(parts[1])
	if len(fields) < 2 {
		return fmt.Errorf("%w: %v", errSyntaxError, raw)
	}

	valdir := bytes.Split(fields[0], kSlash)
	value, ok := parseUint(valdir[0], 8)
	if !ok || value == 0 {
		return fmt.Errorf("%w: %v", errSyntaxError, valdir[0])
	}

	if !anyOf(valdir[1], kSendRecv, kSendOnly, kRecvOnly, kInactive) {
		return fmt.Errorf("%w: %v", errDirectionString, valdir[1])
	}

	if len(fields) == 3 {
		e.ExtAttr = fields[2]
	}

	e.Value = int(value)
	e.Direction = valdir[1]
	e.URI = fields[1]
	return nil
}

// Marshal creates a string from an ExtMap
func (e ExtMap) Marshal() []byte {
	b := make([]byte, 0, len(kExtmap)+1+e.Len())
	b = append(b, kExtmap...)
	b = append(b, ':')
	b = e.AppendTo(b)
	return b
}

func (e ExtMap) Len() int {
	n := uintLen(uint64(e.Value))
	if len(e.Direction) != 0 {
		n += len(e.Direction) + 1
	}
	if len(e.URI) != 0 {
		n += len(e.URI) + 1
	}
	if len(e.ExtAttr) != 0 {
		n += len(e.ExtAttr) + 1
	}
	return n
}

func (e ExtMap) AppendTo(b []byte) []byte {
	b = strconv.AppendUint(b, uint64(e.Value), 10)
	if len(e.Direction) != 0 {
		b = append(b, '/')
		b = append(b, e.Direction...)
	}
	if len(e.URI) != 0 {
		b = append(b, ' ')
		b = append(b, e.URI...)
	}
	if len(e.ExtAttr) != 0 {
		b = append(b, ' ')
		b = append(b, e.ExtAttr...)
	}
	return b
}
