package postgres

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func mapSuccessionReviewCycle(row sqlc.HrmsSuccessionReviewCycle) *domain.SuccessionReviewCycle {
	return &domain.SuccessionReviewCycle{ID: row.ID, TenantID: row.TenantID, Code: row.Code, Name: row.Name, Status: row.Status, StartsOn: ptrFromDate(row.StartsOn), EndsOn: ptrFromDate(row.EndsOn), ConfidentialityLevel: row.ConfidentialityLevel, Notes: ptrFromText(row.Notes), Metadata: successionRaw(row.Metadata), Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)}
}

func mapSuccessionReviewCycles(rows []sqlc.HrmsSuccessionReviewCycle) []*domain.SuccessionReviewCycle {
	items := make([]*domain.SuccessionReviewCycle, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapSuccessionReviewCycle(row))
	}
	return items
}

func mapSuccessionCriticalRole(row sqlc.HrmsSuccessionCriticalRole) *domain.SuccessionCriticalRole {
	return successionCriticalRoleFromParts(row.ID, row.TenantID, row.CycleID, row.Code, row.Title, row.DepartmentID, row.DesignationID, row.IncumbentWorkerProfileID, row.EmergencyCoverWorkerProfileID, row.Criticality, row.ImpactLevel, row.VacancyRisk, row.AttritionRisk, row.ReadinessTarget, row.SuccessorRequiredCount, row.RoleSummary, row.Status, row.Metadata, row.Inactive, row.CreatedAt, row.CreatedBy, row.UpdatedAt, row.UpdatedBy, nil, nil, nil, nil, nil, nil, 0, 0)
}

func mapSuccessionCriticalRoleDetail(row sqlc.GetSuccessionCriticalRoleRow) *domain.SuccessionCriticalRole {
	return successionCriticalRoleFromParts(row.ID, row.TenantID, row.CycleID, row.Code, row.Title, row.DepartmentID, row.DesignationID, row.IncumbentWorkerProfileID, row.EmergencyCoverWorkerProfileID, row.Criticality, row.ImpactLevel, row.VacancyRisk, row.AttritionRisk, row.ReadinessTarget, row.SuccessorRequiredCount, row.RoleSummary, row.Status, row.Metadata, row.Inactive, row.CreatedAt, row.CreatedBy, row.UpdatedAt, row.UpdatedBy, ptrFromText(row.DepartmentName), ptrFromText(row.DesignationName), ptrFromText(row.IncumbentName), ptrFromText(row.IncumbentCode), ptrFromText(row.EmergencyCoverName), ptrFromText(row.EmergencyCoverCode), 0, 0)
}

func mapSuccessionCriticalRoleList(rows []sqlc.ListSuccessionCriticalRolesRow) []*domain.SuccessionCriticalRole {
	items := make([]*domain.SuccessionCriticalRole, 0, len(rows))
	for _, row := range rows {
		items = append(items, successionCriticalRoleFromParts(row.ID, row.TenantID, row.CycleID, row.Code, row.Title, row.DepartmentID, row.DesignationID, row.IncumbentWorkerProfileID, row.EmergencyCoverWorkerProfileID, row.Criticality, row.ImpactLevel, row.VacancyRisk, row.AttritionRisk, row.ReadinessTarget, row.SuccessorRequiredCount, row.RoleSummary, row.Status, row.Metadata, row.Inactive, row.CreatedAt, row.CreatedBy, row.UpdatedAt, row.UpdatedBy, ptrFromText(row.DepartmentName), ptrFromText(row.DesignationName), ptrFromText(row.IncumbentName), ptrFromText(row.IncumbentCode), ptrFromText(row.EmergencyCoverName), ptrFromText(row.EmergencyCoverCode), row.SuccessorCount, row.ReadyNowCount))
	}
	return items
}

func successionCriticalRoleFromParts(id uuid.UUID, tenantID uuid.UUID, cycleID pgtype.UUID, code string, title string, departmentID pgtype.UUID, designationID pgtype.UUID, incumbentWorkerProfileID pgtype.UUID, emergencyCoverWorkerProfileID pgtype.UUID, criticality string, impactLevel string, vacancyRisk string, attritionRisk string, readinessTarget string, successorRequiredCount int32, roleSummary pgtype.Text, status string, metadata []byte, inactive bool, createdAt pgtype.Timestamptz, createdBy pgtype.UUID, updatedAt pgtype.Timestamptz, updatedBy pgtype.UUID, departmentName *string, designationName *string, incumbentName *string, incumbentCode *string, emergencyCoverName *string, emergencyCoverCode *string, successorCount int64, readyNowCount int64) *domain.SuccessionCriticalRole {
	return &domain.SuccessionCriticalRole{ID: id, TenantID: tenantID, CycleID: ptrFromUUID(cycleID), Code: code, Title: title, DepartmentID: ptrFromUUID(departmentID), DesignationID: ptrFromUUID(designationID), IncumbentWorkerProfileID: ptrFromUUID(incumbentWorkerProfileID), EmergencyCoverWorkerProfileID: ptrFromUUID(emergencyCoverWorkerProfileID), Criticality: criticality, ImpactLevel: impactLevel, VacancyRisk: vacancyRisk, AttritionRisk: attritionRisk, ReadinessTarget: readinessTarget, SuccessorRequiredCount: successorRequiredCount, RoleSummary: ptrFromText(roleSummary), Status: status, Metadata: successionRaw(metadata), Inactive: inactive, CreatedAt: timeFromTimestamptz(createdAt), CreatedBy: ptrFromUUID(createdBy), UpdatedAt: timeFromTimestamptz(updatedAt), UpdatedBy: ptrFromUUID(updatedBy), DepartmentName: departmentName, DesignationName: designationName, IncumbentName: incumbentName, IncumbentCode: incumbentCode, EmergencyCoverName: emergencyCoverName, EmergencyCoverCode: emergencyCoverCode, SuccessorCount: successorCount, ReadyNowCount: readyNowCount}
}

func mapSuccessionNomination(row sqlc.HrmsSuccessionSuccessorNomination) *domain.SuccessionSuccessorNomination {
	return successionNominationFromParts(row.ID, row.TenantID, row.CriticalRoleID, row.WorkerProfileID, row.NominatedBy, row.ReadinessLevel, row.ReadinessMonths, row.PotentialRating, row.PerformanceRating, row.RetentionRisk, row.MobilityPreference, row.NominationStatus, row.DevelopmentNotes, row.Metadata, row.Inactive, row.CreatedAt, row.CreatedBy, row.UpdatedAt, row.UpdatedBy, nil, nil, nil)
}

func mapSuccessionNominationList(rows []sqlc.ListSuccessionSuccessorNominationsRow) []*domain.SuccessionSuccessorNomination {
	items := make([]*domain.SuccessionSuccessorNomination, 0, len(rows))
	for _, row := range rows {
		title := row.CriticalRoleTitle
		items = append(items, successionNominationFromParts(row.ID, row.TenantID, row.CriticalRoleID, row.WorkerProfileID, row.NominatedBy, row.ReadinessLevel, row.ReadinessMonths, row.PotentialRating, row.PerformanceRating, row.RetentionRisk, row.MobilityPreference, row.NominationStatus, row.DevelopmentNotes, row.Metadata, row.Inactive, row.CreatedAt, row.CreatedBy, row.UpdatedAt, row.UpdatedBy, &row.WorkerDisplayName, ptrFromText(row.WorkerCode), &title))
	}
	return items
}

func successionNominationFromParts(id uuid.UUID, tenantID uuid.UUID, criticalRoleID uuid.UUID, workerProfileID uuid.UUID, nominatedBy pgtype.UUID, readinessLevel string, readinessMonths int32, potentialRating pgtype.Text, performanceRating pgtype.Text, retentionRisk string, mobilityPreference pgtype.Text, nominationStatus string, developmentNotes pgtype.Text, metadata []byte, inactive bool, createdAt pgtype.Timestamptz, createdBy pgtype.UUID, updatedAt pgtype.Timestamptz, updatedBy pgtype.UUID, workerDisplayName *string, workerCode *string, criticalRoleTitle *string) *domain.SuccessionSuccessorNomination {
	return &domain.SuccessionSuccessorNomination{ID: id, TenantID: tenantID, CriticalRoleID: criticalRoleID, WorkerProfileID: workerProfileID, NominatedBy: ptrFromUUID(nominatedBy), ReadinessLevel: readinessLevel, ReadinessMonths: readinessMonths, PotentialRating: ptrFromText(potentialRating), PerformanceRating: ptrFromText(performanceRating), RetentionRisk: retentionRisk, MobilityPreference: ptrFromText(mobilityPreference), NominationStatus: nominationStatus, DevelopmentNotes: ptrFromText(developmentNotes), Metadata: successionRaw(metadata), Inactive: inactive, CreatedAt: timeFromTimestamptz(createdAt), CreatedBy: ptrFromUUID(createdBy), UpdatedAt: timeFromTimestamptz(updatedAt), UpdatedBy: ptrFromUUID(updatedBy), WorkerDisplayName: workerDisplayName, WorkerCode: workerCode, CriticalRoleTitle: criticalRoleTitle}
}

func mapSuccessionDevelopmentAction(row sqlc.HrmsSuccessionDevelopmentAction) *domain.SuccessionDevelopmentAction {
	return successionDevelopmentActionFromParts(row.ID, row.TenantID, row.NominationID, row.CriticalRoleID, row.WorkerProfileID, row.ActionType, row.Title, row.LearningCourseID, row.LearningPathID, row.OwnerUserID, row.DueDate, row.Status, row.Notes, row.Metadata, row.Inactive, row.CreatedAt, row.CreatedBy, row.UpdatedAt, row.UpdatedBy, nil, nil, nil, nil, nil)
}

func mapSuccessionDevelopmentActionList(rows []sqlc.ListSuccessionDevelopmentActionsRow) []*domain.SuccessionDevelopmentAction {
	items := make([]*domain.SuccessionDevelopmentAction, 0, len(rows))
	for _, row := range rows {
		items = append(items, successionDevelopmentActionFromParts(row.ID, row.TenantID, row.NominationID, row.CriticalRoleID, row.WorkerProfileID, row.ActionType, row.Title, row.LearningCourseID, row.LearningPathID, row.OwnerUserID, row.DueDate, row.Status, row.Notes, row.Metadata, row.Inactive, row.CreatedAt, row.CreatedBy, row.UpdatedAt, row.UpdatedBy, &row.WorkerDisplayName, ptrFromText(row.WorkerCode), ptrFromText(row.CriticalRoleTitle), ptrFromText(row.LearningCourseTitle), ptrFromText(row.LearningPathTitle)))
	}
	return items
}

func successionDevelopmentActionFromParts(id uuid.UUID, tenantID uuid.UUID, nominationID pgtype.UUID, criticalRoleID pgtype.UUID, workerProfileID uuid.UUID, actionType string, title string, learningCourseID pgtype.UUID, learningPathID pgtype.UUID, ownerUserID pgtype.UUID, dueDate pgtype.Date, status string, notes pgtype.Text, metadata []byte, inactive bool, createdAt pgtype.Timestamptz, createdBy pgtype.UUID, updatedAt pgtype.Timestamptz, updatedBy pgtype.UUID, workerDisplayName *string, workerCode *string, criticalRoleTitle *string, learningCourseTitle *string, learningPathTitle *string) *domain.SuccessionDevelopmentAction {
	return &domain.SuccessionDevelopmentAction{ID: id, TenantID: tenantID, NominationID: ptrFromUUID(nominationID), CriticalRoleID: ptrFromUUID(criticalRoleID), WorkerProfileID: workerProfileID, ActionType: actionType, Title: title, LearningCourseID: ptrFromUUID(learningCourseID), LearningPathID: ptrFromUUID(learningPathID), OwnerUserID: ptrFromUUID(ownerUserID), DueDate: ptrFromDate(dueDate), Status: status, Notes: ptrFromText(notes), Metadata: successionRaw(metadata), Inactive: inactive, CreatedAt: timeFromTimestamptz(createdAt), CreatedBy: ptrFromUUID(createdBy), UpdatedAt: timeFromTimestamptz(updatedAt), UpdatedBy: ptrFromUUID(updatedBy), WorkerDisplayName: workerDisplayName, WorkerCode: workerCode, CriticalRoleTitle: criticalRoleTitle, LearningCourseTitle: learningCourseTitle, LearningPathTitle: learningPathTitle}
}

func mapSuccessionEvent(row sqlc.HrmsSuccessionEvent) *domain.SuccessionEvent {
	return &domain.SuccessionEvent{ID: row.ID, TenantID: row.TenantID, SourceType: row.SourceType, SourceID: ptrFromUUID(row.SourceID), Action: row.Action, FromStatus: ptrFromText(row.FromStatus), ToStatus: ptrFromText(row.ToStatus), Remarks: ptrFromText(row.Remarks), Metadata: successionRaw(row.Metadata), CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy)}
}

func mapSuccessionEvents(rows []sqlc.HrmsSuccessionEvent) []*domain.SuccessionEvent {
	items := make([]*domain.SuccessionEvent, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapSuccessionEvent(row))
	}
	return items
}

func mapSuccessionSummary(rows []sqlc.GetSuccessionSummaryRow) []*domain.SuccessionSummaryRow {
	items := make([]*domain.SuccessionSummaryRow, 0, len(rows))
	for _, row := range rows {
		items = append(items, &domain.SuccessionSummaryRow{Metric: row.Metric, MetricCount: row.MetricCount})
	}
	return items
}

func successionRaw(value []byte) json.RawMessage {
	if len(value) == 0 {
		return json.RawMessage(`{}`)
	}
	return json.RawMessage(value)
}
