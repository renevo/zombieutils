package logutil

import (
	"bufio"
	"log/slog"
	"regexp"
	"strings"
)

type Writer struct {
	IsErr bool
}

// 2024-07-31T14:33:47 54.462 INF

var (
	timestampRegex = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\s\d+\.\d+\s`)
)

func (l Writer) Write(b []byte) (int, error) {
	var offset int

	outputFn := slog.Info
	if l.IsErr {
		outputFn = slog.Error
	}

	for {
		advance, ln, err := bufio.ScanLines(b[offset:], true)
		if err != nil {
			return offset, err
		}
		lineOutputFn := outputFn
		msg := strings.TrimSpace(string(ln))

		// remove the timestamp from the msgs
		msg = timestampRegex.ReplaceAllString(msg, "")

		if len(msg) > 0 {
			// log level conversions
			switch {
			case strings.HasPrefix(msg, "INF:"):
				lineOutputFn = slog.Info
				msg = strings.TrimSpace(strings.TrimPrefix(msg, "INF:"))

			case strings.HasPrefix(msg, "INF "):
				lineOutputFn = slog.Info
				msg = strings.TrimSpace(strings.TrimPrefix(msg, "INF"))

			case strings.HasPrefix(msg, "WRN:"):
				lineOutputFn = slog.Warn
				msg = strings.TrimSpace(strings.TrimPrefix(msg, "WRN:"))

			case strings.HasPrefix(msg, "WRN "):
				lineOutputFn = slog.Warn
				msg = strings.TrimSpace(strings.TrimPrefix(msg, "WRN"))

			case strings.HasPrefix(msg, "ERR:"):
				lineOutputFn = slog.Error
				msg = strings.TrimSpace(strings.TrimPrefix(msg, "ERR:"))

			case strings.HasPrefix(msg, "ERR "):
				lineOutputFn = slog.Error
				msg = strings.TrimSpace(strings.TrimPrefix(msg, "ERR"))

			case strings.HasPrefix(msg, "WARNING:"):
				lineOutputFn = slog.Warn
				msg = strings.TrimSpace(strings.TrimPrefix(msg, "WARNING:"))

			case strings.HasPrefix(msg, "ERROR:"):
				lineOutputFn = slog.Warn
				msg = strings.TrimSpace(strings.TrimPrefix(msg, "ERROR:"))
			}

			// noisy shit that doesn't matter for our server
			if strings.HasPrefix(msg, "Shader") ||
				strings.HasPrefix(msg, "#pragma") ||
				strings.Contains(msg, "There is no texture data available to upload.") ||
				strings.Contains(msg, ", you may have forgotten turning Fallback off?") ||
				strings.Contains(msg, "shader is not supported on this GPU") ||
				strings.HasPrefix(msg, "Did you use #pragma only_renderers") ||
				strings.HasPrefix(msg, "Unsupported: '") ||
				strings.HasPrefix(msg, "No mesh data available for mesh") ||
				strings.HasPrefix(msg, "Couldn't create a Convex Mesh from source mesh") ||
				strings.HasPrefix(msg, "Texture with ID ") ||
				strings.EqualFold(msg, ": No such file or directory") ||
				strings.EqualFold(msg, "-- type 'quit' to exit --") {
				goto ADVANCE
			}

			lineOutputFn(msg)
		}

	ADVANCE:
		offset += advance
		if offset >= len(b) {
			break
		}
	}

	return len(b), nil
}
