package postgres

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func mapOKRCycle(row sqlc.HrmsOkrCycle) *domain.OKRCycle {
	return &domain.OKRCycle{ID: row.ID, TenantID: row.TenantID, Name: row.Name, CycleCode: row.CycleCode, Description: ptrFromText(row.Description), StartDate: okrDateValue(row.StartDate), EndDate: okrDateValue(row.EndDate), Status: row.Status, ReviewCadence: row.ReviewCadence, Metadata: okrRaw(row.Metadata), Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)}
}

func mapOKRCycles(rows []sqlc.HrmsOkrCycle) []*domain.OKRCycle {
	items := make([]*domain.OKRCycle, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapOKRCycle(row))
	}
	return items
}

func mapObjective(row sqlc.HrmsObjective) *domain.Objective {
	return okrObjectiveFromParts(row.ID, row.TenantID, row.CycleID, row.ParentObjectiveID, row.OwnerType, row.OwnerWorkerProfileID, row.OwnerDepartmentID, row.OwnerProjectID, row.Title, row.Description, row.Status, row.Priority, row.ProgressPercent, row.Weight, row.StartDate, row.DueDate, row.Metadata, row.Inactive, row.CreatedAt, row.CreatedBy, row.UpdatedAt, row.UpdatedBy, nil, nil, nil, nil, nil, nil, nil, 0, 0)
}

func mapObjectiveRows(rows []sqlc.ListObjectivesRow) []*domain.Objective {
	items := make([]*domain.Objective, 0, len(rows))
	for _, row := range rows {
		items = append(items, okrObjectiveFromParts(row.ID, row.TenantID, row.CycleID, row.ParentObjectiveID, row.OwnerType, row.OwnerWorkerProfileID, row.OwnerDepartmentID, row.OwnerProjectID, row.Title, row.Description, row.Status, row.Priority, row.ProgressPercent, row.Weight, row.StartDate, row.DueDate, row.Metadata, row.Inactive, row.CreatedAt, row.CreatedBy, row.UpdatedAt, row.UpdatedBy, &row.CycleName, ptrFromText(row.ParentObjectiveTitle), ptrFromText(row.OwnerWorkerName), ptrFromText(row.OwnerWorkerCode), ptrFromText(row.OwnerDepartmentName), ptrFromText(row.OwnerProjectName), ptrFromText(row.OwnerProjectCode), row.KeyResultCount, floatFromNumeric(row.AverageKeyResultProgress)))
	}
	return items
}

func okrObjectiveFromParts(id uuid.UUID, tenantID uuid.UUID, cycleID uuid.UUID, parentObjectiveID pgtype.UUID, ownerType string, ownerWorkerProfileID pgtype.UUID, ownerDepartmentID pgtype.UUID, ownerProjectID pgtype.UUID, title string, description pgtype.Text, status string, priority string, progress pgtype.Numeric, weight pgtype.Numeric, startDate pgtype.Date, dueDate pgtype.Date, metadata []byte, inactive bool, createdAt pgtype.Timestamptz, createdBy pgtype.UUID, updatedAt pgtype.Timestamptz, updatedBy pgtype.UUID, cycleName *string, parentObjectiveTitle *string, ownerWorkerName *string, ownerWorkerCode *string, ownerDepartmentName *string, ownerProjectName *string, ownerProjectCode *string, keyResultCount int32, averageProgress float64) *domain.Objective {
	return &domain.Objective{ID: id, TenantID: tenantID, CycleID: cycleID, ParentObjectiveID: ptrFromUUID(parentObjectiveID), OwnerType: ownerType, OwnerWorkerProfileID: ptrFromUUID(ownerWorkerProfileID), OwnerDepartmentID: ptrFromUUID(ownerDepartmentID), OwnerProjectID: ptrFromUUID(ownerProjectID), Title: title, Description: ptrFromText(description), Status: status, Priority: priority, ProgressPercent: floatFromNumeric(progress), Weight: floatFromNumeric(weight), StartDate: ptrFromDate(startDate), DueDate: ptrFromDate(dueDate), Metadata: okrRaw(metadata), Inactive: inactive, CreatedAt: timeFromTimestamptz(createdAt), CreatedBy: ptrFromUUID(createdBy), UpdatedAt: timeFromTimestamptz(updatedAt), UpdatedBy: ptrFromUUID(updatedBy), CycleName: cycleName, ParentObjectiveTitle: parentObjectiveTitle, OwnerWorkerName: ownerWorkerName, OwnerWorkerCode: ownerWorkerCode, OwnerDepartmentName: ownerDepartmentName, OwnerProjectName: ownerProjectName, OwnerProjectCode: ownerProjectCode, KeyResultCount: keyResultCount, AverageKeyResultProgress: averageProgress}
}

func mapKeyResult(row sqlc.HrmsKeyResult) *domain.KeyResult {
	return okrKeyResultFromParts(row.ID, row.TenantID, row.ObjectiveID, row.Title, row.Description, row.MetricType, row.StartValue, row.TargetValue, row.CurrentValue, row.ProgressPercent, row.Confidence, row.Status, row.Weight, row.UnitLabel, row.DueDate, row.Metadata, row.Inactive, row.CreatedAt, row.CreatedBy, row.UpdatedAt, row.UpdatedBy, nil, nil, pgtype.Date{}, nil)
}

func mapKeyResultRows(rows []sqlc.ListKeyResultsRow) []*domain.KeyResult {
	items := make([]*domain.KeyResult, 0, len(rows))
	for _, row := range rows {
		items = append(items, okrKeyResultFromParts(row.ID, row.TenantID, row.ObjectiveID, row.Title, row.Description, row.MetricType, row.StartValue, row.TargetValue, row.CurrentValue, row.ProgressPercent, row.Confidence, row.Status, row.Weight, row.UnitLabel, row.DueDate, row.Metadata, row.Inactive, row.CreatedAt, row.CreatedBy, row.UpdatedAt, row.UpdatedBy, &row.ObjectiveTitle, &row.CycleName, row.LatestCheckinDate, ptrFromText(row.LatestNote)))
	}
	return items
}

func okrKeyResultFromParts(id uuid.UUID, tenantID uuid.UUID, objectiveID uuid.UUID, title string, description pgtype.Text, metricType string, startValue pgtype.Numeric, targetValue pgtype.Numeric, currentValue pgtype.Numeric, progress pgtype.Numeric, confidence string, status string, weight pgtype.Numeric, unitLabel pgtype.Text, dueDate pgtype.Date, metadata []byte, inactive bool, createdAt pgtype.Timestamptz, createdBy pgtype.UUID, updatedAt pgtype.Timestamptz, updatedBy pgtype.UUID, objectiveTitle *string, cycleName *string, latestCheckinDate pgtype.Date, latestNote *string) *domain.KeyResult {
	return &domain.KeyResult{ID: id, TenantID: tenantID, ObjectiveID: objectiveID, Title: title, Description: ptrFromText(description), MetricType: metricType, StartValue: floatFromNumeric(startValue), TargetValue: floatFromNumeric(targetValue), CurrentValue: floatFromNumeric(currentValue), ProgressPercent: floatFromNumeric(progress), Confidence: confidence, Status: status, Weight: floatFromNumeric(weight), UnitLabel: ptrFromText(unitLabel), DueDate: ptrFromDate(dueDate), Metadata: okrRaw(metadata), Inactive: inactive, CreatedAt: timeFromTimestamptz(createdAt), CreatedBy: ptrFromUUID(createdBy), UpdatedAt: timeFromTimestamptz(updatedAt), UpdatedBy: ptrFromUUID(updatedBy), ObjectiveTitle: objectiveTitle, CycleName: cycleName, LatestCheckinDate: ptrFromDate(latestCheckinDate), LatestNote: latestNote}
}

func mapKeyResultCheckIn(row sqlc.HrmsKeyResultCheckin) *domain.KeyResultCheckIn {
	return okrCheckInFromParts(row.ID, row.TenantID, row.KeyResultID, row.CheckinDate, row.Value, row.ProgressPercent, row.Confidence, row.Status, row.Note, row.Metadata, row.Inactive, row.CreatedAt, row.CreatedBy, nil, uuid.Nil, nil)
}

func mapKeyResultCheckInRows(rows []sqlc.ListKeyResultCheckInsRow) []*domain.KeyResultCheckIn {
	items := make([]*domain.KeyResultCheckIn, 0, len(rows))
	for _, row := range rows {
		items = append(items, okrCheckInFromParts(row.ID, row.TenantID, row.KeyResultID, row.CheckinDate, row.Value, row.ProgressPercent, row.Confidence, row.Status, row.Note, row.Metadata, row.Inactive, row.CreatedAt, row.CreatedBy, &row.KeyResultTitle, row.ObjectiveID, &row.ObjectiveTitle))
	}
	return items
}

func okrCheckInFromParts(id uuid.UUID, tenantID uuid.UUID, keyResultID uuid.UUID, checkInDate pgtype.Date, value pgtype.Numeric, progress pgtype.Numeric, confidence string, status string, note pgtype.Text, metadata []byte, inactive bool, createdAt pgtype.Timestamptz, createdBy pgtype.UUID, keyResultTitle *string, objectiveID uuid.UUID, objectiveTitle *string) *domain.KeyResultCheckIn {
	item := &domain.KeyResultCheckIn{ID: id, TenantID: tenantID, KeyResultID: keyResultID, CheckInDate: okrDateValue(checkInDate), Value: floatFromNumeric(value), ProgressPercent: floatFromNumeric(progress), Confidence: confidence, Status: status, Note: ptrFromText(note), Metadata: okrRaw(metadata), Inactive: inactive, CreatedAt: timeFromTimestamptz(createdAt), CreatedBy: ptrFromUUID(createdBy), KeyResultTitle: keyResultTitle, ObjectiveTitle: objectiveTitle}
	if objectiveID != uuid.Nil {
		item.ObjectiveID = &objectiveID
	}
	return item
}

func mapOKRSummaryRows(rows []sqlc.GetOKRSummaryRow) []*domain.OKRSummaryRow {
	items := make([]*domain.OKRSummaryRow, 0, len(rows))
	for _, row := range rows {
		items = append(items, &domain.OKRSummaryRow{OwnerType: row.OwnerType, ObjectiveCount: row.ObjectiveCount, KeyResultCount: row.KeyResultCount, AverageProgress: floatFromNumeric(row.AverageProgress), AtRiskCount: row.AtRiskCount, CompletedCount: row.CompletedCount})
	}
	return items
}

func okrDateValue(value pgtype.Date) time.Time {
	if value.Valid {
		return value.Time
	}
	return time.Time{}
}

func okrRaw(value []byte) json.RawMessage {
	raw := json.RawMessage(value)
	if len(raw) == 0 {
		return json.RawMessage(`{}`)
	}
	return raw
}
