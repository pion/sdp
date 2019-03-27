package sdp

// EmailAddress describes a structured representations for the "e=" line
// which specifies email contact information for the person responsible for
// the conference.
type EmailAddress struct {
	Value string
}

func (e *EmailAddress) Clone() *EmailAddress {
	return &EmailAddress{Value: e.Value}
}

func (e *EmailAddress) Unmarshal(raw string) error {
	e.Value = raw
	return nil
}

func (e *EmailAddress) Marshal() string {
	return emailKey + e.Value + endline
}
