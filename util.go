package sdp

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/pion/randutil"
)

const (
	attributeKey = "a="
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

func (c Codec) String() string {
	return fmt.Sprintf("%d %s/%d/%s (%s) [%s]", c.PayloadType, c.Name, c.ClockRate, c.EncodingParameters, c.Fmtp, strings.Join(c.RTCPFeedback, ", "))
}

func parseRtpmap(rtpmap string) (Codec, error) {
	var codec Codec
	parsingFailed := errExtractCodecRtpmap

	// a=rtpmap:<payload type> <encoding name>/<clock rate>[/<encoding parameters>]
	split := strings.Split(rtpmap, " ")
	if len(split) != 2 {
		return codec, parsingFailed
	}

	ptSplit := strings.Split(split[0], ":")
	if len(ptSplit) != 2 {
		return codec, parsingFailed
	}

	ptInt, err := strconv.Atoi(ptSplit[1])
	if err != nil {
		return codec, parsingFailed
	}

	codec.PayloadType = uint8(ptInt)

	split = strings.Split(split[1], "/")
	codec.Name = split[0]
	parts := len(split)
	if parts > 1 {
		rate, err := strconv.Atoi(split[1])
		if err != nil {
			return codec, parsingFailed
		}
		codec.ClockRate = uint32(rate)
	}
	if parts > 2 {
		codec.EncodingParameters = split[2]
	}

	return codec, nil
}

func parseFmtp(fmtp string) (Codec, error) {
	var codec Codec
	parsingFailed := errExtractCodecFmtp

	// a=fmtp:<format> <format specific parameters>
	split := strings.Split(fmtp, " ")
	if len(split) != 2 {
		return codec, parsingFailed
	}

	formatParams := split[1]

	split = strings.Split(split[0], ":")
	if len(split) != 2 {
		return codec, parsingFailed
	}

	ptInt, err := strconv.Atoi(split[1])
	if err != nil {
		return codec, parsingFailed
	}

	codec.PayloadType = uint8(ptInt)
	codec.Fmtp = formatParams

	return codec, nil
}

func parseRtcpFb(rtcpFb string) (Codec, error) {
	var codec Codec
	parsingFailed := errExtractCodecRtcpFb

	// a=ftcp-fb:<payload type> <RTCP feedback type> [<RTCP feedback parameter>]
	split := strings.SplitN(rtcpFb, " ", 2)
	if len(split) != 2 {
		return codec, parsingFailed
	}

	ptSplit := strings.Split(split[0], ":")
	if len(ptSplit) != 2 {
		return codec, parsingFailed
	}

	ptInt, err := strconv.Atoi(ptSplit[1])
	if err != nil {
		return codec, parsingFailed
	}

	codec.PayloadType = uint8(ptInt)
	codec.RTCPFeedback = append(codec.RTCPFeedback, split[1])

	return codec, nil
}

type codecFn func(string) (Codec, error)

var codecParsers = map[string]codecFn{
	"rtpmap":  parseRtpmap,
	"fmtp":    parseFmtp,
	"rtcp-fb": parseRtcpFb,
}

func (s *SessionDescription) buildCodecMap() map[uint8]Codec {
	if s.cachedCodecs == nil {
		s.cachedCodecs = make(map[uint8]Codec)
		for _, m := range s.MediaDescriptions {
			for _, a := range m.Attributes {
				attr := a.String()
				if fn, ok := codecParsers[strings.SplitN(attr, ":", 2)[0]]; ok {
					codec, err := fn(attr)
					if err != nil {
						continue
					}

					saved := s.cachedCodecs[codec.PayloadType]

					if saved.PayloadType == 0 {
						saved.PayloadType = codec.PayloadType
					}
					if saved.Name == "" {
						saved.Name = codec.Name
					}
					if saved.ClockRate == 0 {
						saved.ClockRate = codec.ClockRate
					}
					if saved.EncodingParameters == "" {
						saved.EncodingParameters = codec.EncodingParameters
					}
					if saved.Fmtp == "" {
						saved.Fmtp = codec.Fmtp
					}

					saved.RTCPFeedback = append(saved.RTCPFeedback, codec.RTCPFeedback...)

					s.cachedCodecs[saved.PayloadType] = saved
				}
			}
		}
	}
	return s.cachedCodecs
}

func equivalentFmtp(want, got string) bool {
	wantSplit := strings.Split(want, ";")
	gotSplit := strings.Split(got, ";")

	if len(wantSplit) != len(gotSplit) {
		return false
	}

	sort.Strings(wantSplit)
	sort.Strings(gotSplit)

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
	return !((wanted.Name != "" && !strings.EqualFold(wanted.Name, got.Name)) ||
		(wanted.ClockRate != 0 && wanted.ClockRate != got.ClockRate) ||
		(wanted.EncodingParameters != "" && wanted.EncodingParameters != got.EncodingParameters) ||
		wanted.Fmtp != "" && !equivalentFmtp(wanted.Fmtp, got.Fmtp))
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
	baseLexer
	desc *SessionDescription
}
