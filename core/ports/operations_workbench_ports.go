package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

type OperationsWorkbenchRepo interface {
	ListOperationsWorkbenchCards(ctx context.Context, filter domain.OperationsWorkbenchFilter) ([]*domain.OperationsWorkbenchCard, error)
}

type OperationsWorkbenchQuery struct {
	TenantID uuid.UUID
	Lane     *string
	Category *string
	Severity *string
	Search   *string
	Limit    int32
	Offset   int32
}

type OperationsWorkbenchActionCommand struct {
	TenantID uuid.UUID  `json:"tenant_id"`
	CardKey  string     `json:"card_key"`
	Action   string     `json:"action"`
	Remarks  *string    `json:"remarks,omitempty"`
	ActorID  *uuid.UUID `json:"-"`
}
