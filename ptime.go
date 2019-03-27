package sdp

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pions/webrtc/pkg/rtcerr"

	"github.com/pkg/errors"
)

type Ptime struct {
	Value uint64
}

func (p *Ptime) Clone() Attribute {
	return &Ptime{Value: p.Value}
}

func (p *Ptime) Unmarshal(raw string) error {
	parts := strings.Split(raw, ":")
	if len(parts) != 2 {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", raw)}, pkgName)
	}

	value, err := strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", parts[1])}, pkgName)
	}

	p.Value = value
	return nil
}

func (p *Ptime) Marshal() string {
	return attributeKey + p.Name() + ":" + fmt.Sprintf("%d", p.Value) + endline
}

func (p *Ptime) Name() string {
	return AttributeNamePtime
}
