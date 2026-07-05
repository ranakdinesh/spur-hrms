package services

import (
	"context"
	"encoding/base64"
	"strings"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) CreateLearningCourse(ctx context.Context, cmd ports.LearningCourseCommand) (*domain.LearningCourse, error) {
	item, err := s.prepareLearningCourse(ctx, cmd)
	if err != nil {
		return nil, err
	}
	result, err := s.learning.CreateLearningCourse(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create learning course", err, serviceTenantIDField(cmd.TenantID), serviceStringField("code", item.Code))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) UpdateLearningCourse(ctx context.Context, cmd ports.LearningCourseCommand) (*domain.LearningCourse, error) {
	if cmd.ID == uuid.Nil {
		return nil, domain.ErrInvalidLearningCourse
	}
	if _, err := s.learning.GetLearningCourse(ctx, cmd.TenantID, cmd.ID); err != nil {
		return nil, err
	}
	item, err := s.prepareLearningCourse(ctx, cmd)
	if err != nil {
		return nil, err
	}
	item.ID = cmd.ID
	result, err := s.learning.UpdateLearningCourse(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("update learning course", err, serviceTenantIDField(cmd.TenantID), serviceStringField("course_id", cmd.ID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) GetLearningCourse(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.LearningCourse, error) {
	if tenantID == uuid.Nil || id == uuid.Nil {
		return nil, domain.ErrInvalidLearningCourse
	}
	return s.learning.GetLearningCourse(ctx, tenantID, id)
}

func (s *TenantService) ListLearningCourses(ctx context.Context, filter domain.LearningCourseFilter) ([]*domain.LearningCourse, error) {
	if filter.TenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	normalizeLearningPage(&filter.Limit, &filter.Offset)
	return s.learning.ListLearningCourses(ctx, filter)
}

func (s *TenantService) DeleteLearningCourse(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil || id == uuid.Nil {
		return domain.ErrInvalidLearningCourse
	}
	return s.learning.DeleteLearningCourse(ctx, tenantID, id, actorID)
}

func (s *TenantService) CreateLearningPath(ctx context.Context, cmd ports.LearningPathCommand) (*domain.LearningPath, error) {
	item, err := s.prepareLearningPath(ctx, cmd)
	if err != nil {
		return nil, err
	}
	result, err := s.learning.CreateLearningPath(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create learning path", err, serviceTenantIDField(cmd.TenantID), serviceStringField("code", item.Code))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) UpdateLearningPath(ctx context.Context, cmd ports.LearningPathCommand) (*domain.LearningPath, error) {
	if cmd.ID == uuid.Nil {
		return nil, domain.ErrInvalidLearningPath
	}
	if _, err := s.learning.GetLearningPath(ctx, cmd.TenantID, cmd.ID); err != nil {
		return nil, err
	}
	item, err := s.prepareLearningPath(ctx, cmd)
	if err != nil {
		return nil, err
	}
	item.ID = cmd.ID
	result, err := s.learning.UpdateLearningPath(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("update learning path", err, serviceTenantIDField(cmd.TenantID), serviceStringField("path_id", cmd.ID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) GetLearningPath(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.LearningPath, error) {
	if tenantID == uuid.Nil || id == uuid.Nil {
		return nil, domain.ErrInvalidLearningPath
	}
	return s.learning.GetLearningPath(ctx, tenantID, id)
}

func (s *TenantService) ListLearningPaths(ctx context.Context, filter domain.LearningPathFilter) ([]*domain.LearningPath, error) {
	if filter.TenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	normalizeLearningPage(&filter.Limit, &filter.Offset)
	return s.learning.ListLearningPaths(ctx, filter)
}

func (s *TenantService) DeleteLearningPath(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil || id == uuid.Nil {
		return domain.ErrInvalidLearningPath
	}
	return s.learning.DeleteLearningPath(ctx, tenantID, id, actorID)
}

func (s *TenantService) UpsertLearningPathCourse(ctx context.Context, cmd ports.LearningPathCourseCommand) (*domain.LearningPathCourse, error) {
	if cmd.TenantID == uuid.Nil || cmd.PathID == uuid.Nil || cmd.CourseID == uuid.Nil {
		return nil, domain.ErrInvalidLearningPathCourse
	}
	if _, err := s.learning.GetLearningPath(ctx, cmd.TenantID, cmd.PathID); err != nil {
		return nil, err
	}
	if _, err := s.learning.GetLearningCourse(ctx, cmd.TenantID, cmd.CourseID); err != nil {
		return nil, err
	}
	item := &domain.LearningPathCourse{TenantID: cmd.TenantID, PathID: cmd.PathID, CourseID: cmd.CourseID, SortOrder: cmd.SortOrder, Required: cmd.Required}
	return s.learning.UpsertLearningPathCourse(ctx, item, cmd.ActorID)
}

func (s *TenantService) ListLearningPathCourses(ctx context.Context, tenantID uuid.UUID, pathID uuid.UUID) ([]*domain.LearningPathCourse, error) {
	if tenantID == uuid.Nil || pathID == uuid.Nil {
		return nil, domain.ErrInvalidLearningPathCourse
	}
	return s.learning.ListLearningPathCourses(ctx, tenantID, pathID)
}

func (s *TenantService) DeleteLearningPathCourse(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil || id == uuid.Nil {
		return domain.ErrInvalidLearningPathCourse
	}
	return s.learning.DeleteLearningPathCourse(ctx, tenantID, id, actorID)
}

func (s *TenantService) CreateLearningEnrollment(ctx context.Context, cmd ports.LearningEnrollmentCommand) (*domain.LearningEnrollment, error) {
	if _, err := s.learning.GetLearningCourse(ctx, cmd.TenantID, cmd.CourseID); err != nil {
		return nil, err
	}
	if cmd.PathID != nil {
		if _, err := s.learning.GetLearningPath(ctx, cmd.TenantID, *cmd.PathID); err != nil {
			return nil, err
		}
	}
	if _, err := s.workerProfiles.GetWorkerProfile(ctx, cmd.TenantID, cmd.WorkerProfileID); err != nil {
		return nil, err
	}
	if cmd.AssignedBy == nil {
		cmd.AssignedBy = cmd.ActorID
	}
	item, err := domain.NewLearningEnrollment(domain.LearningEnrollmentInput{TenantID: cmd.TenantID, CourseID: cmd.CourseID, PathID: cmd.PathID, WorkerProfileID: cmd.WorkerProfileID, AssignmentSource: cmd.AssignmentSource, Status: cmd.Status, NominatedBy: cmd.NominatedBy, AssignedBy: cmd.AssignedBy, DueDate: cmd.DueDate, Notes: cmd.Notes, Metadata: cmd.Metadata})
	if err != nil {
		s.logError("validate learning enrollment", err, serviceTenantIDField(cmd.TenantID), serviceStringField("worker_profile_id", cmd.WorkerProfileID.String()))
		return nil, err
	}
	result, err := s.learning.CreateLearningEnrollment(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create learning enrollment", err, serviceTenantIDField(cmd.TenantID), serviceStringField("course_id", cmd.CourseID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) UpdateLearningEnrollmentStatus(ctx context.Context, cmd ports.LearningEnrollmentStatusCommand) (*domain.LearningEnrollment, error) {
	if cmd.TenantID == uuid.Nil || cmd.ID == uuid.Nil {
		return nil, domain.ErrInvalidLearningEnrollment
	}
	status := domain.NormalizeLearningEnrollmentStatus(cmd.Status)
	if status == "" {
		return nil, domain.ErrInvalidLearningEnrollment
	}
	result, err := s.learning.UpdateLearningEnrollmentStatus(ctx, cmd.TenantID, cmd.ID, status, cmd.Score, cmd.CertificateURL, cmd.CertificateFileName, cmd.CertificateContentType, cmd.Notes, cmd.ActorID)
	if err != nil {
		s.logError("update learning enrollment status", err, serviceTenantIDField(cmd.TenantID), serviceStringField("enrollment_id", cmd.ID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) UploadLearningCertificate(ctx context.Context, cmd ports.LearningCertificateCommand) (*domain.LearningEnrollment, error) {
	if cmd.TenantID == uuid.Nil || cmd.EnrollmentID == uuid.Nil || strings.TrimSpace(cmd.FileName) == "" || strings.TrimSpace(cmd.ContentBase64) == "" {
		return nil, domain.ErrInvalidLearningEnrollment
	}
	enrollment, err := s.learning.GetLearningEnrollment(ctx, cmd.TenantID, cmd.EnrollmentID)
	if err != nil {
		return nil, err
	}
	content, err := base64.StdEncoding.DecodeString(cmd.ContentBase64)
	if err != nil {
		s.logError("decode learning certificate", err, serviceTenantIDField(cmd.TenantID), serviceStringField("enrollment_id", cmd.EnrollmentID.String()))
		return nil, domain.ErrInvalidLearningEnrollment
	}
	if s.learningCertificateStorage == nil {
		return nil, domain.ErrStorageProviderSettingsNotFound
	}
	contentType := strings.TrimSpace(cmd.ContentType)
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	url, err := s.learningCertificateStorage.StoreLearningCertificate(ctx, ports.StoreLearningCertificateInput{TenantID: cmd.TenantID, EnrollmentID: cmd.EnrollmentID, WorkerProfileID: enrollment.WorkerProfileID, FileName: strings.TrimSpace(cmd.FileName), ContentType: contentType, Content: content})
	if err != nil {
		s.logError("store learning certificate", err, serviceTenantIDField(cmd.TenantID), serviceStringField("enrollment_id", cmd.EnrollmentID.String()))
		return nil, err
	}
	return s.UpdateLearningEnrollmentStatus(ctx, ports.LearningEnrollmentStatusCommand{TenantID: cmd.TenantID, ID: cmd.EnrollmentID, Status: domain.LearningEnrollmentCompleted, Score: cmd.Score, CertificateURL: &url, CertificateFileName: &cmd.FileName, CertificateContentType: &contentType, Notes: cmd.Notes, ActorID: cmd.ActorID})
}

func (s *TenantService) ListLearningEnrollments(ctx context.Context, filter domain.LearningEnrollmentFilter) ([]*domain.LearningEnrollment, error) {
	if filter.TenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	normalizeLearningPage(&filter.Limit, &filter.Offset)
	return s.learning.ListLearningEnrollments(ctx, filter)
}

func (s *TenantService) DeleteLearningEnrollment(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil || id == uuid.Nil {
		return domain.ErrInvalidLearningEnrollment
	}
	return s.learning.DeleteLearningEnrollment(ctx, tenantID, id, actorID)
}

func (s *TenantService) CreateLearningRecommendation(ctx context.Context, cmd ports.LearningRecommendationCommand) (*domain.LearningRecommendation, error) {
	item, err := domain.NewLearningRecommendation(domain.LearningRecommendationInput{TenantID: cmd.TenantID, WorkerProfileID: cmd.WorkerProfileID, SkillID: cmd.SkillID, CourseID: cmd.CourseID, PathID: cmd.PathID, SourceType: cmd.SourceType, Reason: cmd.Reason, Priority: cmd.Priority, ConfidenceScore: cmd.ConfidenceScore, Status: cmd.Status, Metadata: cmd.Metadata})
	if err != nil {
		s.logError("validate learning recommendation", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	result, err := s.learning.CreateLearningRecommendation(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create learning recommendation", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) UpdateLearningRecommendationStatus(ctx context.Context, cmd ports.LearningRecommendationStatusCommand) (*domain.LearningRecommendation, error) {
	if cmd.TenantID == uuid.Nil || cmd.ID == uuid.Nil {
		return nil, domain.ErrInvalidLearningRecommendation
	}
	status := domain.NormalizeLearningRecommendationStatus(cmd.Status)
	if status == "" {
		return nil, domain.ErrInvalidLearningRecommendation
	}
	return s.learning.UpdateLearningRecommendationStatus(ctx, cmd.TenantID, cmd.ID, status, cmd.ActorID)
}

func (s *TenantService) ListLearningRecommendations(ctx context.Context, filter domain.LearningRecommendationFilter) ([]*domain.LearningRecommendation, error) {
	if filter.TenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	normalizeLearningPage(&filter.Limit, &filter.Offset)
	return s.learning.ListLearningRecommendations(ctx, filter)
}

func (s *TenantService) GenerateSkillGapLearningRecommendations(ctx context.Context, tenantID uuid.UUID, actorID *uuid.UUID) ([]*domain.LearningRecommendation, error) {
	if tenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	return s.learning.GenerateSkillGapLearningRecommendations(ctx, tenantID, actorID)
}

func (s *TenantService) GetLearningSummary(ctx context.Context, tenantID uuid.UUID) ([]*domain.LearningSummaryRow, error) {
	if tenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	return s.learning.GetLearningSummary(ctx, tenantID)
}

func (s *TenantService) prepareLearningCourse(ctx context.Context, cmd ports.LearningCourseCommand) (*domain.LearningCourse, error) {
	if cmd.SkillID != nil {
		if _, err := s.skills.GetSkill(ctx, cmd.TenantID, *cmd.SkillID); err != nil {
			return nil, err
		}
	}
	if cmd.ComplianceRuleID != nil {
		if _, err := s.compliance.GetComplianceRule(ctx, cmd.TenantID, *cmd.ComplianceRuleID); err != nil {
			return nil, err
		}
	}
	item, err := domain.NewLearningCourse(domain.LearningCourseInput{TenantID: cmd.TenantID, Code: cmd.Code, Title: cmd.Title, Description: cmd.Description, CourseType: cmd.CourseType, DeliveryMode: cmd.DeliveryMode, Provider: cmd.Provider, DurationMinutes: cmd.DurationMinutes, SkillID: cmd.SkillID, ComplianceRuleID: cmd.ComplianceRuleID, Mandatory: cmd.Mandatory, AIReadiness: cmd.AIReadiness, CertificateRequired: cmd.CertificateRequired, BudgetAmount: cmd.BudgetAmount, CurrencyCode: cmd.CurrencyCode, IsActive: cmd.IsActive, Metadata: cmd.Metadata})
	if err != nil {
		s.logError("validate learning course", err, serviceTenantIDField(cmd.TenantID), serviceStringField("code", cmd.Code))
		return nil, err
	}
	return item, nil
}

func (s *TenantService) prepareLearningPath(ctx context.Context, cmd ports.LearningPathCommand) (*domain.LearningPath, error) {
	if cmd.SkillID != nil {
		if _, err := s.skills.GetSkill(ctx, cmd.TenantID, *cmd.SkillID); err != nil {
			return nil, err
		}
	}
	item, err := domain.NewLearningPath(domain.LearningPathInput{TenantID: cmd.TenantID, Code: cmd.Code, Title: cmd.Title, Description: cmd.Description, PathType: cmd.PathType, TargetRole: cmd.TargetRole, SkillID: cmd.SkillID, AIReadiness: cmd.AIReadiness, IsActive: cmd.IsActive, Metadata: cmd.Metadata})
	if err != nil {
		s.logError("validate learning path", err, serviceTenantIDField(cmd.TenantID), serviceStringField("code", cmd.Code))
		return nil, err
	}
	return item, nil
}

func normalizeLearningPage(limit *int32, offset *int32) {
	if *limit <= 0 || *limit > 200 {
		*limit = 100
	}
	if *offset < 0 {
		*offset = 0
	}
}
