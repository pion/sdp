package sdp

import (
	"fmt"
	"strings"

	"github.com/pions/webrtc/pkg/rtcerr"
	"github.com/pkg/errors"
)

const (
	AttributeNameCategory         = "cat"
	AttributeNameKeywds           = "keywds"
	AttributeNameTool             = "tool"
	AttributeNamePtime            = "ptime"
	AttributeNameMaxPtime         = "maxptime"
	AttributeNameRtpMap           = "rtpmap"
	AttributeNameRecvOnly         = "recvonly"
	AttributeNameSendRecv         = "sendrecv"
	AttributeNameSendOnly         = "sendonly"
	AttributeNameInactive         = "inactive"
	AttributeNameOrient           = "orient"
	AttributeNameType             = "type"
	AttributeNameCharset          = "charset"
	AttributeNameSdpLang          = "sdplang"
	AttributeNameLang             = "lang"
	AttributeNameFramerate        = "framerate"
	AttributeNameQuality          = "quality"
	AttributeNameFmtp             = "fmtp"
	AttributeNameCandidate        = "candidate"
	AttributeNameRemoteCandidates = "remote-candidates"
	AttributeNameEndOfCandidates  = "end-of-candidates"
	AttributeNameIceLite          = "ice-lite"
	AttributeNameIceMismatch      = "ice-mismatch"
	AttributeNameIceUfrag         = "ice-ufrag"
	AttributeNameIcePwd           = "ice-pwd"
	AttributeNameIceOptions       = "ice-options"
	AttributeNameIdentity         = "identity"
	AttributeNameGroup            = "group"
	AttributeNameSsrc             = "ssrc"
	AttributeNameRtcpFb           = "rtcp-fb"
	AttributeNameSsrcGroup        = "ssrc-group"
	AttributeNameMsID             = "msid"
	AttributeNameSetup            = "setup"
	AttributeNameMID              = "mid"
	AttributeNameRtcp             = "rtcp"
	AttributeNameRtcpMux          = "rtcp-mux"
	AttributeNameRtcpRsize        = "rtcp-rsize"
	AttributeNameRtcpMuxOnly      = "rtcp-mux-only"
	AttributeNameFingerprint      = "fingerprint"
	AttributeNameBundleOnly       = "bundle-only"
	AttributeNameSctpPort         = "sctp-port"
	AttributeNameMaxMessageSize   = "max-message-size"
	AttributeNameTlsID            = "tls-id"
	AttributeNameExtMap           = "extmap"
	AttributeNameRID              = "rid"
	AttributeNameImageAttr        = "imageattr"
	AttributeNameSimulcast        = "simulcast"
)

// Attribute describes the "a=" field which represents the primary means for
// extending SDP.
type Attribute interface {
	Clone() Attribute
	Unmarshal(raw string) error
	Marshal() string
	Name() string
}

func NewAttribute(name fmt.Stringer) Attribute {
	switch name.String() {
	case AttributeNameCategory:
		return &Category{}
	case AttributeNameKeywds:
		return &Keywds{}
	case AttributeNameTool:
		return &Tool{}
	case AttributeNamePtime:
		return &Ptime{}
	case AttributeNameMaxPtime:
		return &MaxPtime{}
	case AttributeNameRtpMap:
		return &RtpMap{}
	case AttributeNameRecvOnly:
		return &RecvOnly{}
	case AttributeNameSendRecv:
		return &SendRecv{}
	case AttributeNameSendOnly:
		return &SendOnly{}
	case AttributeNameInactive:
		return &Inactive{}
	case AttributeNameOrient:
		return &Orient{}
	case AttributeNameType:
		return &Type{}
	case AttributeNameCharset:
		return &Charset{}
	case AttributeNameSdpLang:
		return &SdpLang{}
	case AttributeNameLang:
		return &Lang{}
	case AttributeNameFramerate:
		return &Framerate{}
	case AttributeNameQuality:
		return &Quality{}
	case AttributeNameFmtp:
		return &Fmtp{}
	case AttributeNameCandidate:
		return &Candidate{}
	case AttributeNameRemoteCandidates:
		return &RemoteCandidates{}
	case AttributeNameEndOfCandidates:
		return &EndOfCandidates{}
	case AttributeNameIceLite:
		return &IceLite{}
	case AttributeNameIceMismatch:
		return &IceMismatch{}
	case AttributeNameIceUfrag:
		return &IceUfrag{}
	case AttributeNameIcePwd:
		return &IcePwd{}
	case AttributeNameIceOptions:
		return &IceOptions{}
	case AttributeNameIdentity:
		return &Identity{}
	case AttributeNameGroup:
		return &Group{}
	case AttributeNameSsrc:
		return &Ssrc{}
	case AttributeNameRtcpFb:
		return &RtcpFb{}
	case AttributeNameSsrcGroup:
		return &SsrcGroup{}
	case AttributeNameMsID:
		return &MsID{}
	case AttributeNameSetup:
		return &Setup{}
	case AttributeNameMID:
		return &MID{}
	case AttributeNameRtcp:
		return &Rtcp{}
	case AttributeNameRtcpMux:
		return &RtcpMux{}
	case AttributeNameRtcpRsize:
		return &RtcpRsize{}
	case AttributeNameRtcpMuxOnly:
		return &RtcpMuxOnly{}
	case AttributeNameFingerprint:
		return &Fingerprint{}
	case AttributeNameBundleOnly:
		return &BundleOnly{}
	case AttributeNameSctpPort:
		return &SctpPort{}
	case AttributeNameMaxMessageSize:
		return &MaxMessageSize{}
	case AttributeNameTlsID:
		return &TlsID{}
	case AttributeNameExtMap:
		return &ExtMap{}
	case AttributeNameRID:
		return &RID{}
	case AttributeNameImageAttr:
		return &ImageAttr{}
	case AttributeNameSimulcast:
		return &Simulcast{}
	default:
		return nil
	}
}

type SessionAttributes []Attribute

func (a *SessionAttributes) Clone() *SessionAttributes {
	attrs := &SessionAttributes{}
	for _, attr := range *a {
		*attrs = append(*attrs, attr.Clone())
	}
	return attrs
}

func (a *SessionAttributes) Add(attr Attribute) {
	if attr == nil {
		return
	}
	*a = append(*a, attr)
}

func (a *SessionAttributes) Get(name string) ([]Attribute, bool) {
	var attrs []Attribute
	for _, attr := range *a {
		if attr.Name() == name {
			attrs = append(attrs, attr)
		}
	}
	if len(attrs) > 0 {
		return attrs, true
	}
	return nil, false
}

// func (a *SessionAttributes) Groups() ([]*Group, bool) {
// 	var groups []*Group
// 	if attrs, ok := (*a).Get(AttributeNameGroup); ok {
// 		groups = append(groups, attrs[0].(*Group))
// 	}
// 	if len(groups) > 0 {
// 		return groups, true
// 	}
// 	return nil, false
// }

// func (a *SessionAttributes) HasGroup(semantic Semantic) bool {
// 	if attrs, ok := (*a).Get(AttributeNameGroup); ok {
// 		if attrs[0].(*Group).Semantic == semantic {
// 			return true
// 		}
// 	}
// 	return false
// }

func (a *SessionAttributes) GetGroup(semantic Semantic) *Group {
	if attrs, ok := (*a).Get(AttributeNameGroup); ok {
		if attrs[0].(*Group).Semantic == semantic {
			return attrs[0].(*Group)
		}
	}
	return nil
}

func (a *SessionAttributes) Unmarshal(raw string) error {
	parts := strings.Split(raw, ":")
	switch parts[0] {
	case AttributeNameCategory,
		AttributeNameKeywds,
		AttributeNameTool,
		AttributeNameRecvOnly,
		AttributeNameSendRecv,
		AttributeNameSendOnly,
		AttributeNameInactive,
		AttributeNameType,
		AttributeNameCharset,
		AttributeNameSdpLang,
		AttributeNameLang,
		AttributeNameIceLite,
		AttributeNameIceUfrag,
		AttributeNameIcePwd,
		AttributeNameIceOptions,
		AttributeNameGroup,
		AttributeNameSetup,
		AttributeNameFingerprint,
		AttributeNameIdentity,
		AttributeNameExtMap:
	default:
		return errors.Wrap(&rtcerr.NotSupportedError{Err: fmt.Errorf("%v", parts[0])}, pkgName)
	}

	attribute := NewAttribute(&stringer{Value: parts[0]})
	if err := attribute.Unmarshal(raw); err != nil {
		return err
	}

	*a = append(*a, attribute)
	return nil
}

type MediaAttributes []Attribute

func (a *MediaAttributes) Clone() *MediaAttributes {
	attrs := &MediaAttributes{}
	for _, attr := range *a {
		*attrs = append(*attrs, attr.Clone())
	}
	return attrs
}

func (a *MediaAttributes) Add(attr Attribute) {
	if attr == nil {
		return
	}
	*a = append(*a, attr)
}

func (a *MediaAttributes) Get(name string) ([]Attribute, bool) {
	var attrs []Attribute
	for _, attr := range *a {
		if attr.Name() == name {
			attrs = append(attrs, attr)
		}
	}

	if len(attrs) > 0 {
		return attrs, true
	}
	return nil, false
}

func (a *MediaAttributes) Unmarshal(raw string) error {
	parts := strings.Split(raw, ":")
	switch parts[0] {
	case AttributeNamePtime,
		AttributeNameMaxPtime,
		AttributeNameRtpMap,
		AttributeNameRecvOnly,
		AttributeNameSendRecv,
		AttributeNameSendOnly,
		AttributeNameInactive,
		AttributeNameOrient,
		AttributeNameSdpLang,
		AttributeNameLang,
		AttributeNameFramerate,
		AttributeNameQuality,
		AttributeNameFmtp,
		AttributeNameCandidate,
		AttributeNameRemoteCandidates,
		AttributeNameEndOfCandidates,
		AttributeNameIceMismatch,
		AttributeNameIceUfrag,
		AttributeNameIcePwd,
		AttributeNameMID,
		AttributeNameBundleOnly,
		AttributeNameSctpPort,
		AttributeNameMaxMessageSize,
		AttributeNameSetup,
		AttributeNameTlsID,
		AttributeNameExtMap,
		AttributeNameSsrc,
		AttributeNameSsrcGroup,
		AttributeNameRtcp,
		AttributeNameRtcpMux,
		AttributeNameRtcpRsize,
		AttributeNameRtcpMuxOnly,
		AttributeNameRtcpFb,
		AttributeNameMsID,
		AttributeNameImageAttr,
		AttributeNameRID,
		AttributeNameFingerprint,
		AttributeNameSimulcast:
	default:
		return errors.Wrap(&rtcerr.NotSupportedError{Err: fmt.Errorf("%v", parts[0])}, pkgName)
	}

	attribute := NewAttribute(&stringer{Value: parts[0]})
	if err := attribute.Unmarshal(raw); err != nil {
		return err
	}

	*a = append(*a, attribute)
	return nil
}
