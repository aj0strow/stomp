package stomp

import (
	"bytes"
	"errors"
	"sort"
	"strings"
)

var InvalidFrame = errors.New("stomp: invalid frame")

type Frame struct {
	Command string
	Headers map[string]string
	Body    []byte
}

// Marshal structured frame into bytes.
func MarshalFrame(frame *Frame) []byte {
	buf := &bytes.Buffer{}
	buf.WriteString(frame.Command)
	buf.WriteString("\n")

	// Sort header names in alphabetic order.
	names := make([]string, len(frame.Headers))
	i := 0
	for name := range frame.Headers {
		names[i] = name
		i++
	}
	sort.Strings(names)

	for _, name := range names {
		buf.WriteString(EncodeHeader(name))
		buf.WriteString(": ")
		buf.WriteString(EncodeHeader(frame.Headers[name]))
		buf.WriteString("\n")
	}
	buf.WriteString("\n")
	buf.Write(frame.Body)
	buf.WriteRune(0)
	return buf.Bytes()
}

// Unmarshal input bytes into a structured STOMP frame. 
func UnmarshalFrame(input []byte) (*Frame, error) {
	buf := bytes.NewBuffer(input)
	command, err := buf.ReadBytes('\n')
	if err != nil {
		return nil, InvalidFrame
	}
	command = bytes.TrimSpace(command)
	headers := map[string]string{}
	for {
		header, err := buf.ReadBytes('\n')
		if err != nil {
			return nil, InvalidFrame
		}
		if len(header) == 1 {
			break
		}
		parts := bytes.SplitN(header, []byte(":"), 2)
		if len(parts) != 2 {
			return nil, InvalidFrame
		}
		name := bytes.TrimSpace(parts[0])
		value := bytes.TrimSpace(parts[1])
		headers[DecodeHeader(string(name))] = DecodeHeader(string(value))
	}
	body, err := buf.ReadBytes(0)
	if err != nil {
		return nil, InvalidFrame
	}
	body = bytes.TrimSuffix(body, []byte{0})
	frame := &Frame{
		Command: string(command),
		Headers: headers,
		Body:    body,
	}
	return frame, nil
}

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
