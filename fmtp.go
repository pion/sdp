package sdp

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pions/webrtc/pkg/rtcerr"

	"github.com/pkg/errors"
)

type Fmtp struct {
	Format int
	Params [][]string
}

func (f *Fmtp) Clone() Attribute {
	fmtp := &Fmtp{}
	fmtp.Format = f.Format

	for _, params := range f.Params {
		newParams := append([]string(nil), params...)
		fmtp.Params = append(fmtp.Params, newParams)
	}

	return fmtp
}

func (f *Fmtp) Unmarshal(raw string) error {
	parts := strings.Split(raw, ":")
	if len(parts) != 2 {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", raw)}, pkgName)
	}

	fmtp := strings.SplitN(parts[1], " ", 2)
	if len(parts) != 2 {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", raw)}, pkgName)
	}

	format, err := strconv.ParseInt(fmtp[0], 10, 64)
	if err != nil {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", fmtp[0])}, pkgName)
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

	f.Format = int(format)
	f.Params = params
	return nil
}

func (f *Fmtp) Marshal() string {
	return attributeKey + f.Name() + ":" + f.string() + endline
}

func (f *Fmtp) string() string {
	var params []string
	for _, keyVal := range f.Params {
		params = append(params, keyVal[0]+"="+keyVal[1])
	}

	return fmt.Sprintf(
		"%d %v",
		f.Format,
		strings.Join(params, ";"),
	)
}

func (f *Fmtp) Name() string {
	return AttributeNameFmtp
}
