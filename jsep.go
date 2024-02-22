// SPDX-FileCopyrightText: 2023 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

package sdp

import (
	"bytes"
	"fmt"
	"strconv"
	"time"
)

// Constants for SDP attributes used in JSEP
const (
	AttrKeyCandidate        = "candidate"
	AttrKeyEndOfCandidates  = "end-of-candidates"
	AttrKeyIdentity         = "identity"
	AttrKeyGroup            = "group"
	AttrKeySSRC             = "ssrc"
	AttrKeySSRCGroup        = "ssrc-group"
	AttrKeyMsid             = "msid"
	AttrKeyMsidSemantic     = "msid-semantic"
	AttrKeyConnectionSetup  = "setup"
	AttrKeyMID              = "mid"
	AttrKeyICELite          = "ice-lite"
	AttrKeyRTCPMux          = "rtcp-mux"
	AttrKeyRTCPRsize        = "rtcp-rsize"
	AttrKeyInactive         = "inactive"
	AttrKeyRecvOnly         = "recvonly"
	AttrKeySendOnly         = "sendonly"
	AttrKeySendRecv         = "sendrecv"
	AttrKeyExtMap           = "extmap"
	AttrKeyExtMapAllowMixed = "extmap-allow-mixed"
)

// Constants for semantic tokens used in JSEP
const (
	SemanticTokenLipSynchronization     = "LS"
	SemanticTokenFlowIdentification     = "FID"
	SemanticTokenForwardErrorCorrection = "FEC"
	SemanticTokenWebRTCMediaStreams     = "WMS"
)

// Constants for extmap key
const (
	ExtMapValueTransportCC = 3
)

var (
	ExtMapValueTransportCCURI = []byte("http://www.ietf.org/id/draft-holmer-rmcat-transport-wide-cc-extensions-01")
)

// API to match draft-ietf-rtcweb-jsep
// Move to webrtc or its own package?

// NewJSEPSessionDescription creates a new SessionDescription with
// some settings that are required by the JSEP spec.
//
// Note: Since v2.4.0, session ID has been fixed to use crypto random according to
//
//	JSEP spec, so that NewJSEPSessionDescription now returns error as a second
//	return value.
func NewJSEPSessionDescription(identity bool) (*SessionDescription, error) {
	sid, err := newSessionID()
	if err != nil {
		return nil, err
	}
	d := &SessionDescription{
		Version: 0,
		Origin: Origin{
			Username:       kDash,
			SessionID:      sid,
			SessionVersion: uint64(time.Now().Unix()),
			NetworkType:    kIn,
			AddressType:    kIp4,
			UnicastAddress: kUnroutableAddr,
		},
		SessionName: SessionName(kDash),
		TimeDescriptions: []TimeDescription{
			{
				Timing: Timing{
					StartTime: 0,
					StopTime:  0,
				},
				RepeatTimes: nil,
			},
		},
		Attributes: []Attribute{
			// 	"Attribute(ice-options:trickle)", // TODO: implement trickle ICE
		},
	}

	if identity {
		d.WithPropertyAttribute(kIdenity)
	}

	return d, nil
}

// WithPropertyAttribute adds a property attribute 'a=key' to the session description
func (s *SessionDescription) WithPropertyAttribute(key []byte) *SessionDescription {
	s.Attributes = append(s.Attributes, NewPropertyAttribute(key))
	return s
}

// WithValueAttribute adds a value attribute 'a=key:value' to the session description
func (s *SessionDescription) WithValueAttribute(key, value []byte) *SessionDescription {
	s.Attributes = append(s.Attributes, NewAttribute(key, value))
	return s
}

// WithFingerprint adds a fingerprint to the session description
func (s *SessionDescription) WithFingerprint(algorithm, value []byte) *SessionDescription {
	return s.WithValueAttribute(kFingerprint, bytes.Join([][]byte{algorithm, value}, kSpace))
}

// WithMedia adds a media description to the session description
func (s *SessionDescription) WithMedia(md MediaDescription) *SessionDescription {
	s.MediaDescriptions = append(s.MediaDescriptions, md)
	return s
}

// NewJSEPMediaDescription creates a new MediaName with
// some settings that are required by the JSEP spec.
func NewJSEPMediaDescription(codecType []byte, _ []string) MediaDescription {
	return MediaDescription{
		MediaName: MediaName{
			Media:  codecType,
			Port:   RangedPort{Value: 9},
			Protos: [][]byte{kUdp, kTls, kRtp, kSavp},
		},
		ConnectionInformation: ConnectionInformation{
			NetworkType: kIn,
			AddressType: kIp4,
			Address: Address{
				Address: kUnroutableAddr,
			},
		},
	}
}

// WithPropertyAttribute adds a property attribute 'a=key' to the media description
func (d *MediaDescription) WithPropertyAttribute(key []byte) *MediaDescription {
	d.Attributes = append(d.Attributes, NewPropertyAttribute(key))
	return d
}

// WithValueAttribute adds a value attribute 'a=key:value' to the media description
func (d *MediaDescription) WithValueAttribute(key, value []byte) *MediaDescription {
	d.Attributes = append(d.Attributes, NewAttribute(key, value))
	return d
}

// WithFingerprint adds a fingerprint to the media description
func (d *MediaDescription) WithFingerprint(algorithm, value []byte) *MediaDescription {
	return d.WithValueAttribute(kFingerprint, bytes.Join([][]byte{algorithm, value}, kSpace))
}

// WithICECredentials adds ICE credentials to the media description
func (d *MediaDescription) WithICECredentials(username, password []byte) *MediaDescription {
	return d.
		WithValueAttribute(kIceUfrag, username).
		WithValueAttribute(kIcePwd, password)
}

// WithCodec adds codec information to the media description
func (d *MediaDescription) WithCodec(payloadType uint8, name string, clockrate uint32, channels uint16, fmtp string) *MediaDescription {
	d.MediaName.Formats = append(d.MediaName.Formats, []byte(strconv.FormatUint(uint64(payloadType), 10)))
	rtpmap := fmt.Sprintf("%d %s/%d", payloadType, name, clockrate)
	if channels > 0 {
		rtpmap += fmt.Sprintf("/%d", channels)
	}
	// TODO
	d.WithValueAttribute(kRtpmap, []byte(rtpmap))
	if fmtp != "" {
		d.WithValueAttribute(kFmtp, []byte(fmt.Sprintf("%d %s", payloadType, fmtp)))
	}
	return d
}

// WithMediaSource adds media source information to the media description
func (d *MediaDescription) WithMediaSource(ssrc uint32, cname, streamLabel, label string) *MediaDescription {
	return d.
		WithValueAttribute(kSsrc, []byte(fmt.Sprintf("%d cname:%s", ssrc, cname))). // Deprecated but not phased out?
		WithValueAttribute(kSsrc, []byte(fmt.Sprintf("%d msid:%s %s", ssrc, streamLabel, label))).
		WithValueAttribute(kSsrc, []byte(fmt.Sprintf("%d mslabel:%s", ssrc, streamLabel))). // Deprecated but not phased out?
		WithValueAttribute(kSsrc, []byte(fmt.Sprintf("%d label:%s", ssrc, label)))          // Deprecated but not phased out?
}

// WithCandidate adds an ICE candidate to the media description
// Deprecated: use WithICECandidate instead
func (d *MediaDescription) WithCandidate(value []byte) *MediaDescription {
	return d.WithValueAttribute(kCandidate, value)
}

// WithExtMap adds an extmap to the media description
func (d *MediaDescription) WithExtMap(e ExtMap) *MediaDescription {
	return d.WithPropertyAttribute(e.Marshal())
}

// WithTransportCCExtMap adds an extmap to the media description
func (d *MediaDescription) WithTransportCCExtMap() *MediaDescription {
	e := ExtMap{
		Value: ExtMapValueTransportCC,
		URI:   ExtMapValueTransportCCURI,
	}
	return d.WithExtMap(e)
}
