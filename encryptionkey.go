package sdp

// EncryptionKey describes the "k=" which conveys encryption key information.
type EncryptionKey struct {
	Value string
}

func (e *EncryptionKey) Clone() *EncryptionKey {
	return &EncryptionKey{Value: e.Value}
}

func (e *EncryptionKey) Unmarshal(raw string) error {
	e.Value = raw
	return nil
}

func (e EncryptionKey) Marshal() string {
	return encryptionKey + e.Value + endline
}
