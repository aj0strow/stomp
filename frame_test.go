package stomp

import (
	"bytes"
	"fmt"
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

func TestFrameMarshaling(t *testing.T) {
	tests := []struct {
		Frame *Frame
		Bytes []byte
	}{
		{
			&Frame{
				Command: "CONNECT",
				Headers: map[string]string{
					"version":    "1.2",
					"heart-beat": "0,30",
				},
			},
			[]byte("CONNECT\nheart-beat: 0,30\nversion: 1.2\n\n\x00"),
		},
		{
			&Frame{
				Command: "MESSAGE",
				Headers: map[string]string{
					"subscription": "1",
					"destination":  "/quotes/AAPL",
				},
				Body: []byte(`{"price":145.62}`),
			},
			[]byte("MESSAGE\ndestination: /quotes/AAPL\nsubscription: 1\n\n{\"price\":145.62}\x00"),
		},
	}
	for _, tt := range tests {
		msg := MarshalFrame(tt.Frame)
		if !bytes.Equal(msg, tt.Bytes) {
			t.Errorf("bad frame output:\nframe: %# v\nhave: %q\nwant: %q\n", tt.Frame, msg, tt.Bytes)
		}
		frame, err := UnmarshalFrame(tt.Bytes)
		if err != nil {
			t.Fatal(err)
		}
		for _, line := range frameDiff(tt.Frame, frame) {
			t.Errorf(" %s\n", line)
		}
	}
}

func frameDiff(frame *Frame, other *Frame) []string {
	var ds []string
	if frame.Command != other.Command {
		ds = append(ds, fmt.Sprintf("Command: %s != %s", frame.Command, other.Command))
	}
	for name := range frame.Headers {
		if frame.Headers[name] != other.Headers[name] {
			ds = append(ds, fmt.Sprintf("Headers[%s]: %s != %s", name, frame.Headers[name], other.Headers[name]))
		}
	}
	for name := range other.Headers {
		if _, ok := frame.Headers[name]; !ok {
			ds = append(ds, fmt.Sprintf("Headers[%s]: %s != %s", name, frame.Headers[name], other.Headers[name]))
		}
	}
	if !bytes.Equal(frame.Body, other.Body) {
		ds = append(ds, fmt.Sprintf("Body: %q != %q", frame.Body, other.Body))
	}
	return ds
}
