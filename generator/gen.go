package generator

import (
	"fmt"
	"regexp"
)

type Generator struct {
	lastCode []byte
}

func Make(code string) (*Generator, error) {
	if code == "" {
		return &Generator{
			lastCode: []byte("0000"),
		}, nil
	}

	ok, err := regexp.MatchString(`^\w{4}$`, code)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("bad code: %s", code)
	}

	return &Generator{
		lastCode: []byte(code),
	}, nil
}

func (g *Generator) LastCode() []byte {
	return g.lastCode
}

func addChar(b byte) (byte, bool) {
	switch b {
	case '\x39':
		return '\x41', false
	case '\x5a':
		return '\x61', false
	case '\x7a':
		return '\x30', true
	}
	return byte(int(b) + 1), false
}

func (g *Generator) NextCode() ([]byte, bool) {
	var (
		next bool
		i    int
	)

	for i = 3; i >= 0; i-- {
		g.lastCode[i], next = addChar(g.lastCode[i])
		if !next {
			break
		}
	}
	if i < 0 {
		return g.lastCode, false
	}
	return g.lastCode, true
}
