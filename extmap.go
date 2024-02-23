// SPDX-FileCopyrightText: 2023 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

package sdp

import (
	"fmt"
	"strconv"
	"strings"
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
	Direction Direction
	URI       URI
	ExtAttr   string
}

// Clone converts this object to an Attribute
func (e ExtMap) Clone() Attribute {
	return Attribute{Key: "extmap", Value: string(e.AppendTo(nil))}
}

// Unmarshal creates an Extmap from a string
func (e *ExtMap) Unmarshal(raw string) error {
	parts := strings.SplitN(raw, ":", 2)
	if len(parts) != 2 {
		return fmt.Errorf("%w: %v", errSyntaxError, raw)
	}

	fields := strings.Fields(parts[1])
	if len(fields) < 2 {
		return fmt.Errorf("%w: %v", errSyntaxError, raw)
	}

	valdir := strings.Split(fields[0], "/")
	value, _, err := parseUint8(valdir[0])
	if err != nil || value == 0 {
		return fmt.Errorf("%w: %v", errSyntaxError, valdir[0])
	}

	var direction Direction
	if len(valdir) == 2 {
		direction, err = NewDirection(valdir[1])
		if err != nil {
			return err
		}
	}

	if len(fields) == 3 {
		e.ExtAttr = fields[2]
	}

	e.Value = int(value)
	e.Direction = direction
	e.URI = URI(fields[1])
	return nil
}

// Marshal creates a string from an ExtMap
func (e ExtMap) Marshal() string {
	b := make([]byte, 0, len("extmap")+1+e.Len())
	b = append(b, "extmap"...)
	b = append(b, ':')
	b = e.AppendTo(b)
	return string(b)
}

func (e ExtMap) Len() int {
	n := uintLen(uint64(e.Value))
	if e.Direction != unknown {
		n += len(e.Direction.String()) + 1
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
	if e.Direction != unknown {
		b = append(b, '/')
		b = append(b, e.Direction.String()...)
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

// Name returns the constant name of this object
func (e ExtMap) Name() string {
	return "extmap"
}
