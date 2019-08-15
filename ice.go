package sdp

import (
	"fmt"
	"strconv"
	"strings"
)

// ICECandidate is used to (un)marshal ICE candidates.
type ICECandidate struct {
	Foundation          string
	Component           uint16
	Priority            uint32
	Address             string
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

// https://tools.ietf.org/id/draft-ietf-mmusic-ice-sip-sdp-14.html

// ICEMediaParameters media-level ICE parameters
type ICEMediaParameters struct {
	UserFragment string
	Password     string
	Mismatch     bool
	Pacing       uint16
	Options      string
}

// ICESessionParameters session-level ICE parameters
type ICESessionParameters struct {
	UserFragment string //a=ice-ufrag:8hhY
	Password     string //a=ice-pwd:asd88fgpdd777uzjYhagZg
	Lite         bool   //a=ice-lite
	Options      string //a=ice-options:ice2
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

// Marshal returns the string representation of the ICECandidate
func (c ICECandidate) Marshal() string {
	val := fmt.Sprintf("%s %d %s %d %s %d typ %s",
		c.Foundation,
		c.Component,
		c.Protocol,
		c.Priority,
		c.Address,
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

// Unmarshal popuulates the ICECandidate from its string representation
func (c *ICECandidate) Unmarshal(raw string) error {
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

	// Address
	c.Address = split[4]

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

// GetICEParams extracts session-level ICE parameters
func (s *SessionDescription) GetICEParams() ICESessionParameters {
	ret := ICESessionParameters{}

	var value string
	var found bool

	_, found = s.Attribute("ice-lite")
	if found {
		ret.Lite = true
	} else {
		ret.Lite = false
	}
	value, found = s.Attribute("ice-ufrag")
	if found {
		ret.UserFragment = value
	} else {
		ret.UserFragment = ""
	}
	value, found = s.Attribute("ice-pwd")
	if found {
		ret.Password = value
	} else {
		ret.Password = ""
	}
	value, found = s.Attribute("ice-options")
	if found {
		ret.Options = value
	} else {
		ret.Options = ""
	}

	return ret
}

// GetICEParams extracts media-section-level ICE parameters
func (s *MediaDescription) GetICEParams() ICEMediaParameters {
	ret := ICEMediaParameters{}

	var value string
	var found bool

	_, found = s.Attribute("ice-mismatch")
	if found {
		ret.Mismatch = true
	} else {
		ret.Mismatch = false
	}
	value, found = s.Attribute("ice-ufrag")
	if found {
		ret.UserFragment = value
	} else {
		ret.UserFragment = ""
	}
	value, found = s.Attribute("ice-pacing")
	if found {
		intval, err := strconv.Atoi(value)

		if err != nil {
			//todo: proper logging here
			ret.Pacing = 0
		} else {
			ret.Pacing = (uint16)(intval)
		}

	} else {
		ret.UserFragment = ""
	}
	value, found = s.Attribute("ice-pwd")
	if found {
		ret.Password = value
	} else {
		ret.Password = ""
	}
	value, found = s.Attribute("ice-options")
	if found {
		ret.Options = value
	} else {
		ret.Options = ""
	}

	return ret
}
