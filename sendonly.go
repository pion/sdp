package sdp

type SendOnly struct{}

func (s *SendOnly) Clone() Attribute {
	return &SendOnly{}
}

func (s *SendOnly) Unmarshal(raw string) error {
	return nil
}

func (s *SendOnly) Marshal() string {
	return attributeKey + s.Name() + endline
}

func (s *SendOnly) Name() string {
	return AttributeNameSendOnly
}
