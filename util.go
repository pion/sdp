// SPDX-FileCopyrightText: 2023 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

package sdp

import (
	"bytes"
	"errors"
	"io"
	"sort"
	"strconv"

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
	Name               []byte
	ClockRate          uint32
	EncodingParameters []byte
	Fmtp               []byte
	RTCPFeedback       [][]byte
}

const (
	unknown = iota
)

func (c Codec) Len() int {
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

func (c Codec) AppendTo(b []byte) []byte {
	b = growByteSlice(b, c.Len())
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

// func (c Codec) String() string {
// 	return fmt.Sprintf("%d %s/%d/%s (%s) [%s]", c.PayloadType, c.Name, c.ClockRate, c.EncodingParameters, c.Fmtp, bytes.Join(c.RTCPFeedback, ", "))
// }

func parseRtpmap(rtpmap Attribute) (codec Codec, err error) {
	// <payload type> <encoding name>/<clock rate>[/<encoding parameters>]
	i := bytes.IndexRune(rtpmap.Value, ' ')
	if i == -1 {
		return codec, errExtractCodecRtpmap
	}

	ptInt, ok := parseUint(rtpmap.Value[:i], 8)
	if !ok {
		return codec, errExtractCodecRtpmap
	}
	codec.PayloadType = uint8(ptInt)

	split := bytes.Split(rtpmap.Value[i+1:], kSlash)
	codec.Name = split[0]
	parts := len(split)
	if parts > 1 {
		rate, ok := parseUint(split[1], 32)
		if !ok {
			return codec, errExtractCodecRtpmap
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
	i := bytes.IndexRune(fmtp.Value, ' ')
	if i == -1 {
		return codec, errExtractCodecFmtp
	}

	ptInt, ok := parseUint(fmtp.Value[i+1:], 8)
	if !ok {
		return codec, errExtractCodecFmtp
	}
	codec.PayloadType = uint8(ptInt)

	codec.Fmtp = fmtp.Value[:i]

	return codec, nil
}

func parseRtcpFb(rtcpFb Attribute) (codec Codec, err error) {
	// <payload type> <RTCP feedback type> [<RTCP feedback parameter>]
	i := bytes.IndexRune(rtcpFb.Value, ' ')
	if i == -1 {
		return codec, errExtractCodecRtcpFb
	}

	ptInt, ok := parseUint(rtcpFb.Value[:i], 8)
	if !ok {
		return codec, errExtractCodecRtcpFb
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
			Name:        kPcmu,
			ClockRate:   8000,
		},
		8: {
			PayloadType: 8,
			Name:        kPcma,
			ClockRate:   8000,
		},
	}

	for _, m := range s.MediaDescriptions {
		for _, a := range m.Attributes {
			if bytes.Equal(a.Key, kRtpmap) {
				codec, err := parseRtpmap(a)
				if err == nil {
					mergeCodecs(codec, codecs)
				}
			} else if bytes.Equal(a.Key, kFmtp) {
				codec, err := parseFmtp(a)
				if err == nil {
					mergeCodecs(codec, codecs)
				}
			} else if bytes.Equal(a.Key, kRtcpFb) {
				codec, err := parseRtcpFb(a)
				if err == nil {
					mergeCodecs(codec, codecs)
				}
			}
		}
	}

	return codecs
}

func equivalentFmtp(want, got []byte) bool {
	wantSplit := bytes.Split(want, kSemicolon)
	gotSplit := bytes.Split(got, kSemicolon)

	if len(wantSplit) != len(gotSplit) {
		return false
	}

	sort.Slice(wantSplit, func(i, j int) bool { return bytes.Compare(wantSplit[i], wantSplit[j]) == -1 })
	sort.Slice(gotSplit, func(i, j int) bool { return bytes.Compare(gotSplit[i], gotSplit[j]) == -1 })

	for i, wantPart := range wantSplit {
		wantPart = bytes.TrimSpace(wantPart)
		gotPart := bytes.TrimSpace(gotSplit[i])
		if !bytes.Equal(gotPart, wantPart) {
			return false
		}
	}

	return true
}

func codecsMatch(wanted, got Codec) bool {
	if wanted.Name != nil && !bytes.EqualFold(wanted.Name, got.Name) {
		return false
	}
	if wanted.ClockRate != 0 && wanted.ClockRate != got.ClockRate {
		return false
	}
	if wanted.EncodingParameters != nil && !bytes.Equal(wanted.EncodingParameters, got.EncodingParameters) {
		return false
	}
	if wanted.Fmtp != nil && !equivalentFmtp(wanted.Fmtp, got.Fmtp) {
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

type stateFn func(*lexer) (stateFn, error)

type lexer struct {
	desc *SessionDescription
	baseLexer
}

type attrName byte

const invalidAttrName attrName = 0

type attrNameToState func(name attrName) stateFn

func (l *lexer) handleType(fn attrNameToState) (stateFn, error) {
	name, err := l.readFieldName()
	if err == io.EOF {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	if res := fn(name); res != nil {
		return res, nil
	}

	return nil, l.syntaxError()
}

func uintLen(n uint64) int {
	return log10(n) + 1
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

// increase capacity of byte slice slice to accommodate at least n slices
func growByteSliceSlice(b [][]byte, n int) [][]byte {
	if cap(b)-len(b) >= n {
		return b
	}
	bc := make([][]byte, len(b), len(b)+n)
	copy(bc, b)
	return bc
}
