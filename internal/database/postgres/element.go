package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/x0tf/server/internal/shared"
)

// elementService represents the postgres element service implementation
type elementService struct {
	pool *pgxpool.Pool
}

// Count counts the total amount of elements
func (service *elementService) Count() (int, error) {
	query := "SELECT COUNT(*) FROM elements"

	row := service.pool.QueryRow(context.Background(), query)

	var count int
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

// Element retrieves an element with a specific key out of a specific namespace
func (service *elementService) Element(namespace, key string) (*shared.Element, error) {
	query := "SELECT * FROM elements WHERE namespace = $1 AND key = $2"

	element, err := rowToElement(service.pool.QueryRow(context.Background(), query, namespace, key))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return element, nil
}

// Elements retrieves a list of elements using the given limit and offset
func (service *elementService) Elements(limit, offset int) ([]*shared.Element, error) {
	query := fmt.Sprintf("SELECT * FROM elements ORDER BY created LIMIT %d OFFSET %d", limit, offset)

	rows, err := service.pool.Query(context.Background(), query)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*shared.Element{}, nil
		}
		return nil, err
	}

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

// ElementsInNamespace retrieves a list of elements inside a specific namespace using the given limit and offset
func (service *elementService) ElementsInNamespace(namespace string, limit, offset int) ([]*shared.Element, error) {
	query := fmt.Sprintf("SELECT * FROM elements WHERE namespace = $1 ORDER BY created LIMIT %d OFFSET %d", limit, offset)

	rows, err := service.pool.Query(context.Background(), query, namespace)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*shared.Element{}, nil
		}
		return nil, err
	}

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
func (service *elementService) CreateOrReplace(element *shared.Element) error {
	query := `
		INSERT INTO elements (namespace, key, type, internal_data, public_data, views, max_views, valid_from, valid_until, created)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (namespace, key) DO UPDATE
			SET type = excluded.type,
				internal_data = excluded.internal_data,
				public_data = excluded.public_data,
				views = excluded.views,
				max_views = excluded.max_views,
				valid_from = excluded.valid_from,
				valid_until = excluded.valid_until,
				created = excluded.created
	`

	_, err := service.pool.Exec(
		context.Background(),
		query,
		element.Namespace,
		element.Key,
		element.Type,
		element.InternalData,
		element.PublicData,
		element.Views,
		element.MaxViews,
		element.ValidFrom,
		element.ValidUntil,
		element.Created,
	)
	return err
}

// Delete deletes an element with a specific key out of a specific namespace
func (service *elementService) Delete(namespace, key string) error {
	query := "DELETE FROM elements WHERE namespace = $1 AND key = $2"

	_, err := service.pool.Exec(context.Background(), query, namespace, key)
	return err
}

// DeleteInNamespace deletes all elements out of a specific namespace
func (service *elementService) DeleteInNamespace(namespace string) error {
	query := "DELETE FROM elements WHERE namespace = $1"

	_, err := service.pool.Exec(context.Background(), query, namespace)
	return err
}

// rowToElement reads a pgx row into an element instance
func rowToElement(row pgx.Row) (*shared.Element, error) {
	var namespace string
	var key string
	var typ shared.ElementType
	var internalData map[string]interface{}
	var publicData map[string]interface{}
	var views int
	var maxViews int
	var validFrom int64
	var validUntil int64
	var created int64

	if err := row.Scan(&namespace, &key, &typ, &internalData, &publicData, &views, &maxViews, &validFrom, &validUntil, &created); err != nil {
		return nil, err
	}

	return &shared.Element{
		Namespace:    namespace,
		Key:          key,
		Type:         typ,
		InternalData: internalData,
		PublicData:   publicData,
		Views:        views,
		MaxViews:     maxViews,
		ValidFrom:    validFrom,
		ValidUntil:   validUntil,
		Created:      created,
	}, nil
}
