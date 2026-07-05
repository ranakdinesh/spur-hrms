package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (s *Store) CreateSkillCategory(ctx context.Context, item *domain.SkillCategory, actorID *uuid.UUID) (*domain.SkillCategory, error) {
	row, err := s.getQueries(ctx).CreateSkillCategory(ctx, sqlc.CreateSkillCategoryParams{TenantID: uuidFromPtr(item.TenantID), ParentID: uuidFromPtr(item.ParentID), Code: item.Code, Name: item.Name, Description: textFromPtr(item.Description), SortOrder: item.SortOrder, Metadata: jsonBytesFromRaw(item.Metadata), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create skill category", err, optionalTenantIDField(item.TenantID), stringField("code", item.Code))
	}
	return mapSkillCategory(row), nil
}

func (s *Store) UpdateSkillCategory(ctx context.Context, item *domain.SkillCategory, actorID *uuid.UUID) (*domain.SkillCategory, error) {
	row, err := s.getQueries(ctx).UpdateSkillCategory(ctx, sqlc.UpdateSkillCategoryParams{TenantID: uuidFromPtr(item.TenantID), ID: item.ID, ParentID: uuidFromPtr(item.ParentID), Code: item.Code, Name: item.Name, Description: textFromPtr(item.Description), SortOrder: item.SortOrder, Metadata: jsonBytesFromRaw(item.Metadata), UpdatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "update skill category", err, optionalTenantIDField(item.TenantID), stringField("skill_category_id", item.ID.String()))
	}
	return mapSkillCategory(row), nil
}

func (s *Store) GetSkillCategory(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.SkillCategory, error) {
	row, err := s.getQueries(ctx).GetSkillCategory(ctx, sqlc.GetSkillCategoryParams{TenantID: uuidFromPtr(&tenantID), ID: id})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrSkillCategoryNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "get skill category", err, tenantIDField(tenantID), stringField("skill_category_id", id.String()))
	}
	return mapSkillCategory(row), nil
}

func (s *Store) ListSkillCategories(ctx context.Context, filter domain.SkillCategoryFilter) ([]*domain.SkillCategory, error) {
	rows, err := s.getQueries(ctx).ListSkillCategories(ctx, sqlc.ListSkillCategoriesParams{TenantID: uuidFromPtr(&filter.TenantID), SourceScope: textFromPtr(filter.SourceScope), ParentID: uuidFromPtr(filter.ParentID), Search: textFromPtr(filter.Search)})
	if err != nil {
		return nil, s.logDBError(ctx, "list skill categories", err, tenantIDField(filter.TenantID))
	}
	return mapSkillCategoryList(rows), nil
}

func (s *Store) DeleteSkillCategory(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteSkillCategory(ctx, sqlc.SoftDeleteSkillCategoryParams{TenantID: uuidFromPtr(&tenantID), ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete skill category", err, tenantIDField(tenantID), stringField("skill_category_id", id.String()))
	}
	return nil
}

func (s *Store) CreateSkill(ctx context.Context, item *domain.Skill, actorID *uuid.UUID) (*domain.Skill, error) {
	row, err := s.getQueries(ctx).CreateSkill(ctx, sqlc.CreateSkillParams{TenantID: uuidFromPtr(item.TenantID), CategoryID: uuidFromPtr(item.CategoryID), Code: item.Code, Name: item.Name, Description: textFromPtr(item.Description), SkillType: item.SkillType, CertificateRequired: item.CertificateRequired, AssessmentRequired: item.AssessmentRequired, IsActive: item.IsActive, Metadata: jsonBytesFromRaw(item.Metadata), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create skill", err, optionalTenantIDField(item.TenantID), stringField("code", item.Code))
	}
	return mapSkill(row), nil
}

func (s *Store) UpdateSkill(ctx context.Context, item *domain.Skill, actorID *uuid.UUID) (*domain.Skill, error) {
	row, err := s.getQueries(ctx).UpdateSkill(ctx, sqlc.UpdateSkillParams{TenantID: uuidFromPtr(item.TenantID), ID: item.ID, CategoryID: uuidFromPtr(item.CategoryID), Code: item.Code, Name: item.Name, Description: textFromPtr(item.Description), SkillType: item.SkillType, CertificateRequired: item.CertificateRequired, AssessmentRequired: item.AssessmentRequired, IsActive: item.IsActive, Metadata: jsonBytesFromRaw(item.Metadata), UpdatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "update skill", err, optionalTenantIDField(item.TenantID), stringField("skill_id", item.ID.String()))
	}
	return mapSkill(row), nil
}

func (s *Store) GetSkill(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.Skill, error) {
	row, err := s.getQueries(ctx).GetSkill(ctx, sqlc.GetSkillParams{TenantID: uuidFromPtr(&tenantID), ID: id})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrSkillNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "get skill", err, tenantIDField(tenantID), stringField("skill_id", id.String()))
	}
	return mapSkill(row), nil
}

func (s *Store) GetSkillByCode(ctx context.Context, tenantID uuid.UUID, code string) (*domain.Skill, error) {
	row, err := s.getQueries(ctx).GetSkillByCode(ctx, sqlc.GetSkillByCodeParams{TenantID: uuidFromPtr(&tenantID), Lower: code})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrSkillNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "get skill by code", err, tenantIDField(tenantID), stringField("code", code))
	}
	return mapSkill(row), nil
}

func (s *Store) ListSkills(ctx context.Context, filter domain.SkillFilter) ([]*domain.Skill, error) {
	rows, err := s.getQueries(ctx).ListSkills(ctx, sqlc.ListSkillsParams{TenantID: uuidFromPtr(&filter.TenantID), CategoryID: uuidFromPtr(filter.CategoryID), SkillType: textFromPtr(filter.SkillType), SourceScope: textFromPtr(filter.SourceScope), IsActive: boolFromSkillPtr(filter.IsActive), Search: textFromPtr(filter.Search)})
	if err != nil {
		return nil, s.logDBError(ctx, "list skills", err, tenantIDField(filter.TenantID))
	}
	return mapSkillList(rows), nil
}

func (s *Store) DeleteSkill(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteSkill(ctx, sqlc.SoftDeleteSkillParams{TenantID: uuidFromPtr(&tenantID), ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete skill", err, tenantIDField(tenantID), stringField("skill_id", id.String()))
	}
	return nil
}

func (s *Store) CreateWorkerSkill(ctx context.Context, item *domain.WorkerSkill, actorID *uuid.UUID) (*domain.WorkerSkill, error) {
	row, err := s.getQueries(ctx).CreateWorkerSkill(ctx, workerSkillCreateParams(item, actorID))
	if err != nil {
		return nil, s.logDBError(ctx, "create worker skill", err, tenantIDField(item.TenantID), stringField("worker_profile_id", item.WorkerProfileID.String()), stringField("skill_id", item.SkillID.String()))
	}
	return mapWorkerSkill(row), nil
}

func (s *Store) UpdateWorkerSkill(ctx context.Context, item *domain.WorkerSkill, actorID *uuid.UUID) (*domain.WorkerSkill, error) {
	row, err := s.getQueries(ctx).UpdateWorkerSkill(ctx, workerSkillUpdateParams(item, actorID))
	if err != nil {
		return nil, s.logDBError(ctx, "update worker skill", err, tenantIDField(item.TenantID), stringField("worker_skill_id", item.ID.String()))
	}
	return mapWorkerSkill(row), nil
}

func (s *Store) GetWorkerSkill(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.WorkerSkill, error) {
	row, err := s.getQueries(ctx).GetWorkerSkill(ctx, sqlc.GetWorkerSkillParams{TenantID: tenantID, ID: id})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrWorkerSkillNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "get worker skill", err, tenantIDField(tenantID), stringField("worker_skill_id", id.String()))
	}
	return mapWorkerSkill(row), nil
}

func (s *Store) ListWorkerSkills(ctx context.Context, filter domain.WorkerSkillFilter) ([]*domain.WorkerSkill, error) {
	rows, err := s.getQueries(ctx).ListWorkerSkills(ctx, sqlc.ListWorkerSkillsParams{TenantID: filter.TenantID, WorkerProfileID: uuidFromPtr(filter.WorkerProfileID), SkillID: uuidFromPtr(filter.SkillID), CategoryID: uuidFromPtr(filter.CategoryID), Proficiency: textFromPtr(filter.Proficiency), VerificationStatus: textFromPtr(filter.VerificationStatus), CertificateExpiringBefore: dateFromPtr(filter.CertificateExpiringBefore), Search: textFromPtr(filter.Search)})
	if err != nil {
		return nil, s.logDBError(ctx, "list worker skills", err, tenantIDField(filter.TenantID))
	}
	return mapWorkerSkillList(rows), nil
}

func (s *Store) UpdateWorkerSkillVerification(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, notes *string, actorID *uuid.UUID) (*domain.WorkerSkill, error) {
	row, err := s.getQueries(ctx).UpdateWorkerSkillVerification(ctx, sqlc.UpdateWorkerSkillVerificationParams{TenantID: tenantID, ID: id, VerificationStatus: status, UpdatedBy: uuidFromPtr(actorID), Notes: textFromPtr(notes)})
	if err != nil {
		return nil, s.logDBError(ctx, "update worker skill verification", err, tenantIDField(tenantID), stringField("worker_skill_id", id.String()), stringField("status", status))
	}
	return mapWorkerSkill(row), nil
}

func (s *Store) DeleteWorkerSkill(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteWorkerSkill(ctx, sqlc.SoftDeleteWorkerSkillParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete worker skill", err, tenantIDField(tenantID), stringField("worker_skill_id", id.String()))
	}
	return nil
}

func (s *Store) CreateWorkerSkillAssessment(ctx context.Context, item *domain.WorkerSkillAssessment, actorID *uuid.UUID) (*domain.WorkerSkillAssessment, error) {
	row, err := s.getQueries(ctx).CreateWorkerSkillAssessment(ctx, sqlc.CreateWorkerSkillAssessmentParams{TenantID: item.TenantID, WorkerSkillID: item.WorkerSkillID, AssessmentType: item.AssessmentType, ResultStatus: item.ResultStatus, Score: numericFromSkillFloat(item.Score), MaxScore: numericFromSkillFloat(item.MaxScore), AssessedBy: uuidFromPtr(item.AssessedBy), AssessedOn: dateFromTime(item.AssessedOn), EvidenceUrl: textFromPtr(item.EvidenceURL), Notes: textFromPtr(item.Notes), Metadata: jsonBytesFromRaw(item.Metadata), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create worker skill assessment", err, tenantIDField(item.TenantID), stringField("worker_skill_id", item.WorkerSkillID.String()))
	}
	return mapWorkerSkillAssessment(row), nil
}

func (s *Store) ListWorkerSkillAssessments(ctx context.Context, tenantID uuid.UUID, workerSkillID *uuid.UUID) ([]*domain.WorkerSkillAssessment, error) {
	rows, err := s.getQueries(ctx).ListWorkerSkillAssessments(ctx, sqlc.ListWorkerSkillAssessmentsParams{TenantID: tenantID, WorkerSkillID: uuidFromPtr(workerSkillID)})
	if err != nil {
		return nil, s.logDBError(ctx, "list worker skill assessments", err, tenantIDField(tenantID))
	}
	return mapWorkerSkillAssessments(rows), nil
}

func (s *Store) GetSkillsSummary(ctx context.Context, tenantID uuid.UUID) ([]*domain.SkillsSummaryRow, error) {
	rows, err := s.getQueries(ctx).GetSkillsSummary(ctx, tenantID)
	if err != nil {
		return nil, s.logDBError(ctx, "get skills summary", err, tenantIDField(tenantID))
	}
	return mapSkillsSummary(rows), nil
}

func workerSkillCreateParams(item *domain.WorkerSkill, actorID *uuid.UUID) sqlc.CreateWorkerSkillParams {
	return sqlc.CreateWorkerSkillParams{TenantID: item.TenantID, WorkerProfileID: item.WorkerProfileID, SkillID: item.SkillID, SkillNameSnapshot: item.SkillNameSnapshot, Proficiency: item.Proficiency, YearsExperience: numericFromSkillFloat(item.YearsExperience), LastUsedOn: dateFromPtr(item.LastUsedOn), VerificationStatus: item.VerificationStatus, CertificateUrl: textFromPtr(item.CertificateURL), CertificateExpiresOn: dateFromPtr(item.CertificateExpiresOn), AssessmentScore: numericFromSkillFloat(item.AssessmentScore), AssessedOn: dateFromPtr(item.AssessedOn), Notes: textFromPtr(item.Notes), Metadata: jsonBytesFromRaw(item.Metadata), CreatedBy: uuidFromPtr(actorID)}
}

func workerSkillUpdateParams(item *domain.WorkerSkill, actorID *uuid.UUID) sqlc.UpdateWorkerSkillParams {
	return sqlc.UpdateWorkerSkillParams{TenantID: item.TenantID, ID: item.ID, SkillID: item.SkillID, SkillNameSnapshot: item.SkillNameSnapshot, Proficiency: item.Proficiency, YearsExperience: numericFromSkillFloat(item.YearsExperience), LastUsedOn: dateFromPtr(item.LastUsedOn), VerificationStatus: item.VerificationStatus, CertificateUrl: textFromPtr(item.CertificateURL), CertificateExpiresOn: dateFromPtr(item.CertificateExpiresOn), AssessmentScore: numericFromSkillFloat(item.AssessmentScore), AssessedOn: dateFromPtr(item.AssessedOn), Notes: textFromPtr(item.Notes), Metadata: jsonBytesFromRaw(item.Metadata), UpdatedBy: uuidFromPtr(actorID)}
}
