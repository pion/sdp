// SPDX-FileCopyrightText: 2023 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

package sdp

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
)

var errDocumentStart = errors.New("already on document start")

const eof byte = 0

type syntaxError struct {
	s string
	i int
}

func (e syntaxError) Error() string {
	if e.i < 0 {
		e.i = 0
	}
	return fmt.Sprintf("sdp: syntax error at pos %d: %s", e.i, strconv.QuoteToASCII(e.s[e.i:e.i+1]))
}

type baseLexer struct {
	value []byte
	pos   int
	attrs []Attribute
}

func (l *baseLexer) reset() {
	l.pos = 0
}

func (l baseLexer) syntaxError() error {
	return syntaxError{s: string(l.value), i: l.pos - 1}
}

func (l *baseLexer) unreadByte() error {
	if l.pos <= 0 {
		return errDocumentStart
	}
	l.pos--
	return nil
}

func (l *baseLexer) readByte() byte {
	if l.pos >= len(l.value) {
		return eof
	}
	l.pos++
	return l.value[l.pos-1]
}

func (l *baseLexer) nextLine() error {
	for {
		ch := l.readByte()
		if ch == eof {
			return nil
		}
		if !isNewline(ch) {
			return l.unreadByte()
		}
	}
}

func (l *baseLexer) readWhitespace() error {
	for {
		ch := l.readByte()
		if ch == eof {
			return nil
		}
		if !isWhitespace(ch) {
			return l.unreadByte()
		}
	}
}

func (l *baseLexer) readUint64Field() (i uint64, err error) {
	for {
		ch := l.readByte()
		if ch == eof {
			if i == 0 {
				return i, io.EOF
			}
			break
		}

		if isNewline(ch) {
			if err := l.unreadByte(); err != nil {
				return i, err
			}
			break
		}

		if isWhitespace(ch) {
			if err := l.readWhitespace(); err != nil {
				return i, err
			}
			break
		}

		if ch < '0' || ch > '9' {
			return i, l.syntaxError()
		}

		i = i*10 + uint64(ch-'0')
	}

	return i, nil
}

// Returns next field on this line or empty string if no more fields on line
func (l *baseLexer) readField() ([]byte, error) {
	start := l.pos
	var stop int
	for {
		stop = l.pos
		ch := l.readByte()
		if ch == eof {
			if stop == start {
				return nil, io.EOF
			}
			break
		}

		if isNewline(ch) {
			if err := l.unreadByte(); err != nil {
				return nil, err
			}
			break
		}

		if isWhitespace(ch) {
			if err := l.readWhitespace(); err != nil {
				return nil, err
			}
			break
		}
	}
	return l.value[start:stop], nil
}

// Returns symbols until line end
func (l *baseLexer) readLine() ([]byte, error) {
	start := l.pos
	trim := 1
	for {
		switch l.readByte() {
		case eof:
			return nil, io.EOF
		case '\r':
			trim++
		case '\n':
			return l.value[start : l.pos-trim], nil
		}
	}
}

func (l *baseLexer) readUntil(until byte) ([]byte, error) {
	start := l.pos
	for {
		switch l.readByte() {
		case eof:
			return nil, io.EOF
		case until:
			return l.value[start:l.pos], nil
		}
	}
}

func (l *baseLexer) readFieldName() (attrName, error) {
	for {
		ch := l.readByte()
		if ch == eof {
			return invalidAttrName, io.EOF
		}

		if isNewline(ch) {
			continue
		}

		err := l.unreadByte()
		if err != nil {
			return invalidAttrName, err
		}

		name, err := l.readUntil('=')
		if err != nil {
			return invalidAttrName, err
		}

		if len(name) == 2 {
			return attrName(name[0]), nil
		}

		return invalidAttrName, l.syntaxError()
	}
}

func isNewline(ch byte) bool { return ch == '\n' || ch == '\r' }

func isWhitespace(ch byte) bool { return ch == ' ' || ch == '\t' }

func anyOf(element []byte, data ...[]byte) bool {
	for _, v := range data {
		if bytes.Equal(element, v) {
			return true
		}
	}
	return false
}
