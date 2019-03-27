package sdp

type RtcpMux struct{}

func (r *RtcpMux) Clone() Attribute {
	return &RtcpMux{}
}

func (r *RtcpMux) Unmarshal(raw string) error {
	return nil
}

func (r *RtcpMux) Marshal() string {
	return attributeKey + r.Name() + endline
}

func (r *RtcpMux) Name() string {
	return AttributeNameRtcpMux
}
