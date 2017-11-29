package stomp

import (
	"strings"
)

// Encode the header name or value.
func EncodeHeader(input string) string {
	r := strings.NewReplacer(
		"\r", "\\r",
		"\n", "\\n",
		":", "\\c",
		"\\", "\\\\",
	)
	return r.Replace(input)
}

// Decode the header name or value.
func DecodeHeader(input string) string {
	r := strings.NewReplacer(
		"\\r", "\r",
		"\\n", "\n",
		"\\c", ":",
		"\\\\", "\\",
	)
	return r.Replace(input)
}
