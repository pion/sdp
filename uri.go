package sdp

import (
	"net/url"

	"github.com/pions/webrtc/pkg/rtcerr"
	"github.com/pkg/errors"
)

type URL struct {
	Value url.URL
}

func (u *URL) Clone() *URL {
	uri, _ := url.Parse(u.Value.String())
	return &URL{Value: *uri}
}

func (u *URL) Unmarshal(raw string) error {
	uri, err := url.Parse(raw)
	if err != nil {
		return errors.Wrap(&rtcerr.UnknownError{Err: err}, pkgName)
	}

	u.Value = *uri
	return nil
}

func (u *URL) Marshal() string {
	return uriKey + u.Value.String() + endline
}
