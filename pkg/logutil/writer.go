package logutil

import (
	"bufio"
	"strings"
)

type Writer func(...interface{})

func (l Writer) Write(b []byte) (int, error) {

	var offset int

	for {
		advance, ln, err := bufio.ScanLines(b[offset:], true)
		if err != nil {
			return offset, err
		}

		msg := strings.TrimSpace(string(ln))

		if len(msg) > 0 {
			l(msg)
		}

		offset += advance
		if offset >= len(b) {
			break
		}
	}

	return len(b), nil
}
