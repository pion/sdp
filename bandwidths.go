package sdp

type Bandwidths []Bandwidth

func (b *Bandwidths) Clone() *Bandwidths {
	bandwidths := &Bandwidths{}
	for _, bandwidth := range *b {
		*bandwidths = append(*bandwidths, *bandwidth.Clone())
	}
	return bandwidths
}
