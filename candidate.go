package sdp

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pions/webrtc/pkg/rtcerr"

	"github.com/pkg/errors"
)

// Candidate is defined in https://tools.ietf.org/html/rfc5245.
type Candidate struct {
	Foundation string
	Component  uint
	Protocol   Protocol
	Priority   uint
	Addr       string
	Port       uint16
	Type       CandidateType
	Params     [][]string
}

func (c *Candidate) Clone() Attribute {
	candidate := &Candidate{}
	candidate.Foundation = c.Foundation
	candidate.Component = c.Component
	candidate.Protocol = c.Protocol
	candidate.Priority = c.Priority
	candidate.Addr = c.Addr
	candidate.Port = c.Port
	candidate.Type = c.Type

	for _, params := range c.Params {
		newParams := append([]string(nil), params...)
		candidate.Params = append(candidate.Params, newParams)
	}

	return candidate
}

func (c *Candidate) Unmarshal(raw string) error {
	parts := strings.Split(raw, ":")
	if len(parts) != 2 {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", raw)}, pkgName)
	}

	split := strings.Fields(parts[1])
	if len(split) < 8 {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", raw)}, pkgName)
	}

	component, err := strconv.ParseUint(split[1], 10, 64)
	if err != nil {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", split[1])}, pkgName)
	}

	protocol := NewProtocol(split[2])
	if protocol == Protocol(unknown) {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", split[2])}, pkgName)
	}

	priority, err := strconv.ParseUint(split[3], 10, 64)
	if err != nil {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", split[3])}, pkgName)
	}

	port, err := parsePort(split[5])
	if err != nil {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", split[5])}, pkgName)
	}

	if split[6] != "typ" {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", raw)}, pkgName)
	}

	candidateType := NewCandidateType(split[7])
	if candidateType == CandidateType(unknown) {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", split[7])}, pkgName)
	}

	var params [][]string
	for i := 8; i < len(split); i += 2 {
		if len(split) < i+1 {
			return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", raw)}, pkgName)
		}
		params = append(params, split[i:i+2])
	}

	c.Foundation = split[0]
	c.Component = uint(component)
	c.Protocol = protocol
	c.Priority = uint(priority)
	c.Addr = split[4] // TODO verify address - issue #3
	c.Port = uint16(port)
	c.Type = candidateType
	c.Params = params
	return nil
}

func (c *Candidate) Marshal() string {
	return attributeKey + c.Name() + ":" + c.string() + endline
}

func (c *Candidate) string() string {
	var params []string
	for _, keyVal := range c.Params {
		params = append(params, keyVal[0]+" "+keyVal[1])
	}

	if len(params) > 0 {
		params = append([]string{""}, params...)
	}

	return fmt.Sprintf(
		"%v %d %v %d %v %d typ %v%v",
		c.Foundation,
		c.Component,
		c.Protocol.String(),
		c.Priority,
		c.Addr,
		c.Port,
		c.Type.String(),
		strings.Join(params, " "),
	)
}

func (c *Candidate) Name() string {
	return AttributeNameCandidate
}
