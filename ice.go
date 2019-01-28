package sdp

import (
	"fmt"
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

	return nil
}
