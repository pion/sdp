package sdp

import (
	"fmt"
	"strings"

	"github.com/pions/webrtc/pkg/rtcerr"

	"github.com/pkg/errors"
)

// Fingerprint is defined in https://tools.ietf.org/html/rfc4572.
type Fingerprint struct {
	HashFunc    HashFunc
	Fingerprint string
}

func (f *Fingerprint) Clone() Attribute {
	return &Fingerprint{
		HashFunc:    f.HashFunc,
		Fingerprint: f.Fingerprint,
	}
}

func (f *Fingerprint) Unmarshal(raw string) error {
	parts := strings.SplitN(raw, ":", 2)
	if len(parts) != 2 {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", raw)}, pkgName)
	}

	split := strings.Fields(parts[1])
	if len(split) != 2 {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", raw)}, pkgName)
	}

	hashFunc := NewHashFunc(split[0])
	if hashFunc == HashFunc(unknown) {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", split[0])}, pkgName)
	}

	f.HashFunc = hashFunc
	f.Fingerprint = split[1]
	return nil
}

func (f Fingerprint) Marshal() string {
	return attributeKey + f.Name() + ":" + f.String() + endline
}

func (f *Fingerprint) String() string {
	return fmt.Sprintf(
		"%v %v",
		f.HashFunc.String(),
		f.Fingerprint,
	)
}

func (f *Fingerprint) Name() string {
	return AttributeNameFingerprint
}
