package generator

import (
	"errors"
	"fmt"
	"log"
	"math"
	"regexp"
)

const (
	CmdGet = iota
	CmdStop
	CmdCount
	CmdSet
)

var (
	ErrorEndSequence = errors.New("end codes")
)

type Data struct {
	Cmd      int
	Code     string
	ChanCode chan *Result
}

type Result struct {
	Code  []byte
	Error error
}

type Generator struct {
	lastCode []byte
	start    bool
}

func Make(code string) (*Generator, error) {
	gen := new(Generator)
	err := gen.set(code)
	return gen, err
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

func (g *Generator) set(code string) error {
	if code == "" {
		g.lastCode = []byte("0000")
		g.start = true
		return nil
	}

	ok, err := regexp.MatchString(`^\w{4}$`, code)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("bad code: %s", code)
	}

	g.lastCode = []byte(code)
	g.start = false
	return nil
}

func (g *Generator) nextCode() ([]byte, bool) {
	if g.start {
		g.start = false
		return g.lastCode, true
	}

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

func countChar(b byte) (count float64) {
	switch {
	case b >= '\x30' && b <= '\x39':
		count = float64(b - '\x30')
	case b >= '\x41' && b <= '\x4A':
		count = float64(b - '\x41' + 10)
	default:
		count = float64(b - '\x61' + 36)
	}
	return
}

func (g *Generator) freeCount() uint32 {
	var count float64 = math.Pow(62.0, 4.0) - 1
	if !g.start {
		for i, c := range g.lastCode {
			count -= math.Pow(62.0, 3.0-float64(i)) * countChar(c)
		}
	}
	return uint32(count)
}

func (g *Generator) Start(ch chan *Data) {
	log.Println("Generator started...")
	go func() {
		var ok bool
		for d := range ch {
			res := new(Result)
			switch d.Cmd {
			case CmdGet:
				if res.Code, ok = g.nextCode(); !ok {
					res.Error = ErrorEndSequence
				}
				d.ChanCode <- res
			case CmdCount:
				res.Code = []byte(fmt.Sprintf("%d", g.freeCount()))
				d.ChanCode <- res
			case CmdSet:
				res.Error = g.set(d.Code)
				d.ChanCode <- res
			case CmdStop:
				log.Println("Generator stop")
				d.ChanCode <- res
				return
			}
		}
	}()
}
