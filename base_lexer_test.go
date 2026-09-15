// SPDX-FileCopyrightText: 2026 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

package sdp

import (
	"fmt"
	"io"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLexer(t *testing.T) {
	t.Run("single field", func(t *testing.T) {
		for k, value := range map[string]string{
			"clean":            "aaa",
			"with extra space": "aaa ",
			"with linebreak":   "aaa \n",
			"with linebreak 2": "aaa \r\n",
		} {
			l := &baseLexer{value: value}
			field, err := l.readField()
			assert.NoError(t, err)
			assert.Equalf(t, "aaa", field, "%s: aaa not parsed, got: '%v'", k, field)
		}
	})

	t.Run("syntax error", func(t *testing.T) {
		l := &baseLexer{value: "12NaN"}
		_, err := l.readUint64Field()
		assert.Error(t, err)
	})

	t.Run("many fields", func(t *testing.T) {
		lex := &baseLexer{value: "aaa  123\nf1 f2\nlast"}

		t.Run("first line", func(t *testing.T) {
			field, err := lex.readField()
			assert.NoError(t, err)
			assert.Equal(t, "aaa", field)

			value, err := lex.readUint64Field()
			assert.NoError(t, err)
			assert.Equal(t, value, uint64(123))

			assert.NoError(t, lex.nextLine())
		})

		t.Run("second line", func(t *testing.T) {
			field, err := lex.readField()
			assert.NoError(t, err)
			assert.Equal(t, "f1", field)

			field, err = lex.readField()
			assert.NoError(t, err)
			assert.Equal(t, "f2", field)

			field, err = lex.readField()
			assert.NoError(t, err)
			assert.Empty(t, field)

			assert.NoError(t, lex.nextLine())
		})

		t.Run("last line", func(t *testing.T) {
			field, err := lex.readField()
			assert.NoError(t, err)
			assert.Equal(t, "last", field)
		})
	})
}

func TestSyntaxError_Error(t *testing.T) {
	t.Run("index in range", func(t *testing.T) {
		e := syntaxError{s: "hello", i: 1}
		assert.Equal(t, byte('e'), e.s[e.i])
	})

	t.Run("negative index coerced to zero", func(t *testing.T) {
		e := syntaxError{s: "hello", i: -2}
		assert.NotPanics(t, func() { _ = e.Error() })
	})

	t.Run("escaped newline", func(t *testing.T) {
		e := syntaxError{s: "a\nb", i: 1} // points to '\n'
		assert.Equal(t, byte('\n'), e.s[e.i])
	})
}

func TestUnreadByte_ErrorAtStart(t *testing.T) {
	l := &baseLexer{value: "", pos: 0}
	err := l.unreadByte()
	assert.ErrorIs(t, err, errDocumentStart)
	assert.Equal(t, 0, l.pos, "pos should remain at 0 after failed unread")
}

func TestReadByte_ReturnsEOFAtEnd(t *testing.T) {
	l := &baseLexer{value: "a", pos: 1} // already at end
	b, err := l.readByte()
	assert.Equal(t, byte(0), b)
	assert.ErrorIs(t, err, io.EOF)
	assert.Equal(t, 1, l.pos, "pos should not advance on EOF")
}

func TestNextLine_NoErrorOnEOF(t *testing.T) {
	l := &baseLexer{value: "", pos: 0}
	err := l.nextLine()
	assert.NoError(t, err)
	assert.Equal(t, 0, l.pos, "pos should remain at 0 on empty input")
}

func TestReadWhitespace_NoErrorOnEOF(t *testing.T) {
	l := &baseLexer{value: "", pos: 0}
	err := l.readWhitespace()
	assert.NoError(t, err)
	assert.Equal(t, 0, l.pos, "pos should remain at 0 on empty input")
}

func TestReadUint64Field_Errors(t *testing.T) {
	t.Run("empty input -> EOF", func(t *testing.T) {
		l := &baseLexer{value: "", pos: 0}
		_, err := l.readUint64Field()
		assert.ErrorIs(t, err, io.EOF)
	})

	t.Run("non-digit at start -> syntaxError", func(t *testing.T) {
		l := &baseLexer{value: "x123", pos: 0}
		_, err := l.readUint64Field()
		var se syntaxError
		assert.ErrorAs(t, err, &se)
	})
}

func TestReadField_Errors(t *testing.T) {
	t.Run("empty input -> returns EOF", func(t *testing.T) {
		l := &baseLexer{value: "", pos: 0}
		s, err := l.readField()
		assert.Empty(t, s)
		assert.ErrorIs(t, err, io.EOF)
		assert.Equal(t, 0, l.pos)
	})

	t.Run("starting at end of input -> returns EOF", func(t *testing.T) {
		l := &baseLexer{value: "abc", pos: len("abc")}
		s, err := l.readField()
		assert.Empty(t, s)
		assert.ErrorIs(t, err, io.EOF)
		assert.Equal(t, len("abc"), l.pos)
	})
}

func TestReadRequiredField_PropagatesReadFieldError(t *testing.T) {
	// Start at end/empty so readField() returns EOF.
	l := &lexer{baseLexer: baseLexer{value: "", pos: 0}}

	got, err := l.readRequiredField()
	assert.Empty(t, got)
	assert.ErrorIs(t, err, io.EOF)
}

func TestReadRequiredField_FieldMissingOnLeadingWhitespace(t *testing.T) {
	// Leading whitespace makes readField() return "" with nil error,
	// which should trigger errFieldMissing in readRequiredField().
	l := &lexer{baseLexer: baseLexer{value: "   \t"}}

	got, err := l.readRequiredField()
	assert.Empty(t, got)
	assert.ErrorIs(t, err, errFieldMissing)
}

func TestReadLine_CRLFTrimsCorrectly(t *testing.T) {
	l := &baseLexer{value: "abc\r\nx", pos: 0}
	s, err := l.readLine()

	assert.NoError(t, err)
	assert.Equal(t, "abc", s)
	assert.Equal(t, len("abc\r\n"), l.pos, "pos should be after the newline sequence")
}

func TestReadLine_EOFOnEmptyInput(t *testing.T) {
	l := &baseLexer{value: "", pos: 0}
	s, err := l.readLine()

	assert.Empty(t, s)
	assert.ErrorIs(t, err, io.EOF)
	assert.Equal(t, 0, l.pos, "pos should remain at 0 on empty input")
}

func TestReadLine_EOFWhenNoNewlinePresent(t *testing.T) {
	l := &baseLexer{value: "tail", pos: 0}
	s, err := l.readLine()

	assert.Empty(t, s)
	assert.ErrorIs(t, err, io.EOF)
	assert.Equal(t, len("tail"), l.pos, "pos should advance to end on EOF")
}

// largeSDP builds an offer with the given number of video sections, each with
// the attributes a browser typically sends.
func largeSDP(sections int) string {
	var sdp strings.Builder
	sdp.WriteString("v=0\r\no=- 4611731400430051336 2 IN IP4 127.0.0.1\r\ns=-\r\nt=0 0\r\na=group:BUNDLE")
	for i := range sections {
		fmt.Fprintf(&sdp, " %d", i)
	}
	sdp.WriteString("\r\na=msid-semantic: WMS\r\n")
	for i := range sections {
		sdp.WriteString("m=video 9 UDP/TLS/RTP/SAVPF 96 97 98 99 100 101\r\n")
		sdp.WriteString("c=IN IP4 0.0.0.0\r\na=rtcp:9 IN IP4 0.0.0.0\r\n")
		sdp.WriteString("a=ice-ufrag:AbCdEf01\r\na=ice-pwd:0123456789abcdefghijklmnop\r\na=ice-options:trickle\r\n")
		sdp.WriteString("a=fingerprint:sha-256 6B:8B:F0:65:5F:78:E2:51:3B:AC:6F:F3:3F:46:1B:35:")
		sdp.WriteString("DC:B8:5F:64:1A:24:C2:43:F0:A1:58:D0:A1:2C:19:08\r\n")
		fmt.Fprintf(&sdp, "a=setup:actpass\r\na=mid:%d\r\n", i)
		sdp.WriteString("a=extmap:1 urn:ietf:params:rtp-hdrext:toffset\r\n")
		sdp.WriteString("a=extmap:2 http://www.webrtc.org/experiments/rtp-hdrext/abs-send-time\r\n")
		sdp.WriteString("a=extmap:3 http://www.ietf.org/id/draft-holmer-rmcat-transport-wide-cc-extensions-01\r\n")
		sdp.WriteString("a=extmap:4 urn:ietf:params:rtp-hdrext:sdes:mid\r\n")
		fmt.Fprintf(&sdp, "a=sendonly\r\na=msid:stream%d track%d\r\na=rtcp-mux\r\na=rtcp-rsize\r\n", i, i)
		for payloadType, name := range map[int]string{96: "VP8", 98: "VP9", 100: "H264"} {
			fmt.Fprintf(&sdp, "a=rtpmap:%d %s/90000\r\n", payloadType, name)
			for _, feedback := range []string{"goog-remb", "transport-cc", "ccm fir", "nack", "nack pli"} {
				fmt.Fprintf(&sdp, "a=rtcp-fb:%d %s\r\n", payloadType, feedback)
			}
			fmt.Fprintf(&sdp, "a=rtpmap:%d rtx/90000\r\na=fmtp:%d apt=%d\r\n", payloadType+1, payloadType+1, payloadType)
		}
		fmt.Fprintf(&sdp, "a=ssrc-group:FID %d %d\r\n", 1000+2*i, 1001+2*i)
		fmt.Fprintf(&sdp, "a=ssrc:%d cname:cname%d\r\na=ssrc:%d cname:cname%d\r\n", 1000+2*i, i, 1001+2*i, i)
	}

	return sdp.String()
}

func heapInuse() uint64 {
	runtime.GC()
	runtime.GC()
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	return memStats.HeapInuse
}

// Keeping one parsed attribute must not keep the whole input alive.
func TestUnmarshalDoesNotRetainInput(t *testing.T) {
	const parses = 16
	input := largeSDP(200)

	before := heapInuse()
	mids := make([]string, 0, parses)
	for range parses {
		// a fresh copy per parse, so only the parse result can keep it alive
		sd := SessionDescription{}
		assert.NoError(t, sd.UnmarshalString(strings.Clone(input)))
		mid, ok := sd.MediaDescriptions[0].Attribute(AttrKeyMID)
		assert.True(t, ok)
		mids = append(mids, mid)
	}
	after := heapInuse()

	t.Logf("heap in use: before %d KiB, after %d KiB (%d parses of a %d KiB input)",
		before/1024, after/1024, parses, len(input)/1024)
	// before the lexer copied its output, this grew by parses*len(input)
	assert.Less(t, after, before+uint64(len(input))/2, "parsed strings pin the input")
	runtime.KeepAlive(mids)
}

func BenchmarkUnmarshalLarge(b *testing.B) {
	input := largeSDP(200)
	b.SetBytes(int64(len(input)))
	b.ReportAllocs()
	for range b.N {
		var sd SessionDescription
		assert.NoError(b, sd.UnmarshalString(input))
	}
}
