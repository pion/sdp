package sdp

// Session describes a structured representations for the "s=" field
// and is the textual session name.
type Session struct {
	Value string
}

func (s *Session) Clone() *Session {
	return &Session{Value: s.Value}
}

func (s *Session) Unmarshal(raw string) error {
	s.Value = raw
	return nil
}

func (s *Session) Marshal() string {
	return sessionKey + s.Value + endline
}
