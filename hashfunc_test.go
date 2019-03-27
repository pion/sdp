package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHashFunc(t *testing.T) {
	tests := []struct {
		value    string
		expected HashFunc
	}{
		{"unknown", HashFunc(unknown)},
		{"sha-1", HashFuncSHA1},
		{"sha-224", HashFuncSHA224},
		{"sha-256", HashFuncSHA256},
		{"sha-384", HashFuncSHA384},
		{"sha-512", HashFuncSHA512},
		{"md5", HashFuncMD5},
		{"md2", HashFuncMD2},
	}

	for i, u := range tests {
		assert.Equal(t, u.expected, NewHashFunc(u.value), "%d: %+v", i, u)
	}
}

func TestHashFunc_String(t *testing.T) {
	tests := []struct {
		actual   HashFunc
		expected string
	}{
		{HashFunc(unknown), unknownStr},
		{HashFuncSHA1, "sha-1"},
		{HashFuncSHA224, "sha-224"},
		{HashFuncSHA256, "sha-256"},
		{HashFuncSHA384, "sha-384"},
		{HashFuncSHA512, "sha-512"},
		{HashFuncMD5, "md5"},
		{HashFuncMD2, "md2"},
	}

	for i, u := range tests {
		assert.Equal(t, u.expected, u.actual.String(), "%d: %+v", i, u)
	}
}
