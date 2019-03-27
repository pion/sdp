package sdp

import (
	"net"
	"strconv"
	"strings"
)

// Address desribes a structured address token from within the "c=" field.
type Address struct {
	IP    net.IP
	TTL   *int
	Range *int
}

func (a *Address) Clone() *Address {
	addr := &Address{}
	addr.IP = append(net.IP(nil), a.IP...)

	if a.TTL != nil {
		tmp := *a.TTL
		addr.TTL = &tmp
	}

	if a.Range != nil {
		tmp := *a.Range
		addr.Range = &tmp
	}

	return addr
}

func (a *Address) String() string {
	var parts []string
	parts = append(parts, a.IP.String())
	if a.TTL != nil {
		parts = append(parts, strconv.Itoa(*a.TTL))
	}

	if a.Range != nil {
		parts = append(parts, strconv.Itoa(*a.Range))
	}

	return strings.Join(parts, "/")
}
