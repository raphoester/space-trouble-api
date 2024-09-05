package id

import (
	"github.com/raphoester/chaos"
)

type ChaoticGenerator struct {
	calls int
	seed  string
}

func NewChaoticFactory(seed string) *ChaoticGenerator {
	return &ChaoticGenerator{seed: seed}
}

func (g *ChaoticGenerator) Generate() ID {
	g.calls++
	id := chaos.UUID(g.seed, g.calls)
	return ID(id.String())
}
