package sdp

import (
	"fmt"
	"strings"

	"github.com/pions/webrtc/pkg/rtcerr"

	"github.com/pkg/errors"
)

type Keywds struct {
	Path []string
}

func (k *Keywds) Clone() Attribute {
	kw := &Keywds{}
	kw.Path = append([]string(nil), k.Path...)
	return kw
}

func (k *Keywds) Unmarshal(raw string) error {
	parts := strings.Split(raw, ":")
	if len(parts) != 2 {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", raw)}, pkgName)
	}

	k.Path = strings.Split(parts[1], ".")
	return nil
}

func (k *Keywds) Marshal() string {
	return attributeKey + k.Name() + ":" + strings.Join(k.Path, ".") + endline
}

func (k *Keywds) Name() string {
	return AttributeNameKeywds
}
