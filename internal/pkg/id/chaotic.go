package id

import (
	"github.com/raphoester/chaos"
)

type ChaoticGenerator struct {
	chaos *chaos.Chaos
}

func NewChaoticFactory(seed string) *ChaoticGenerator {
	return &ChaoticGenerator{
		chaos: chaos.New(seed),
	}
}

func (g *ChaoticGenerator) Generate() ID {
	id := g.chaos.UUID()
	return ID(id.String())
}
