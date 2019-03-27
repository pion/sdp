package sdp

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pions/webrtc/pkg/rtcerr"

	"github.com/pkg/errors"
)

// Media describes the "m=" field storage structure.
type Media struct {
	Type    MediaType
	Port    RangedPort
	Protos  []string
	Formats []string
}

func (m *Media) Clone() *Media {
	media := &Media{}
	media.Type = m.Type
	media.Port = *m.Port.Copy()
	media.Protos = append([]string(nil), m.Protos...)
	media.Formats = append([]string(nil), m.Formats...)
	return media
}

func (m *Media) Unmarshal(raw string) error {
	fields := strings.Fields(raw)
	if len(fields) < 4 {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("m=%v", fields)}, pkgName)
	}

	// <media>
	// Set according to currently registered with IANA
	// https://tools.ietf.org/html/rfc4566#section-5.14
	if i := indexOf(fields[0], mediaTypes); i == -1 {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("port %v", fields[0])}, pkgName)
	}

	// <port>
	var err error
	var port RangedPort
	parts := strings.Split(fields[1], "/")
	port.Value, err = parsePort(parts[0])
	if err != nil {
		return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("port %v", parts[0])}, pkgName)
	}

	if len(parts) > 1 {
		portRange, err := strconv.Atoi(parts[1])
		if err != nil {
			return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", parts)}, pkgName)
		}
		port.Range = &portRange
	}

	// <proto>
	// Set according to currently registered with IANA
	// https://tools.ietf.org/html/rfc4566#section-5.14
	var protos []string
	for _, proto := range strings.Split(fields[2], "/") {
		if i := indexOf(proto, []string{"UDP", "RTP", "AVP", "SAVP", "SAVPF", "TLS", "DTLS", "SCTP"}); i == -1 {
			return errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", fields[2])}, pkgName)
		}
		protos = append(protos, proto)
	}

	// <fmt>...
	var formats []string
	for i := 3; i < len(fields); i++ {
		formats = append(formats, fields[i])
	}

	m.Type = NewMediaType(fields[0])
	m.Port = port
	m.Protos = protos
	m.Formats = formats
	return nil
}

func (m *Media) Marshal() string {
	return mediaKey + strings.Join([]string{
		m.Type.String(),
		m.Port.String(),
		strings.Join(m.Protos, "/"),
		strings.Join(m.Formats, " "),
	}, " ") + endline
}
