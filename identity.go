package sdp

import (
	"fmt"
	"strings"

	"github.com/pions/webrtc/pkg/rtcerr"

	"github.com/pkg/errors"
)

// Identity is defined in https://tools.ietf.org/html/draft-ietf-rtcweb-security-arch-15.
type Identity struct {
	Assertion string
	Params    [][]string
}

func (i *Identity) Clone() Attribute {
	identity := &Identity{}
	identity.Assertion = i.Assertion

	for _, params := range i.Params {
		newParams := append([]string(nil), params...)
		identity.Params = append(identity.Params, newParams)
	}

	return identity
}

func (i *Identity) Unmarshal(raw string) error {
	parts := strings.Split(raw, ":")
	if len(parts) != 2 {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", raw)}, pkgName)
	}

	fmtp := strings.SplitN(parts[1], " ", 2)
	if len(parts) != 2 {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", raw)}, pkgName)
	}

	var params [][]string
	rawParams := strings.Split(fmtp[1], ";")

	for _, rawParam := range rawParams {
		if rawParam == "" {
			continue
		}

		keyVal := strings.Split(strings.TrimSpace(rawParam), "=")
		if len(keyVal) != 2 {
			return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", rawParams)}, pkgName)
		}
		params = append(params, keyVal)
	}

	i.Assertion = fmtp[0]
	i.Params = params
	return nil
}

func (i Identity) Marshal() string {
	return attributeKey + i.Name() + ":" + i.String() + endline
}

func (i Identity) String() string {
	var params []string
	for _, keyVal := range i.Params {
		params = append(params, keyVal[0]+"="+keyVal[1])
	}

	return fmt.Sprintf(
		"%v %v",
		i.Assertion,
		strings.Join(params, ";"),
	)
}

func (i *Identity) Name() string {
	return AttributeNameIdentity
}
