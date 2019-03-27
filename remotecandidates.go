package sdp

import (
	"fmt"
	"strings"

	"github.com/pions/webrtc/pkg/rtcerr"

	"github.com/pkg/errors"
)

// RemoteCandidates is defined in https://tools.ietf.org/html/rfc5245.
type RemoteCandidates []RemoteCandidate

func (r *RemoteCandidates) Clone() Attribute {
	remoteCandidates := &RemoteCandidates{}
	for _, remoteCandidate := range *r {
		*remoteCandidates = append(*remoteCandidates, *remoteCandidate.Clone())
	}
	return remoteCandidates
}

func (r *RemoteCandidates) Unmarshal(raw string) error {
	parts := strings.SplitN(raw, ":", 2)
	if len(parts) != 2 {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", raw)}, pkgName)
	}

	fields := strings.Fields(parts[1])
	if len(fields)%3 != 0 {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", parts[1])}, pkgName)
	}

	var remoteCandidates RemoteCandidates
	for i := 0; i < len(fields); i += 3 {
		remoteCandidate := RemoteCandidate{}
		if err := remoteCandidate.Unmarshal(strings.Join(fields[i:i+3], " ")); err != nil {
			return err
		}
		remoteCandidates = append(remoteCandidates, remoteCandidate)
	}

	*r = remoteCandidates
	return nil
}

func (r *RemoteCandidates) Marshal() string {
	return attributeKey + r.Name() + ":" + r.string() + endline
}

func (r *RemoteCandidates) string() string {
	remoteCandidates := make([]string, 0)
	for _, z := range *r {
		remoteCandidates = append(remoteCandidates, z.Marshal())
	}
	return strings.Join(remoteCandidates, " ")
}

func (r *RemoteCandidates) Name() string {
	return AttributeNameRemoteCandidates
}
