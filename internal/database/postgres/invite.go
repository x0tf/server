package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/x0tf/server/internal/shared"
)

// inviteService represents the postgres invite service implementation
type inviteService struct {
	pool *pgxpool.Pool
}

// Count counts the total amount of invites
func (service *inviteService) Count() (int, error) {
	query := "SELECT COUNT(*) FROM invites"

	row := service.pool.QueryRow(context.Background(), query)

	var count int
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

// Invite retrieves an invite with a specific code
func (service *inviteService) Invite(code string) (*shared.Invite, error) {
	query := "SELECT * FROM invites WHERE code = $1"

	invite, err := rowToInvite(service.pool.QueryRow(context.Background(), query, code))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return invite, nil
}

// Invites retrieves a list of invites using the given limit and offset
func (service *inviteService) Invites(limit, offset int) ([]*shared.Invite, error) {
	query := fmt.Sprintf("SELECT * FROM invites ORDER BY created LIMIT %d OFFSET %d", limit, offset)

	rows, err := service.pool.Query(context.Background(), query)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*shared.Invite{}, nil
		}
		return nil, err
	}

	var invites []*shared.Invite
	for rows.Next() {
		invite, err := rowToInvite(rows)
		if err != nil {
			return nil, err
		}
		invites = append(invites, invite)
	}

	return invites, nil
}

// CreateOrReplace creates or replaces an invite
func (service *inviteService) CreateOrReplace(invite *shared.Invite) error {
	query := `
		INSERT INTO invites (code, uses, max_uses, created)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE
			SET uses = excluded.uses,
				max_uses = excluded.max_uses,
				created = excluded.created
	`

	_, err := service.pool.Exec(context.Background(), query, invite.Code, invite.Uses, invite.MaxUses, invite.Created)
	return err
}

// Delete deletes an invite with a specific code
func (service *inviteService) Delete(code string) error {
	query := "DELETE FROM invites WHERE code = $1"

	_, err := service.pool.Exec(context.Background(), query, code)
	return err
}

// rowToInvite reads a pgx row into an invite instance
func rowToInvite(row pgx.Row) (*shared.Invite, error) {
	var code string
	var uses int
	var maxUses int
	var created int64

	if err := row.Scan(&code, &uses, &maxUses, &created); err != nil {
		return nil, err
	}

	return &shared.Invite{
		Code:    code,
		Uses:    uses,
		MaxUses: maxUses,
		Created: created,
	}, nil
}
