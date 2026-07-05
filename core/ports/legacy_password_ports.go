package ports

import (
	"context"

	"github.com/google/uuid"
)

// LegacyPasswordMigrationPort is implemented by the identity module/host app.
// HRMS must not look up or mutate identity-owned auth tables directly.
type LegacyPasswordMigrationPort interface {
	VerifyAndMigrateLegacyPassword(ctx context.Context, cmd LegacyPasswordMigrationCommand) (*LegacyPasswordMigrationResult, error)
}

type LegacyPasswordMigrationCommand struct {
	TenantID   uuid.UUID
	Identifier string
	Password   string
	ActorID    *uuid.UUID
}

type LegacyPasswordMigrationResult struct {
	Migrated              bool       `json:"migrated"`
	UserID                *uuid.UUID `json:"user_id,omitempty"`
	RequiresPasswordReset bool       `json:"requires_password_reset,omitempty"`
	Message               string     `json:"message,omitempty"`
}
