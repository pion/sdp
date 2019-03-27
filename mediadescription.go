package sdp

// MediaDescription represents a media type.
//
// https://tools.ietf.org/html/rfc4566#section-5.14
type MediaDescription struct {
	// m=<media> <port>/<number of ports> <proto> <fmt> ...
	Media Media

	// i=<session description>
	Information *Information

	// c=<nettype> <addrtype> <connection-address>
	Connection *Connection

	// b=<bwtype>:<bandwidth>
	Bandwidths Bandwidths

	// k=<method>
	// k=<method>:<encryption key>
	EncryptionKey *EncryptionKey

	// a=<attribute>
	// a=<attribute>:<value>
	Attributes MediaAttributes
}

func (m *MediaDescription) Clone() *MediaDescription {
	mediaDesc := &MediaDescription{}
	mediaDesc.Media = *m.Media.Clone()

	if m.Information != nil {
		mediaDesc.Information = m.Information.Clone()
	}

	if m.Connection != nil {
		mediaDesc.Connection = m.Connection.Clone()
	}

	if len(m.Bandwidths) > 0 {
		mediaDesc.Bandwidths = *m.Bandwidths.Clone()
	}

	if m.EncryptionKey != nil {
		mediaDesc.EncryptionKey = m.EncryptionKey.Clone()
	}

	if len(m.Attributes) > 0 {
		mediaDesc.Attributes = *m.Attributes.Clone()
	}

	return mediaDesc
}

func (m *MediaDescription) MID() string {
	if attrs, ok := m.Attributes.Get(AttributeNameMID); ok {
		return attrs[0].(*MID).Value
	}
	return ""
}

func (m *MediaDescription) Rejected() bool {
	_, hasBundleOnly := m.Attributes.Get(AttributeNameBundleOnly)
	if m.Media.Port.Value == 0 && !hasBundleOnly {
		return true
	}
	return false
}
