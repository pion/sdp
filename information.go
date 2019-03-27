package sdp

// Information describes the "i=" field which provides textual information
// about the session.
type Information struct {
	Value string
}

func (i *Information) Clone() *Information {
	return &Information{Value: i.Value}
}

func (i *Information) Unmarshal(raw string) error {
	i.Value = raw
	return nil
}

func (i *Information) Marshal() string {
	return infoKey + i.Value + endline
}
