package sdp

type RtcpRsize struct{}

func (r *RtcpRsize) Clone() Attribute {
	return &RtcpRsize{}
}

func (r *RtcpRsize) Unmarshal(raw string) error {
	return nil
}

func (r *RtcpRsize) Marshal() string {
	return attributeKey + r.Name() + endline
}

func (r *RtcpRsize) Name() string {
	return AttributeNameRtcpRsize
}
