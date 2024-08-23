package id

import (
	"math/rand"

	"github.com/google/uuid"
	"github.com/raphoester/space-trouble-api/internal/pkg/chaos"
)

type ChaoticGenerator struct {
	rnd *rand.Rand
}

func NewChaoticFactory(seed string) *ChaoticGenerator {
	rnd := rand.New(rand.NewSource(int64(chaos.Int(seed, 100))))
	return &ChaoticGenerator{rnd: rnd}
}

func (g *ChaoticGenerator) Generate() ID {
	id, err := uuid.NewRandomFromReader(g.rnd)
	if err != nil {
		panic(err)
	}
	return ID(id.String())
}
