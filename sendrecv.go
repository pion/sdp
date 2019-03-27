package sdp

type SendRecv struct{}

func (s *SendRecv) Clone() Attribute {
	return &SendRecv{}
}

func (s *SendRecv) Unmarshal(raw string) error {
	return nil
}

func (s *SendRecv) Marshal() string {
	return attributeKey + s.Name() + endline
}

func (s *SendRecv) Name() string {
	return AttributeNameSendRecv
}
