package sdp

import (
	"bufio"
	"fmt"
	"io"
	"strconv"

	"github.com/pkg/errors"

	"github.com/pions/webrtc/pkg/rtcerr"
)

const (
	pkgName       = "sdp"
	endline       = "\r\n"
	versionKey    = "v="
	originKey     = "o="
	sessionKey    = "s="
	infoKey       = "i="
	uriKey        = "u="
	emailKey      = "e="
	phoneKey      = "p="
	connectionKey = "c="
	bandwidthKey  = "b="
	timingKey     = "t="
	repeatTimeKey = "r="
	timeZonesKey  = "z="
	encryptionKey = "k="
	attributeKey  = "a="
	mediaKey      = "m="
)

const (
	unknown    = iota
	unknownStr = "unknown"
)

type stringer struct {
	Value string
}

func (s stringer) String() string {
	return s.Value
}

type lexer struct {
	desc  *SessionDescription
	input *bufio.Reader
}

type stateFn func(*lexer) (stateFn, error)

func readType(input *bufio.Reader) (string, error) {
	key, err := input.ReadString('=')
	if err != nil {
		return key, errors.Wrap(err, pkgName)
	}

	if len(key) != 2 {
		return key, errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", key)}, pkgName)
	}

	return key, nil
}

func readValue(input *bufio.Reader) (string, error) {
	line, err := input.ReadString('\n')
	if err != nil && err != io.EOF {
		return line, errors.Wrap(&rtcerr.UnknownError{Err: err}, pkgName)
	}

	if len(line) == 0 {
		return line, errors.Wrap(&rtcerr.UnknownError{Err: io.EOF}, pkgName)
	}

	if line[len(line)-1] == '\n' {
		drop := 1
		if len(line) > 1 && line[len(line)-2] == '\r' {
			drop = 2
		}
		line = line[:len(line)-drop]
	}

	return line, nil
}

func indexOf(element string, data []string) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1
}

func parseTimeUnits(value string) (int64, error) {
	// Some time offsets in the protocol can be provided with a shorthand
	// notation. This code ensures to convert it to NTP timestamp format.
	//      d - days (86400 seconds)
	//      h - hours (3600 seconds)
	//      m - minutes (60 seconds)
	//      s - seconds (allowed for completeness)
	switch value[len(value)-1:] {
	case "d":
		num, err := strconv.ParseInt(value[:len(value)-1], 10, 64)
		if err != nil {
			return 0, errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", value)}, pkgName)
		}
		return num * 86400, nil
	case "h":
		num, err := strconv.ParseInt(value[:len(value)-1], 10, 64)
		if err != nil {
			return 0, errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", value)}, pkgName)
		}
		return num * 3600, nil
	case "m":
		num, err := strconv.ParseInt(value[:len(value)-1], 10, 64)
		if err != nil {
			return 0, errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", value)}, pkgName)
		}
		return num * 60, nil
	}

	num, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("%v", value)}, pkgName)
	}

	return num, nil
}

func parsePort(value string) (int, error) {
	port, err := strconv.Atoi(value)
	if err != nil {
		return 0, errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("port %v", port)}, pkgName)
	}

	if port < 0 || port > 65536 {
		return 0, errors.Wrap(&rtcerr.SyntaxError{Err: fmt.Errorf("port %v", port)}, pkgName)
	}

	return port, nil
}
