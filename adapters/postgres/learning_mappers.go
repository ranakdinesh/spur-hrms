package postgres

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func mapLearningCourse(row sqlc.HrmsLearningCourse) *domain.LearningCourse {
	return learningCourseFromParts(row.ID, row.TenantID, row.Code, row.Title, row.Description, row.CourseType, row.DeliveryMode, row.Provider, row.DurationMinutes, row.SkillID, row.ComplianceRuleID, row.Mandatory, row.AiReadiness, row.CertificateRequired, row.BudgetAmount, row.CurrencyCode, row.IsActive, row.Metadata, row.Inactive, row.CreatedAt, row.CreatedBy, row.UpdatedAt, row.UpdatedBy, nil, nil, nil)
}

func mapLearningCourseList(rows []sqlc.ListLearningCoursesRow) []*domain.LearningCourse {
	items := make([]*domain.LearningCourse, 0, len(rows))
	for _, row := range rows {
		items = append(items, learningCourseFromParts(row.ID, row.TenantID, row.Code, row.Title, row.Description, row.CourseType, row.DeliveryMode, row.Provider, row.DurationMinutes, row.SkillID, row.ComplianceRuleID, row.Mandatory, row.AiReadiness, row.CertificateRequired, row.BudgetAmount, row.CurrencyCode, row.IsActive, row.Metadata, row.Inactive, row.CreatedAt, row.CreatedBy, row.UpdatedAt, row.UpdatedBy, ptrFromText(row.SkillName), ptrFromText(row.SkillCode), ptrFromText(row.ComplianceRuleTitle)))
	}
	return items
}

func learningCourseFromParts(id uuid.UUID, tenantID uuid.UUID, code string, title string, description pgtype.Text, courseType string, deliveryMode string, provider pgtype.Text, durationMinutes int32, skillID pgtype.UUID, complianceRuleID pgtype.UUID, mandatory bool, aiReadiness bool, certificateRequired bool, budgetAmount pgtype.Numeric, currencyCode string, isActive bool, metadataBytes []byte, inactive bool, createdAt pgtype.Timestamptz, createdBy pgtype.UUID, updatedAt pgtype.Timestamptz, updatedBy pgtype.UUID, skillName *string, skillCode *string, complianceRuleTitle *string) *domain.LearningCourse {
	return &domain.LearningCourse{ID: id, TenantID: tenantID, Code: code, Title: title, Description: ptrFromText(description), CourseType: courseType, DeliveryMode: deliveryMode, Provider: ptrFromText(provider), DurationMinutes: durationMinutes, SkillID: ptrFromUUID(skillID), ComplianceRuleID: ptrFromUUID(complianceRuleID), Mandatory: mandatory, AIReadiness: aiReadiness, CertificateRequired: certificateRequired, BudgetAmount: ptrFromNumeric(budgetAmount), CurrencyCode: currencyCode, IsActive: isActive, Metadata: learningRaw(metadataBytes), Inactive: inactive, CreatedAt: timeFromTimestamptz(createdAt), CreatedBy: ptrFromUUID(createdBy), UpdatedAt: timeFromTimestamptz(updatedAt), UpdatedBy: ptrFromUUID(updatedBy), SkillName: skillName, SkillCode: skillCode, ComplianceRuleTitle: complianceRuleTitle}
}

func mapLearningPath(row sqlc.HrmsLearningPath) *domain.LearningPath {
	return learningPathFromParts(row.ID, row.TenantID, row.Code, row.Title, row.Description, row.PathType, row.TargetRole, row.SkillID, row.AiReadiness, row.IsActive, row.Metadata, row.Inactive, row.CreatedAt, row.CreatedBy, row.UpdatedAt, row.UpdatedBy, nil, 0, 0)
}

func mapLearningPathList(rows []sqlc.ListLearningPathsRow) []*domain.LearningPath {
	items := make([]*domain.LearningPath, 0, len(rows))
	for _, row := range rows {
		items = append(items, learningPathFromParts(row.ID, row.TenantID, row.Code, row.Title, row.Description, row.PathType, row.TargetRole, row.SkillID, row.AiReadiness, row.IsActive, row.Metadata, row.Inactive, row.CreatedAt, row.CreatedBy, row.UpdatedAt, row.UpdatedBy, ptrFromText(row.SkillName), row.CourseCount, row.TotalMinutes))
	}
	return items
}

func learningPathFromParts(id uuid.UUID, tenantID uuid.UUID, code string, title string, description pgtype.Text, pathType string, targetRole pgtype.Text, skillID pgtype.UUID, aiReadiness bool, isActive bool, metadataBytes []byte, inactive bool, createdAt pgtype.Timestamptz, createdBy pgtype.UUID, updatedAt pgtype.Timestamptz, updatedBy pgtype.UUID, skillName *string, courseCount int32, totalMinutes int32) *domain.LearningPath {
	return &domain.LearningPath{ID: id, TenantID: tenantID, Code: code, Title: title, Description: ptrFromText(description), PathType: pathType, TargetRole: ptrFromText(targetRole), SkillID: ptrFromUUID(skillID), AIReadiness: aiReadiness, IsActive: isActive, Metadata: learningRaw(metadataBytes), Inactive: inactive, CreatedAt: timeFromTimestamptz(createdAt), CreatedBy: ptrFromUUID(createdBy), UpdatedAt: timeFromTimestamptz(updatedAt), UpdatedBy: ptrFromUUID(updatedBy), SkillName: skillName, CourseCount: courseCount, TotalMinutes: totalMinutes}
}

func mapLearningPathCourse(row sqlc.HrmsLearningPathCourse) *domain.LearningPathCourse {
	return &domain.LearningPathCourse{ID: row.ID, TenantID: row.TenantID, PathID: row.PathID, CourseID: row.CourseID, SortOrder: row.SortOrder, Required: row.Required, Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)}
}

func mapLearningPathCourseList(rows []sqlc.ListLearningPathCoursesRow) []*domain.LearningPathCourse {
	items := make([]*domain.LearningPathCourse, 0, len(rows))
	for _, row := range rows {
		items = append(items, &domain.LearningPathCourse{ID: row.ID, TenantID: row.TenantID, PathID: row.PathID, CourseID: row.CourseID, SortOrder: row.SortOrder, Required: row.Required, Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy), CourseCode: learningStringPtr(row.CourseCode), CourseTitle: learningStringPtr(row.CourseTitle), CourseType: learningStringPtr(row.CourseType), DeliveryMode: learningStringPtr(row.DeliveryMode), DurationMinutes: row.DurationMinutes})
	}
	return items
}

func mapLearningEnrollment(row sqlc.HrmsLearningEnrollment) *domain.LearningEnrollment {
	return learningEnrollmentFromParts(row.ID, row.TenantID, row.CourseID, row.PathID, row.WorkerProfileID, row.AssignmentSource, row.Status, row.NominatedBy, row.AssignedBy, row.DueDate, row.StartedAt, row.CompletedAt, row.Score, row.CertificateUrl, row.CertificateFileName, row.CertificateContentType, row.Notes, row.Metadata, row.Inactive, row.CreatedAt, row.CreatedBy, row.UpdatedAt, row.UpdatedBy, nil, nil, nil, false, false, nil, nil, nil)
}

func mapLearningEnrollmentList(rows []sqlc.ListLearningEnrollmentsRow) []*domain.LearningEnrollment {
	items := make([]*domain.LearningEnrollment, 0, len(rows))
	for _, row := range rows {
		items = append(items, learningEnrollmentFromParts(row.ID, row.TenantID, row.CourseID, row.PathID, row.WorkerProfileID, row.AssignmentSource, row.Status, row.NominatedBy, row.AssignedBy, row.DueDate, row.StartedAt, row.CompletedAt, row.Score, row.CertificateUrl, row.CertificateFileName, row.CertificateContentType, row.Notes, row.Metadata, row.Inactive, row.CreatedAt, row.CreatedBy, row.UpdatedAt, row.UpdatedBy, &row.CourseTitle, &row.CourseCode, &row.CourseType, row.AiReadiness, row.Mandatory, ptrFromText(row.PathTitle), &row.WorkerDisplayName, ptrFromText(row.WorkerCode)))
	}
	return items
}

func learningEnrollmentFromParts(id uuid.UUID, tenantID uuid.UUID, courseID uuid.UUID, pathID pgtype.UUID, workerProfileID uuid.UUID, assignmentSource string, status string, nominatedBy pgtype.UUID, assignedBy pgtype.UUID, dueDate pgtype.Date, startedAt pgtype.Timestamptz, completedAt pgtype.Timestamptz, score pgtype.Numeric, certificateURL pgtype.Text, certificateFileName pgtype.Text, certificateContentType pgtype.Text, notes pgtype.Text, metadataBytes []byte, inactive bool, createdAt pgtype.Timestamptz, createdBy pgtype.UUID, updatedAt pgtype.Timestamptz, updatedBy pgtype.UUID, courseTitle *string, courseCode *string, courseType *string, aiReadiness bool, mandatory bool, pathTitle *string, workerDisplayName *string, workerCode *string) *domain.LearningEnrollment {
	return &domain.LearningEnrollment{ID: id, TenantID: tenantID, CourseID: courseID, PathID: ptrFromUUID(pathID), WorkerProfileID: workerProfileID, AssignmentSource: assignmentSource, Status: status, NominatedBy: ptrFromUUID(nominatedBy), AssignedBy: ptrFromUUID(assignedBy), DueDate: ptrFromDate(dueDate), StartedAt: ptrFromTimestamptz(startedAt), CompletedAt: ptrFromTimestamptz(completedAt), Score: ptrFromNumeric(score), CertificateURL: ptrFromText(certificateURL), CertificateFileName: ptrFromText(certificateFileName), CertificateContentType: ptrFromText(certificateContentType), Notes: ptrFromText(notes), Metadata: learningRaw(metadataBytes), Inactive: inactive, CreatedAt: timeFromTimestamptz(createdAt), CreatedBy: ptrFromUUID(createdBy), UpdatedAt: timeFromTimestamptz(updatedAt), UpdatedBy: ptrFromUUID(updatedBy), CourseTitle: courseTitle, CourseCode: courseCode, CourseType: courseType, AIReadiness: aiReadiness, Mandatory: mandatory, PathTitle: pathTitle, WorkerDisplayName: workerDisplayName, WorkerCode: workerCode}
}

func mapLearningRecommendation(row sqlc.HrmsLearningRecommendation) *domain.LearningRecommendation {
	return learningRecommendationFromParts(row.ID, row.TenantID, row.WorkerProfileID, row.SkillID, row.CourseID, row.PathID, row.SourceType, row.Reason, row.Priority, row.ConfidenceScore, row.Status, row.Metadata, row.Inactive, row.CreatedAt, row.CreatedBy, row.UpdatedAt, row.UpdatedBy, nil, nil, nil, nil, nil)
}

func mapLearningRecommendationList(rows []sqlc.ListLearningRecommendationsRow) []*domain.LearningRecommendation {
	items := make([]*domain.LearningRecommendation, 0, len(rows))
	for _, row := range rows {
		items = append(items, learningRecommendationFromParts(row.ID, row.TenantID, row.WorkerProfileID, row.SkillID, row.CourseID, row.PathID, row.SourceType, row.Reason, row.Priority, row.ConfidenceScore, row.Status, row.Metadata, row.Inactive, row.CreatedAt, row.CreatedBy, row.UpdatedAt, row.UpdatedBy, ptrFromText(row.WorkerDisplayName), ptrFromText(row.WorkerCode), ptrFromText(row.SkillName), ptrFromText(row.CourseTitle), ptrFromText(row.PathTitle)))
	}
	return items
}

func mapLearningRecommendations(rows []sqlc.HrmsLearningRecommendation) []*domain.LearningRecommendation {
	items := make([]*domain.LearningRecommendation, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapLearningRecommendation(row))
	}
	return items
}

func learningRecommendationFromParts(id uuid.UUID, tenantID uuid.UUID, workerProfileID pgtype.UUID, skillID pgtype.UUID, courseID pgtype.UUID, pathID pgtype.UUID, sourceType string, reason string, priority string, confidenceScore pgtype.Numeric, status string, metadataBytes []byte, inactive bool, createdAt pgtype.Timestamptz, createdBy pgtype.UUID, updatedAt pgtype.Timestamptz, updatedBy pgtype.UUID, workerDisplayName *string, workerCode *string, skillName *string, courseTitle *string, pathTitle *string) *domain.LearningRecommendation {
	return &domain.LearningRecommendation{ID: id, TenantID: tenantID, WorkerProfileID: ptrFromUUID(workerProfileID), SkillID: ptrFromUUID(skillID), CourseID: ptrFromUUID(courseID), PathID: ptrFromUUID(pathID), SourceType: sourceType, Reason: reason, Priority: priority, ConfidenceScore: ptrFromNumeric(confidenceScore), Status: status, Metadata: learningRaw(metadataBytes), Inactive: inactive, CreatedAt: timeFromTimestamptz(createdAt), CreatedBy: ptrFromUUID(createdBy), UpdatedAt: timeFromTimestamptz(updatedAt), UpdatedBy: ptrFromUUID(updatedBy), WorkerDisplayName: workerDisplayName, WorkerCode: workerCode, SkillName: skillName, CourseTitle: courseTitle, PathTitle: pathTitle}
}

func mapLearningSummary(rows []sqlc.GetLearningSummaryRow) []*domain.LearningSummaryRow {
	items := make([]*domain.LearningSummaryRow, 0, len(rows))
	for _, row := range rows {
		items = append(items, &domain.LearningSummaryRow{Metric: row.Metric, MetricCount: row.MetricCount})
	}
	return items
}

func learningRaw(value []byte) json.RawMessage {
	if len(value) == 0 {
		return json.RawMessage(`{}`)
	}
	return json.RawMessage(value)
}

func learningStringPtr(value string) *string {
	return &value
}
