package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (s *Store) CreateBenefitPlan(ctx context.Context, item *domain.BenefitPlan, actorID *uuid.UUID) (*domain.BenefitPlan, error) {
	row, err := s.getQueries(ctx).CreateBenefitPlan(ctx, sqlc.CreateBenefitPlanParams{TenantID: item.TenantID, Code: item.Code, Name: item.Name, PlanType: item.PlanType, Description: textFromPtr(item.Description), ProviderName: textFromPtr(item.ProviderName), PolicyNumber: textFromPtr(item.PolicyNumber), CoverageAmount: numericFromFloatPtr(item.CoverageAmount), EmployerContribution: compNumeric(item.EmployerContribution), EmployeeContribution: compNumeric(item.EmployeeContribution), CurrencyCode: item.CurrencyCode, EligibilityRule: jsonBytesFromRaw(item.EligibilityRule), InsuranceMetadata: jsonBytesFromRaw(item.InsuranceMetadata), EffectiveFrom: dateFromPtr(item.EffectiveFrom), EffectiveTo: dateFromPtr(item.EffectiveTo), IsActive: item.IsActive, CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create benefit plan", err, tenantIDField(item.TenantID), stringField("code", item.Code))
	}
	return mapBenefitPlan(row), nil
}

func (s *Store) UpdateBenefitPlan(ctx context.Context, item *domain.BenefitPlan, actorID *uuid.UUID) (*domain.BenefitPlan, error) {
	row, err := s.getQueries(ctx).UpdateBenefitPlan(ctx, sqlc.UpdateBenefitPlanParams{TenantID: item.TenantID, ID: item.ID, Code: item.Code, Name: item.Name, PlanType: item.PlanType, Description: textFromPtr(item.Description), ProviderName: textFromPtr(item.ProviderName), PolicyNumber: textFromPtr(item.PolicyNumber), CoverageAmount: numericFromFloatPtr(item.CoverageAmount), EmployerContribution: compNumeric(item.EmployerContribution), EmployeeContribution: compNumeric(item.EmployeeContribution), CurrencyCode: item.CurrencyCode, EligibilityRule: jsonBytesFromRaw(item.EligibilityRule), InsuranceMetadata: jsonBytesFromRaw(item.InsuranceMetadata), EffectiveFrom: dateFromPtr(item.EffectiveFrom), EffectiveTo: dateFromPtr(item.EffectiveTo), IsActive: item.IsActive, UpdatedBy: uuidFromPtr(actorID)})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrBenefitPlanNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "update benefit plan", err, tenantIDField(item.TenantID), stringField("benefit_plan_id", item.ID.String()))
	}
	return mapBenefitPlan(row), nil
}

func (s *Store) ListBenefitPlans(ctx context.Context, filter domain.BenefitFilter) ([]*domain.BenefitPlan, error) {
	rows, err := s.getQueries(ctx).ListBenefitPlans(ctx, sqlc.ListBenefitPlansParams{TenantID: filter.TenantID, Limit: limitOrDefault(filter.Limit), Offset: filter.Offset, ActiveOnly: boolFromPtr(filter.ActiveOnly), PlanType: textFromPtr(filter.PlanType), Search: textFromPtr(filter.Search)})
	if err != nil {
		return nil, s.logDBError(ctx, "list benefit plans", err, tenantIDField(filter.TenantID))
	}
	return mapBenefitPlans(rows), nil
}

func (s *Store) GetBenefitPlan(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.BenefitPlan, error) {
	row, err := s.getQueries(ctx).GetBenefitPlan(ctx, sqlc.GetBenefitPlanParams{TenantID: tenantID, ID: id})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrBenefitPlanNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "get benefit plan", err, tenantIDField(tenantID), stringField("benefit_plan_id", id.String()))
	}
	return mapBenefitPlan(row), nil
}

func (s *Store) DeleteBenefitPlan(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteBenefitPlan(ctx, sqlc.SoftDeleteBenefitPlanParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete benefit plan", err, tenantIDField(tenantID), stringField("benefit_plan_id", id.String()))
	}
	return nil
}

func (s *Store) CreateBenefitEnrollmentWindow(ctx context.Context, item *domain.BenefitEnrollmentWindow, actorID *uuid.UUID) (*domain.BenefitEnrollmentWindow, error) {
	row, err := s.getQueries(ctx).CreateBenefitEnrollmentWindow(ctx, sqlc.CreateBenefitEnrollmentWindowParams{TenantID: item.TenantID, PlanID: item.PlanID, Name: item.Name, OpensOn: dateFromTime(item.OpensOn), ClosesOn: dateFromTime(item.ClosesOn), EffectiveFrom: dateFromPtr(item.EffectiveFrom), EffectiveTo: dateFromPtr(item.EffectiveTo), Status: item.Status, Metadata: jsonBytesFromRaw(item.Metadata), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create benefit enrollment window", err, tenantIDField(item.TenantID), stringField("plan_id", item.PlanID.String()))
	}
	return mapBenefitEnrollmentWindow(row), nil
}

func (s *Store) UpdateBenefitEnrollmentWindow(ctx context.Context, item *domain.BenefitEnrollmentWindow, actorID *uuid.UUID) (*domain.BenefitEnrollmentWindow, error) {
	row, err := s.getQueries(ctx).UpdateBenefitEnrollmentWindow(ctx, sqlc.UpdateBenefitEnrollmentWindowParams{TenantID: item.TenantID, ID: item.ID, PlanID: item.PlanID, Name: item.Name, OpensOn: dateFromTime(item.OpensOn), ClosesOn: dateFromTime(item.ClosesOn), EffectiveFrom: dateFromPtr(item.EffectiveFrom), EffectiveTo: dateFromPtr(item.EffectiveTo), Status: item.Status, Metadata: jsonBytesFromRaw(item.Metadata), UpdatedBy: uuidFromPtr(actorID)})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrBenefitEnrollmentWindowNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "update benefit enrollment window", err, tenantIDField(item.TenantID), stringField("benefit_window_id", item.ID.String()))
	}
	return mapBenefitEnrollmentWindow(row), nil
}

func (s *Store) ListBenefitEnrollmentWindows(ctx context.Context, filter domain.BenefitFilter) ([]*domain.BenefitEnrollmentWindow, error) {
	rows, err := s.getQueries(ctx).ListBenefitEnrollmentWindows(ctx, sqlc.ListBenefitEnrollmentWindowsParams{TenantID: filter.TenantID, Limit: limitOrDefault(filter.Limit), Offset: filter.Offset, PlanID: uuidFromPtr(filter.PlanID), Status: textFromPtr(filter.Status)})
	if err != nil {
		return nil, s.logDBError(ctx, "list benefit enrollment windows", err, tenantIDField(filter.TenantID))
	}
	return mapBenefitEnrollmentWindows(rows), nil
}

func (s *Store) GetBenefitEnrollmentWindow(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.BenefitEnrollmentWindow, error) {
	row, err := s.getQueries(ctx).GetBenefitEnrollmentWindow(ctx, sqlc.GetBenefitEnrollmentWindowParams{TenantID: tenantID, ID: id})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrBenefitEnrollmentWindowNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "get benefit enrollment window", err, tenantIDField(tenantID), stringField("benefit_window_id", id.String()))
	}
	return mapBenefitEnrollmentWindow(row), nil
}

func (s *Store) DeleteBenefitEnrollmentWindow(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteBenefitEnrollmentWindow(ctx, sqlc.SoftDeleteBenefitEnrollmentWindowParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete benefit enrollment window", err, tenantIDField(tenantID), stringField("benefit_window_id", id.String()))
	}
	return nil
}

func (s *Store) CreateBenefitDependent(ctx context.Context, item *domain.BenefitDependent, actorID *uuid.UUID) (*domain.BenefitDependent, error) {
	row, err := s.getQueries(ctx).CreateBenefitDependent(ctx, sqlc.CreateBenefitDependentParams{TenantID: item.TenantID, EmployeeUserID: item.EmployeeUserID, FullName: item.FullName, Relationship: item.Relationship, DateOfBirth: dateFromPtr(item.DateOfBirth), Gender: textFromPtr(item.Gender), NomineePercentage: numericFromFloatPtr(item.NomineePercentage), IsNominee: item.IsNominee, Metadata: jsonBytesFromRaw(item.Metadata), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create benefit dependent", err, tenantIDField(item.TenantID), stringField("employee_user_id", item.EmployeeUserID.String()))
	}
	return mapBenefitDependent(row), nil
}

func (s *Store) UpdateBenefitDependent(ctx context.Context, item *domain.BenefitDependent, actorID *uuid.UUID) (*domain.BenefitDependent, error) {
	row, err := s.getQueries(ctx).UpdateBenefitDependent(ctx, sqlc.UpdateBenefitDependentParams{TenantID: item.TenantID, ID: item.ID, EmployeeUserID: item.EmployeeUserID, FullName: item.FullName, Relationship: item.Relationship, DateOfBirth: dateFromPtr(item.DateOfBirth), Gender: textFromPtr(item.Gender), NomineePercentage: numericFromFloatPtr(item.NomineePercentage), IsNominee: item.IsNominee, Metadata: jsonBytesFromRaw(item.Metadata), UpdatedBy: uuidFromPtr(actorID)})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrBenefitDependentNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "update benefit dependent", err, tenantIDField(item.TenantID), stringField("benefit_dependent_id", item.ID.String()))
	}
	return mapBenefitDependent(row), nil
}

func (s *Store) ListBenefitDependents(ctx context.Context, filter domain.BenefitFilter) ([]*domain.BenefitDependent, error) {
	rows, err := s.getQueries(ctx).ListBenefitDependents(ctx, sqlc.ListBenefitDependentsParams{TenantID: filter.TenantID, Limit: limitOrDefault(filter.Limit), Offset: filter.Offset, EmployeeUserID: uuidFromPtr(filter.EmployeeUserID), NomineesOnly: boolFromPtr(filter.NomineesOnly)})
	if err != nil {
		return nil, s.logDBError(ctx, "list benefit dependents", err, tenantIDField(filter.TenantID))
	}
	return mapBenefitDependents(rows), nil
}

func (s *Store) GetBenefitDependent(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.BenefitDependent, error) {
	row, err := s.getQueries(ctx).GetBenefitDependent(ctx, sqlc.GetBenefitDependentParams{TenantID: tenantID, ID: id})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrBenefitDependentNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "get benefit dependent", err, tenantIDField(tenantID), stringField("benefit_dependent_id", id.String()))
	}
	return mapBenefitDependent(row), nil
}

func (s *Store) DeleteBenefitDependent(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteBenefitDependent(ctx, sqlc.SoftDeleteBenefitDependentParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete benefit dependent", err, tenantIDField(tenantID), stringField("benefit_dependent_id", id.String()))
	}
	return nil
}

func (s *Store) CreateBenefitEnrollment(ctx context.Context, item *domain.BenefitEnrollment, actorID *uuid.UUID) (*domain.BenefitEnrollment, error) {
	row, err := s.getQueries(ctx).CreateBenefitEnrollment(ctx, sqlc.CreateBenefitEnrollmentParams{TenantID: item.TenantID, PlanID: item.PlanID, WindowID: uuidFromPtr(item.WindowID), EmployeeUserID: item.EmployeeUserID, Status: item.Status, CoverageLevel: textFromPtr(item.CoverageLevel), SelectedAmount: numericFromFloatPtr(item.SelectedAmount), EmployeeContribution: compNumeric(item.EmployeeContribution), EmployerContribution: compNumeric(item.EmployerContribution), EffectiveFrom: dateFromPtr(item.EffectiveFrom), EffectiveTo: dateFromPtr(item.EffectiveTo), SubmittedAt: timestamptzFromPtr(item.SubmittedAt), ReviewRemarks: textFromPtr(item.ReviewRemarks), Metadata: jsonBytesFromRaw(item.Metadata), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create benefit enrollment", err, tenantIDField(item.TenantID), stringField("employee_user_id", item.EmployeeUserID.String()))
	}
	return mapBenefitEnrollment(row), nil
}

func (s *Store) UpdateBenefitEnrollment(ctx context.Context, item *domain.BenefitEnrollment, actorID *uuid.UUID) (*domain.BenefitEnrollment, error) {
	row, err := s.getQueries(ctx).UpdateBenefitEnrollment(ctx, sqlc.UpdateBenefitEnrollmentParams{TenantID: item.TenantID, ID: item.ID, PlanID: item.PlanID, WindowID: uuidFromPtr(item.WindowID), EmployeeUserID: item.EmployeeUserID, Status: item.Status, CoverageLevel: textFromPtr(item.CoverageLevel), SelectedAmount: numericFromFloatPtr(item.SelectedAmount), EmployeeContribution: compNumeric(item.EmployeeContribution), EmployerContribution: compNumeric(item.EmployerContribution), EffectiveFrom: dateFromPtr(item.EffectiveFrom), EffectiveTo: dateFromPtr(item.EffectiveTo), SubmittedAt: timestamptzFromPtr(item.SubmittedAt), ReviewRemarks: textFromPtr(item.ReviewRemarks), Metadata: jsonBytesFromRaw(item.Metadata), UpdatedBy: uuidFromPtr(actorID)})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrBenefitEnrollmentNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "update benefit enrollment", err, tenantIDField(item.TenantID), stringField("benefit_enrollment_id", item.ID.String()))
	}
	return mapBenefitEnrollment(row), nil
}

func (s *Store) UpdateBenefitEnrollmentStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, reviewerID *uuid.UUID, remarks *string) (*domain.BenefitEnrollment, error) {
	row, err := s.getQueries(ctx).UpdateBenefitEnrollmentStatus(ctx, sqlc.UpdateBenefitEnrollmentStatusParams{TenantID: tenantID, ID: id, Status: status, ReviewedBy: uuidFromPtr(reviewerID), ReviewRemarks: textFromPtr(remarks)})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrBenefitEnrollmentNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "update benefit enrollment status", err, tenantIDField(tenantID), stringField("benefit_enrollment_id", id.String()), stringField("status", status))
	}
	return mapBenefitEnrollment(row), nil
}

func (s *Store) ListBenefitEnrollments(ctx context.Context, filter domain.BenefitFilter) ([]*domain.BenefitEnrollment, error) {
	rows, err := s.getQueries(ctx).ListBenefitEnrollments(ctx, sqlc.ListBenefitEnrollmentsParams{TenantID: filter.TenantID, Limit: limitOrDefault(filter.Limit), Offset: filter.Offset, EmployeeUserID: uuidFromPtr(filter.EmployeeUserID), PlanID: uuidFromPtr(filter.PlanID), Status: textFromPtr(filter.Status)})
	if err != nil {
		return nil, s.logDBError(ctx, "list benefit enrollments", err, tenantIDField(filter.TenantID))
	}
	return mapBenefitEnrollments(rows), nil
}

func (s *Store) GetBenefitEnrollment(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.BenefitEnrollment, error) {
	row, err := s.getQueries(ctx).GetBenefitEnrollment(ctx, sqlc.GetBenefitEnrollmentParams{TenantID: tenantID, ID: id})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrBenefitEnrollmentNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "get benefit enrollment", err, tenantIDField(tenantID), stringField("benefit_enrollment_id", id.String()))
	}
	return mapBenefitEnrollment(row), nil
}

func (s *Store) DeleteBenefitEnrollment(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteBenefitEnrollment(ctx, sqlc.SoftDeleteBenefitEnrollmentParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete benefit enrollment", err, tenantIDField(tenantID), stringField("benefit_enrollment_id", id.String()))
	}
	return nil
}

func (s *Store) CreateBenefitClaimType(ctx context.Context, item *domain.BenefitClaimType, actorID *uuid.UUID) (*domain.BenefitClaimType, error) {
	row, err := s.getQueries(ctx).CreateBenefitClaimType(ctx, sqlc.CreateBenefitClaimTypeParams{TenantID: item.TenantID, PlanID: uuidFromPtr(item.PlanID), Code: item.Code, Name: item.Name, Description: textFromPtr(item.Description), AnnualLimit: numericFromFloatPtr(item.AnnualLimit), PerClaimLimit: numericFromFloatPtr(item.PerClaimLimit), RequiresAttachment: item.RequiresAttachment, Taxable: item.Taxable, PayrollComponentCode: textFromPtr(item.PayrollComponentCode), EligibilityRule: jsonBytesFromRaw(item.EligibilityRule), IsActive: item.IsActive, CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create benefit claim type", err, tenantIDField(item.TenantID), stringField("code", item.Code))
	}
	return mapBenefitClaimType(row), nil
}

func (s *Store) UpdateBenefitClaimType(ctx context.Context, item *domain.BenefitClaimType, actorID *uuid.UUID) (*domain.BenefitClaimType, error) {
	row, err := s.getQueries(ctx).UpdateBenefitClaimType(ctx, sqlc.UpdateBenefitClaimTypeParams{TenantID: item.TenantID, ID: item.ID, PlanID: uuidFromPtr(item.PlanID), Code: item.Code, Name: item.Name, Description: textFromPtr(item.Description), AnnualLimit: numericFromFloatPtr(item.AnnualLimit), PerClaimLimit: numericFromFloatPtr(item.PerClaimLimit), RequiresAttachment: item.RequiresAttachment, Taxable: item.Taxable, PayrollComponentCode: textFromPtr(item.PayrollComponentCode), EligibilityRule: jsonBytesFromRaw(item.EligibilityRule), IsActive: item.IsActive, UpdatedBy: uuidFromPtr(actorID)})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrBenefitClaimTypeNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "update benefit claim type", err, tenantIDField(item.TenantID), stringField("benefit_claim_type_id", item.ID.String()))
	}
	return mapBenefitClaimType(row), nil
}

func (s *Store) ListBenefitClaimTypes(ctx context.Context, filter domain.BenefitFilter) ([]*domain.BenefitClaimType, error) {
	rows, err := s.getQueries(ctx).ListBenefitClaimTypes(ctx, sqlc.ListBenefitClaimTypesParams{TenantID: filter.TenantID, Limit: limitOrDefault(filter.Limit), Offset: filter.Offset, ActiveOnly: boolFromPtr(filter.ActiveOnly), PlanID: uuidFromPtr(filter.PlanID), Search: textFromPtr(filter.Search)})
	if err != nil {
		return nil, s.logDBError(ctx, "list benefit claim types", err, tenantIDField(filter.TenantID))
	}
	return mapBenefitClaimTypes(rows), nil
}

func (s *Store) GetBenefitClaimType(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.BenefitClaimType, error) {
	row, err := s.getQueries(ctx).GetBenefitClaimType(ctx, sqlc.GetBenefitClaimTypeParams{TenantID: tenantID, ID: id})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrBenefitClaimTypeNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "get benefit claim type", err, tenantIDField(tenantID), stringField("benefit_claim_type_id", id.String()))
	}
	return mapBenefitClaimType(row), nil
}

func (s *Store) DeleteBenefitClaimType(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteBenefitClaimType(ctx, sqlc.SoftDeleteBenefitClaimTypeParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete benefit claim type", err, tenantIDField(tenantID), stringField("benefit_claim_type_id", id.String()))
	}
	return nil
}

func (s *Store) CreateBenefitClaim(ctx context.Context, item *domain.BenefitClaim, actorID *uuid.UUID) (*domain.BenefitClaim, error) {
	row, err := s.getQueries(ctx).CreateBenefitClaim(ctx, benefitClaimCreateParams(item, actorID))
	if err != nil {
		return nil, s.logDBError(ctx, "create benefit claim", err, tenantIDField(item.TenantID), stringField("employee_user_id", item.EmployeeUserID.String()))
	}
	return mapBenefitClaim(row), nil
}

func (s *Store) UpdateBenefitClaim(ctx context.Context, item *domain.BenefitClaim, actorID *uuid.UUID) (*domain.BenefitClaim, error) {
	row, err := s.getQueries(ctx).UpdateBenefitClaim(ctx, benefitClaimUpdateParams(item, actorID))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrBenefitClaimNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "update benefit claim", err, tenantIDField(item.TenantID), stringField("benefit_claim_id", item.ID.String()))
	}
	return mapBenefitClaim(row), nil
}

func (s *Store) UpdateBenefitClaimStatus(ctx context.Context, item *domain.BenefitClaim, actorID *uuid.UUID) (*domain.BenefitClaim, error) {
	row, err := s.getQueries(ctx).UpdateBenefitClaimStatus(ctx, sqlc.UpdateBenefitClaimStatusParams{TenantID: item.TenantID, ID: item.ID, Status: item.Status, ApprovedAmount: numericFromFloatPtr(item.ApprovedAmount), PaymentStatus: item.PaymentStatus, PaymentReference: textFromPtr(item.PaymentReference), PaidAt: timestamptzFromPtr(item.PaidAt), ReviewedBy: uuidFromPtr(actorID), ReviewRemarks: textFromPtr(item.ReviewRemarks), PayrollExportStatus: item.PayrollExportStatus, PayrollExportedAt: timestamptzFromPtr(item.PayrollExportedAt), PayrollExportReference: textFromPtr(item.PayrollExportReference)})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrBenefitClaimNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "update benefit claim status", err, tenantIDField(item.TenantID), stringField("benefit_claim_id", item.ID.String()), stringField("status", item.Status))
	}
	return mapBenefitClaim(row), nil
}

func (s *Store) ListBenefitClaims(ctx context.Context, filter domain.BenefitFilter) ([]*domain.BenefitClaim, error) {
	rows, err := s.getQueries(ctx).ListBenefitClaims(ctx, sqlc.ListBenefitClaimsParams{TenantID: filter.TenantID, Limit: limitOrDefault(filter.Limit), Offset: filter.Offset, EmployeeUserID: uuidFromPtr(filter.EmployeeUserID), ClaimTypeID: uuidFromPtr(filter.ClaimTypeID), Status: textFromPtr(filter.Status), PaymentStatus: textFromPtr(filter.PaymentStatus), PayrollExportStatus: textFromPtr(filter.PayrollExportStatus)})
	if err != nil {
		return nil, s.logDBError(ctx, "list benefit claims", err, tenantIDField(filter.TenantID))
	}
	return mapBenefitClaims(rows), nil
}

func (s *Store) GetBenefitClaim(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.BenefitClaim, error) {
	row, err := s.getQueries(ctx).GetBenefitClaim(ctx, sqlc.GetBenefitClaimParams{TenantID: tenantID, ID: id})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrBenefitClaimNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "get benefit claim", err, tenantIDField(tenantID), stringField("benefit_claim_id", id.String()))
	}
	return mapBenefitClaim(row), nil
}

func (s *Store) DeleteBenefitClaim(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteBenefitClaim(ctx, sqlc.SoftDeleteBenefitClaimParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete benefit claim", err, tenantIDField(tenantID), stringField("benefit_claim_id", id.String()))
	}
	return nil
}

func (s *Store) GetBenefitClaimLimitUsage(ctx context.Context, tenantID uuid.UUID, employeeUserID uuid.UUID, claimTypeID uuid.UUID, startDate string, endDate string) (float64, error) {
	value, err := s.getQueries(ctx).GetBenefitClaimLimitUsage(ctx, sqlc.GetBenefitClaimLimitUsageParams{TenantID: tenantID, EmployeeUserID: employeeUserID, ClaimTypeID: claimTypeID, ExpenseDate: dateFromString(startDate), ExpenseDate_2: dateFromString(endDate)})
	if err != nil {
		return 0, s.logDBError(ctx, "get benefit claim limit usage", err, tenantIDField(tenantID), stringField("employee_user_id", employeeUserID.String()))
	}
	return floatFromNumeric(value), nil
}

func (s *Store) CreateBenefitClaimAttachment(ctx context.Context, item *domain.BenefitClaimAttachment, actorID *uuid.UUID) (*domain.BenefitClaimAttachment, error) {
	row, err := s.getQueries(ctx).CreateBenefitClaimAttachment(ctx, sqlc.CreateBenefitClaimAttachmentParams{TenantID: item.TenantID, ClaimID: item.ClaimID, FileName: item.FileName, ContentType: item.ContentType, StoragePath: item.StoragePath, ChecksumSha256: textFromPtr(item.ChecksumSHA256), SizeBytes: item.SizeBytes, UploadedBy: uuidFromPtr(item.UploadedBy), Metadata: jsonBytesFromRaw(item.Metadata), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create benefit claim attachment", err, tenantIDField(item.TenantID), stringField("benefit_claim_id", item.ClaimID.String()))
	}
	return mapBenefitClaimAttachment(row), nil
}

func (s *Store) ListBenefitClaimAttachments(ctx context.Context, tenantID uuid.UUID, claimID uuid.UUID) ([]*domain.BenefitClaimAttachment, error) {
	rows, err := s.getQueries(ctx).ListBenefitClaimAttachments(ctx, sqlc.ListBenefitClaimAttachmentsParams{TenantID: tenantID, ClaimID: claimID})
	if err != nil {
		return nil, s.logDBError(ctx, "list benefit claim attachments", err, tenantIDField(tenantID), stringField("benefit_claim_id", claimID.String()))
	}
	return mapBenefitClaimAttachments(rows), nil
}

func (s *Store) CreateBenefitEvent(ctx context.Context, item *domain.BenefitEvent, actorID *uuid.UUID) (*domain.BenefitEvent, error) {
	row, err := s.getQueries(ctx).CreateBenefitEvent(ctx, sqlc.CreateBenefitEventParams{TenantID: item.TenantID, SourceType: item.SourceType, SourceID: item.SourceID, Action: item.Action, FromStatus: textFromPtr(item.FromStatus), ToStatus: textFromPtr(item.ToStatus), ActorUserID: uuidFromPtr(item.ActorUserID), Remarks: textFromPtr(item.Remarks), Metadata: jsonBytesFromRaw(item.Metadata), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create benefit event", err, tenantIDField(item.TenantID), stringField("source_id", item.SourceID.String()))
	}
	return mapBenefitEvent(row), nil
}

func (s *Store) ListBenefitEvents(ctx context.Context, filter domain.BenefitFilter) ([]*domain.BenefitEvent, error) {
	rows, err := s.getQueries(ctx).ListBenefitEvents(ctx, sqlc.ListBenefitEventsParams{TenantID: filter.TenantID, Limit: limitOrDefault(filter.Limit), Offset: filter.Offset, SourceType: textFromPtr(filter.SourceType), SourceID: uuidFromPtr(filter.SourceID)})
	if err != nil {
		return nil, s.logDBError(ctx, "list benefit events", err, tenantIDField(filter.TenantID))
	}
	return mapBenefitEvents(rows), nil
}

func (s *Store) GetBenefitsSummary(ctx context.Context, tenantID uuid.UUID) ([]*domain.BenefitsSummaryRow, error) {
	rows, err := s.getQueries(ctx).GetBenefitsSummary(ctx, tenantID)
	if err != nil {
		return nil, s.logDBError(ctx, "get benefits summary", err, tenantIDField(tenantID))
	}
	return mapBenefitsSummaryRows(rows), nil
}

func benefitClaimCreateParams(item *domain.BenefitClaim, actorID *uuid.UUID) sqlc.CreateBenefitClaimParams {
	return sqlc.CreateBenefitClaimParams{TenantID: item.TenantID, ClaimNumber: item.ClaimNumber, ClaimTypeID: item.ClaimTypeID, PlanID: uuidFromPtr(item.PlanID), EmployeeUserID: item.EmployeeUserID, DependentID: uuidFromPtr(item.DependentID), ExpenseDate: dateFromTime(item.ExpenseDate), SubmittedAt: timestamptzFromPtr(item.SubmittedAt), ClaimAmount: compNumeric(item.ClaimAmount), ApprovedAmount: numericFromFloatPtr(item.ApprovedAmount), CurrencyCode: item.CurrencyCode, Status: item.Status, PaymentStatus: item.PaymentStatus, PaymentReference: textFromPtr(item.PaymentReference), PaidAt: timestamptzFromPtr(item.PaidAt), ReviewedBy: uuidFromPtr(item.ReviewedBy), ReviewedAt: timestamptzFromPtr(item.ReviewedAt), ReviewRemarks: textFromPtr(item.ReviewRemarks), PayrollExportStatus: item.PayrollExportStatus, PayrollExportedAt: timestamptzFromPtr(item.PayrollExportedAt), PayrollExportReference: textFromPtr(item.PayrollExportReference), Notes: textFromPtr(item.Notes), Metadata: jsonBytesFromRaw(item.Metadata), CreatedBy: uuidFromPtr(actorID)}
}

func benefitClaimUpdateParams(item *domain.BenefitClaim, actorID *uuid.UUID) sqlc.UpdateBenefitClaimParams {
	return sqlc.UpdateBenefitClaimParams{TenantID: item.TenantID, ID: item.ID, ClaimTypeID: item.ClaimTypeID, PlanID: uuidFromPtr(item.PlanID), EmployeeUserID: item.EmployeeUserID, DependentID: uuidFromPtr(item.DependentID), ExpenseDate: dateFromTime(item.ExpenseDate), SubmittedAt: timestamptzFromPtr(item.SubmittedAt), ClaimAmount: compNumeric(item.ClaimAmount), ApprovedAmount: numericFromFloatPtr(item.ApprovedAmount), CurrencyCode: item.CurrencyCode, Status: item.Status, PaymentStatus: item.PaymentStatus, PaymentReference: textFromPtr(item.PaymentReference), PaidAt: timestamptzFromPtr(item.PaidAt), ReviewedBy: uuidFromPtr(item.ReviewedBy), ReviewedAt: timestamptzFromPtr(item.ReviewedAt), ReviewRemarks: textFromPtr(item.ReviewRemarks), PayrollExportStatus: item.PayrollExportStatus, PayrollExportedAt: timestamptzFromPtr(item.PayrollExportedAt), PayrollExportReference: textFromPtr(item.PayrollExportReference), Notes: textFromPtr(item.Notes), Metadata: jsonBytesFromRaw(item.Metadata), UpdatedBy: uuidFromPtr(actorID)}
}
