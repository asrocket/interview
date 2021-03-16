package generator

import (
	"math/rand"
	"time"
)

var symbols = []byte{'q', 'w', 'e', 'r', 't', 'y', 'u', 'i', 'o', 'p', 'a', 's', 'd', 'f',
	'g', 'h', 'j', 'k', 'l', 'z', 'x', 'c', 'v', 'b', 'n', 'm'}

type LineGenerator interface {
	Generate() []byte
}

type defaultGenerator struct {
	random *rand.Rand
	maxLen int
	buf    []byte
}

func (g *defaultGenerator) Generate() []byte {
	lineLen := g.random.Intn(g.maxLen)
	g.random.Read(g.buf[:lineLen])
	for j := 0; j < lineLen; j++ {
		g.buf[j] = symbols[int(g.buf[j])%len(symbols)]
	}
	return append(g.buf[:lineLen], '\n')
}

func NewLineGenerator(maxLen int) LineGenerator {
	return &defaultGenerator{
		random: rand.New(rand.NewSource(time.Now().Unix())),
		buf:    make([]byte, maxLen),
		maxLen: maxLen,
	}
}
