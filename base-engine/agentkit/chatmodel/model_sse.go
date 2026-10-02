package chatmodel

import (
	"bufio"
	"bytes"
	"io"
	"strings"
	"unicode/utf8"
)

const maxSSELineBytes = 1024 * 1024

func parseModelSSE(reader io.Reader, onData func([]byte) error) error {
	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, 4096), maxSSELineBytes)
	var data [][]byte
	total := 0
	for scanner.Scan() {
		line := scanner.Bytes()
		total += len(line) + 1
		if total > maxModelResponseBytes || !utf8.Valid(line) {
			return ErrModelProtocol
		}
		if len(line) == 0 {
			if len(data) == 0 {
				continue
			}
			payload := bytes.Join(data, []byte("\n"))
			data = nil
			if strings.TrimSpace(string(payload)) == "[DONE]" {
				return nil
			}
			if err := onData(payload); err != nil {
				return err
			}
			continue
		}
		if bytes.HasPrefix(line, []byte("data:")) {
			value := bytes.TrimPrefix(line, []byte("data:"))
			data = append(data, bytes.Clone(bytes.TrimPrefix(value, []byte(" "))))
		}
	}
	if scanner.Err() != nil {
		return ErrModelProtocol
	}
	return ErrModelProtocol
}
