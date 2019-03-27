package sdp

// PhoneNumber describes a structured representations for the "p=" line
// specify phone contact information for the person responsible for the
// conference.
type PhoneNumber struct {
	Value string
}

func (p *PhoneNumber) Clone() *PhoneNumber {
	return &PhoneNumber{Value: p.Value}
}

func (p *PhoneNumber) Unmarshal(raw string) error {
	p.Value = raw
	return nil
}

func (p *PhoneNumber) Marshal() string {
	return phoneKey + p.Value + endline
}
