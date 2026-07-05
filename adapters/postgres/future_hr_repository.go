package postgres

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (s *Store) GetPeopleAnalyticsWorkspace(ctx context.Context, tenantID uuid.UUID) (*domain.PeopleAnalyticsWorkspace, error) {
	payload, err := s.getQueries(ctx).GetPeopleAnalyticsWorkspace(ctx, tenantID)
	if err != nil {
		return nil, s.logDBError(ctx, "get people analytics workspace", err, tenantIDField(tenantID))
	}
	return &domain.PeopleAnalyticsWorkspace{Workspace: json.RawMessage(payload)}, nil
}

func (s *Store) UpsertPrivacyConsent(ctx context.Context, item *domain.PrivacyConsent, actorID *uuid.UUID) (*domain.PrivacyConsent, error) {
	row, err := s.getQueries(ctx).UpsertPrivacyConsent(ctx, sqlc.UpsertPrivacyConsentParams{
		TenantID:        item.TenantID,
		ConsentKey:      item.ConsentKey,
		ConsentArea:     item.ConsentArea,
		Status:          item.Status,
		LawfulBasis:     item.LawfulBasis,
		Channel:         item.Channel,
		Source:          item.Source,
		Purpose:         item.Purpose,
		Evidence:        jsonBytesFromRaw(item.Evidence),
		Metadata:        jsonBytesFromRaw(item.Metadata),
		EmployeeUserID:  uuidFromPtr(item.EmployeeUserID),
		WorkerProfileID: uuidFromPtr(item.WorkerProfileID),
		GrantedAt:       timestamptzFromPtr(item.GrantedAt),
		RevokedAt:       timestamptzFromPtr(item.RevokedAt),
		ExpiresAt:       timestamptzFromPtr(item.ExpiresAt),
		ActorID:         uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "upsert privacy consent", err, tenantIDField(item.TenantID), stringField("consent_key", item.ConsentKey))
	}
	return mapPrivacyConsent(row), nil
}

func (s *Store) ListPrivacyConsents(ctx context.Context, filter domain.PrivacyEcosystemFilter) ([]*domain.PrivacyConsent, error) {
	rows, err := s.getQueries(ctx).ListPrivacyConsents(ctx, sqlc.ListPrivacyConsentsParams{TenantID: filter.TenantID, Limit: limitOrDefault(filter.Limit), Offset: filter.Offset, Status: textFromPtr(filter.Status), ConsentArea: textFromPtr(filter.ConsentArea)})
	if err != nil {
		return nil, s.logDBError(ctx, "list privacy consents", err, tenantIDField(filter.TenantID))
	}
	return mapPrivacyConsents(rows), nil
}

func (s *Store) CreateDataErasureRequest(ctx context.Context, item *domain.DataErasureRequest, actorID *uuid.UUID) (*domain.DataErasureRequest, error) {
	row, err := s.getQueries(ctx).CreateDataErasureRequest(ctx, sqlc.CreateDataErasureRequestParams{
		TenantID:        item.TenantID,
		RequestKey:      item.RequestKey,
		RequestType:     item.RequestType,
		Status:          item.Status,
		Priority:        item.Priority,
		Reason:          item.Reason,
		Scope:           jsonBytesFromRaw(item.Scope),
		AuditSummary:    jsonBytesFromRaw(item.AuditSummary),
		SubjectUserID:   uuidFromPtr(item.SubjectUserID),
		WorkerProfileID: uuidFromPtr(item.WorkerProfileID),
		RequestedBy:     uuidFromPtr(item.RequestedBy),
		RetainedReason:  textFromPtr(item.RetainedReason),
		DueAt:           timestamptzFromPtr(item.DueAt),
		ActorID:         uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "create data erasure request", err, tenantIDField(item.TenantID), stringField("request_key", item.RequestKey))
	}
	return mapDataErasureRequest(row), nil
}

func (s *Store) ListDataErasureRequests(ctx context.Context, filter domain.PrivacyEcosystemFilter) ([]*domain.DataErasureRequest, error) {
	rows, err := s.getQueries(ctx).ListDataErasureRequests(ctx, sqlc.ListDataErasureRequestsParams{TenantID: filter.TenantID, Limit: limitOrDefault(filter.Limit), Offset: filter.Offset, Status: textFromPtr(filter.Status), Priority: textFromPtr(filter.Priority)})
	if err != nil {
		return nil, s.logDBError(ctx, "list data erasure requests", err, tenantIDField(filter.TenantID))
	}
	return mapDataErasureRequests(rows), nil
}

func (s *Store) UpdateDataErasureRequestStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, retainedReason *string, auditSummary json.RawMessage, actorID *uuid.UUID) (*domain.DataErasureRequest, error) {
	row, err := s.getQueries(ctx).UpdateDataErasureRequestStatus(ctx, sqlc.UpdateDataErasureRequestStatusParams{TenantID: tenantID, ID: id, Status: status, RetainedReason: textFromPtr(retainedReason), AuditSummary: jsonBytesFromRaw(auditSummary), ActorID: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "update data erasure request status", err, tenantIDField(tenantID), stringField("request_id", id.String()))
	}
	return mapDataErasureRequest(row), nil
}

func (s *Store) UpsertEcosystemIntegrationHook(ctx context.Context, item *domain.EcosystemIntegrationHook, actorID *uuid.UUID) (*domain.EcosystemIntegrationHook, error) {
	row, err := s.getQueries(ctx).UpsertEcosystemIntegrationHook(ctx, sqlc.UpsertEcosystemIntegrationHookParams{
		TenantID:        item.TenantID,
		HookKey:         item.HookKey,
		Provider:        item.Provider,
		Channel:         item.Channel,
		Direction:       item.Direction,
		Status:          item.Status,
		DisplayName:     item.DisplayName,
		EventTypes:      item.EventTypes,
		ConsentRequired: item.ConsentRequired,
		MobileSafe:      item.MobileSafe,
		Config:          jsonBytesFromRaw(item.Config),
		EndpointUrl:     textFromPtr(item.EndpointURL),
		SecretRef:       textFromPtr(item.SecretRef),
		LastCheckedAt:   timestamptzFromPtr(item.LastCheckedAt),
		LastError:       textFromPtr(item.LastError),
		ActorID:         uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "upsert ecosystem integration hook", err, tenantIDField(item.TenantID), stringField("hook_key", item.HookKey))
	}
	return mapEcosystemIntegrationHook(row), nil
}

func (s *Store) ListEcosystemIntegrationHooks(ctx context.Context, filter domain.PrivacyEcosystemFilter) ([]*domain.EcosystemIntegrationHook, error) {
	rows, err := s.getQueries(ctx).ListEcosystemIntegrationHooks(ctx, sqlc.ListEcosystemIntegrationHooksParams{TenantID: filter.TenantID, Limit: limitOrDefault(filter.Limit), Offset: filter.Offset, Channel: textFromPtr(filter.Channel), Status: textFromPtr(filter.Status)})
	if err != nil {
		return nil, s.logDBError(ctx, "list ecosystem integration hooks", err, tenantIDField(filter.TenantID))
	}
	return mapEcosystemIntegrationHooks(rows), nil
}

func (s *Store) UpsertMobileAPIConstraint(ctx context.Context, item *domain.MobileAPIConstraint, actorID *uuid.UUID) (*domain.MobileAPIConstraint, error) {
	row, err := s.getQueries(ctx).UpsertMobileAPIConstraint(ctx, sqlc.UpsertMobileAPIConstraintParams{
		TenantID:              item.TenantID,
		ConstraintKey:         item.ConstraintKey,
		Workflow:              item.Workflow,
		OfflineSupported:      item.OfflineSupported,
		LowBandwidthMode:      item.LowBandwidthMode,
		RequiresLocation:      item.RequiresLocation,
		RequiresDeviceBinding: item.RequiresDeviceBinding,
		MaxPayloadKb:          item.MaxPayloadKB,
		Status:                item.Status,
		Config:                jsonBytesFromRaw(item.Config),
		MinAndroidVersion:     textFromPtr(item.MinAndroidVersion),
		MinIosVersion:         textFromPtr(item.MinIOSVersion),
		Notes:                 textFromPtr(item.Notes),
		ActorID:               uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "upsert mobile api constraint", err, tenantIDField(item.TenantID), stringField("constraint_key", item.ConstraintKey))
	}
	return mapMobileAPIConstraint(row), nil
}

func (s *Store) ListMobileAPIConstraints(ctx context.Context, filter domain.PrivacyEcosystemFilter) ([]*domain.MobileAPIConstraint, error) {
	rows, err := s.getQueries(ctx).ListMobileAPIConstraints(ctx, sqlc.ListMobileAPIConstraintsParams{TenantID: filter.TenantID, Limit: limitOrDefault(filter.Limit), Offset: filter.Offset, Workflow: textFromPtr(filter.Workflow), Status: textFromPtr(filter.Status)})
	if err != nil {
		return nil, s.logDBError(ctx, "list mobile api constraints", err, tenantIDField(filter.TenantID))
	}
	return mapMobileAPIConstraints(rows), nil
}
