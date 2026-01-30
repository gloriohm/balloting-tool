package io

import (
	"bufio"
	"io"
	"strings"
)

func detectDelimiter(r io.Reader) (rune, *bufio.Reader, error) {
	br := bufio.NewReader(r)

	peek, err := br.Peek(4096)
	if err != nil && err != io.EOF {
		return 0, nil, err
	}

	line := firstNonEmptyLine(peek)
	if line == "" {
		return ',', br, nil
	}

	commas := countRune(line, ',')
	semis := countRune(line, ';')

	if semis > commas {
		return ';', br, nil
	}

	return ',', br, nil
}

func firstNonEmptyLine(b []byte) string {
	start := 0
	for i, c := range b {
		if c == '\n' {
			line := strings.TrimSpace(string(b[start:i]))
			if line != "" {
				return line
			}
			start = i + 1
		}
	}
	return strings.TrimSpace(string(b))
}

func countRune(s string, r rune) int {
	n := 0
	for _, c := range s {
		if c == r {
			n++
		}
	}
	return n
}
