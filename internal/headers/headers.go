package headers

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

const crlf = "\r\n"

type Headers map[string]string

func NewHeaders() Headers {
	return map[string]string{}
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return 0, false, nil
	}
	if idx == 0 {
		// the empty line
		// headers are done, consume the CRLF
		return 2, true, nil
	}

	line := data[:idx]

	parts := bytes.SplitN(line, []byte(":"), 2)
	if len(parts) != 2 {
		return 0, false, fmt.Errorf("malformed header line (missing ':'): %s", string(line))
	}

	keyRaw := string(parts[0])

	if keyRaw != strings.TrimRight(keyRaw, " ") {
		return 0, false, fmt.Errorf("invalid header name: %s", keyRaw)
	}

	re := regexp.MustCompile(`^[A-Za-z0-9!#$%&'*+\-.^_` + "`" + `|~]+$`)

	value := string(bytes.TrimSpace(parts[1]))
	key := strings.ToLower(strings.TrimSpace(keyRaw))

	if !re.MatchString(key) {
		return 0, false, fmt.Errorf("invalid character in header")
	}

	h.Set(key, value)
	return idx + 2, false, nil
}

func (h Headers) Set(key, value string) {
	key = strings.ToLower(key)
	v, ok := h[key]
	if ok {
		value = strings.Join([]string{
			v,
			value,
		}, ", ")
	}

	h[key] = value
}

func (h Headers) Get(key string) (string, error) {
	k := strings.ToLower(key)
	v, ok := h[k]
	if !ok {
		return "", fmt.Errorf("no such key")
	}

	return v, nil
}

func (h Headers) Override(key, value string) {
	key = strings.ToLower(key)
	h[key] = value
}

func (h Headers) Remove(key string) {
	key = strings.ToLower(key)
	delete(h, key)
}
