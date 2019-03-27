package sdp

type Inactive struct{}

func (i *Inactive) Clone() Attribute {
	return &Inactive{}
}

func (i *Inactive) Unmarshal(raw string) error {
	return nil
}

func (i *Inactive) Marshal() string {
	return attributeKey + i.Name() + endline
}

func (i *Inactive) Name() string {
	return AttributeNameInactive
}
