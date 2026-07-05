package services

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) CreateBenefitPlan(ctx context.Context, cmd ports.BenefitPlanCommand) (*domain.BenefitPlan, error) {
	item, err := s.benefitPlanFromCommand(cmd)
	if err != nil {
		s.log.Warn().Err(err).Str("operation", "create benefit plan").Str("tenant_id", cmd.TenantID.String()).Msg("invalid benefit plan")
		return nil, err
	}
	item.IsActive = true
	result, err := s.benefitsClaims.CreateBenefitPlan(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create benefit plan", err, serviceTenantIDField(cmd.TenantID), serviceStringField("code", item.Code))
		return nil, err
	}
	s.recordBenefitEvent(ctx, result.TenantID, "plan", result.ID, "created", nil, &result.IsActive, cmd.ActorID, nil)
	return result, nil
}

func (s *TenantService) UpdateBenefitPlan(ctx context.Context, cmd ports.BenefitPlanCommand) (*domain.BenefitPlan, error) {
	item, err := s.benefitPlanFromCommand(cmd)
	if err != nil {
		s.log.Warn().Err(err).Str("operation", "update benefit plan").Str("tenant_id", cmd.TenantID.String()).Msg("invalid benefit plan")
		return nil, err
	}
	result, err := s.benefitsClaims.UpdateBenefitPlan(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("update benefit plan", err, serviceTenantIDField(cmd.TenantID), serviceStringField("benefit_plan_id", cmd.ID.String()))
		return nil, err
	}
	s.recordBenefitEvent(ctx, result.TenantID, "plan", result.ID, "updated", nil, &result.IsActive, cmd.ActorID, nil)
	return result, nil
}

func (s *TenantService) ListBenefitPlans(ctx context.Context, filter domain.BenefitFilter) ([]*domain.BenefitPlan, error) {
	items, err := s.benefitsClaims.ListBenefitPlans(ctx, filter)
	if err != nil {
		s.logError("list benefit plans", err, serviceTenantIDField(filter.TenantID))
	}
	return items, err
}

func (s *TenantService) DeleteBenefitPlan(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.benefitsClaims.DeleteBenefitPlan(ctx, tenantID, id, actorID); err != nil {
		s.logError("delete benefit plan", err, serviceTenantIDField(tenantID), serviceStringField("benefit_plan_id", id.String()))
		return err
	}
	return nil
}

func (s *TenantService) CreateBenefitEnrollmentWindow(ctx context.Context, cmd ports.BenefitEnrollmentWindowCommand) (*domain.BenefitEnrollmentWindow, error) {
	item, err := s.benefitWindowFromCommand(ctx, cmd)
	if err != nil {
		s.log.Warn().Err(err).Str("operation", "create benefit enrollment window").Str("tenant_id", cmd.TenantID.String()).Msg("invalid benefit window")
		return nil, err
	}
	result, err := s.benefitsClaims.CreateBenefitEnrollmentWindow(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create benefit enrollment window", err, serviceTenantIDField(cmd.TenantID), serviceStringField("plan_id", cmd.PlanID.String()))
		return nil, err
	}
	s.recordBenefitStatusEvent(ctx, result.TenantID, "window", result.ID, "created", nil, &result.Status, cmd.ActorID, nil)
	return result, nil
}

func (s *TenantService) UpdateBenefitEnrollmentWindow(ctx context.Context, cmd ports.BenefitEnrollmentWindowCommand) (*domain.BenefitEnrollmentWindow, error) {
	item, err := s.benefitWindowFromCommand(ctx, cmd)
	if err != nil {
		s.log.Warn().Err(err).Str("operation", "update benefit enrollment window").Str("tenant_id", cmd.TenantID.String()).Msg("invalid benefit window")
		return nil, err
	}
	result, err := s.benefitsClaims.UpdateBenefitEnrollmentWindow(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("update benefit enrollment window", err, serviceTenantIDField(cmd.TenantID), serviceStringField("benefit_window_id", cmd.ID.String()))
		return nil, err
	}
	s.recordBenefitStatusEvent(ctx, result.TenantID, "window", result.ID, "updated", nil, &result.Status, cmd.ActorID, nil)
	return result, nil
}

func (s *TenantService) ListBenefitEnrollmentWindows(ctx context.Context, filter domain.BenefitFilter) ([]*domain.BenefitEnrollmentWindow, error) {
	items, err := s.benefitsClaims.ListBenefitEnrollmentWindows(ctx, filter)
	if err != nil {
		s.logError("list benefit enrollment windows", err, serviceTenantIDField(filter.TenantID))
	}
	return items, err
}

func (s *TenantService) DeleteBenefitEnrollmentWindow(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	return s.benefitsClaims.DeleteBenefitEnrollmentWindow(ctx, tenantID, id, actorID)
}

func (s *TenantService) CreateBenefitDependent(ctx context.Context, cmd ports.BenefitDependentCommand) (*domain.BenefitDependent, error) {
	item, err := s.benefitDependentFromCommand(ctx, cmd)
	if err != nil {
		s.log.Warn().Err(err).Str("operation", "create benefit dependent").Str("tenant_id", cmd.TenantID.String()).Msg("invalid benefit dependent")
		return nil, err
	}
	result, err := s.benefitsClaims.CreateBenefitDependent(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create benefit dependent", err, serviceTenantIDField(cmd.TenantID), serviceStringField("employee_user_id", cmd.EmployeeUserID.String()))
		return nil, err
	}
	s.recordBenefitEvent(ctx, result.TenantID, "dependent", result.ID, "created", nil, nil, cmd.ActorID, nil)
	return result, nil
}

func (s *TenantService) UpdateBenefitDependent(ctx context.Context, cmd ports.BenefitDependentCommand) (*domain.BenefitDependent, error) {
	item, err := s.benefitDependentFromCommand(ctx, cmd)
	if err != nil {
		s.log.Warn().Err(err).Str("operation", "update benefit dependent").Str("tenant_id", cmd.TenantID.String()).Msg("invalid benefit dependent")
		return nil, err
	}
	result, err := s.benefitsClaims.UpdateBenefitDependent(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("update benefit dependent", err, serviceTenantIDField(cmd.TenantID), serviceStringField("benefit_dependent_id", cmd.ID.String()))
		return nil, err
	}
	s.recordBenefitEvent(ctx, result.TenantID, "dependent", result.ID, "updated", nil, nil, cmd.ActorID, nil)
	return result, nil
}

func (s *TenantService) ListBenefitDependents(ctx context.Context, filter domain.BenefitFilter) ([]*domain.BenefitDependent, error) {
	items, err := s.benefitsClaims.ListBenefitDependents(ctx, filter)
	if err != nil {
		s.logError("list benefit dependents", err, serviceTenantIDField(filter.TenantID))
	}
	return items, err
}

func (s *TenantService) DeleteBenefitDependent(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	return s.benefitsClaims.DeleteBenefitDependent(ctx, tenantID, id, actorID)
}

func (s *TenantService) CreateBenefitEnrollment(ctx context.Context, cmd ports.BenefitEnrollmentCommand) (*domain.BenefitEnrollment, error) {
	item, err := s.benefitEnrollmentFromCommand(ctx, cmd)
	if err != nil {
		s.log.Warn().Err(err).Str("operation", "create benefit enrollment").Str("tenant_id", cmd.TenantID.String()).Msg("invalid benefit enrollment")
		return nil, err
	}
	result, err := s.benefitsClaims.CreateBenefitEnrollment(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create benefit enrollment", err, serviceTenantIDField(cmd.TenantID), serviceStringField("employee_user_id", cmd.EmployeeUserID.String()))
		return nil, err
	}
	s.recordBenefitStatusEvent(ctx, result.TenantID, "enrollment", result.ID, "created", nil, &result.Status, cmd.ActorID, nil)
	return result, nil
}

func (s *TenantService) UpdateBenefitEnrollment(ctx context.Context, cmd ports.BenefitEnrollmentCommand) (*domain.BenefitEnrollment, error) {
	item, err := s.benefitEnrollmentFromCommand(ctx, cmd)
	if err != nil {
		s.log.Warn().Err(err).Str("operation", "update benefit enrollment").Str("tenant_id", cmd.TenantID.String()).Msg("invalid benefit enrollment")
		return nil, err
	}
	result, err := s.benefitsClaims.UpdateBenefitEnrollment(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("update benefit enrollment", err, serviceTenantIDField(cmd.TenantID), serviceStringField("benefit_enrollment_id", cmd.ID.String()))
		return nil, err
	}
	s.recordBenefitStatusEvent(ctx, result.TenantID, "enrollment", result.ID, "updated", nil, &result.Status, cmd.ActorID, nil)
	return result, nil
}

func (s *TenantService) UpdateBenefitEnrollmentStatus(ctx context.Context, cmd ports.BenefitEnrollmentStatusCommand) (*domain.BenefitEnrollment, error) {
	before, _ := s.benefitsClaims.GetBenefitEnrollment(ctx, cmd.TenantID, cmd.ID)
	status := normalizeBenefitServiceStatus(cmd.Status, domain.BenefitEnrollmentStatusSubmitted)
	result, err := s.benefitsClaims.UpdateBenefitEnrollmentStatus(ctx, cmd.TenantID, cmd.ID, status, cmd.ActorID, cmd.Remarks)
	if err != nil {
		s.logError("update benefit enrollment status", err, serviceTenantIDField(cmd.TenantID), serviceStringField("benefit_enrollment_id", cmd.ID.String()))
		return nil, err
	}
	var from *string
	if before != nil {
		from = &before.Status
	}
	s.recordBenefitStatusEvent(ctx, result.TenantID, "enrollment", result.ID, "status_updated", from, &result.Status, cmd.ActorID, cmd.Remarks)
	return result, nil
}

func (s *TenantService) ListBenefitEnrollments(ctx context.Context, filter domain.BenefitFilter) ([]*domain.BenefitEnrollment, error) {
	items, err := s.benefitsClaims.ListBenefitEnrollments(ctx, filter)
	if err != nil {
		s.logError("list benefit enrollments", err, serviceTenantIDField(filter.TenantID))
	}
	return items, err
}

func (s *TenantService) DeleteBenefitEnrollment(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	return s.benefitsClaims.DeleteBenefitEnrollment(ctx, tenantID, id, actorID)
}

func (s *TenantService) CreateBenefitClaimType(ctx context.Context, cmd ports.BenefitClaimTypeCommand) (*domain.BenefitClaimType, error) {
	item, err := s.benefitClaimTypeFromCommand(ctx, cmd)
	if err != nil {
		s.log.Warn().Err(err).Str("operation", "create benefit claim type").Str("tenant_id", cmd.TenantID.String()).Msg("invalid benefit claim type")
		return nil, err
	}
	item.IsActive = true
	result, err := s.benefitsClaims.CreateBenefitClaimType(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create benefit claim type", err, serviceTenantIDField(cmd.TenantID), serviceStringField("code", item.Code))
		return nil, err
	}
	s.recordBenefitEvent(ctx, result.TenantID, "claim_type", result.ID, "created", nil, &result.IsActive, cmd.ActorID, nil)
	return result, nil
}

func (s *TenantService) UpdateBenefitClaimType(ctx context.Context, cmd ports.BenefitClaimTypeCommand) (*domain.BenefitClaimType, error) {
	item, err := s.benefitClaimTypeFromCommand(ctx, cmd)
	if err != nil {
		s.log.Warn().Err(err).Str("operation", "update benefit claim type").Str("tenant_id", cmd.TenantID.String()).Msg("invalid benefit claim type")
		return nil, err
	}
	result, err := s.benefitsClaims.UpdateBenefitClaimType(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("update benefit claim type", err, serviceTenantIDField(cmd.TenantID), serviceStringField("benefit_claim_type_id", cmd.ID.String()))
		return nil, err
	}
	s.recordBenefitEvent(ctx, result.TenantID, "claim_type", result.ID, "updated", nil, &result.IsActive, cmd.ActorID, nil)
	return result, nil
}

func (s *TenantService) ListBenefitClaimTypes(ctx context.Context, filter domain.BenefitFilter) ([]*domain.BenefitClaimType, error) {
	items, err := s.benefitsClaims.ListBenefitClaimTypes(ctx, filter)
	if err != nil {
		s.logError("list benefit claim types", err, serviceTenantIDField(filter.TenantID))
	}
	return items, err
}

func (s *TenantService) DeleteBenefitClaimType(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	return s.benefitsClaims.DeleteBenefitClaimType(ctx, tenantID, id, actorID)
}

func (s *TenantService) CreateBenefitClaim(ctx context.Context, cmd ports.BenefitClaimCommand) (*domain.BenefitClaim, error) {
	item, err := s.benefitClaimFromCommand(ctx, cmd)
	if err != nil {
		s.log.Warn().Err(err).Str("operation", "create benefit claim").Str("tenant_id", cmd.TenantID.String()).Msg("invalid benefit claim")
		return nil, err
	}
	if item.ClaimNumber == "" {
		item.ClaimNumber = fmt.Sprintf("CLM-%s-%s", time.Now().UTC().Format("20060102"), strings.ToUpper(uuid.NewString()[:8]))
	}
	if item.Status == domain.BenefitClaimStatusSubmitted && item.SubmittedAt == nil {
		now := time.Now().UTC()
		item.SubmittedAt = &now
	}
	claimType, err := s.benefitsClaims.GetBenefitClaimType(ctx, item.TenantID, item.ClaimTypeID)
	if err != nil {
		return nil, err
	}
	if err := s.validateBenefitClaimLimits(ctx, item, claimType); err != nil {
		return nil, err
	}
	result, err := s.benefitsClaims.CreateBenefitClaim(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create benefit claim", err, serviceTenantIDField(cmd.TenantID), serviceStringField("employee_user_id", cmd.EmployeeUserID.String()))
		return nil, err
	}
	s.recordBenefitStatusEvent(ctx, result.TenantID, "claim", result.ID, "created", nil, &result.Status, cmd.ActorID, nil)
	return result, nil
}

func (s *TenantService) UpdateBenefitClaim(ctx context.Context, cmd ports.BenefitClaimCommand) (*domain.BenefitClaim, error) {
	item, err := s.benefitClaimFromCommand(ctx, cmd)
	if err != nil {
		s.log.Warn().Err(err).Str("operation", "update benefit claim").Str("tenant_id", cmd.TenantID.String()).Msg("invalid benefit claim")
		return nil, err
	}
	existing, err := s.benefitsClaims.GetBenefitClaim(ctx, cmd.TenantID, cmd.ID)
	if err != nil {
		return nil, err
	}
	item.ClaimNumber = existing.ClaimNumber
	item.ReviewedBy = existing.ReviewedBy
	item.ReviewedAt = existing.ReviewedAt
	claimType, err := s.benefitsClaims.GetBenefitClaimType(ctx, item.TenantID, item.ClaimTypeID)
	if err != nil {
		return nil, err
	}
	if err := s.validateBenefitClaimLimits(ctx, item, claimType); err != nil {
		return nil, err
	}
	result, err := s.benefitsClaims.UpdateBenefitClaim(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("update benefit claim", err, serviceTenantIDField(cmd.TenantID), serviceStringField("benefit_claim_id", cmd.ID.String()))
		return nil, err
	}
	s.recordBenefitStatusEvent(ctx, result.TenantID, "claim", result.ID, "updated", &existing.Status, &result.Status, cmd.ActorID, nil)
	return result, nil
}

func (s *TenantService) UpdateBenefitClaimStatus(ctx context.Context, cmd ports.BenefitClaimStatusCommand) (*domain.BenefitClaim, error) {
	claim, err := s.benefitsClaims.GetBenefitClaim(ctx, cmd.TenantID, cmd.ID)
	if err != nil {
		return nil, err
	}
	from := claim.Status
	now := time.Now().UTC()
	claim.Status = normalizeBenefitServiceStatus(cmd.Status, domain.BenefitClaimStatusUnderReview)
	claim.ReviewRemarks = cmd.Remarks
	switch claim.Status {
	case domain.BenefitClaimStatusApproved:
		if cmd.ApprovedAmount != nil {
			claim.ApprovedAmount = cmd.ApprovedAmount
		} else if claim.ApprovedAmount == nil {
			amount := claim.ClaimAmount
			claim.ApprovedAmount = &amount
		}
		claim.PaymentStatus = domain.BenefitPaymentStatusPending
		claim.PayrollExportStatus = domain.BenefitPayrollExportStatusReady
	case domain.BenefitClaimStatusRejected, domain.BenefitClaimStatusCancelled:
		claim.PaymentStatus = domain.BenefitPaymentStatusNotPayable
		claim.PayrollExportStatus = domain.BenefitPayrollExportStatusBlocked
	case domain.BenefitClaimStatusPaid:
		if cmd.ApprovedAmount != nil {
			claim.ApprovedAmount = cmd.ApprovedAmount
		} else if claim.ApprovedAmount == nil {
			amount := claim.ClaimAmount
			claim.ApprovedAmount = &amount
		}
		claim.PaymentStatus = domain.BenefitPaymentStatusPaid
		claim.PaidAt = &now
		claim.PaymentReference = cmd.PaymentReference
		claim.PayrollExportStatus = domain.BenefitPayrollExportStatusExported
		claim.PayrollExportedAt = &now
		claim.PayrollExportReference = cmd.PayrollExportReference
	case domain.BenefitClaimStatusSubmitted:
		claim.SubmittedAt = &now
	}
	result, err := s.benefitsClaims.UpdateBenefitClaimStatus(ctx, claim, cmd.ActorID)
	if err != nil {
		s.logError("update benefit claim status", err, serviceTenantIDField(cmd.TenantID), serviceStringField("benefit_claim_id", cmd.ID.String()))
		return nil, err
	}
	s.recordBenefitStatusEvent(ctx, result.TenantID, "claim", result.ID, "status_updated", &from, &result.Status, cmd.ActorID, cmd.Remarks)
	return result, nil
}

func (s *TenantService) ListBenefitClaims(ctx context.Context, filter domain.BenefitFilter) ([]*domain.BenefitClaim, error) {
	items, err := s.benefitsClaims.ListBenefitClaims(ctx, filter)
	if err != nil {
		s.logError("list benefit claims", err, serviceTenantIDField(filter.TenantID))
	}
	return items, err
}

func (s *TenantService) DeleteBenefitClaim(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	return s.benefitsClaims.DeleteBenefitClaim(ctx, tenantID, id, actorID)
}

func (s *TenantService) CreateBenefitClaimAttachment(ctx context.Context, cmd ports.BenefitClaimAttachmentCommand) (*domain.BenefitClaimAttachment, error) {
	if s.objectStorage == nil {
		return nil, domain.ErrStorageProviderSettingsNotFound
	}
	claim, err := s.benefitsClaims.GetBenefitClaim(ctx, cmd.TenantID, cmd.ClaimID)
	if err != nil {
		return nil, err
	}
	content, err := base64.StdEncoding.DecodeString(strings.TrimSpace(cmd.ContentBase64))
	if err != nil || len(content) == 0 {
		return nil, domain.ErrInvalidBenefitClaimAttachment
	}
	settings, err := s.resolveWorkflowStorageSettings(ctx, cmd.TenantID)
	if err != nil {
		return nil, err
	}
	entityID := uuid.New()
	contentType := strings.TrimSpace(cmd.ContentType)
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	storagePath, err := s.objectStorage.PutObject(ctx, settings, ports.StoreObjectInput{TenantID: cmd.TenantID, Category: ports.StorageCategoryBenefitClaim, OwnerID: cmd.ClaimID, EntityID: entityID, FileName: cmd.FileName, ContentType: contentType, Content: content})
	if err != nil {
		s.logError("store benefit claim attachment", err, serviceTenantIDField(cmd.TenantID), serviceStringField("benefit_claim_id", cmd.ClaimID.String()))
		return nil, err
	}
	sum := sha256.Sum256(content)
	checksum := hex.EncodeToString(sum[:])
	item, err := domain.NewBenefitClaimAttachment(domain.BenefitClaimAttachment{TenantID: cmd.TenantID, ClaimID: claim.ID, FileName: cmd.FileName, ContentType: contentType, StoragePath: storagePath, ChecksumSHA256: &checksum, SizeBytes: int64(len(content)), UploadedBy: cmd.ActorID, Metadata: cmd.Metadata})
	if err != nil {
		return nil, err
	}
	result, err := s.benefitsClaims.CreateBenefitClaimAttachment(ctx, item, cmd.ActorID)
	if err != nil {
		return nil, err
	}
	s.recordBenefitEvent(ctx, result.TenantID, "attachment", result.ID, "created", nil, nil, cmd.ActorID, nil)
	return result, nil
}

func (s *TenantService) ListBenefitClaimAttachments(ctx context.Context, tenantID uuid.UUID, claimID uuid.UUID) ([]*domain.BenefitClaimAttachment, error) {
	return s.benefitsClaims.ListBenefitClaimAttachments(ctx, tenantID, claimID)
}

func (s *TenantService) ListBenefitEvents(ctx context.Context, filter domain.BenefitFilter) ([]*domain.BenefitEvent, error) {
	return s.benefitsClaims.ListBenefitEvents(ctx, filter)
}

func (s *TenantService) GetBenefitsSummary(ctx context.Context, tenantID uuid.UUID) ([]*domain.BenefitsSummaryRow, error) {
	return s.benefitsClaims.GetBenefitsSummary(ctx, tenantID)
}

func (s *TenantService) benefitPlanFromCommand(cmd ports.BenefitPlanCommand) (*domain.BenefitPlan, error) {
	from, err := parseOptionalDate(cmd.EffectiveFrom)
	if err != nil {
		return nil, err
	}
	to, err := parseOptionalDate(cmd.EffectiveTo)
	if err != nil {
		return nil, err
	}
	return domain.NewBenefitPlan(domain.BenefitPlan{ID: cmd.ID, TenantID: cmd.TenantID, Code: cmd.Code, Name: cmd.Name, PlanType: cmd.PlanType, Description: cmd.Description, ProviderName: cmd.ProviderName, PolicyNumber: cmd.PolicyNumber, CoverageAmount: cmd.CoverageAmount, EmployerContribution: cmd.EmployerContribution, EmployeeContribution: cmd.EmployeeContribution, CurrencyCode: cmd.CurrencyCode, EligibilityRule: cmd.EligibilityRule, InsuranceMetadata: cmd.InsuranceMetadata, EffectiveFrom: from, EffectiveTo: to, IsActive: cmd.IsActive})
}

func (s *TenantService) benefitWindowFromCommand(ctx context.Context, cmd ports.BenefitEnrollmentWindowCommand) (*domain.BenefitEnrollmentWindow, error) {
	if _, err := s.benefitsClaims.GetBenefitPlan(ctx, cmd.TenantID, cmd.PlanID); err != nil {
		return nil, err
	}
	opens, err := parseBenefitRequiredDate(cmd.OpensOn)
	if err != nil {
		return nil, err
	}
	closes, err := parseBenefitRequiredDate(cmd.ClosesOn)
	if err != nil {
		return nil, err
	}
	from, err := parseOptionalDate(cmd.EffectiveFrom)
	if err != nil {
		return nil, err
	}
	to, err := parseOptionalDate(cmd.EffectiveTo)
	if err != nil {
		return nil, err
	}
	return domain.NewBenefitEnrollmentWindow(domain.BenefitEnrollmentWindow{ID: cmd.ID, TenantID: cmd.TenantID, PlanID: cmd.PlanID, Name: cmd.Name, OpensOn: *opens, ClosesOn: *closes, EffectiveFrom: from, EffectiveTo: to, Status: cmd.Status, Metadata: cmd.Metadata})
}

func (s *TenantService) benefitDependentFromCommand(ctx context.Context, cmd ports.BenefitDependentCommand) (*domain.BenefitDependent, error) {
	if _, err := s.employees.GetEmployeeByUserID(ctx, cmd.TenantID, cmd.EmployeeUserID); err != nil {
		return nil, err
	}
	dob, err := parseOptionalDate(cmd.DateOfBirth)
	if err != nil {
		return nil, err
	}
	return domain.NewBenefitDependent(domain.BenefitDependent{ID: cmd.ID, TenantID: cmd.TenantID, EmployeeUserID: cmd.EmployeeUserID, FullName: cmd.FullName, Relationship: cmd.Relationship, DateOfBirth: dob, Gender: cmd.Gender, NomineePercentage: cmd.NomineePercentage, IsNominee: cmd.IsNominee, Metadata: cmd.Metadata})
}

func (s *TenantService) benefitEnrollmentFromCommand(ctx context.Context, cmd ports.BenefitEnrollmentCommand) (*domain.BenefitEnrollment, error) {
	if _, err := s.benefitsClaims.GetBenefitPlan(ctx, cmd.TenantID, cmd.PlanID); err != nil {
		return nil, err
	}
	if cmd.WindowID != nil {
		if _, err := s.benefitsClaims.GetBenefitEnrollmentWindow(ctx, cmd.TenantID, *cmd.WindowID); err != nil {
			return nil, err
		}
	}
	if _, err := s.employees.GetEmployeeByUserID(ctx, cmd.TenantID, cmd.EmployeeUserID); err != nil {
		return nil, err
	}
	from, err := parseOptionalDate(cmd.EffectiveFrom)
	if err != nil {
		return nil, err
	}
	to, err := parseOptionalDate(cmd.EffectiveTo)
	if err != nil {
		return nil, err
	}
	submittedAt, err := parseBenefitOptionalDateTime(cmd.SubmittedAt)
	if err != nil {
		return nil, err
	}
	return domain.NewBenefitEnrollment(domain.BenefitEnrollment{ID: cmd.ID, TenantID: cmd.TenantID, PlanID: cmd.PlanID, WindowID: cmd.WindowID, EmployeeUserID: cmd.EmployeeUserID, Status: cmd.Status, CoverageLevel: cmd.CoverageLevel, SelectedAmount: cmd.SelectedAmount, EmployeeContribution: cmd.EmployeeContribution, EmployerContribution: cmd.EmployerContribution, EffectiveFrom: from, EffectiveTo: to, SubmittedAt: submittedAt, ReviewRemarks: cmd.ReviewRemarks, Metadata: cmd.Metadata})
}

func (s *TenantService) benefitClaimTypeFromCommand(ctx context.Context, cmd ports.BenefitClaimTypeCommand) (*domain.BenefitClaimType, error) {
	if cmd.PlanID != nil {
		if _, err := s.benefitsClaims.GetBenefitPlan(ctx, cmd.TenantID, *cmd.PlanID); err != nil {
			return nil, err
		}
	}
	return domain.NewBenefitClaimType(domain.BenefitClaimType{ID: cmd.ID, TenantID: cmd.TenantID, PlanID: cmd.PlanID, Code: cmd.Code, Name: cmd.Name, Description: cmd.Description, AnnualLimit: cmd.AnnualLimit, PerClaimLimit: cmd.PerClaimLimit, RequiresAttachment: cmd.RequiresAttachment, Taxable: cmd.Taxable, PayrollComponentCode: cmd.PayrollComponentCode, EligibilityRule: cmd.EligibilityRule, IsActive: cmd.IsActive})
}

func (s *TenantService) benefitClaimFromCommand(ctx context.Context, cmd ports.BenefitClaimCommand) (*domain.BenefitClaim, error) {
	if _, err := s.employees.GetEmployeeByUserID(ctx, cmd.TenantID, cmd.EmployeeUserID); err != nil {
		return nil, err
	}
	expenseDate, err := parseBenefitRequiredDate(cmd.ExpenseDate)
	if err != nil {
		return nil, err
	}
	submittedAt, err := parseBenefitOptionalDateTime(cmd.SubmittedAt)
	if err != nil {
		return nil, err
	}
	paidAt, err := parseBenefitOptionalDateTime(cmd.PaidAt)
	if err != nil {
		return nil, err
	}
	payrollExportedAt, err := parseBenefitOptionalDateTime(cmd.PayrollExportedAt)
	if err != nil {
		return nil, err
	}
	return domain.NewBenefitClaim(domain.BenefitClaim{ID: cmd.ID, TenantID: cmd.TenantID, ClaimNumber: cmd.ClaimNumber, ClaimTypeID: cmd.ClaimTypeID, PlanID: cmd.PlanID, EmployeeUserID: cmd.EmployeeUserID, DependentID: cmd.DependentID, ExpenseDate: *expenseDate, SubmittedAt: submittedAt, ClaimAmount: cmd.ClaimAmount, ApprovedAmount: cmd.ApprovedAmount, CurrencyCode: cmd.CurrencyCode, Status: cmd.Status, PaymentStatus: cmd.PaymentStatus, PaymentReference: cmd.PaymentReference, PaidAt: paidAt, ReviewRemarks: cmd.ReviewRemarks, PayrollExportStatus: cmd.PayrollExportStatus, PayrollExportedAt: payrollExportedAt, PayrollExportReference: cmd.PayrollExportReference, Notes: cmd.Notes, Metadata: cmd.Metadata})
}

func (s *TenantService) validateBenefitClaimLimits(ctx context.Context, claim *domain.BenefitClaim, claimType *domain.BenefitClaimType) error {
	if claimType.PerClaimLimit != nil && claim.ClaimAmount > *claimType.PerClaimLimit {
		return domain.ErrInvalidBenefitClaim
	}
	if claimType.AnnualLimit == nil {
		return nil
	}
	start := time.Date(claim.ExpenseDate.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(claim.ExpenseDate.Year(), 12, 31, 0, 0, 0, 0, time.UTC)
	used, err := s.benefitsClaims.GetBenefitClaimLimitUsage(ctx, claim.TenantID, claim.EmployeeUserID, claim.ClaimTypeID, start.Format("2006-01-02"), end.Format("2006-01-02"))
	if err != nil {
		return err
	}
	if used+claim.ClaimAmount > *claimType.AnnualLimit {
		return domain.ErrInvalidBenefitClaim
	}
	return nil
}

func parseBenefitRequiredDate(value string) (*time.Time, error) {
	parsed, err := parseOptionalDate(value)
	if err != nil || parsed == nil {
		return nil, domain.ErrInvalidBenefitClaim
	}
	return parsed, nil
}

func parseBenefitOptionalDateTime(value string) (*time.Time, error) {
	clean := strings.TrimSpace(value)
	if clean == "" {
		return nil, nil
	}
	if parsed, err := time.Parse(time.RFC3339, clean); err == nil {
		return &parsed, nil
	}
	parsed, err := time.Parse("2006-01-02", clean)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

func normalizeBenefitServiceStatus(value string, fallback string) string {
	clean := strings.ToLower(strings.TrimSpace(value))
	if clean == "" {
		return fallback
	}
	return clean
}

func (s *TenantService) recordBenefitStatusEvent(ctx context.Context, tenantID uuid.UUID, sourceType string, sourceID uuid.UUID, action string, fromStatus *string, toStatus *string, actorID *uuid.UUID, remarks *string) {
	_, err := s.benefitsClaims.CreateBenefitEvent(ctx, &domain.BenefitEvent{TenantID: tenantID, SourceType: sourceType, SourceID: sourceID, Action: action, FromStatus: fromStatus, ToStatus: toStatus, ActorUserID: actorID, Remarks: remarks}, actorID)
	if err != nil {
		s.log.Warn().Err(err).Str("operation", "record benefit event").Str("tenant_id", tenantID.String()).Str("source_id", sourceID.String()).Msg("benefit event skipped")
	}
}

func (s *TenantService) recordBenefitEvent(ctx context.Context, tenantID uuid.UUID, sourceType string, sourceID uuid.UUID, action string, fromActive *bool, toActive *bool, actorID *uuid.UUID, remarks *string) {
	var fromStatus, toStatus *string
	if fromActive != nil {
		value := fmt.Sprintf("active:%t", *fromActive)
		fromStatus = &value
	}
	if toActive != nil {
		value := fmt.Sprintf("active:%t", *toActive)
		toStatus = &value
	}
	s.recordBenefitStatusEvent(ctx, tenantID, sourceType, sourceID, action, fromStatus, toStatus, actorID, remarks)
}
