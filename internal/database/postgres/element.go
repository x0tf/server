package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/x0tf/server/internal/shared"
)

// ElementService represents the postgres element service
type ElementService struct {
	pool *pgxpool.Pool
}

// NewElementService creates a new postgres element service
func NewElementService(dsn string) (*ElementService, error) {
	// Open a postgres connection pool
	pool, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	// Create and return the element service
	return &ElementService{
		pool: pool,
	}, nil
}

// InitializeTable initializes the element table
func (service *ElementService) InitializeTable() error {
	query := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			namespace VARCHAR(32) NOT NULL,
			key VARCHAR(32) NOT NULL,
			type SMALLINT NOT NULL,
			data TEXT NOT NULL,
			PRIMARY KEY (namespace, key)
		)
    `, tableElements)
	_, err := service.pool.Exec(context.Background(), query)
	return err
}

// Element searches for a single element with a specific key in a specific namespace
func (service *ElementService) Element(sourceNamespace, sourceKey string) (*shared.Element, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE namespace = $1 AND key = $2", tableElements)
	element, err := rowToElement(service.pool.QueryRow(context.Background(), query, sourceNamespace, sourceKey))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return element, nil
}

// Elements searches for all elements
func (service *ElementService) Elements() ([]*shared.Element, error) {
	query := fmt.Sprintf("SELECT * FROM %s", tableElements)
	rows, err := service.pool.Query(context.Background(), query)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	var elements []*shared.Element
	for rows.Next() {
		element, err := rowToElement(rows)
		if err != nil {
			return nil, err
		}
		elements = append(elements, element)
	}
	return elements, nil
}

// ElementsInNamespace searches for all elements in a specific namespace
func (service *ElementService) ElementsInNamespace(namespace string) ([]*shared.Element, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE namespace = $1", tableElements)
	rows, err := service.pool.Query(context.Background(), query, namespace)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	var elements []*shared.Element
	for rows.Next() {
		element, err := rowToElement(rows)
		if err != nil {
			return nil, err
		}
		elements = append(elements, element)
	}
	return elements, nil
}

// CreateOrReplace creates or replaces an element
func (service *ElementService) CreateOrReplace(element *shared.Element) error {
	query := fmt.Sprintf(`
		INSERT INTO %s (namespace, key, type, data)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (namespace, key) DO UPDATE
			SET type = excluded.type,
				data = excluded.data
    `, tableElements)
	_, err := service.pool.Exec(context.Background(), query, element.Namespace, element.Key, element.Type, element.Data)
	return err
}

// Delete deletes an element
func (service *ElementService) Delete(namespace, key string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE namespace = $1 AND key = $2", tableElements)
	_, err := service.pool.Exec(context.Background(), query, namespace, key)
	return err
}

// DeleteInNamespace deletes every element in a namespace
func (service *ElementService) DeleteInNamespace(namespace string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE namespace = $1", tableElements)
	_, err := service.pool.Exec(context.Background(), query, namespace)
	return err
}

// Close closes the postgres element service
func (service *ElementService) Close() {
	service.pool.Close()
}

// rowToElement creates an element from a postgres row
func rowToElement(row pgx.Row) (*shared.Element, error) {
	var namespace string
	var key string
	var typ shared.ElementType
	var data string

	err := row.Scan(&namespace, &key, &typ, &data)
	if err != nil {
		return nil, err
	}

	return &shared.Element{
		Namespace: namespace,
		Key:       key,
		Type:      typ,
		Data:      data,
	}, nil
}
