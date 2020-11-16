package sdp

import "testing"

func TestLexer(t *testing.T) {
	t.Run("single field", func(t *testing.T) {
		for k, s := range map[string]string{
			"clean":            "aaa",
			"with extra space": "aaa ",
			"with linebreak":   "aaa \n",
			"with linebreak 2": "aaa \r\n",
		} {
			t.Run(k, func(t *testing.T) {
				l := &baseLexer{data: []byte(s)}
				field, err := l.readStringField()
				if err != nil {
					t.Fatal(err)
				}
				if field != "aaa" {
					t.Errorf("aaa not parsed, got: '%v'", field)
				}
			})
		}
	})

	t.Run("two fields", func(t *testing.T) {
		l := &baseLexer{data: []byte(`aaa  123`)}

		field, err := l.readStringField()
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
			t.Errorf("value not parsed, got: '%v'", value)
		}
	})

}
