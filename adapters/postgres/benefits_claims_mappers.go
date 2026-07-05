package postgres

import (
	"encoding/json"

	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func mapBenefitPlan(row sqlc.HrmsBenefitPlan) *domain.BenefitPlan {
	return &domain.BenefitPlan{
		ID: row.ID, TenantID: row.TenantID, Code: row.Code, Name: row.Name, PlanType: row.PlanType,
		Description: ptrFromText(row.Description), ProviderName: ptrFromText(row.ProviderName), PolicyNumber: ptrFromText(row.PolicyNumber),
		CoverageAmount: floatPtrFromNumeric(row.CoverageAmount), EmployerContribution: floatFromNumeric(row.EmployerContribution), EmployeeContribution: floatFromNumeric(row.EmployeeContribution),
		CurrencyCode: row.CurrencyCode, EligibilityRule: benefitRaw(row.EligibilityRule), InsuranceMetadata: benefitRaw(row.InsuranceMetadata),
		EffectiveFrom: ptrFromDate(row.EffectiveFrom), EffectiveTo: ptrFromDate(row.EffectiveTo), IsActive: row.IsActive, Inactive: row.Inactive,
		CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy),
	}
}

func mapBenefitPlans(rows []sqlc.HrmsBenefitPlan) []*domain.BenefitPlan {
	items := make([]*domain.BenefitPlan, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapBenefitPlan(row))
	}
	return items
}

func mapBenefitEnrollmentWindow(row sqlc.HrmsBenefitEnrollmentWindow) *domain.BenefitEnrollmentWindow {
	return &domain.BenefitEnrollmentWindow{
		ID: row.ID, TenantID: row.TenantID, PlanID: row.PlanID, Name: row.Name, OpensOn: timeFromDate(row.OpensOn), ClosesOn: timeFromDate(row.ClosesOn),
		EffectiveFrom: ptrFromDate(row.EffectiveFrom), EffectiveTo: ptrFromDate(row.EffectiveTo), Status: row.Status, Metadata: benefitRaw(row.Metadata), Inactive: row.Inactive,
		CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy),
	}
}

func mapBenefitEnrollmentWindows(rows []sqlc.HrmsBenefitEnrollmentWindow) []*domain.BenefitEnrollmentWindow {
	items := make([]*domain.BenefitEnrollmentWindow, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapBenefitEnrollmentWindow(row))
	}
	return items
}

func mapBenefitDependent(row sqlc.HrmsBenefitDependent) *domain.BenefitDependent {
	return &domain.BenefitDependent{
		ID: row.ID, TenantID: row.TenantID, EmployeeUserID: row.EmployeeUserID, FullName: row.FullName, Relationship: row.Relationship,
		DateOfBirth: ptrFromDate(row.DateOfBirth), Gender: ptrFromText(row.Gender), NomineePercentage: floatPtrFromNumeric(row.NomineePercentage),
		IsNominee: row.IsNominee, Metadata: benefitRaw(row.Metadata), Inactive: row.Inactive,
		CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy),
	}
}

func mapBenefitDependents(rows []sqlc.HrmsBenefitDependent) []*domain.BenefitDependent {
	items := make([]*domain.BenefitDependent, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapBenefitDependent(row))
	}
	return items
}

func mapBenefitEnrollment(row sqlc.HrmsBenefitEnrollment) *domain.BenefitEnrollment {
	return &domain.BenefitEnrollment{
		ID: row.ID, TenantID: row.TenantID, PlanID: row.PlanID, WindowID: ptrFromUUID(row.WindowID), EmployeeUserID: row.EmployeeUserID,
		Status: row.Status, CoverageLevel: ptrFromText(row.CoverageLevel), SelectedAmount: floatPtrFromNumeric(row.SelectedAmount),
		EmployeeContribution: floatFromNumeric(row.EmployeeContribution), EmployerContribution: floatFromNumeric(row.EmployerContribution),
		EffectiveFrom: ptrFromDate(row.EffectiveFrom), EffectiveTo: ptrFromDate(row.EffectiveTo), SubmittedAt: ptrFromTimestamptz(row.SubmittedAt),
		ReviewedBy: ptrFromUUID(row.ReviewedBy), ReviewedAt: ptrFromTimestamptz(row.ReviewedAt), ReviewRemarks: ptrFromText(row.ReviewRemarks),
		Metadata: benefitRaw(row.Metadata), Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy),
		UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy),
	}
}

func mapBenefitEnrollments(rows []sqlc.HrmsBenefitEnrollment) []*domain.BenefitEnrollment {
	items := make([]*domain.BenefitEnrollment, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapBenefitEnrollment(row))
	}
	return items
}

func mapBenefitClaimType(row sqlc.HrmsBenefitClaimType) *domain.BenefitClaimType {
	return &domain.BenefitClaimType{
		ID: row.ID, TenantID: row.TenantID, PlanID: ptrFromUUID(row.PlanID), Code: row.Code, Name: row.Name, Description: ptrFromText(row.Description),
		AnnualLimit: floatPtrFromNumeric(row.AnnualLimit), PerClaimLimit: floatPtrFromNumeric(row.PerClaimLimit), RequiresAttachment: row.RequiresAttachment,
		Taxable: row.Taxable, PayrollComponentCode: ptrFromText(row.PayrollComponentCode), EligibilityRule: benefitRaw(row.EligibilityRule),
		IsActive: row.IsActive, Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy),
		UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy),
	}
}

func mapBenefitClaimTypes(rows []sqlc.HrmsBenefitClaimType) []*domain.BenefitClaimType {
	items := make([]*domain.BenefitClaimType, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapBenefitClaimType(row))
	}
	return items
}

func mapBenefitClaim(row sqlc.HrmsBenefitClaim) *domain.BenefitClaim {
	return &domain.BenefitClaim{
		ID: row.ID, TenantID: row.TenantID, ClaimNumber: row.ClaimNumber, ClaimTypeID: row.ClaimTypeID, PlanID: ptrFromUUID(row.PlanID),
		EmployeeUserID: row.EmployeeUserID, DependentID: ptrFromUUID(row.DependentID), ExpenseDate: timeFromDate(row.ExpenseDate), SubmittedAt: ptrFromTimestamptz(row.SubmittedAt),
		ClaimAmount: floatFromNumeric(row.ClaimAmount), ApprovedAmount: floatPtrFromNumeric(row.ApprovedAmount), CurrencyCode: row.CurrencyCode,
		Status: row.Status, PaymentStatus: row.PaymentStatus, PaymentReference: ptrFromText(row.PaymentReference), PaidAt: ptrFromTimestamptz(row.PaidAt),
		ReviewedBy: ptrFromUUID(row.ReviewedBy), ReviewedAt: ptrFromTimestamptz(row.ReviewedAt), ReviewRemarks: ptrFromText(row.ReviewRemarks),
		PayrollExportStatus: row.PayrollExportStatus, PayrollExportedAt: ptrFromTimestamptz(row.PayrollExportedAt), PayrollExportReference: ptrFromText(row.PayrollExportReference),
		Notes: ptrFromText(row.Notes), Metadata: benefitRaw(row.Metadata), Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt),
		CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy),
	}
}

func mapBenefitClaims(rows []sqlc.HrmsBenefitClaim) []*domain.BenefitClaim {
	items := make([]*domain.BenefitClaim, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapBenefitClaim(row))
	}
	return items
}

func mapBenefitClaimAttachment(row sqlc.HrmsBenefitClaimAttachment) *domain.BenefitClaimAttachment {
	return &domain.BenefitClaimAttachment{
		ID: row.ID, TenantID: row.TenantID, ClaimID: row.ClaimID, FileName: row.FileName, ContentType: row.ContentType,
		StoragePath: row.StoragePath, ChecksumSHA256: ptrFromText(row.ChecksumSha256), SizeBytes: row.SizeBytes, UploadedBy: ptrFromUUID(row.UploadedBy),
		Metadata: benefitRaw(row.Metadata), Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy),
		UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy),
	}
}

func mapBenefitClaimAttachments(rows []sqlc.HrmsBenefitClaimAttachment) []*domain.BenefitClaimAttachment {
	items := make([]*domain.BenefitClaimAttachment, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapBenefitClaimAttachment(row))
	}
	return items
}

func mapBenefitEvent(row sqlc.HrmsBenefitEvent) *domain.BenefitEvent {
	return &domain.BenefitEvent{
		ID: row.ID, TenantID: row.TenantID, SourceType: row.SourceType, SourceID: row.SourceID, Action: row.Action,
		FromStatus: ptrFromText(row.FromStatus), ToStatus: ptrFromText(row.ToStatus), ActorUserID: ptrFromUUID(row.ActorUserID), Remarks: ptrFromText(row.Remarks),
		Metadata: benefitRaw(row.Metadata), Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy),
		UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy),
	}
}

func mapBenefitEvents(rows []sqlc.HrmsBenefitEvent) []*domain.BenefitEvent {
	items := make([]*domain.BenefitEvent, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapBenefitEvent(row))
	}
	return items
}

func mapBenefitsSummaryRows(rows []sqlc.GetBenefitsSummaryRow) []*domain.BenefitsSummaryRow {
	items := make([]*domain.BenefitsSummaryRow, 0, len(rows))
	for _, row := range rows {
		items = append(items, &domain.BenefitsSummaryRow{Metric: row.Metric, MetricCount: row.MetricCount, Amount: floatFromNumeric(row.Amount)})
	}
	return items
}

func benefitRaw(value []byte) json.RawMessage {
	if len(value) == 0 || !json.Valid(value) {
		return json.RawMessage(`{}`)
	}
	return json.RawMessage(value)
}
