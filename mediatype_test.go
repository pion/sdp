package sdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMediaType(t *testing.T) {
	tests := []struct {
		value    string
		expected MediaType
	}{
		{"unknown", MediaType(unknown)},
		{"audio", MediaTypeAudio},
		{"video", MediaTypeVideo},
		{"application", MediaTypeApplication},
		{"text", MediaTypeText},
		{"message", MediaTypeMessage},
	}

	for i, u := range tests {
		assert.Equal(t, u.expected, NewMediaType(u.value), "%d: %+v", i, u)
	}
}

func TestMediaType_String(t *testing.T) {
	tests := []struct {
		actual   MediaType
		expected string
	}{
		{MediaType(unknown), unknownStr},
		{MediaTypeAudio, "audio"},
		{MediaTypeVideo, "video"},
		{MediaTypeApplication, "application"},
		{MediaTypeText, "text"},
		{MediaTypeMessage, "message"},
	}

	for i, u := range tests {
		assert.Equal(t, u.expected, u.actual.String(), "%d: %+v", i, u)
	}
}
