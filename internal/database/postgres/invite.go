package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/x0tf/server/internal/shared"
	"strings"
)

// InviteService represents the postgres invite service
type InviteService struct {
	pool *pgxpool.Pool
}

// NewInviteService creates a new postgres invite service
func NewInviteService(dsn string) (*InviteService, error) {
	// Open a postgres connection pool
	pool, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	// Create and return the invite service
	return &InviteService{
		pool: pool,
	}, nil
}

// InitializeTable initializes the invite table
func (service *InviteService) InitializeTable() error {
	query := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			token VARCHAR(32) NOT NULL,
			PRIMARY KEY (token)
		)
    `, tableInvites)
	_, err := service.pool.Exec(context.Background(), query)
	return err
}

// IsValid searches for a single invite with a specific token
func (service *InviteService) IsValid(token string) (bool, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE token = $1", tableInvites)
	if err := service.pool.QueryRow(context.Background(), query, token).Scan(nil); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Invites searches for all invites
func (service *InviteService) Invites() ([]shared.Invite, error) {
	query := fmt.Sprintf("SELECT * FROM %s", tableInvites)
	rows, err := service.pool.Query(context.Background(), query)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	var invites []shared.Invite
	for rows.Next() {
		var invite shared.Invite
		if err = rows.Scan(&invite); err != nil {
			return nil, err
		}
		invites = append(invites, invite)
	}
	return invites, nil
}

// Create creates an invite
func (service *InviteService) Create(invite shared.Invite) error {
	query := fmt.Sprintf("INSERT INTO %s (token) VALUES ($1)", tableInvites)
	_, err := service.pool.Exec(context.Background(), query, invite)
	if err != nil && strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
		err = nil
	}
	return err
}

// Delete deletes an invite
func (service *InviteService) Delete(token string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE token = $1", tableInvites)
	_, err := service.pool.Exec(context.Background(), query, token)
	return err
}

// Close closes the postgres invite service
func (service *InviteService) Close() {
	service.pool.Close()
}
