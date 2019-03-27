package sdp

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pions/webrtc/pkg/rtcerr"

	"github.com/pkg/errors"
)

// Bandwidth describes an optional field which denotes the proposed bandwidth
// to be used by the session or media.
type Bandwidth struct {
	Experimental bool
	Type         string
	Bandwidth    uint64
}

func (b *Bandwidth) Clone() *Bandwidth {
	return &Bandwidth{
		Experimental: b.Experimental,
		Type:         b.Type,
		Bandwidth:    b.Bandwidth,
	}
}

func (b *Bandwidth) Unmarshal(raw string) error {
	parts := strings.Split(raw, ":")
	if len(parts) != 2 {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("b=%v", parts)}, pkgName)
	}

	experimental := strings.HasPrefix(parts[0], "X-")
	if experimental {
		parts[0] = strings.TrimPrefix(parts[0], "X-")
	} else {
		// Set according to currently registered with IANA
		// https://tools.ietf.org/html/rfc4566#section-5.8
		if i := indexOf(parts[0], []string{"CT", "AS"}); i == -1 {
			return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", parts[0])}, pkgName)
		}
	}

	bandwidth, err := strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", parts[1])}, pkgName)
	}

	b.Experimental = experimental
	b.Type = parts[0]
	b.Bandwidth = bandwidth
	return nil
}

func (b *Bandwidth) Marshal() string {
	return bandwidthKey + b.string() + endline
}

func (b *Bandwidth) string() string {
	var output string
	if b.Experimental {
		output += "X-"
	}
	output += b.Type + ":" + strconv.FormatUint(b.Bandwidth, 10)
	return output
}
