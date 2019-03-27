package sdp

type RecvOnly struct{}

func (r *RecvOnly) Clone() Attribute {
	return &RecvOnly{}
}

func (r *RecvOnly) Unmarshal(raw string) error {
	return nil
}

func (r *RecvOnly) Marshal() string {
	return attributeKey + r.Name() + endline
}

func (r *RecvOnly) Name() string {
	return AttributeNameRecvOnly
}
