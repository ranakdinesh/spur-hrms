package postgres

import (
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func mapPolicyType(row sqlc.HrmsPolicyType) *domain.PolicyType {
	return &domain.PolicyType{
		ID:        row.ID,
		TenantID:  ptrFromUUID(row.TenantID),
		Name:      row.Name,
		IsSystem:  row.IsSystem,
		Inactive:  row.Inactive,
		CreatedAt: timeFromTimestamptz(row.CreatedAt),
		CreatedBy: ptrFromUUID(row.CreatedBy),
		UpdatedAt: timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy: ptrFromUUID(row.UpdatedBy),
	}
}

func mapPolicyTypes(rows []sqlc.HrmsPolicyType) []*domain.PolicyType {
	items := make([]*domain.PolicyType, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapPolicyType(row))
	}
	return items
}

func mapCompanyPolicy(row sqlc.HrmsCompanyPolicy) *domain.CompanyPolicy {
	return &domain.CompanyPolicy{
		ID:           row.ID,
		TenantID:     row.TenantID,
		PolicyTypeID: ptrFromUUID(row.PolicyTypeID),
		Title:        row.Title,
		FilePath:     ptrFromText(row.FilePath),
		Description:  ptrFromText(row.Description),
		Inactive:     row.Inactive,
		CreatedAt:    timeFromTimestamptz(row.CreatedAt),
		CreatedBy:    ptrFromUUID(row.CreatedBy),
		UpdatedAt:    timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:    ptrFromUUID(row.UpdatedBy),
	}
}

func mapCompanyPolicies(rows []sqlc.HrmsCompanyPolicy) []*domain.CompanyPolicy {
	items := make([]*domain.CompanyPolicy, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapCompanyPolicy(row))
	}
	return items
}
