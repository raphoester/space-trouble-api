package id

import "github.com/google/uuid"

type FixedIDGenerator struct{}

func (f *FixedIDGenerator) Generate() ID {
	return ID(uuid.MustParse("fea5b3b4-4b3b-4b3b-4b3b-4b3b4b3b4b3b").String())
}
