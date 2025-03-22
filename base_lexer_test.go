// SPDX-FileCopyrightText: 2023 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

package sdp

import (
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
			assert.Equalf(t, field, "aaa", "%s: aaa not parsed, got: '%v'", k, field)
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
			assert.Equal(t, field, "aaa")

			value, err := lex.readUint64Field()
			assert.NoError(t, err)
			assert.Equal(t, value, uint64(123))

			assert.NoError(t, lex.nextLine())
		})

		t.Run("second line", func(t *testing.T) {
			field, err := lex.readField()
			assert.NoError(t, err)
			assert.Equal(t, field, "f1")

			field, err = lex.readField()
			assert.NoError(t, err)
			assert.Equal(t, field, "f2")

			field, err = lex.readField()
			assert.NoError(t, err)
			assert.Equal(t, field, "")

			assert.NoError(t, lex.nextLine())
		})

		t.Run("last line", func(t *testing.T) {
			field, err := lex.readField()
			assert.NoError(t, err)
			assert.Equal(t, field, "last")
		})
	})
}
