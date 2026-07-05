package ports

import "github.com/google/uuid"

type HRCommandCenterQuery struct {
	TenantID uuid.UUID
	Limit    int32
}
