package sdp

// IceLite is defined in https://tools.ietf.org/html/rfc5245.
type IceLite struct{}

func (i *IceLite) Clone() Attribute {
	return &IceLite{}
}

func (i *IceLite) Unmarshal(raw string) error {
	return nil
}

func (i *IceLite) Marshal() string {
	return attributeKey + i.Name() + endline
}

func (i *IceLite) Name() string {
	return AttributeNameIceLite
}
