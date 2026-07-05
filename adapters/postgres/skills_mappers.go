package postgres

import (
	"encoding/json"
	"math"
	"math/big"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func mapSkillCategory(row sqlc.HrmsSkillCategory) *domain.SkillCategory {
	return skillCategoryFromParts(row.ID, row.TenantID, row.ParentID, row.Code, row.Name, row.Description, row.SourceScope, row.SortOrder, row.Metadata, row.Inactive, row.CreatedAt, row.CreatedBy, row.UpdatedAt, row.UpdatedBy, nil)
}

func mapSkillCategoryList(rows []sqlc.ListSkillCategoriesRow) []*domain.SkillCategory {
	items := make([]*domain.SkillCategory, 0, len(rows))
	for _, row := range rows {
		items = append(items, skillCategoryFromParts(row.ID, row.TenantID, row.ParentID, row.Code, row.Name, row.Description, row.SourceScope, row.SortOrder, row.Metadata, row.Inactive, row.CreatedAt, row.CreatedBy, row.UpdatedAt, row.UpdatedBy, ptrFromText(row.ParentName)))
	}
	return items
}

func skillCategoryFromParts(id uuid.UUID, tenantID pgtype.UUID, parentID pgtype.UUID, code string, name string, description pgtype.Text, sourceScope string, sortOrder int32, metadataBytes []byte, inactive bool, createdAt pgtype.Timestamptz, createdBy pgtype.UUID, updatedAt pgtype.Timestamptz, updatedBy pgtype.UUID, parentName *string) *domain.SkillCategory {
	metadata := json.RawMessage(metadataBytes)
	if len(metadata) == 0 {
		metadata = json.RawMessage(`{}`)
	}
	return &domain.SkillCategory{ID: id, TenantID: ptrFromUUID(tenantID), ParentID: ptrFromUUID(parentID), Code: code, Name: name, Description: ptrFromText(description), SourceScope: sourceScope, SortOrder: sortOrder, Metadata: metadata, Inactive: inactive, CreatedAt: timeFromTimestamptz(createdAt), CreatedBy: ptrFromUUID(createdBy), UpdatedAt: timeFromTimestamptz(updatedAt), UpdatedBy: ptrFromUUID(updatedBy), ParentName: parentName}
}

func mapSkill(row sqlc.HrmsSkill) *domain.Skill {
	return skillFromParts(row.ID, row.TenantID, row.CategoryID, row.Code, row.Name, row.Description, row.SkillType, row.SourceScope, row.CertificateRequired, row.AssessmentRequired, row.IsActive, row.Metadata, row.Inactive, row.CreatedAt, row.CreatedBy, row.UpdatedAt, row.UpdatedBy, nil, nil)
}

func mapSkillList(rows []sqlc.ListSkillsRow) []*domain.Skill {
	items := make([]*domain.Skill, 0, len(rows))
	for _, row := range rows {
		items = append(items, skillFromParts(row.ID, row.TenantID, row.CategoryID, row.Code, row.Name, row.Description, row.SkillType, row.SourceScope, row.CertificateRequired, row.AssessmentRequired, row.IsActive, row.Metadata, row.Inactive, row.CreatedAt, row.CreatedBy, row.UpdatedAt, row.UpdatedBy, ptrFromText(row.CategoryName), ptrFromText(row.CategoryCode)))
	}
	return items
}

func skillFromParts(id uuid.UUID, tenantID pgtype.UUID, categoryID pgtype.UUID, code string, name string, description pgtype.Text, skillType string, sourceScope string, certificateRequired bool, assessmentRequired bool, isActive bool, metadataBytes []byte, inactive bool, createdAt pgtype.Timestamptz, createdBy pgtype.UUID, updatedAt pgtype.Timestamptz, updatedBy pgtype.UUID, categoryName *string, categoryCode *string) *domain.Skill {
	metadata := json.RawMessage(metadataBytes)
	if len(metadata) == 0 {
		metadata = json.RawMessage(`{}`)
	}
	return &domain.Skill{ID: id, TenantID: ptrFromUUID(tenantID), CategoryID: ptrFromUUID(categoryID), Code: code, Name: name, Description: ptrFromText(description), SkillType: skillType, SourceScope: sourceScope, CertificateRequired: certificateRequired, AssessmentRequired: assessmentRequired, IsActive: isActive, Metadata: metadata, Inactive: inactive, CreatedAt: timeFromTimestamptz(createdAt), CreatedBy: ptrFromUUID(createdBy), UpdatedAt: timeFromTimestamptz(updatedAt), UpdatedBy: ptrFromUUID(updatedBy), CategoryName: categoryName, CategoryCode: categoryCode}
}

func mapWorkerSkill(row sqlc.HrmsWorkerSkill) *domain.WorkerSkill {
	return workerSkillFromParts(row.ID, row.TenantID, row.WorkerProfileID, row.SkillID, row.SkillNameSnapshot, row.Proficiency, row.YearsExperience, row.LastUsedOn, row.VerificationStatus, row.CertificateUrl, row.CertificateExpiresOn, row.AssessmentScore, row.AssessedOn, row.EndorsedBy, row.EndorsedAt, row.VerifiedBy, row.VerifiedAt, row.Notes, row.Metadata, row.Inactive, row.CreatedAt, row.CreatedBy, row.UpdatedAt, row.UpdatedBy, nil, nil, nil, nil, nil, nil, false, false, nil, nil)
}

func mapWorkerSkillList(rows []sqlc.ListWorkerSkillsRow) []*domain.WorkerSkill {
	items := make([]*domain.WorkerSkill, 0, len(rows))
	for _, row := range rows {
		workerName := row.WorkerDisplayName
		skillCode := row.SkillCode
		skillName := row.SkillName
		skillType := row.SkillType
		sourceScope := row.SkillSourceScope
		items = append(items, workerSkillFromParts(row.ID, row.TenantID, row.WorkerProfileID, row.SkillID, row.SkillNameSnapshot, row.Proficiency, row.YearsExperience, row.LastUsedOn, row.VerificationStatus, row.CertificateUrl, row.CertificateExpiresOn, row.AssessmentScore, row.AssessedOn, row.EndorsedBy, row.EndorsedAt, row.VerifiedBy, row.VerifiedAt, row.Notes, row.Metadata, row.Inactive, row.CreatedAt, row.CreatedBy, row.UpdatedAt, row.UpdatedBy, &workerName, ptrFromText(row.WorkerCode), &skillCode, &skillName, &skillType, &sourceScope, row.CertificateRequired, row.AssessmentRequired, ptrFromText(row.CategoryName), ptrFromText(row.CategoryCode)))
	}
	return items
}

func workerSkillFromParts(id uuid.UUID, tenantID uuid.UUID, workerProfileID uuid.UUID, skillID uuid.UUID, skillNameSnapshot string, proficiency string, yearsExperience pgtype.Numeric, lastUsedOn pgtype.Date, verificationStatus string, certificateURL pgtype.Text, certificateExpiresOn pgtype.Date, assessmentScore pgtype.Numeric, assessedOn pgtype.Date, endorsedBy pgtype.UUID, endorsedAt pgtype.Timestamptz, verifiedBy pgtype.UUID, verifiedAt pgtype.Timestamptz, notes pgtype.Text, metadataBytes []byte, inactive bool, createdAt pgtype.Timestamptz, createdBy pgtype.UUID, updatedAt pgtype.Timestamptz, updatedBy pgtype.UUID, workerDisplayName *string, workerCode *string, skillCode *string, skillName *string, skillType *string, skillSourceScope *string, certificateRequired bool, assessmentRequired bool, categoryName *string, categoryCode *string) *domain.WorkerSkill {
	metadata := json.RawMessage(metadataBytes)
	if len(metadata) == 0 {
		metadata = json.RawMessage(`{}`)
	}
	return &domain.WorkerSkill{ID: id, TenantID: tenantID, WorkerProfileID: workerProfileID, SkillID: skillID, SkillNameSnapshot: skillNameSnapshot, Proficiency: proficiency, YearsExperience: ptrFromNumeric(yearsExperience), LastUsedOn: ptrFromDate(lastUsedOn), VerificationStatus: verificationStatus, CertificateURL: ptrFromText(certificateURL), CertificateExpiresOn: ptrFromDate(certificateExpiresOn), AssessmentScore: ptrFromNumeric(assessmentScore), AssessedOn: ptrFromDate(assessedOn), EndorsedBy: ptrFromUUID(endorsedBy), EndorsedAt: timePtrFromSkillTimestamptz(endorsedAt), VerifiedBy: ptrFromUUID(verifiedBy), VerifiedAt: timePtrFromSkillTimestamptz(verifiedAt), Notes: ptrFromText(notes), Metadata: metadata, Inactive: inactive, CreatedAt: timeFromTimestamptz(createdAt), CreatedBy: ptrFromUUID(createdBy), UpdatedAt: timeFromTimestamptz(updatedAt), UpdatedBy: ptrFromUUID(updatedBy), WorkerDisplayName: workerDisplayName, WorkerCode: workerCode, SkillCode: skillCode, SkillName: skillName, SkillType: skillType, SkillSourceScope: skillSourceScope, CertificateRequired: certificateRequired, AssessmentRequired: assessmentRequired, CategoryName: categoryName, CategoryCode: categoryCode}
}

func mapWorkerSkillAssessment(row sqlc.HrmsWorkerSkillAssessment) *domain.WorkerSkillAssessment {
	metadata := json.RawMessage(row.Metadata)
	if len(metadata) == 0 {
		metadata = json.RawMessage(`{}`)
	}
	return &domain.WorkerSkillAssessment{ID: row.ID, TenantID: row.TenantID, WorkerSkillID: row.WorkerSkillID, AssessmentType: row.AssessmentType, ResultStatus: row.ResultStatus, Score: ptrFromNumeric(row.Score), MaxScore: ptrFromNumeric(row.MaxScore), AssessedBy: ptrFromUUID(row.AssessedBy), AssessedOn: timeFromDate(row.AssessedOn), EvidenceURL: ptrFromText(row.EvidenceUrl), Notes: ptrFromText(row.Notes), Metadata: metadata, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy)}
}

func mapWorkerSkillAssessments(rows []sqlc.HrmsWorkerSkillAssessment) []*domain.WorkerSkillAssessment {
	items := make([]*domain.WorkerSkillAssessment, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapWorkerSkillAssessment(row))
	}
	return items
}

func mapSkillsSummary(rows []sqlc.GetSkillsSummaryRow) []*domain.SkillsSummaryRow {
	items := make([]*domain.SkillsSummaryRow, 0, len(rows))
	for _, row := range rows {
		items = append(items, &domain.SkillsSummaryRow{Status: row.Status, WorkerSkillCount: row.WorkerSkillCount, WorkerCount: row.WorkerCount, SkillCount: row.SkillCount, ExpiringCertificateCount: row.ExpiringCertificateCount})
	}
	return items
}

func numericFromSkillFloat(value *float64) pgtype.Numeric {
	if value == nil {
		return pgtype.Numeric{Valid: false}
	}
	scaled := int64(math.Round(*value * 100))
	return pgtype.Numeric{Int: big.NewInt(scaled), Exp: -2, Valid: true}
}

func boolFromSkillPtr(value *bool) pgtype.Bool {
	if value == nil {
		return pgtype.Bool{Valid: false}
	}
	return pgtype.Bool{Bool: *value, Valid: true}
}

func timePtrFromSkillTimestamptz(value pgtype.Timestamptz) *time.Time {
	if !value.Valid {
		return nil
	}
	clean := value.Time.UTC()
	return &clean
}
