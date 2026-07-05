package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (s *Store) CreateLearningCourse(ctx context.Context, item *domain.LearningCourse, actorID *uuid.UUID) (*domain.LearningCourse, error) {
	row, err := s.getQueries(ctx).CreateLearningCourse(ctx, learningCourseCreateParams(item, actorID))
	if err != nil {
		return nil, s.logDBError(ctx, "create learning course", err, tenantIDField(item.TenantID), stringField("code", item.Code))
	}
	return mapLearningCourse(row), nil
}

func (s *Store) UpdateLearningCourse(ctx context.Context, item *domain.LearningCourse, actorID *uuid.UUID) (*domain.LearningCourse, error) {
	row, err := s.getQueries(ctx).UpdateLearningCourse(ctx, learningCourseUpdateParams(item, actorID))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrLearningCourseNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "update learning course", err, tenantIDField(item.TenantID), stringField("course_id", item.ID.String()))
	}
	return mapLearningCourse(row), nil
}

func (s *Store) GetLearningCourse(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.LearningCourse, error) {
	row, err := s.getQueries(ctx).GetLearningCourse(ctx, sqlc.GetLearningCourseParams{TenantID: tenantID, ID: id})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrLearningCourseNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "get learning course", err, tenantIDField(tenantID), stringField("course_id", id.String()))
	}
	return mapLearningCourse(row), nil
}

func (s *Store) ListLearningCourses(ctx context.Context, filter domain.LearningCourseFilter) ([]*domain.LearningCourse, error) {
	rows, err := s.getQueries(ctx).ListLearningCourses(ctx, sqlc.ListLearningCoursesParams{TenantID: filter.TenantID, CourseType: textFromPtr(filter.CourseType), SkillID: uuidFromPtr(filter.SkillID), Mandatory: boolFromSkillPtr(filter.Mandatory), AiReadiness: boolFromSkillPtr(filter.AIReadiness), IsActive: boolFromSkillPtr(filter.IsActive), Search: textFromPtr(filter.Search), Limit: filter.Limit, Offset: filter.Offset})
	if err != nil {
		return nil, s.logDBError(ctx, "list learning courses", err, tenantIDField(filter.TenantID))
	}
	return mapLearningCourseList(rows), nil
}

func (s *Store) DeleteLearningCourse(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteLearningCourse(ctx, sqlc.SoftDeleteLearningCourseParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete learning course", err, tenantIDField(tenantID), stringField("course_id", id.String()))
	}
	return nil
}

func (s *Store) CreateLearningPath(ctx context.Context, item *domain.LearningPath, actorID *uuid.UUID) (*domain.LearningPath, error) {
	row, err := s.getQueries(ctx).CreateLearningPath(ctx, sqlc.CreateLearningPathParams{TenantID: item.TenantID, Code: item.Code, Title: item.Title, Description: textFromPtr(item.Description), PathType: item.PathType, TargetRole: textFromPtr(item.TargetRole), SkillID: uuidFromPtr(item.SkillID), AiReadiness: item.AIReadiness, IsActive: item.IsActive, Metadata: jsonBytesFromRaw(item.Metadata), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create learning path", err, tenantIDField(item.TenantID), stringField("code", item.Code))
	}
	return mapLearningPath(row), nil
}

func (s *Store) UpdateLearningPath(ctx context.Context, item *domain.LearningPath, actorID *uuid.UUID) (*domain.LearningPath, error) {
	row, err := s.getQueries(ctx).UpdateLearningPath(ctx, sqlc.UpdateLearningPathParams{TenantID: item.TenantID, ID: item.ID, Code: item.Code, Title: item.Title, Description: textFromPtr(item.Description), PathType: item.PathType, TargetRole: textFromPtr(item.TargetRole), SkillID: uuidFromPtr(item.SkillID), AiReadiness: item.AIReadiness, IsActive: item.IsActive, Metadata: jsonBytesFromRaw(item.Metadata), UpdatedBy: uuidFromPtr(actorID)})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrLearningPathNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "update learning path", err, tenantIDField(item.TenantID), stringField("path_id", item.ID.String()))
	}
	return mapLearningPath(row), nil
}

func (s *Store) GetLearningPath(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.LearningPath, error) {
	row, err := s.getQueries(ctx).GetLearningPath(ctx, sqlc.GetLearningPathParams{TenantID: tenantID, ID: id})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrLearningPathNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "get learning path", err, tenantIDField(tenantID), stringField("path_id", id.String()))
	}
	return mapLearningPath(row), nil
}

func (s *Store) ListLearningPaths(ctx context.Context, filter domain.LearningPathFilter) ([]*domain.LearningPath, error) {
	rows, err := s.getQueries(ctx).ListLearningPaths(ctx, sqlc.ListLearningPathsParams{TenantID: filter.TenantID, PathType: textFromPtr(filter.PathType), SkillID: uuidFromPtr(filter.SkillID), AiReadiness: boolFromSkillPtr(filter.AIReadiness), IsActive: boolFromSkillPtr(filter.IsActive), Search: textFromPtr(filter.Search), Limit: filter.Limit, Offset: filter.Offset})
	if err != nil {
		return nil, s.logDBError(ctx, "list learning paths", err, tenantIDField(filter.TenantID))
	}
	return mapLearningPathList(rows), nil
}

func (s *Store) DeleteLearningPath(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteLearningPath(ctx, sqlc.SoftDeleteLearningPathParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete learning path", err, tenantIDField(tenantID), stringField("path_id", id.String()))
	}
	return nil
}

func (s *Store) UpsertLearningPathCourse(ctx context.Context, item *domain.LearningPathCourse, actorID *uuid.UUID) (*domain.LearningPathCourse, error) {
	row, err := s.getQueries(ctx).UpsertLearningPathCourse(ctx, sqlc.UpsertLearningPathCourseParams{TenantID: item.TenantID, PathID: item.PathID, CourseID: item.CourseID, SortOrder: item.SortOrder, Required: item.Required, CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "upsert learning path course", err, tenantIDField(item.TenantID), stringField("path_id", item.PathID.String()), stringField("course_id", item.CourseID.String()))
	}
	return mapLearningPathCourse(row), nil
}

func (s *Store) ListLearningPathCourses(ctx context.Context, tenantID uuid.UUID, pathID uuid.UUID) ([]*domain.LearningPathCourse, error) {
	rows, err := s.getQueries(ctx).ListLearningPathCourses(ctx, sqlc.ListLearningPathCoursesParams{TenantID: tenantID, PathID: pathID})
	if err != nil {
		return nil, s.logDBError(ctx, "list learning path courses", err, tenantIDField(tenantID), stringField("path_id", pathID.String()))
	}
	return mapLearningPathCourseList(rows), nil
}

func (s *Store) DeleteLearningPathCourse(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteLearningPathCourse(ctx, sqlc.SoftDeleteLearningPathCourseParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete learning path course", err, tenantIDField(tenantID), stringField("path_course_id", id.String()))
	}
	return nil
}

func (s *Store) CreateLearningEnrollment(ctx context.Context, item *domain.LearningEnrollment, actorID *uuid.UUID) (*domain.LearningEnrollment, error) {
	row, err := s.getQueries(ctx).CreateLearningEnrollment(ctx, sqlc.CreateLearningEnrollmentParams{TenantID: item.TenantID, CourseID: item.CourseID, PathID: uuidFromPtr(item.PathID), WorkerProfileID: item.WorkerProfileID, AssignmentSource: item.AssignmentSource, Status: item.Status, NominatedBy: uuidFromPtr(item.NominatedBy), AssignedBy: uuidFromPtr(item.AssignedBy), DueDate: dateFromPtr(item.DueDate), Notes: textFromPtr(item.Notes), Metadata: jsonBytesFromRaw(item.Metadata), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create learning enrollment", err, tenantIDField(item.TenantID), stringField("worker_profile_id", item.WorkerProfileID.String()), stringField("course_id", item.CourseID.String()))
	}
	return mapLearningEnrollment(row), nil
}

func (s *Store) UpdateLearningEnrollmentStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, score *float64, certificateURL *string, certificateFileName *string, certificateContentType *string, notes *string, actorID *uuid.UUID) (*domain.LearningEnrollment, error) {
	row, err := s.getQueries(ctx).UpdateLearningEnrollmentStatus(ctx, sqlc.UpdateLearningEnrollmentStatusParams{TenantID: tenantID, ID: id, Status: status, Score: numericFromSkillFloat(score), CertificateUrl: textFromPtr(certificateURL), CertificateFileName: textFromPtr(certificateFileName), CertificateContentType: textFromPtr(certificateContentType), Notes: textFromPtr(notes), UpdatedBy: uuidFromPtr(actorID)})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrLearningEnrollmentNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "update learning enrollment status", err, tenantIDField(tenantID), stringField("enrollment_id", id.String()), stringField("status", status))
	}
	return mapLearningEnrollment(row), nil
}

func (s *Store) GetLearningEnrollment(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.LearningEnrollment, error) {
	row, err := s.getQueries(ctx).GetLearningEnrollment(ctx, sqlc.GetLearningEnrollmentParams{TenantID: tenantID, ID: id})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrLearningEnrollmentNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "get learning enrollment", err, tenantIDField(tenantID), stringField("enrollment_id", id.String()))
	}
	return mapLearningEnrollment(row), nil
}

func (s *Store) ListLearningEnrollments(ctx context.Context, filter domain.LearningEnrollmentFilter) ([]*domain.LearningEnrollment, error) {
	rows, err := s.getQueries(ctx).ListLearningEnrollments(ctx, sqlc.ListLearningEnrollmentsParams{TenantID: filter.TenantID, WorkerProfileID: uuidFromPtr(filter.WorkerProfileID), CourseID: uuidFromPtr(filter.CourseID), Status: textFromPtr(filter.Status), AssignmentSource: textFromPtr(filter.AssignmentSource), Search: textFromPtr(filter.Search), Limit: filter.Limit, Offset: filter.Offset})
	if err != nil {
		return nil, s.logDBError(ctx, "list learning enrollments", err, tenantIDField(filter.TenantID))
	}
	return mapLearningEnrollmentList(rows), nil
}

func (s *Store) DeleteLearningEnrollment(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteLearningEnrollment(ctx, sqlc.SoftDeleteLearningEnrollmentParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete learning enrollment", err, tenantIDField(tenantID), stringField("enrollment_id", id.String()))
	}
	return nil
}

func (s *Store) CreateLearningRecommendation(ctx context.Context, item *domain.LearningRecommendation, actorID *uuid.UUID) (*domain.LearningRecommendation, error) {
	row, err := s.getQueries(ctx).CreateLearningRecommendation(ctx, sqlc.CreateLearningRecommendationParams{TenantID: item.TenantID, WorkerProfileID: uuidFromPtr(item.WorkerProfileID), SkillID: uuidFromPtr(item.SkillID), CourseID: uuidFromPtr(item.CourseID), PathID: uuidFromPtr(item.PathID), SourceType: item.SourceType, Reason: item.Reason, Priority: item.Priority, ConfidenceScore: numericFromSkillFloat(item.ConfidenceScore), Status: item.Status, Metadata: jsonBytesFromRaw(item.Metadata), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create learning recommendation", err, tenantIDField(item.TenantID), stringField("source_type", item.SourceType))
	}
	return mapLearningRecommendation(row), nil
}

func (s *Store) UpdateLearningRecommendationStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, actorID *uuid.UUID) (*domain.LearningRecommendation, error) {
	row, err := s.getQueries(ctx).UpdateLearningRecommendationStatus(ctx, sqlc.UpdateLearningRecommendationStatusParams{TenantID: tenantID, ID: id, Status: status, UpdatedBy: uuidFromPtr(actorID)})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrLearningRecommendationNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "update learning recommendation status", err, tenantIDField(tenantID), stringField("recommendation_id", id.String()), stringField("status", status))
	}
	return mapLearningRecommendation(row), nil
}

func (s *Store) ListLearningRecommendations(ctx context.Context, filter domain.LearningRecommendationFilter) ([]*domain.LearningRecommendation, error) {
	rows, err := s.getQueries(ctx).ListLearningRecommendations(ctx, sqlc.ListLearningRecommendationsParams{TenantID: filter.TenantID, WorkerProfileID: uuidFromPtr(filter.WorkerProfileID), SkillID: uuidFromPtr(filter.SkillID), SourceType: textFromPtr(filter.SourceType), Status: textFromPtr(filter.Status), Search: textFromPtr(filter.Search), Limit: filter.Limit, Offset: filter.Offset})
	if err != nil {
		return nil, s.logDBError(ctx, "list learning recommendations", err, tenantIDField(filter.TenantID))
	}
	return mapLearningRecommendationList(rows), nil
}

func (s *Store) GenerateSkillGapLearningRecommendations(ctx context.Context, tenantID uuid.UUID, actorID *uuid.UUID) ([]*domain.LearningRecommendation, error) {
	rows, err := s.getQueries(ctx).GenerateSkillGapLearningRecommendations(ctx, sqlc.GenerateSkillGapLearningRecommendationsParams{TenantID: tenantID, CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "generate skill gap learning recommendations", err, tenantIDField(tenantID))
	}
	return mapLearningRecommendations(rows), nil
}

func (s *Store) GetLearningSummary(ctx context.Context, tenantID uuid.UUID) ([]*domain.LearningSummaryRow, error) {
	rows, err := s.getQueries(ctx).GetLearningSummary(ctx, tenantID)
	if err != nil {
		return nil, s.logDBError(ctx, "get learning summary", err, tenantIDField(tenantID))
	}
	return mapLearningSummary(rows), nil
}

func learningCourseCreateParams(item *domain.LearningCourse, actorID *uuid.UUID) sqlc.CreateLearningCourseParams {
	return sqlc.CreateLearningCourseParams{TenantID: item.TenantID, Code: item.Code, Title: item.Title, Description: textFromPtr(item.Description), CourseType: item.CourseType, DeliveryMode: item.DeliveryMode, Provider: textFromPtr(item.Provider), DurationMinutes: item.DurationMinutes, SkillID: uuidFromPtr(item.SkillID), ComplianceRuleID: uuidFromPtr(item.ComplianceRuleID), Mandatory: item.Mandatory, AiReadiness: item.AIReadiness, CertificateRequired: item.CertificateRequired, BudgetAmount: numericFromSkillFloat(item.BudgetAmount), CurrencyCode: item.CurrencyCode, IsActive: item.IsActive, Metadata: jsonBytesFromRaw(item.Metadata), CreatedBy: uuidFromPtr(actorID)}
}

func learningCourseUpdateParams(item *domain.LearningCourse, actorID *uuid.UUID) sqlc.UpdateLearningCourseParams {
	return sqlc.UpdateLearningCourseParams{TenantID: item.TenantID, ID: item.ID, Code: item.Code, Title: item.Title, Description: textFromPtr(item.Description), CourseType: item.CourseType, DeliveryMode: item.DeliveryMode, Provider: textFromPtr(item.Provider), DurationMinutes: item.DurationMinutes, SkillID: uuidFromPtr(item.SkillID), ComplianceRuleID: uuidFromPtr(item.ComplianceRuleID), Mandatory: item.Mandatory, AiReadiness: item.AIReadiness, CertificateRequired: item.CertificateRequired, BudgetAmount: numericFromSkillFloat(item.BudgetAmount), CurrencyCode: item.CurrencyCode, IsActive: item.IsActive, Metadata: jsonBytesFromRaw(item.Metadata), UpdatedBy: uuidFromPtr(actorID)}
}
