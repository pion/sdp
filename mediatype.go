package sdp

// MediaType describes the <media> section of the mline in media description.
//
// https://tools.ietf.org/html/rfc4566#section-8.2.1
type MediaType int

const (
	// MediaTypeAudio indicates that the media type is audio.
	MediaTypeAudio MediaType = iota + 1

	// MediaTypeVideo indicates that the media type is video
	MediaTypeVideo

	// MediaTypeApplication indicates that the media type in application.
	MediaTypeApplication

	// MediaTypeText indicates that the media type is text.
	MediaTypeText

	// MediaTypeMessage indicates that the media type is message.
	MediaTypeMessage
)

const (
	mediaTypeAudioStr       = "audio"
	mediaTypeVideoStr       = "video"
	mediaTypeApplicationStr = "application"
	mediaTypeTextStr        = "text"
	mediaTypeMessageStr     = "message"
)

var mediaTypes = []string{
	mediaTypeAudioStr,
	mediaTypeVideoStr,
	mediaTypeApplicationStr,
	mediaTypeTextStr,
	mediaTypeMessageStr,
}

// NewMediaType defines a procedure for creating a new media type from a raw
// string.
func NewMediaType(raw string) MediaType {
	switch raw {
	case mediaTypeAudioStr:
		return MediaTypeAudio
	case mediaTypeVideoStr:
		return MediaTypeVideo
	case mediaTypeApplicationStr:
		return MediaTypeApplication
	case mediaTypeTextStr:
		return MediaTypeText
	case mediaTypeMessageStr:
		return MediaTypeMessage
	default:
		return MediaType(unknown)
	}
}

func (t MediaType) String() string {
	switch t {
	case MediaTypeAudio:
		return mediaTypeAudioStr
	case MediaTypeVideo:
		return mediaTypeVideoStr
	case MediaTypeApplication:
		return mediaTypeApplicationStr
	case MediaTypeText:
		return mediaTypeTextStr
	case MediaTypeMessage:
		return mediaTypeMessageStr
	default:
		return unknownStr
	}
}
