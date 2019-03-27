package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewProtocol(t *testing.T) {
	tests := []struct {
		value    string
		expected Protocol
	}{
		{"unknown", Protocol(unknown)},
		{"UDP", ProtocolUDP},
		{"TCP", ProtocolTCP},
	}

	for i, u := range tests {
		assert.Equal(t, u.expected, NewProtocol(u.value), "%d: %+v", i, u)
	}
}

func TestProtocol_String(t *testing.T) {
	tests := []struct {
		actual   Protocol
		expected string
	}{
		{Protocol(unknown), unknownStr},
		{ProtocolUDP, "UDP"},
		{ProtocolTCP, "TCP"},
	}

	for i, u := range tests {
		assert.Equal(t, u.expected, u.actual.String(), "%d: %+v", i, u)
	}
}
