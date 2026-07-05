package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/db"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/internal/logging"
	"github.com/rs/zerolog"
)

type Store struct {
	pool *pgxpool.Pool
	q    *sqlc.Queries
	log  *zerolog.Logger
}

func New(pool *pgxpool.Pool, log ...*zerolog.Logger) *Store {
	return &Store{pool: pool, q: sqlc.New(pool), log: logging.Component(logging.First(log...), "postgres")}
}

func (s *Store) RunAsSystem(ctx context.Context, fn func(context.Context) error) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return s.logDBError(ctx, "begin system transaction", fmt.Errorf("hrms: begin system transaction: %w", err))
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, "SELECT set_config('app.tenant_id', '', true), set_config('app.is_super_admin', 'true', true)"); err != nil {
		return s.logDBError(ctx, "set system rls context", fmt.Errorf("hrms: set system rls context: %w", err))
	}

	if err := fn(db.WithTx(ctx, tx)); err != nil {
		s.logDBError(ctx, "run system transaction callback", err)
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return s.logDBError(ctx, "commit system transaction", fmt.Errorf("hrms: commit system transaction: %w", err))
	}
	return nil
}

func (s *Store) getQueries(ctx context.Context) *sqlc.Queries {
	if tx := db.GetTx(ctx); tx != nil {
		return s.q.WithTx(tx)
	}
	return s.q
}

type logField func(*zerolog.Event)

func tenantIDField(tenantID uuid.UUID) logField {
	return func(event *zerolog.Event) {
		if tenantID != uuid.Nil {
			event.Str("tenant_id", tenantID.String())
		}
	}
}

func optionalTenantIDField(tenantID *uuid.UUID) logField {
	return func(event *zerolog.Event) {
		if tenantID != nil && *tenantID != uuid.Nil {
			event.Str("tenant_id", tenantID.String())
		}
	}
}

func stringField(key string, value string) logField {
	return func(event *zerolog.Event) {
		if value != "" {
			event.Str(key, value)
		}
	}
}

func (s *Store) logDBError(ctx context.Context, operation string, err error, fields ...logField) error {
	if s != nil && s.log != nil && err != nil {
		event := s.log.Error().Err(err).Str("operation", operation)
		if db.GetTx(ctx) != nil {
			event.Bool("transaction_scoped", true)
		}
		for _, field := range fields {
			if field != nil {
				field(event)
			}
		}
		event.Msg("hrms postgres error")
	}
	return err
}
