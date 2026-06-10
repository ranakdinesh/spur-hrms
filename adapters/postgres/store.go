package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"y/core/domain"
)

type Store struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Store {
	return &Store{pool: pool}
}

func (s *Store) Create(ctx context.Context, e *domain.Hrms) (*domain.Hrms, error) {
	// TODO: use sqlc generated function after running: sqlc generate
	// Example:
	// q := sqlc.New(s.pool)
	// row, err := q.CreateHrms(ctx, sqlc.CreateHrmsParams{
	//     ID:       e.ID,
	//     TenantID: e.TenantID,
	// })
	// if err != nil { return nil, fmt.Errorf("hrms: create: %w", err) }
	// return mapRow(row), nil
	_ = fmt.Sprintf // suppress import error
	return e, nil
}

func (s *Store) GetByID(ctx context.Context, id uuid.UUID, tenantID uuid.UUID) (*domain.Hrms, error) {
	// TODO: implement
	return nil, nil
}

func (s *Store) List(ctx context.Context, tenantID uuid.UUID) ([]*domain.Hrms, error) {
	// TODO: implement
	return []*domain.Hrms{}, nil
}

func (s *Store) Delete(ctx context.Context, id uuid.UUID, tenantID uuid.UUID) error {
	// TODO: implement
	return nil
}
