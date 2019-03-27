package sdp

// IceMismatch is defined in https://tools.ietf.org/html/rfc5245.
type IceMismatch struct{}

func (i *IceMismatch) Clone() Attribute {
	return &IceMismatch{}
}

func (i *IceMismatch) Unmarshal(raw string) error {
	return nil
}

func (i *IceMismatch) Marshal() string {
	return attributeKey + i.Name() + endline
}

func (i *IceMismatch) Name() string {
	return AttributeNameIceMismatch
}
