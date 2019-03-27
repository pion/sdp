package sdp

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pions/webrtc/pkg/rtcerr"

	"github.com/pkg/errors"
)

// Origin defines the structure for the "o=" field which provides the
// originator of the session plus a session identifier and version number.
type Origin struct {
	Username       string
	SessionID      uint64
	SessionVersion uint64
	NetworkType    string
	AddressType    string
	UnicastAddress string
}

func (o *Origin) Clone() *Origin {
	return &Origin{
		Username:       o.Username,
		SessionID:      o.SessionID,
		SessionVersion: o.SessionVersion,
		NetworkType:    o.NetworkType,
		AddressType:    o.AddressType,
		UnicastAddress: o.UnicastAddress,
	}
}

func (o *Origin) Unmarshal(raw string) error {
	fields := strings.Fields(raw)
	if len(fields) != 6 {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("o=%v", fields)}, pkgName)
	}

	sessionID, err := strconv.ParseUint(fields[1], 10, 64)
	if err != nil {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", fields[1])}, pkgName)
	}

	sessionVersion, err := strconv.ParseUint(fields[2], 10, 64)
	if err != nil {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", fields[2])}, pkgName)
	}

	// Set according to currently registered with IANA
	// https://tools.ietf.org/html/rfc4566#section-8.2.6
	if i := indexOf(fields[3], []string{"IN"}); i == -1 {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", fields[3])}, pkgName)
	}

	// Set according to currently registered with IANA
	// https://tools.ietf.org/html/rfc4566#section-8.2.7
	if i := indexOf(fields[4], []string{"IP4", "IP6"}); i == -1 {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", fields[4])}, pkgName)
	}

	// TODO validated UnicastAddress - issue #2

	o.Username = fields[0]
	o.SessionID = sessionID
	o.SessionVersion = sessionVersion
	o.NetworkType = fields[3]
	o.AddressType = fields[4]
	o.UnicastAddress = fields[5]
	return nil
}

func (o *Origin) Marshal() string {
	return originKey + fmt.Sprintf(
		"%v %d %d %v %v %v",
		o.Username,
		o.SessionID,
		o.SessionVersion,
		o.NetworkType,
		o.AddressType,
		o.UnicastAddress,
	) + endline
}
