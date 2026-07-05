package ports

import "github.com/google/uuid"

type OperationCatalogQuery struct {
	TenantID    uuid.UUID
	IncludeAll  bool
	Permissions []string
}
