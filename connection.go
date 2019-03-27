package sdp

import (
	"fmt"
	"github.com/pions/webrtc/pkg/rtcerr"
	"net"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// Connection defines the representation for the "c=" field containing
// connection data.
type Connection struct {
	NetworkType string
	AddressType string
	Address     *Address
}

func (c *Connection) Clone() *Connection {
	conn := &Connection{}
	conn.NetworkType = c.NetworkType
	conn.AddressType = c.AddressType

	if c.Address != nil {
		conn.Address = c.Address.Clone()
	}

	return conn
}

func (c *Connection) Unmarshal(raw string) error {
	fields := strings.Fields(raw)
	if len(fields) < 2 {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("c=%v", fields)}, pkgName)
	}

	// Set according to currently registered with IANA
	// https://tools.ietf.org/html/rfc4566#section-8.2.6
	if i := indexOf(fields[0], []string{"IN"}); i == -1 {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", fields[0])}, pkgName)
	}

	// Set according to currently registered with IANA
	// https://tools.ietf.org/html/rfc4566#section-8.2.7
	if i := indexOf(fields[1], []string{"IP4", "IP6"}); i == -1 {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", fields[1])}, pkgName)
	}

	var connAddr *Address
	if len(fields) > 2 {
		connAddr = &Address{}

		parts := strings.Split(fields[2], "/")
		connAddr.IP = net.ParseIP(parts[0])
		if connAddr.IP == nil {
			return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", fields[2])}, pkgName)
		}

		isIP6 := connAddr.IP.To4() == nil
		if len(parts) > 1 {
			val, err := strconv.ParseInt(parts[1], 10, 32)
			if err != nil {
				return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", fields[2])}, pkgName)
			}

			if isIP6 {
				multi := int(val)
				connAddr.Range = &multi
			} else {
				ttl := int(val)
				connAddr.TTL = &ttl
			}
		}

		if len(parts) > 2 {
			val, err := strconv.ParseInt(parts[2], 10, 32)
			if err != nil {
				return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", fields[2])}, pkgName)
			}

			multi := int(val)
			connAddr.Range = &multi
		}

	}

	c.NetworkType = fields[0]
	c.AddressType = fields[1]
	c.Address = connAddr
	return nil
}

func (c *Connection) Marshal() string {
	return connectionKey + fmt.Sprintf(
		"%v %v %v",
		c.NetworkType,
		c.AddressType,
		c.Address.String(),
	) + endline
}
