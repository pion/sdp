package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConnectionRole(t *testing.T) {
	tests := []struct {
		value    string
		expected ConnectionRole
	}{
		{"unknown", ConnectionRole(unknown)},
		{"active", ConnectionRoleActive},
		{"passive", ConnectionRolePassive},
		{"actpass", ConnectionRoleActPass},
		{"holdconn", ConnectionRoleHoldConn},
	}

	for i, u := range tests {
		assert.Equal(t, u.expected, NewConnectionRole(u.value), "%d: %+v", i, u)
	}
}

func TestConnectionRole_String(t *testing.T) {
	tests := []struct {
		actual   ConnectionRole
		expected string
	}{
		{ConnectionRole(unknown), unknownStr},
		{ConnectionRoleActive, "active"},
		{ConnectionRolePassive, "passive"},
		{ConnectionRoleActPass, "actpass"},
		{ConnectionRoleHoldConn, "holdconn"},
	}

	for i, u := range tests {
		assert.Equal(t, u.expected, u.actual.String(), "%d: %+v", i, u)
	}
}
