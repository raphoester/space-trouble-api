package main

import (
	"os"

	"github.com/raphoester/space-trouble-api/internal/pkg/postgres"
)

func main() {
	dsn := os.Getenv("POSTGRES_DSN")
	pg, err := postgres.New(dsn)
	if err != nil {
		panic(err)
	}

	path := os.Getenv("MIGRATIONS_PATH")
	if err := pg.Migrate(path); err != nil {
		panic(err)
	}
}
