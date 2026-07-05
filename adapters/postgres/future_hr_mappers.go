package postgres

import (
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func mapPrivacyConsent(row sqlc.HrmsPrivacyConsent) *domain.PrivacyConsent {
	return &domain.PrivacyConsent{
		ID:              row.ID,
		TenantID:        row.TenantID,
		EmployeeUserID:  ptrFromUUID(row.EmployeeUserID),
		WorkerProfileID: ptrFromUUID(row.WorkerProfileID),
		ConsentKey:      row.ConsentKey,
		ConsentArea:     row.ConsentArea,
		Status:          row.Status,
		LawfulBasis:     row.LawfulBasis,
		Channel:         row.Channel,
		Source:          row.Source,
		Purpose:         row.Purpose,
		GrantedAt:       ptrFromTimestamptz(row.GrantedAt),
		RevokedAt:       ptrFromTimestamptz(row.RevokedAt),
		ExpiresAt:       ptrFromTimestamptz(row.ExpiresAt),
		Evidence:        jsonRawDefault(row.Evidence, `{}`),
		Metadata:        jsonRawDefault(row.Metadata, `{}`),
		Inactive:        row.Inactive,
		CreatedAt:       timeFromTimestamptz(row.CreatedAt),
		CreatedBy:       ptrFromUUID(row.CreatedBy),
		UpdatedAt:       timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:       ptrFromUUID(row.UpdatedBy),
	}
}

func mapPrivacyConsents(rows []sqlc.HrmsPrivacyConsent) []*domain.PrivacyConsent {
	items := make([]*domain.PrivacyConsent, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapPrivacyConsent(row))
	}
	return items
}

func mapDataErasureRequest(row sqlc.HrmsDataErasureRequest) *domain.DataErasureRequest {
	return &domain.DataErasureRequest{
		ID:              row.ID,
		TenantID:        row.TenantID,
		RequestKey:      row.RequestKey,
		SubjectUserID:   ptrFromUUID(row.SubjectUserID),
		WorkerProfileID: ptrFromUUID(row.WorkerProfileID),
		RequestType:     row.RequestType,
		Status:          row.Status,
		Priority:        row.Priority,
		RequestedBy:     ptrFromUUID(row.RequestedBy),
		Reason:          row.Reason,
		Scope:           jsonRawDefault(row.Scope, `{}`),
		RetainedReason:  ptrFromText(row.RetainedReason),
		DueAt:           ptrFromTimestamptz(row.DueAt),
		CompletedAt:     ptrFromTimestamptz(row.CompletedAt),
		ReviewedBy:      ptrFromUUID(row.ReviewedBy),
		ReviewedAt:      ptrFromTimestamptz(row.ReviewedAt),
		AuditSummary:    jsonRawDefault(row.AuditSummary, `{}`),
		Inactive:        row.Inactive,
		CreatedAt:       timeFromTimestamptz(row.CreatedAt),
		CreatedBy:       ptrFromUUID(row.CreatedBy),
		UpdatedAt:       timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:       ptrFromUUID(row.UpdatedBy),
	}
}

func mapDataErasureRequests(rows []sqlc.HrmsDataErasureRequest) []*domain.DataErasureRequest {
	items := make([]*domain.DataErasureRequest, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapDataErasureRequest(row))
	}
	return items
}

func mapEcosystemIntegrationHook(row sqlc.HrmsEcosystemIntegrationHook) *domain.EcosystemIntegrationHook {
	return &domain.EcosystemIntegrationHook{
		ID:              row.ID,
		TenantID:        row.TenantID,
		HookKey:         row.HookKey,
		Provider:        row.Provider,
		Channel:         row.Channel,
		Direction:       row.Direction,
		Status:          row.Status,
		DisplayName:     row.DisplayName,
		EndpointURL:     ptrFromText(row.EndpointUrl),
		EventTypes:      row.EventTypes,
		SecretRef:       ptrFromText(row.SecretRef),
		ConsentRequired: row.ConsentRequired,
		MobileSafe:      row.MobileSafe,
		LastCheckedAt:   ptrFromTimestamptz(row.LastCheckedAt),
		LastError:       ptrFromText(row.LastError),
		Config:          jsonRawDefault(row.Config, `{}`),
		Inactive:        row.Inactive,
		CreatedAt:       timeFromTimestamptz(row.CreatedAt),
		CreatedBy:       ptrFromUUID(row.CreatedBy),
		UpdatedAt:       timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:       ptrFromUUID(row.UpdatedBy),
	}
}

func mapEcosystemIntegrationHooks(rows []sqlc.HrmsEcosystemIntegrationHook) []*domain.EcosystemIntegrationHook {
	items := make([]*domain.EcosystemIntegrationHook, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapEcosystemIntegrationHook(row))
	}
	return items
}

func mapMobileAPIConstraint(row sqlc.HrmsMobileApiConstraint) *domain.MobileAPIConstraint {
	return &domain.MobileAPIConstraint{
		ID:                    row.ID,
		TenantID:              row.TenantID,
		ConstraintKey:         row.ConstraintKey,
		Workflow:              row.Workflow,
		MinAndroidVersion:     ptrFromText(row.MinAndroidVersion),
		MinIOSVersion:         ptrFromText(row.MinIosVersion),
		OfflineSupported:      row.OfflineSupported,
		LowBandwidthMode:      row.LowBandwidthMode,
		RequiresLocation:      row.RequiresLocation,
		RequiresDeviceBinding: row.RequiresDeviceBinding,
		MaxPayloadKB:          row.MaxPayloadKb,
		Status:                row.Status,
		Notes:                 ptrFromText(row.Notes),
		Config:                jsonRawDefault(row.Config, `{}`),
		Inactive:              row.Inactive,
		CreatedAt:             timeFromTimestamptz(row.CreatedAt),
		CreatedBy:             ptrFromUUID(row.CreatedBy),
		UpdatedAt:             timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:             ptrFromUUID(row.UpdatedBy),
	}
}

func mapMobileAPIConstraints(rows []sqlc.HrmsMobileApiConstraint) []*domain.MobileAPIConstraint {
	items := make([]*domain.MobileAPIConstraint, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapMobileAPIConstraint(row))
	}
	return items
}
