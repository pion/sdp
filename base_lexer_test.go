// SPDX-FileCopyrightText: 2023 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

package sdp

import (
	"io"
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
