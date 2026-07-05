package postgres

import (
	"encoding/json"

	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func mapWorkerType(row sqlc.HrmsWorkerType) *domain.WorkerType {
	statutoryDefaults := json.RawMessage(row.StatutoryDefaults)
	if len(statutoryDefaults) == 0 {
		statutoryDefaults = json.RawMessage(`{}`)
	}
	return &domain.WorkerType{
		ID:                  row.ID,
		TenantID:            row.TenantID,
		Code:                row.Code,
		Name:                row.Name,
		ClassificationGroup: row.ClassificationGroup,
		Description:         ptrFromText(row.Description),
		AttendanceMode:      row.AttendanceMode,
		PayMode:             row.PayMode,
		TDSSection:          row.TdsSection,
		PFApplicable:        row.PfApplicable,
		ESICApplicable:      row.EsicApplicable,
		PTApplicable:        row.PtApplicable,
		LWFApplicable:       row.LwfApplicable,
		CLRAApplicable:      row.ClraApplicable,
		LeaveApplicable:     row.LeaveApplicable,
		OvertimeApplicable:  row.OvertimeApplicable,
		RequiresAgreement:   row.RequiresAgreement,
		RequiresInvoice:     row.RequiresInvoice,
		RequiresAttendance:  row.RequiresAttendance,
		StatutoryDefaults:   statutoryDefaults,
		ComplianceNotes:     ptrFromText(row.ComplianceNotes),
		IsSystemDefault:     row.IsSystemDefault,
		SortOrder:           row.SortOrder,
		Inactive:            row.Inactive,
		CreatedAt:           timeFromTimestamptz(row.CreatedAt),
		CreatedBy:           ptrFromUUID(row.CreatedBy),
		UpdatedAt:           timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:           ptrFromUUID(row.UpdatedBy),
	}
}

func mapWorkerTypes(rows []sqlc.HrmsWorkerType) []*domain.WorkerType {
	items := make([]*domain.WorkerType, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapWorkerType(row))
	}
	return items
}

func mapWorkerClassificationRule(row sqlc.HrmsWorkerClassificationRule) *domain.WorkerClassificationRule {
	conditions := json.RawMessage(row.Conditions)
	if len(conditions) == 0 {
		conditions = json.RawMessage(`{}`)
	}
	outcome := json.RawMessage(row.Outcome)
	if len(outcome) == 0 {
		outcome = json.RawMessage(`{}`)
	}
	return &domain.WorkerClassificationRule{
		ID:           row.ID,
		TenantID:     row.TenantID,
		WorkerTypeID: row.WorkerTypeID,
		RuleName:     row.RuleName,
		RuleType:     row.RuleType,
		Priority:     row.Priority,
		Conditions:   conditions,
		Outcome:      outcome,
		Notes:        ptrFromText(row.Notes),
		Inactive:     row.Inactive,
		CreatedAt:    timeFromTimestamptz(row.CreatedAt),
		CreatedBy:    ptrFromUUID(row.CreatedBy),
		UpdatedAt:    timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:    ptrFromUUID(row.UpdatedBy),
	}
}

func mapWorkerClassificationRules(rows []sqlc.HrmsWorkerClassificationRule) []*domain.WorkerClassificationRule {
	items := make([]*domain.WorkerClassificationRule, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapWorkerClassificationRule(row))
	}
	return items
}

func jsonBytesFromRaw(value json.RawMessage) []byte {
	if len(value) == 0 {
		return []byte(`{}`)
	}
	return []byte(value)
}
