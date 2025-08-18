package headers

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

var sep = []byte("\r\n")
var ErrorBadHeader = fmt.Errorf("bad header (no registered nurse)")
var ErrorBadFL = fmt.Errorf("bad field line")
var ErrorBadKey = fmt.Errorf("bad key")

var headerKeyRe = regexp.MustCompile(`^[A-Za-z0-9!#$%&'*+\-.\^_` + "`" + `|~]+$`)

type Headers map[string]string

func (h Headers) Get(key string) (string, bool) {
	val, ok := h[strings.ToLower(key)]
	return val, ok
}

func (h Headers) Set(key, val string) {
	if _, ok := h[strings.ToLower(key)]; ok {
		h[strings.ToLower(key)] = h[strings.ToLower(key)] + ", " + val
	} else {
		h[strings.ToLower(key)] = val
	}
}

func (h Headers) Replace(key, val string) {
	h[strings.ToLower(key)] = val
}
func validKey(key []byte) bool {
	return headerKeyRe.Match(key) && (len(key) >= 1) && !bytes.HasSuffix(key, []byte(" "))
}

func (h Headers) ForEach(f func(key, value string)) {
	for key, value := range h {
		f(key, value)
	}
}
func (h Headers) Parse(data []byte) (int, bool, error) {
	read := 0
	done := false

	for {
		idx := bytes.Index(data, sep)
		if idx == -1 {
			// incomplete line, stop
			break
		}

		// empty line -> headers finished
		if idx == 0 {
			read += len(sep)
			done = true
			break
		}

		line := data[:idx]
		parts := bytes.SplitN(line, []byte(":"), 2)
		if len(parts) != 2 {
			return 0, false, ErrorBadFL
		}

		key := bytes.TrimLeft(parts[0], " ")
		value := bytes.TrimSpace(parts[1])

		if !validKey(key) {
			return 0, false, ErrorBadKey
		}

		// overwrite instead of append
		h.Set(string(key), string(value))

		// advance the slice to the next line
		data = data[idx+len(sep):]
		read += idx + len(sep)
	}

	return read, done, nil
}

func NewHeaders() Headers {
	return make(Headers)
}
