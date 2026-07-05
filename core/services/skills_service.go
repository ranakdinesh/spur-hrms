package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) CreateSkillCategory(ctx context.Context, cmd ports.SkillCategoryCommand) (*domain.SkillCategory, error) {
	item, err := domain.NewSkillCategory(domain.SkillCategoryInput{TenantID: cmd.TenantID, ParentID: cmd.ParentID, Code: cmd.Code, Name: cmd.Name, Description: cmd.Description, SortOrder: cmd.SortOrder, Metadata: cmd.Metadata})
	if err != nil {
		s.logError("validate skill category create", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	if item.ParentID != nil {
		if _, err := s.skills.GetSkillCategory(ctx, cmd.TenantID, *item.ParentID); err != nil {
			s.logError("validate skill category parent", err, serviceTenantIDField(cmd.TenantID), serviceStringField("parent_id", item.ParentID.String()))
			return nil, err
		}
	}
	result, err := s.skills.CreateSkillCategory(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create skill category", err, serviceTenantIDField(cmd.TenantID), serviceStringField("code", item.Code))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) UpdateSkillCategory(ctx context.Context, cmd ports.SkillCategoryCommand) (*domain.SkillCategory, error) {
	if cmd.ID == uuid.Nil {
		return nil, domain.ErrInvalidSkillCategory
	}
	if _, err := s.skills.GetSkillCategory(ctx, cmd.TenantID, cmd.ID); err != nil {
		return nil, err
	}
	item, err := domain.NewSkillCategory(domain.SkillCategoryInput{TenantID: cmd.TenantID, ParentID: cmd.ParentID, Code: cmd.Code, Name: cmd.Name, Description: cmd.Description, SortOrder: cmd.SortOrder, Metadata: cmd.Metadata})
	if err != nil {
		s.logError("validate skill category update", err, serviceTenantIDField(cmd.TenantID), serviceStringField("skill_category_id", cmd.ID.String()))
		return nil, err
	}
	if item.ParentID != nil && *item.ParentID == cmd.ID {
		return nil, domain.ErrInvalidSkillCategory
	}
	item.ID = cmd.ID
	result, err := s.skills.UpdateSkillCategory(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("update skill category", err, serviceTenantIDField(cmd.TenantID), serviceStringField("skill_category_id", cmd.ID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) GetSkillCategory(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.SkillCategory, error) {
	if tenantID == uuid.Nil || id == uuid.Nil {
		return nil, domain.ErrInvalidSkillCategory
	}
	return s.skills.GetSkillCategory(ctx, tenantID, id)
}

func (s *TenantService) ListSkillCategories(ctx context.Context, filter domain.SkillCategoryFilter) ([]*domain.SkillCategory, error) {
	if filter.TenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	return s.skills.ListSkillCategories(ctx, filter)
}

func (s *TenantService) DeleteSkillCategory(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil || id == uuid.Nil {
		return domain.ErrInvalidSkillCategory
	}
	return s.skills.DeleteSkillCategory(ctx, tenantID, id, actorID)
}

func (s *TenantService) CreateSkill(ctx context.Context, cmd ports.SkillCommand) (*domain.Skill, error) {
	item, err := s.prepareSkill(ctx, cmd)
	if err != nil {
		return nil, err
	}
	result, err := s.skills.CreateSkill(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create skill", err, serviceTenantIDField(cmd.TenantID), serviceStringField("code", item.Code))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) UpdateSkill(ctx context.Context, cmd ports.SkillCommand) (*domain.Skill, error) {
	if cmd.ID == uuid.Nil {
		return nil, domain.ErrInvalidSkill
	}
	if _, err := s.skills.GetSkill(ctx, cmd.TenantID, cmd.ID); err != nil {
		return nil, err
	}
	item, err := s.prepareSkill(ctx, cmd)
	if err != nil {
		return nil, err
	}
	item.ID = cmd.ID
	result, err := s.skills.UpdateSkill(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("update skill", err, serviceTenantIDField(cmd.TenantID), serviceStringField("skill_id", cmd.ID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) GetSkill(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.Skill, error) {
	if tenantID == uuid.Nil || id == uuid.Nil {
		return nil, domain.ErrInvalidSkill
	}
	return s.skills.GetSkill(ctx, tenantID, id)
}

func (s *TenantService) ListSkills(ctx context.Context, filter domain.SkillFilter) ([]*domain.Skill, error) {
	if filter.TenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	return s.skills.ListSkills(ctx, filter)
}

func (s *TenantService) DeleteSkill(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil || id == uuid.Nil {
		return domain.ErrInvalidSkill
	}
	return s.skills.DeleteSkill(ctx, tenantID, id, actorID)
}

func (s *TenantService) CreateWorkerSkill(ctx context.Context, cmd ports.WorkerSkillCommand) (*domain.WorkerSkill, error) {
	item, err := s.prepareWorkerSkill(ctx, cmd)
	if err != nil {
		return nil, err
	}
	result, err := s.skills.CreateWorkerSkill(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create worker skill", err, serviceTenantIDField(cmd.TenantID), serviceStringField("worker_profile_id", cmd.WorkerProfileID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) UpdateWorkerSkill(ctx context.Context, cmd ports.WorkerSkillCommand) (*domain.WorkerSkill, error) {
	if cmd.ID == uuid.Nil {
		return nil, domain.ErrInvalidWorkerSkill
	}
	if _, err := s.skills.GetWorkerSkill(ctx, cmd.TenantID, cmd.ID); err != nil {
		return nil, err
	}
	item, err := s.prepareWorkerSkill(ctx, cmd)
	if err != nil {
		return nil, err
	}
	item.ID = cmd.ID
	result, err := s.skills.UpdateWorkerSkill(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("update worker skill", err, serviceTenantIDField(cmd.TenantID), serviceStringField("worker_skill_id", cmd.ID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) GetWorkerSkill(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.WorkerSkill, error) {
	if tenantID == uuid.Nil || id == uuid.Nil {
		return nil, domain.ErrInvalidWorkerSkill
	}
	return s.skills.GetWorkerSkill(ctx, tenantID, id)
}

func (s *TenantService) ListWorkerSkills(ctx context.Context, filter domain.WorkerSkillFilter) ([]*domain.WorkerSkill, error) {
	if filter.TenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	return s.skills.ListWorkerSkills(ctx, filter)
}

func (s *TenantService) UpdateWorkerSkillVerification(ctx context.Context, cmd ports.WorkerSkillVerificationCommand) (*domain.WorkerSkill, error) {
	if cmd.TenantID == uuid.Nil || cmd.ID == uuid.Nil {
		return nil, domain.ErrInvalidWorkerSkill
	}
	status := domain.NormalizeSkillVerificationStatus(cmd.Status)
	if status == "" {
		return nil, domain.ErrInvalidWorkerSkill
	}
	result, err := s.skills.UpdateWorkerSkillVerification(ctx, cmd.TenantID, cmd.ID, status, cmd.Notes, cmd.ActorID)
	if err != nil {
		s.logError("update worker skill verification", err, serviceTenantIDField(cmd.TenantID), serviceStringField("worker_skill_id", cmd.ID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) DeleteWorkerSkill(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil || id == uuid.Nil {
		return domain.ErrInvalidWorkerSkill
	}
	return s.skills.DeleteWorkerSkill(ctx, tenantID, id, actorID)
}

func (s *TenantService) CreateWorkerSkillAssessment(ctx context.Context, cmd ports.WorkerSkillAssessmentCommand) (*domain.WorkerSkillAssessment, error) {
	if _, err := s.skills.GetWorkerSkill(ctx, cmd.TenantID, cmd.WorkerSkillID); err != nil {
		return nil, err
	}
	assessedBy := cmd.AssessedBy
	if assessedBy == nil {
		assessedBy = cmd.ActorID
	}
	item, err := domain.NewWorkerSkillAssessment(domain.WorkerSkillAssessmentInput{TenantID: cmd.TenantID, WorkerSkillID: cmd.WorkerSkillID, AssessmentType: cmd.AssessmentType, ResultStatus: cmd.ResultStatus, Score: cmd.Score, MaxScore: cmd.MaxScore, AssessedBy: assessedBy, AssessedOn: cmd.AssessedOn, EvidenceURL: cmd.EvidenceURL, Notes: cmd.Notes, Metadata: cmd.Metadata})
	if err != nil {
		s.logError("validate worker skill assessment", err, serviceTenantIDField(cmd.TenantID), serviceStringField("worker_skill_id", cmd.WorkerSkillID.String()))
		return nil, err
	}
	result, err := s.skills.CreateWorkerSkillAssessment(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create worker skill assessment", err, serviceTenantIDField(cmd.TenantID), serviceStringField("worker_skill_id", cmd.WorkerSkillID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) ListWorkerSkillAssessments(ctx context.Context, tenantID uuid.UUID, workerSkillID *uuid.UUID) ([]*domain.WorkerSkillAssessment, error) {
	if tenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	return s.skills.ListWorkerSkillAssessments(ctx, tenantID, workerSkillID)
}

func (s *TenantService) GetSkillsSummary(ctx context.Context, tenantID uuid.UUID) ([]*domain.SkillsSummaryRow, error) {
	if tenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	return s.skills.GetSkillsSummary(ctx, tenantID)
}

func (s *TenantService) prepareSkill(ctx context.Context, cmd ports.SkillCommand) (*domain.Skill, error) {
	if cmd.CategoryID != nil && *cmd.CategoryID != uuid.Nil {
		if _, err := s.skills.GetSkillCategory(ctx, cmd.TenantID, *cmd.CategoryID); err != nil {
			return nil, err
		}
	}
	return domain.NewSkill(domain.SkillInput{TenantID: cmd.TenantID, CategoryID: cmd.CategoryID, Code: cmd.Code, Name: cmd.Name, Description: cmd.Description, SkillType: cmd.SkillType, CertificateRequired: cmd.CertificateRequired, AssessmentRequired: cmd.AssessmentRequired, IsActive: cmd.IsActive, Metadata: cmd.Metadata})
}

func (s *TenantService) prepareWorkerSkill(ctx context.Context, cmd ports.WorkerSkillCommand) (*domain.WorkerSkill, error) {
	if _, err := s.GetWorkerProfile(ctx, cmd.TenantID, cmd.WorkerProfileID); err != nil {
		return nil, err
	}
	skill, err := s.skills.GetSkill(ctx, cmd.TenantID, cmd.SkillID)
	if err != nil {
		return nil, err
	}
	return domain.NewWorkerSkill(domain.WorkerSkillInput{TenantID: cmd.TenantID, WorkerProfileID: cmd.WorkerProfileID, SkillID: cmd.SkillID, SkillNameSnapshot: skill.Name, Proficiency: cmd.Proficiency, YearsExperience: cmd.YearsExperience, LastUsedOn: cmd.LastUsedOn, VerificationStatus: cmd.VerificationStatus, CertificateURL: cmd.CertificateURL, CertificateExpiresOn: cmd.CertificateExpiresOn, AssessmentScore: cmd.AssessmentScore, AssessedOn: cmd.AssessedOn, Notes: cmd.Notes, Metadata: cmd.Metadata})
}
