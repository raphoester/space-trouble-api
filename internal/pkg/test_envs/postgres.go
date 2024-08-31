package test_envs

import (
	"fmt"
	"path/filepath"

	"github.com/ory/dockertest"
	"github.com/raphoester/space-trouble-api/internal/pkg/basicutil"
	"github.com/raphoester/space-trouble-api/internal/pkg/postgres"
)

type Postgres struct {
	PG        *postgres.Postgres
	container *dockertest.Resource
	pool      *dockertest.Pool
}

func (p *Postgres) Destroy() error {
	return p.pool.Purge(p.container)
}

func (p *Postgres) Clean() error {
	_, err := p.PG.DB.Exec("TRUNCATE TABLE bookings CASCADE;")
	return err
}

func NewPostgres() (*Postgres, error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, fmt.Errorf("failed creating dockertest pool: %w", err)
	}

	credentialsBase := "example"
	container, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Name:       credentialsBase,
		Tag:        "15.4-alpine3.18",
		Env: []string{
			fmt.Sprintf("POSTGRES_USER=%s", credentialsBase),
			fmt.Sprintf("POSTGRES_PASSWORD=%s", credentialsBase),
			fmt.Sprintf("POSTGRES_DB=%s", credentialsBase),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed running container: %w", err)
	}

	pgPort := container.GetPort("5432/tcp")

	dsn := fmt.Sprintf(
		"host=localhost user=example password=example dbname=example port=%s sslmode=disable", pgPort)

	var db *postgres.Postgres
	// Wait for the container to be ready and get its connection details.
	err = pool.Retry(func() error {
		db, err = postgres.New(dsn)
		if err != nil {
			return err
		}

		if _, err = db.DB.Exec(`SELECT 1`); err != nil {
			return fmt.Errorf("failed to ping db: %w", err)
		}

		return nil
	})

	projectRoot, err := basicutil.FindProjectRoot()
	if err != nil {
		return nil, fmt.Errorf("failed finding project root: %w", err)
	}

	if err = db.Migrate(filepath.Join(projectRoot, "sql")); err != nil {
		return nil, fmt.Errorf("failed running migrations: %w", err)
	}

	return &Postgres{
		PG:        db,
		container: container,
		pool:      pool,
	}, nil
}
