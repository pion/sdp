package sdp

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pions/webrtc/pkg/rtcerr"

	"github.com/pkg/errors"
)

type Quality struct {
	Value int
}

func (q *Quality) Clone() Attribute {
	return &Quality{Value: q.Value}
}

func (q *Quality) Unmarshal(raw string) error {
	parts := strings.Split(raw, ":")
	if len(parts) != 2 {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", raw)}, pkgName)
	}

	value, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", parts[1])}, pkgName)
	}

	q.Value = int(value)
	return nil
}

func (q *Quality) Marshal() string {
	return attributeKey + q.Name() + ":" + fmt.Sprintf("%d", q.Value) + endline
}

func (q *Quality) Name() string {
	return AttributeNameQuality
}
