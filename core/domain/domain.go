package domain

import (
	"time"

	"github.com/google/uuid"
)

// Hrms is the core domain entity for the hrms module.
// Add your business fields below.
type Hrms struct {
	ID        uuid.UUID `json:"id"`
	TenantID  uuid.UUID `json:"tenant_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// TODO: add your domain fields
}

// NewHrms is the factory function. Validates inputs and returns a ready entity.
func NewHrms(tenantID uuid.UUID) (*Hrms, error) {
	now := time.Now().UTC()
	return &Hrms{
		ID:        uuid.New(),
		TenantID:  tenantID,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
