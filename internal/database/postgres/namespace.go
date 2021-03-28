package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/x0tf/server/internal/shared"
)

// namespaceService represents the postgres namespace service implementation
type namespaceService struct {
	pool *pgxpool.Pool
}

// Namespace retrieves a namespace with a specific ID
func (service *namespaceService) Namespace(id string) (*shared.Namespace, error) {
	query := "SELECT * FROM namespaces WHERE id = $1"

	namespace, err := rowToNamespace(service.pool.QueryRow(context.Background(), query, id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return namespace, nil
}

// Namespaces retrieves a list of namespaces using the given limit and offset
func (service *namespaceService) Namespaces(limit, offset int) ([]*shared.Namespace, error) {
	query := fmt.Sprintf("SELECT * FROM namespaces ORDER BY created LIMIT %d OFFSET %d", limit, offset)

	rows, err := service.pool.Query(context.Background(), query)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*shared.Namespace{}, nil
		}
		return nil, err
	}

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
func (service *namespaceService) CreateOrReplace(namespace *shared.Namespace) error {
	query := `
		INSERT INTO namespaces (id, token, active, created)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE
			SET token = excluded.token,
				active = excluded.active,
				created = excluded.created
	`

	_, err := service.pool.Exec(context.Background(), query, namespace.ID, namespace.Token, namespace.Active, namespace.Created)
	return err
}

// Delete deletes a namespace with a specific ID
func (service *namespaceService) Delete(id string) error {
	query := "DELETE FROM namespaces WHERE id = $1"

	_, err := service.pool.Exec(context.Background(), query, id)
	return err
}

// rowToNamespace reads a pgx row into a namespace instance
func rowToNamespace(row pgx.Row) (*shared.Namespace, error) {
	var id string
	var token string
	var active bool
	var created int64

	if err := row.Scan(&id, &token, &active, &created); err != nil {
		return nil, err
	}

	return &shared.Namespace{
		ID:      id,
		Token:   token,
		Active:  active,
		Created: created,
	}, nil
}
