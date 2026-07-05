package postgres

import (
	"encoding/json"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func mapLeavePolicyTemplate(row sqlc.HrmsLeavePolicyTemplate) *domain.LeavePolicyTemplate {
	return &domain.LeavePolicyTemplate{
		ID:            row.ID,
		TenantID:      ptrFromUUID(row.TenantID),
		Name:          row.Name,
		Code:          row.Code,
		Description:   ptrFromText(row.Description),
		IsSystem:      row.IsSystem,
		EffectiveFrom: ptrFromDate(row.EffectiveFrom),
		EffectiveTo:   ptrFromDate(row.EffectiveTo),
		Inactive:      row.Inactive,
		CreatedAt:     timeFromTimestamptz(row.CreatedAt),
		CreatedBy:     ptrFromUUID(row.CreatedBy),
		UpdatedAt:     timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:     ptrFromUUID(row.UpdatedBy),
	}
}

func mapLeavePolicyTemplates(rows []sqlc.HrmsLeavePolicyTemplate) []*domain.LeavePolicyTemplate {
	items := make([]*domain.LeavePolicyTemplate, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapLeavePolicyTemplate(row))
	}
	return items
}

func mapLeavePolicyTemplateRule(row sqlc.HrmsLeavePolicyTemplateRule) (*domain.LeavePolicyTemplateRule, error) {
	config := map[string]any{}
	if len(row.CalculationConfig) > 0 {
		if err := json.Unmarshal(row.CalculationConfig, &config); err != nil {
			return nil, err
		}
	}
	return &domain.LeavePolicyTemplateRule{
		ID:                        row.ID,
		TenantID:                  row.TenantID,
		TemplateID:                row.TemplateID,
		LeaveTypeID:               row.LeaveTypeID,
		FYID:                      ptrFromUUID(row.FyID),
		EmploymentTypeID:          ptrFromUUID(row.EmploymentTypeID),
		DepartmentID:              ptrFromUUID(row.DepartmentID),
		DesignationID:             ptrFromUUID(row.DesignationID),
		ProbationStatus:           ptrFromText(row.ProbationStatus),
		AccrualMethod:             row.AccrualMethod,
		AccrualFrequency:          row.AccrualFrequency,
		CreditDays:                floatFromNumeric(row.CreditDays),
		CreditHours:               floatFromNumeric(row.CreditHours),
		AnnualEntitlement:         floatFromNumeric(row.AnnualEntitlement),
		MinWorkedDays:             row.MinWorkedDays,
		MaxBalance:                floatPtrFromNumeric(row.MaxBalance),
		CarryForwardEnabled:       row.CarryForwardEnabled,
		MaxCarryForward:           floatFromNumeric(row.MaxCarryForward),
		CarryForwardExpiryMonths:  row.CarryForwardExpiryMonths,
		EncashmentEnabled:         row.EncashmentEnabled,
		EncashmentLimit:           floatFromNumeric(row.EncashmentLimit),
		EncashmentPayablePercent:  floatFromNumeric(row.EncashmentPayablePercent),
		NegativeBalanceAllowed:    row.NegativeBalanceAllowed,
		MaxNegativeBalance:        floatFromNumeric(row.MaxNegativeBalance),
		SandwichApplicable:        row.SandwichApplicable,
		IncludeHolidays:           row.IncludeHolidays,
		IncludeWeekoffs:           row.IncludeWeekoffs,
		RequiresDocumentAfterDays: floatPtrFromNumeric(row.RequiresDocumentAfterDays),
		MinRequestDays:            floatFromNumeric(row.MinRequestDays),
		MaxRequestDays:            floatPtrFromNumeric(row.MaxRequestDays),
		MaxRequestsPerYear:        row.MaxRequestsPerYear,
		AccrualDay:                row.AccrualDay,
		LapseUnutilized:           row.LapseUnutilized,
		AllowHalfDay:              row.AllowHalfDay,
		RequiresApproval:          row.RequiresApproval,
		CalculationConfig:         config,
		Priority:                  row.Priority,
		Inactive:                  row.Inactive,
		CreatedAt:                 timeFromTimestamptz(row.CreatedAt),
		CreatedBy:                 ptrFromUUID(row.CreatedBy),
		UpdatedAt:                 timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:                 ptrFromUUID(row.UpdatedBy),
	}, nil
}

func mapLeavePolicyTemplateRules(rows []sqlc.HrmsLeavePolicyTemplateRule) ([]*domain.LeavePolicyTemplateRule, error) {
	items := make([]*domain.LeavePolicyTemplateRule, 0, len(rows))
	for _, row := range rows {
		item, err := mapLeavePolicyTemplateRule(row)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func mapEmployeeLeavePolicyAssignment(row sqlc.HrmsEmployeeLeavePolicyAssignment) *domain.EmployeeLeavePolicyAssignment {
	return &domain.EmployeeLeavePolicyAssignment{
		ID:            row.ID,
		TenantID:      row.TenantID,
		UserID:        row.UserID,
		TemplateID:    row.TemplateID,
		FYID:          ptrFromUUID(row.FyID),
		EffectiveFrom: timeFromDate(row.EffectiveFrom),
		EffectiveTo:   ptrFromDate(row.EffectiveTo),
		Inactive:      row.Inactive,
		CreatedAt:     timeFromTimestamptz(row.CreatedAt),
		CreatedBy:     ptrFromUUID(row.CreatedBy),
		UpdatedAt:     timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:     ptrFromUUID(row.UpdatedBy),
	}
}

func mapEmployeeLeavePolicyAssignments(rows []sqlc.HrmsEmployeeLeavePolicyAssignment) []*domain.EmployeeLeavePolicyAssignment {
	items := make([]*domain.EmployeeLeavePolicyAssignment, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapEmployeeLeavePolicyAssignment(row))
	}
	return items
}

func jsonBytesFromMap(value map[string]any) []byte {
	if value == nil {
		return []byte("{}")
	}
	data, err := json.Marshal(value)
	if err != nil {
		return []byte("{}")
	}
	return data
}

func numericFromFloatPtr(value *float64) pgtype.Numeric {
	if value == nil {
		return pgtype.Numeric{Valid: false}
	}
	return numericFromFloat(*value)
}

func floatPtrFromNumeric(value pgtype.Numeric) *float64 {
	if !value.Valid {
		return nil
	}
	floatValue := floatFromNumeric(value)
	return &floatValue
}
