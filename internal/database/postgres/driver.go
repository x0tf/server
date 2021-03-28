package postgres

import (
	"context"
	"embed"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/johejo/golang-migrate-extra/source/iofs"
)

//go:embed migrations/*.sql
var migrations embed.FS

// postgresDriver represents the postgres database driver
type postgresDriver struct {
	dsn        string
	pool       *pgxpool.Pool
	Namespaces *namespaceService
	Elements   *elementService
	Invites    *inviteService
}

// NewDriver creates a new postgres database driver
func NewDriver(dsn string) (*postgresDriver, error) {
	// Open a postgres connection pool
	pool, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	// Create and return the postgres driver
	return &postgresDriver{
		dsn:  dsn,
		pool: pool,
		Namespaces: &namespaceService{
			pool: pool,
		},
		Elements: &elementService{
			pool: pool,
		},
		Invites: &inviteService{
			pool: pool,
		},
	}, nil
}

// Migrate runs all migrations on the connected database
func (driver *postgresDriver) Migrate() error {
	source, err := iofs.New(migrations, "migrations")
	if err != nil {
		return err
	}

	migrator, err := migrate.NewWithSourceInstance("iofs", source, driver.dsn)
	if err != nil {
		return err
	}

	return migrator.Up()
}

// Close closes the postgres database driver
func (driver *postgresDriver) Close() {
	driver.pool.Close()
}
