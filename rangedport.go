package sdp

import (
	"strconv"
)

// RangedPort supports special format for the media field "m=" port value. If
// it may be necessary to specify multiple transport ports, the protocol allows
// to write it as: <port>/<number of ports> where number of ports is a an
// offsetting range.
type RangedPort struct {
	Value int
	Range *int
}

func (p *RangedPort) Copy() *RangedPort {
	port := &RangedPort{}
	port.Value = p.Value

	if p.Range != nil {
		tmp := *p.Range
		port.Range = &tmp
	}

	return port
}

func (p *RangedPort) String() string {
	output := strconv.Itoa(p.Value)
	if p.Range != nil {
		output += "/" + strconv.Itoa(*p.Range)
	}
	return output
}
