package sdp

type RtcpMuxOnly struct{}

func (r *RtcpMuxOnly) Clone() Attribute {
	return &RtcpMuxOnly{}
}

func (r *RtcpMuxOnly) Unmarshal(raw string) error {
	return nil
}

func (r *RtcpMuxOnly) Marshal() string {
	return attributeKey + r.Name() + endline
}

func (r *RtcpMuxOnly) Name() string {
	return AttributeNameRtcpMuxOnly
}
