package sdp

import (
	"fmt"
	"testing"
)

func TestLexer(t *testing.T) {
	t.Run("single field", func(t *testing.T) {
		for k, value := range map[string]string{
			"clean":            "aaa",
			"with extra space": "aaa ",
			"with linebreak":   "aaa \n",
			"with linebreak 2": "aaa \r\n",
		} {
			l := &baseLexer{value: []byte(value)}
			field, err := l.readField()
			if err != nil {
				t.Fatal(err)
			}
			if field != "aaa" {
				t.Errorf("%s: aaa not parsed, got: '%v'", k, field)
			}
		}
	})

	t.Run("syntax error", func(t *testing.T) {
		l := &baseLexer{value: []byte("12NaN")}
		_, err := l.readUint64Field()
		if err != nil {
			fmt.Println("error message:", err.Error())
		} else {
			t.Fatal("no error")
		}
	})

	t.Run("many fields", func(t *testing.T) {
		l := &baseLexer{value: []byte("aaa  123\nf1 f2\nlast")}

		t.Run("first line", func(t *testing.T) {
			field, err := l.readField()
			if err != nil {
				t.Fatal(err)
			}
			if field != "aaa" {
				t.Errorf("aaa not parsed, got: '%v'", field)
			}

			value, err := l.readUint64Field()
			if err != nil {
				t.Fatal(err)
			}
			if value != 123 {
				t.Errorf("aaa not parsed, got: '%v'", field)
			}

			if err := l.nextLine(); err != nil {
				t.Fatal(err)
			}
		})

		t.Run("second line", func(t *testing.T) {
			field, err := l.readField()
			if err != nil {
				t.Fatal(err)
			}
			if field != "f1" {
				t.Errorf("value not parsed, got: '%v'", field)
			}

			field, err = l.readField()
			if err != nil {
				t.Fatal(err)
			}
			if field != "f2" {
				t.Errorf("value not parsed, got: '%v'", field)
			}

			field, err = l.readField()
			if err != nil {
				t.Fatal(err)
			}
			if field != "" {
				t.Errorf("value not parsed, got: '%v'", field)
			}

			if err := l.nextLine(); err != nil {
				t.Fatal(err)
			}
		})

		t.Run("last line", func(t *testing.T) {
			field, err := l.readField()
			if err != nil {
				t.Fatal(err)
			}
			if field != "last" {
				t.Errorf("value not parsed, got: '%v'", field)
			}
		})
	})
}
