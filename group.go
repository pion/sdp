package sdp

import (
	"fmt"
	"strings"

	"github.com/pions/webrtc/pkg/rtcerr"

	"github.com/pkg/errors"
)

// Group is defined in https://tools.ietf.org/html/rfc5888.
type Group struct {
	Semantic Semantic
	MIDs     []string
}

func (g *Group) Clone() Attribute {
	group := &Group{}
	group.Semantic = g.Semantic
	group.MIDs = append([]string(nil), g.MIDs...)
	return group
}

func (g *Group) Unmarshal(raw string) error {
	parts := strings.Split(raw, ":")
	if len(parts) != 2 {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", raw)}, pkgName)
	}

	split := strings.Fields(parts[1])
	if len(parts) < 1 {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", raw)}, pkgName)
	}

	semantic := NewSemantic(split[0])
	if semantic == Semantic(unknown) {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", split[0])}, pkgName)
	}

	g.Semantic = semantic
	g.MIDs = append(g.MIDs, split[1:]...)
	return nil
}

func (g *Group) Marshal() string {
	return attributeKey + g.Name() + ":" + g.string() + endline
}

func (g *Group) string() string {
	MIDs := strings.Join(g.MIDs, " ")
	if len(g.MIDs) > 0 {
		MIDs = " " + MIDs
	}

	return fmt.Sprintf(
		"%v%v",
		g.Semantic.String(),
		MIDs,
	)
}

func (g *Group) Name() string {
	return AttributeNameGroup
}

// FirstMID returns the first mid
func (g *Group) FirstMID() *string {
	if len(g.MIDs) > 0 {
		tmp := g.MIDs[0]
		return &tmp
	}
	return nil
}

func (g *Group) FindMID(mid string) int {
	for i, each := range g.MIDs {
		if each == mid {
			return i
		}
	}
	return -1
}

func (g *Group) HasMID(mid string) bool {
	if i := g.FindMID(mid); i != -1 {
		return true
	}
	return false
}

func (g *Group) AddMID(mid string) int {
	if i := g.FindMID(mid); i != -1 {
		return i
	}
	g.MIDs = append(g.MIDs, mid)
	return len(g.MIDs) - 1
}

func (g *Group) RemoveMID(mid string) bool {
	if i := g.FindMID(mid); i != -1 {
		copy(g.MIDs[i:], g.MIDs[i+1:])
		g.MIDs[len(g.MIDs)-1] = ""
		g.MIDs = g.MIDs[:len(g.MIDs)-1]
		return true
	}
	return false
}
