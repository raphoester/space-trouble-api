package id

import "github.com/google/uuid"

type RandomGenerator struct {
}

func (r *RandomGenerator) Generate() ID {
	return ID(uuid.New().String())
}
