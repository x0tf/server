package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/x0tf/server/internal/shared"
	"github.com/x0tf/server/internal/static"
)

// NamespaceService represents the postgres namespace service
type NamespaceService struct {
	pool *pgxpool.Pool
}

// NewNamespaceService creates a new postgres namespace service
func NewNamespaceService(dsn string) (*NamespaceService, error) {
	// Open a postgres connection pool
	pool, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	// Create and return the namespace service
	return &NamespaceService{
		pool: pool,
	}, nil
}

// InitializeTable initializes the namespace table
func (service *NamespaceService) InitializeTable() error {
	query := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id VARCHAR(32) NOT NULL,
			token VARCHAR(100) NOT NULL,
			active BOOLEAN NOT NULL,
			PRIMARY KEY (id)
		)
    `, static.PostgresTableNamespaces)
	_, err := service.pool.Exec(context.Background(), query)
	return err
}

// Namespace searches for a namespace by its ID
func (service *NamespaceService) Namespace(sourceID string) (*shared.Namespace, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", static.PostgresTableNamespaces)
	namespace, err := rowToNamespace(service.pool.QueryRow(context.Background(), query, sourceID))
	if err != nil {
		return nil, err
	}
	return namespace, nil
}

// Namespaces searches for all existent namespaces
func (service *NamespaceService) Namespaces() ([]*shared.Namespace, error) {
	query := fmt.Sprintf("SELECT * FROM %s", static.PostgresTableNamespaces)
	rows, err := service.pool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var namespaces []*shared.Namespace
	for rows.Next() {
		namespace, err := rowToNamespace(rows)
		if err != nil {
			return nil, err
		}
		namespaces = append(namespaces, namespace)
	}
	return namespaces, nil
}

// CreateOrReplace creates or replaces a namespace
func (service *NamespaceService) CreateOrReplace(namespace *shared.Namespace) error {
	query := fmt.Sprintf(`
		INSERT INTO %s (id, token, active)
		VALUES ($1, $2, $3)
		ON CONFLICT (id) DO UPDATE
			SET token = excluded.token,
				active = excluded.active
    `, static.PostgresTableNamespaces)
	_, err := service.pool.Exec(context.Background(), query, namespace.ID, namespace.Token, namespace.Active)
	return err
}

// Delete deletes a namespace
func (service *NamespaceService) Delete(id string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", static.PostgresTableNamespaces)
	_, err := service.pool.Exec(context.Background(), query, id)
	return err
}

// Close closes the postgres namespace service
func (service *NamespaceService) Close() {
	service.pool.Close()
}

// rowToNamespace creates a namespace from a postgres row
func rowToNamespace(row pgx.Row) (*shared.Namespace, error) {
	var id string
	var token string
	var active bool

	err := row.Scan(&id, &token, &active)
	if err != nil {
		return nil, err
	}

	return &shared.Namespace{
		ID:     id,
		Token:  token,
		Active: active,
	}, nil
}
