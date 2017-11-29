package stomp

import (
	"testing"
)

func TestHeaderEncoding(t *testing.T) {
	tests := [][]string{
		{"long header\n\r", "long header\\n\\r"},
		{"windows\\path", "windows\\\\path"},
		{`{"os":"windows"}`, `{"os"\c"windows"}`},
	}
	for _, tt := range tests {
		encoded := EncodeHeader(tt[0])
		if encoded != tt[1] {
			t.Errorf("encoding error:\nin: %s\nhave: %s\nwant: %s\n", tt[0], encoded, tt[1])
		}
		decoded := DecodeHeader(encoded)
		if decoded != tt[0] {
			t.Errorf("decoding error:\nin: %s\nhave: %s\nwant: %s\n", tt[1], decoded, tt[0])
		}
	}
}
