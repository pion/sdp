package sdp

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pions/webrtc/pkg/rtcerr"

	"github.com/pkg/errors"
)

type RemoteCandidate struct {
	Component uint
	Addr      string
	Port      uint16
}

func (r *RemoteCandidate) Clone() *RemoteCandidate {
	return &RemoteCandidate{
		Component: r.Component,
		Addr:      r.Addr,
		Port:      r.Port,
	}
}

func (r *RemoteCandidate) Unmarshal(raw string) error {
	split := strings.Fields(raw)
	if len(split) != 3 {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", raw)}, pkgName)
	}

	component, err := strconv.ParseUint(split[0], 10, 64)
	if err != nil {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", split[0])}, pkgName)
	}

	port, err := parsePort(split[2])
	if err != nil {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", split[2])}, pkgName)
	}

	r.Component = uint(component)
	r.Addr = split[1]
	r.Port = uint16(port)
	return nil
}

func (r *RemoteCandidate) Marshal() string {
	return fmt.Sprintf("%d %v %d", r.Component, r.Addr, r.Port)
}
