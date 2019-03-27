package sdp

type EndOfCandidates struct{}

func (e *EndOfCandidates) Clone() Attribute {
	return &EndOfCandidates{}
}

func (e *EndOfCandidates) Unmarshal(raw string) error {
	return nil
}

func (e *EndOfCandidates) Marshal() string {
	return attributeKey + e.Name() + endline
}

func (e *EndOfCandidates) Name() string {
	return AttributeNameEndOfCandidates
}
