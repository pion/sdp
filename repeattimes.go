package sdp

type RepeatTimes []RepeatTime

func (r *RepeatTimes) Clone() *RepeatTimes {
	repeatTimes := &RepeatTimes{}
	for _, repeatTime := range *r {
		*repeatTimes = append(*repeatTimes, *repeatTime.Clone())
	}
	return repeatTimes
}
