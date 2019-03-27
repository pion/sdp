package sdp

type MediaDescriptions []MediaDescription

func (m *MediaDescriptions) Clone() *MediaDescriptions {
	mediaDescs := &MediaDescriptions{}
	for _, mediaDesc := range *m {
		*mediaDescs = append(*mediaDescs, *mediaDesc.Clone())
	}
	return mediaDescs
}

func (m *MediaDescriptions) Add(mediaDesc *MediaDescription) {
	if mediaDesc == nil {
		return
	}
	*m = append(*m, *mediaDesc)
}

func (m *MediaDescriptions) FirstActive() *MediaDescription {
	for _, mediaDesc := range *m {
		if !mediaDesc.Rejected() {
			return &mediaDesc
		}
	}
	return nil
}

func (m *MediaDescriptions) FirstActiveByType(mediaType MediaType) *MediaDescription {
	mediaDescs := map[MediaType]*MediaDescription{
		MediaTypeAudio:       nil,
		MediaTypeVideo:       nil,
		MediaTypeApplication: nil,
		MediaTypeText:        nil,
		MediaTypeMessage:     nil,
	}
	for _, mediaDesc := range *m {
		if !mediaDesc.Rejected() && mediaDescs[mediaDesc.Media.Type] == nil {
			mediaDescs[mediaDesc.Media.Type] = &mediaDesc
		}
	}
	return mediaDescs[mediaType]
}
