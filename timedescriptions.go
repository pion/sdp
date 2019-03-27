package sdp

type TimeDescriptions []TimeDescription

func (t *TimeDescriptions) Clone() *TimeDescriptions {
	timeDescs := &TimeDescriptions{}
	for _, timeDesc := range *t {
		*timeDescs = append(*timeDescs, *timeDesc.Clone())
	}
	return timeDescs
}
