package id

func Parse(id string) ID {
	return ID(id)
}

type ID string

func (i ID) String() string {
	return string(i)
}

type Generator interface {
	Generate() ID
}
