// SPDX-FileCopyrightText: 2023 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

package sdp

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/pion/randutil"
)

var (
	errExtractCodecRtpmap  = errors.New("could not extract codec from rtpmap")
	errExtractCodecFmtp    = errors.New("could not extract codec from fmtp")
	errExtractCodecRtcpFb  = errors.New("could not extract codec from rtcp-fb")
	errPayloadTypeNotFound = errors.New("payload type not found")
	errCodecNotFound       = errors.New("codec not found")
	errSyntaxError         = errors.New("SyntaxError")
)

// ConnectionRole indicates which of the end points should initiate the connection establishment
type ConnectionRole int

const (
	// ConnectionRoleActive indicates the endpoint will initiate an outgoing connection.
	ConnectionRoleActive ConnectionRole = iota + 1

	// ConnectionRolePassive indicates the endpoint will accept an incoming connection.
	ConnectionRolePassive

	// ConnectionRoleActpass indicates the endpoint is willing to accept an incoming connection or to initiate an outgoing connection.
	ConnectionRoleActpass

	// ConnectionRoleHoldconn indicates the endpoint does not want the connection to be established for the time being.
	ConnectionRoleHoldconn
)

func (t ConnectionRole) String() string {
	switch t {
	case ConnectionRoleActive:
		return "active"
	case ConnectionRolePassive:
		return "passive"
	case ConnectionRoleActpass:
		return "actpass"
	case ConnectionRoleHoldconn:
		return "holdconn"
	default:
		return "Unknown"
	}
}

func newSessionID() (uint64, error) {
	// https://tools.ietf.org/html/draft-ietf-rtcweb-jsep-26#section-5.2.1
	// Session ID is recommended to be constructed by generating a 64-bit
	// quantity with the highest bit set to zero and the remaining 63-bits
	// being cryptographically random.
	id, err := randutil.CryptoUint64()
	return id & (^(uint64(1) << 63)), err
}

// Codec represents a codec
type Codec struct {
	PayloadType        uint8
	Name               string
	ClockRate          uint32
	EncodingParameters string
	Fmtp               string
	RTCPFeedback       []string
}

const (
	unknown = iota
)

func (c Codec) ByteLen() int {
	n := uintLen(uint64(c.PayloadType))
	n += len(c.Name)
	n += uintLen(uint64(c.ClockRate))
	n += len(c.EncodingParameters)
	n += len(c.Fmtp)
	for i, f := range c.RTCPFeedback {
		if i > 0 {
			n += 2
		}
		n += len(f)
	}
	return n + 9
}

func (c Codec) MarshalAppend(b []byte) []byte {
	b = growByteSlice(b, c.ByteLen())
	b = strconv.AppendUint(b, uint64(c.PayloadType), 10)
	b = append(b, ' ')
	b = append(b, c.Name...)
	b = append(b, '/')
	b = strconv.AppendUint(b, uint64(c.ClockRate), 10)
	b = append(b, '/')
	b = append(b, c.EncodingParameters...)
	b = append(b, " ("...)
	b = append(b, c.Fmtp...)
	b = append(b, ") ["...)
	for i, f := range c.RTCPFeedback {
		if i > 0 {
			b = append(b, ", "...)
		}
		b = append(b, f...)
	}
	b = append(b, ']')
	return b
}

func (c Codec) String() string {
	return string(c.MarshalAppend(nil))
}

func parseRtpmap(rtpmap Attribute) (codec Codec, err error) {
	// <payload type> <encoding name>/<clock rate>[/<encoding parameters>]
	i := strings.IndexByte(rtpmap.Value, ' ')
	if i == -1 {
		return codec, errExtractCodecRtpmap
	}

	ptInt, _, err := parseUint8(rtpmap.Value[:i])
	if err != nil {
		return codec, fmt.Errorf("%w: %s", errExtractCodecRtpmap, err)
	}
	codec.PayloadType = uint8(ptInt)

	split := strings.Split(rtpmap.Value[i+1:], "/")
	codec.Name = split[0]
	parts := len(split)
	if parts > 1 {
		rate, _, err := parseUint32(split[1])
		if err != nil {
			return codec, fmt.Errorf("%w: %s", errExtractCodecRtpmap, err)
		}
		codec.ClockRate = uint32(rate)
	}
	if parts > 2 {
		codec.EncodingParameters = split[2]
	}

	return codec, nil
}

func parseFmtp(fmtp Attribute) (codec Codec, err error) {
	// <format> <format specific parameters>
	i := strings.IndexByte(fmtp.Value, ' ')
	if i == -1 {
		return codec, errExtractCodecFmtp
	}

	ptInt, _, err := parseUint8(fmtp.Value[i+1:])
	if err != nil {
		return codec, fmt.Errorf("%w: %s", errExtractCodecFmtp, err)
	}
	codec.PayloadType = uint8(ptInt)

	codec.Fmtp = fmtp.Value[:i]

	return codec, nil
}

func parseRtcpFb(rtcpFb Attribute) (codec Codec, err error) {
	// <payload type> <RTCP feedback type> [<RTCP feedback parameter>]
	i := strings.IndexByte(rtcpFb.Value, ' ')
	if i == -1 {
		return codec, errExtractCodecRtcpFb
	}

	ptInt, _, err := parseUint8(rtcpFb.Value[:i])
	if err != nil {
		return codec, fmt.Errorf("%w: %s", errExtractCodecRtcpFb, err)
	}

	codec.PayloadType = uint8(ptInt)
	codec.RTCPFeedback = append(codec.RTCPFeedback, rtcpFb.Value[i+1:])

	return codec, nil
}

func mergeCodecs(codec Codec, codecs map[uint8]Codec) {
	savedCodec := codecs[codec.PayloadType]

	if savedCodec.PayloadType == 0 {
		savedCodec.PayloadType = codec.PayloadType
	}
	if len(savedCodec.Name) == 0 {
		savedCodec.Name = codec.Name
	}
	if savedCodec.ClockRate == 0 {
		savedCodec.ClockRate = codec.ClockRate
	}
	if len(savedCodec.EncodingParameters) == 0 {
		savedCodec.EncodingParameters = codec.EncodingParameters
	}
	if len(savedCodec.Fmtp) == 0 {
		savedCodec.Fmtp = codec.Fmtp
	}
	savedCodec.RTCPFeedback = append(savedCodec.RTCPFeedback, codec.RTCPFeedback...)

	codecs[savedCodec.PayloadType] = savedCodec
}

func (s *SessionDescription) buildCodecMap() map[uint8]Codec {
	codecs := map[uint8]Codec{
		// static codecs that do not require a rtpmap
		0: {
			PayloadType: 0,
			Name:        "pcmu",
			ClockRate:   8000,
		},
		8: {
			PayloadType: 8,
			Name:        "pcma",
			ClockRate:   8000,
		},
	}

	for _, m := range s.MediaDescriptions {
		for _, a := range m.Attributes {
			switch a.Key {
			case "rtpmap":
				codec, err := parseRtpmap(a)
				if err == nil {
					mergeCodecs(codec, codecs)
				}
			case "fmtp":
				codec, err := parseFmtp(a)
				if err == nil {
					mergeCodecs(codec, codecs)
				}
			case "rtcp-fb":
				codec, err := parseRtcpFb(a)
				if err == nil {
					mergeCodecs(codec, codecs)
				}
			}
		}
	}

	return codecs
}

func equivalentFmtp(want, got string) bool {
	wantSplit := strings.Split(want, ";")
	gotSplit := strings.Split(got, ";")

	if len(wantSplit) != len(gotSplit) {
		return false
	}

	sort.Slice(wantSplit, func(i, j int) bool { return strings.Compare(wantSplit[i], wantSplit[j]) == -1 })
	sort.Slice(gotSplit, func(i, j int) bool { return strings.Compare(gotSplit[i], gotSplit[j]) == -1 })

	for i, wantPart := range wantSplit {
		wantPart = strings.TrimSpace(wantPart)
		gotPart := strings.TrimSpace(gotSplit[i])
		if gotPart != wantPart {
			return false
		}
	}

	return true
}

func codecsMatch(wanted, got Codec) bool {
	if wanted.Name != "" && !strings.EqualFold(wanted.Name, got.Name) {
		return false
	}
	if wanted.ClockRate != 0 && wanted.ClockRate != got.ClockRate {
		return false
	}
	if wanted.EncodingParameters != "" && wanted.EncodingParameters != got.EncodingParameters {
		return false
	}
	if wanted.Fmtp != "" && !equivalentFmtp(wanted.Fmtp, got.Fmtp) {
		return false
	}

	return true
}

// GetCodecForPayloadType scans the SessionDescription for the given payload type and returns the codec
func (s *SessionDescription) GetCodecForPayloadType(payloadType uint8) (Codec, error) {
	codecs := s.buildCodecMap()

	codec, ok := codecs[payloadType]
	if ok {
		return codec, nil
	}

	return codec, errPayloadTypeNotFound
}

// GetPayloadTypeForCodec scans the SessionDescription for a codec that matches the provided codec
// as closely as possible and returns its payload type
func (s *SessionDescription) GetPayloadTypeForCodec(wanted Codec) (uint8, error) {
	codecs := s.buildCodecMap()

	for payloadType, codec := range codecs {
		if codecsMatch(wanted, codec) {
			return payloadType, nil
		}
	}

	return 0, errCodecNotFound
}
