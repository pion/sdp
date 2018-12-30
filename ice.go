package sdp

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/pions/webrtc/pkg/ice"
)

// ICECandidate is used to (un)marshal ICE candidates.
type ICECandidate struct {
	Foundation          string
	Component           uint16
	Priority            uint32
	IP                  string
	Protocol            string
	Port                uint16
	Typ                 string
	RelatedAddress      string
	RelatedPort         uint16
	ExtensionAttributes []ICECandidateAttribute
}

// ICECandidateAttribute represents an ICE candidate extension attribute
type ICECandidateAttribute struct {
	Key   string
	Value string
}

// https://tools.ietf.org/html/draft-ietf-mmusic-ice-sip-sdp-24#section-4.1
// candidate-attribute   = "candidate" ":" foundation SP component-id SP
//                            transport SP
//                            priority SP
//                            connection-address SP     ;from RFC 4566
//                            port         ;port from RFC 4566
//                            SP cand-type
//                            [SP rel-addr]
//                            [SP rel-port]
//                            *(SP extension-att-name SP
//                                 extension-att-value)

func (c ICECandidate) marshalString() string {
	val := fmt.Sprintf("%s %d %s %d %s %d typ %s",
		c.Foundation,
		c.Component,
		c.Protocol,
		c.Priority,
		c.IP,
		c.Port,
		c.Typ)

	if len(c.RelatedAddress) > 0 {
		val = fmt.Sprintf("%s raddr %s rport %d",
			val,
			c.RelatedAddress,
			c.RelatedPort)
	}

	for _, attr := range c.ExtensionAttributes {
		val = fmt.Sprintf("%s %s %s",
			val,
			attr.Key,
			attr.Value)
	}
	return val
}

func (c *ICECandidate) unmarshalString(raw string) error {
	split := strings.Fields(raw)
	if len(split) < 8 {
		return fmt.Errorf("attribute not long enough to be ICE candidate (%d)", len(split))
	}

	// Foundation
	c.Foundation = split[0]

	// Component
	component, err := strconv.ParseUint(split[1], 10, 16)
	if err != nil {
		return fmt.Errorf("could not parse component: %v", err)
	}
	c.Component = uint16(component)

	// Protocol
	c.Protocol = split[2]

	// Priority
	priority, err := strconv.ParseUint(split[3], 10, 32)
	if err != nil {
		return fmt.Errorf("could not parse priority: %v", err)
	}
	c.Priority = uint32(priority)

	// IP
	c.IP = split[4]

	// Port
	port, err := strconv.ParseUint(split[5], 10, 16)
	if err != nil {
		return fmt.Errorf("could not parse port: %v", err)
	}
	c.Port = uint16(port)

	c.Typ = split[7]

	if len(split) <= 8 {
		return nil
	}

	split = split[8:]

	if split[0] == "raddr" {
		if len(split) < 4 {
			return fmt.Errorf("could not parse related addresses: incorrect length")
		}

		// RelatedAddress
		c.RelatedAddress = split[1]

		// RelatedPort
		relatedPort, err := strconv.ParseUint(split[3], 10, 16)
		if err != nil {
			return fmt.Errorf("could not parse port: %v", err)
		}
		c.RelatedPort = uint16(relatedPort)

		if len(split) <= 4 {
			return nil
		}

		split = split[4:]
	}

	for i := 0; len(split) > i+1; i += 2 {
		c.ExtensionAttributes = append(c.ExtensionAttributes, ICECandidateAttribute{
			Key:   split[i],
			Value: split[i+1],
		})
	}

	return nil
}

// TODO: Remove the deprecated code below

// ICECandidateUnmarshal takes a candidate strings and returns a ice.Candidate or nil if it fails to parse
// TODO: return error if parsing fails
// Deprecated: use ICECandidate instead
func ICECandidateUnmarshal(raw string) (*ice.Candidate, error) {
	split := strings.Fields(raw)
	if len(split) < 8 {
		return nil, fmt.Errorf("attribute not long enough to be ICE candidate (%d) %s", len(split), raw)
	}

	getValue := func(key string) string {
		rtrnNext := false
		for _, i := range split {
			if rtrnNext {
				return i
			} else if i == key {
				rtrnNext = true
			}
		}
		return ""
	}

	port, err := strconv.Atoi(split[5])
	if err != nil {
		return nil, err
	}

	transport := split[2]

	// TODO verify valid address
	ip := net.ParseIP(split[4])
	if ip == nil {
		return nil, err
	}

	switch getValue("typ") {
	case "host":
		return ice.NewCandidateHost(transport, ip, port)
	case "srflx":
		return ice.NewCandidateServerReflexive(transport, ip, port, "", 0) // TODO: parse related address
	default:
		return nil, fmt.Errorf("Unhandled candidate typ %s", getValue("typ"))
	}
}

func iceCandidateString(c *ice.Candidate, component int) string {
	// TODO: calculate foundation
	switch c.Type {
	case ice.CandidateTypeHost:
		return fmt.Sprintf("foundation %d %s %d %s %d typ host generation 0",
			component, c.NetworkShort(), c.Priority(c.Type.Preference(), uint16(component)), c.IP, c.Port)

	case ice.CandidateTypeServerReflexive:
		return fmt.Sprintf("foundation %d %s %d %s %d typ srflx raddr %s rport %d generation 0",
			component, c.NetworkShort(), c.Priority(c.Type.Preference(), uint16(component)), c.IP, c.Port,
			c.RelatedAddress.Address, c.RelatedAddress.Port)
	}
	return ""
}

// ICECandidateMarshal takes a candidate and returns a string representation
// Deprecated: use ICECandidate instead
func ICECandidateMarshal(c *ice.Candidate) []string {
	out := make([]string, 0)

	out = append(out, iceCandidateString(c, 1))
	out = append(out, iceCandidateString(c, 2))

	return out
}
