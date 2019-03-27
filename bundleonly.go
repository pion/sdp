package sdp

// BundleOnly is defined in https://tools.ietf.org/html/draft-ietf-mmusic-sdp-bundle-negotiation-53.
type BundleOnly struct{}

func (b *BundleOnly) Clone() Attribute {
	return &BundleOnly{}
}

func (b *BundleOnly) Unmarshal(raw string) error {
	return nil
}

func (b BundleOnly) Marshal() string {
	return attributeKey + b.Name() + endline
}

func (b BundleOnly) Name() string {
	return AttributeNameBundleOnly
}
