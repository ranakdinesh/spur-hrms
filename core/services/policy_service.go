package services

import (
	"context"
	"encoding/base64"
	"strings"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) CreatePolicyType(ctx context.Context, cmd ports.PolicyTypeCommand) (*domain.PolicyType, error) {
	item, err := domain.NewPolicyType(domain.PolicyTypeInput{TenantID: cmd.TenantID, Name: cmd.Name, IsSystem: cmd.IsSystem})
	if err != nil {
		s.logError("validate policy type create", err, serviceTenantIDField(cmd.TenantID), serviceStringField("policy_type_name", cmd.Name))
		return nil, err
	}
	result, err := s.policies.CreatePolicyType(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create policy type", err, serviceTenantIDField(cmd.TenantID), serviceStringField("policy_type_name", item.Name))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) ListPolicyTypes(ctx context.Context, tenantID uuid.UUID) ([]*domain.PolicyType, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate policy type list tenant", err)
		return nil, err
	}
	result, err := s.policies.ListPolicyTypes(ctx, tenantID)
	if err != nil {
		s.logError("list policy types", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) GetPolicyType(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.PolicyType, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate policy type get tenant", err)
		return nil, err
	}
	if id == uuid.Nil {
		err := domain.ErrInvalidPolicyTypeID
		s.logError("validate policy type get id", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	result, err := s.policies.GetPolicyType(ctx, tenantID, id)
	if err != nil {
		s.logError("get policy type", err, serviceTenantIDField(tenantID), serviceStringField("policy_type_id", id.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) UpdatePolicyType(ctx context.Context, cmd ports.PolicyTypeCommand) (*domain.PolicyType, error) {
	if cmd.ID == uuid.Nil {
		err := domain.ErrInvalidPolicyTypeID
		s.logError("validate policy type update id", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	existing, err := s.GetPolicyType(ctx, cmd.TenantID, cmd.ID)
	if err != nil {
		return nil, err
	}
	if existing.IsSystem {
		err := domain.ErrSystemPolicyReadOnly
		s.logError("validate policy type update system", err, serviceTenantIDField(cmd.TenantID), serviceStringField("policy_type_id", cmd.ID.String()))
		return nil, err
	}
	item, err := domain.NewPolicyType(domain.PolicyTypeInput{TenantID: cmd.TenantID, Name: cmd.Name})
	if err != nil {
		s.logError("validate policy type update", err, serviceTenantIDField(cmd.TenantID), serviceStringField("policy_type_id", cmd.ID.String()))
		return nil, err
	}
	item.ID = cmd.ID
	result, err := s.policies.UpdatePolicyType(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("update policy type", err, serviceTenantIDField(cmd.TenantID), serviceStringField("policy_type_id", cmd.ID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) DeletePolicyType(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate policy type delete tenant", err)
		return err
	}
	if id == uuid.Nil {
		err := domain.ErrInvalidPolicyTypeID
		s.logError("validate policy type delete id", err, serviceTenantIDField(tenantID))
		return err
	}
	existing, err := s.GetPolicyType(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if existing.IsSystem {
		err := domain.ErrSystemPolicyReadOnly
		s.logError("validate policy type delete system", err, serviceTenantIDField(tenantID), serviceStringField("policy_type_id", id.String()))
		return err
	}
	if err := s.policies.DeletePolicyType(ctx, tenantID, id, actorID); err != nil {
		s.logError("delete policy type", err, serviceTenantIDField(tenantID), serviceStringField("policy_type_id", id.String()))
		return err
	}
	return nil
}

func (s *TenantService) CreateCompanyPolicy(ctx context.Context, cmd ports.CompanyPolicyCommand) (*domain.CompanyPolicy, error) {
	if err := s.validateCompanyPolicyReferences(ctx, cmd.TenantID, cmd.PolicyTypeID); err != nil {
		return nil, err
	}
	fileContent, err := decodePolicyFile(cmd.FileContentBase64)
	if err != nil {
		s.logError("decode company policy file", err, serviceTenantIDField(cmd.TenantID), serviceStringField("policy_title", cmd.Title))
		return nil, err
	}
	if len(fileContent) > 0 && s.policyStorage == nil {
		err := domain.ErrPolicyStorageMissing
		s.logError("store company policy file", err, serviceTenantIDField(cmd.TenantID), serviceStringField("policy_title", cmd.Title))
		return nil, err
	}
	item, err := domain.NewCompanyPolicy(domain.CompanyPolicyInput{TenantID: cmd.TenantID, PolicyTypeID: cmd.PolicyTypeID, Title: cmd.Title, FilePath: cmd.FilePath, Description: cmd.Description})
	if err != nil {
		s.logError("validate company policy create", err, serviceTenantIDField(cmd.TenantID), serviceStringField("policy_title", cmd.Title))
		return nil, err
	}
	result, err := s.policies.CreateCompanyPolicy(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create company policy", err, serviceTenantIDField(cmd.TenantID), serviceStringField("policy_title", item.Title))
		return nil, err
	}
	if len(fileContent) > 0 {
		path, err := s.policyStorage.StorePolicyFile(ctx, ports.StorePolicyFileInput{TenantID: result.TenantID, PolicyID: result.ID, FileName: cmd.FileName, ContentType: cmd.FileContentType, Content: fileContent})
		if err != nil {
			s.logError("store company policy file", err, serviceTenantIDField(result.TenantID), serviceStringField("policy_id", result.ID.String()))
			return nil, err
		}
		result.FilePath = &path
		result, err = s.policies.UpdateCompanyPolicy(ctx, result, cmd.ActorID)
		if err != nil {
			s.logError("attach company policy file", err, serviceTenantIDField(result.TenantID), serviceStringField("policy_id", result.ID.String()))
			return nil, err
		}
	}
	s.notifyCompanyPolicyChanged(ctx, result, "created", cmd.ActorID)
	return result, nil
}

func (s *TenantService) ListCompanyPolicies(ctx context.Context, tenantID uuid.UUID, policyTypeID *uuid.UUID) ([]*domain.CompanyPolicy, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate company policy list tenant", err)
		return nil, err
	}
	result, err := s.policies.ListCompanyPolicies(ctx, tenantID, policyTypeID)
	if err != nil {
		s.logError("list company policies", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) GetCompanyPolicy(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.CompanyPolicy, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate company policy get tenant", err)
		return nil, err
	}
	if id == uuid.Nil {
		err := domain.ErrInvalidPolicyID
		s.logError("validate company policy get id", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	result, err := s.policies.GetCompanyPolicy(ctx, tenantID, id)
	if err != nil {
		s.logError("get company policy", err, serviceTenantIDField(tenantID), serviceStringField("policy_id", id.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) UpdateCompanyPolicy(ctx context.Context, cmd ports.CompanyPolicyCommand) (*domain.CompanyPolicy, error) {
	if cmd.ID == uuid.Nil {
		err := domain.ErrInvalidPolicyID
		s.logError("validate company policy update id", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	if err := s.validateCompanyPolicyReferences(ctx, cmd.TenantID, cmd.PolicyTypeID); err != nil {
		return nil, err
	}
	fileContent, err := decodePolicyFile(cmd.FileContentBase64)
	if err != nil {
		s.logError("decode company policy file", err, serviceTenantIDField(cmd.TenantID), serviceStringField("policy_id", cmd.ID.String()))
		return nil, err
	}
	if len(fileContent) > 0 && s.policyStorage == nil {
		err := domain.ErrPolicyStorageMissing
		s.logError("store company policy file", err, serviceTenantIDField(cmd.TenantID), serviceStringField("policy_id", cmd.ID.String()))
		return nil, err
	}
	item, err := domain.NewCompanyPolicy(domain.CompanyPolicyInput{TenantID: cmd.TenantID, PolicyTypeID: cmd.PolicyTypeID, Title: cmd.Title, FilePath: cmd.FilePath, Description: cmd.Description})
	if err != nil {
		s.logError("validate company policy update", err, serviceTenantIDField(cmd.TenantID), serviceStringField("policy_id", cmd.ID.String()))
		return nil, err
	}
	item.ID = cmd.ID
	if len(fileContent) > 0 {
		path, err := s.policyStorage.StorePolicyFile(ctx, ports.StorePolicyFileInput{TenantID: item.TenantID, PolicyID: item.ID, FileName: cmd.FileName, ContentType: cmd.FileContentType, Content: fileContent})
		if err != nil {
			s.logError("store company policy file", err, serviceTenantIDField(item.TenantID), serviceStringField("policy_id", item.ID.String()))
			return nil, err
		}
		item.FilePath = &path
	}
	result, err := s.policies.UpdateCompanyPolicy(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("update company policy", err, serviceTenantIDField(cmd.TenantID), serviceStringField("policy_id", cmd.ID.String()))
		return nil, err
	}
	s.notifyCompanyPolicyChanged(ctx, result, "updated", cmd.ActorID)
	return result, nil
}

func (s *TenantService) DeleteCompanyPolicy(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	policy, err := s.GetCompanyPolicy(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if err := s.policies.DeleteCompanyPolicy(ctx, tenantID, id, actorID); err != nil {
		s.logError("delete company policy", err, serviceTenantIDField(tenantID), serviceStringField("policy_id", id.String()))
		return err
	}
	s.notifyCompanyPolicyChanged(ctx, policy, "deleted", actorID)
	return nil
}

func (s *TenantService) validateCompanyPolicyReferences(ctx context.Context, tenantID uuid.UUID, policyTypeID *uuid.UUID) error {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate company policy tenant", err)
		return err
	}
	if policyTypeID != nil && *policyTypeID != uuid.Nil {
		if _, err := s.policies.GetPolicyType(ctx, tenantID, *policyTypeID); err != nil {
			s.logError("validate company policy type", err, serviceTenantIDField(tenantID), serviceStringField("policy_type_id", policyTypeID.String()))
			return err
		}
	}
	return nil
}

func (s *TenantService) notifyCompanyPolicyChanged(ctx context.Context, item *domain.CompanyPolicy, action string, actorID *uuid.UUID) {
	if item == nil {
		return
	}
	if s.policyNotifier == nil {
		if s.log != nil {
			s.log.Warn().Str("tenant_id", item.TenantID.String()).Str("policy_id", item.ID.String()).Str("action", action).Msg("hrms: company policy notification hook not configured")
		}
		return
	}
	if err := s.policyNotifier.CompanyPolicyChanged(ctx, ports.CompanyPolicyChangedEvent{TenantID: item.TenantID, PolicyID: item.ID, Title: item.Title, Action: action, ActorID: actorID}); err != nil {
		s.logError("notify company policy changed", err, serviceTenantIDField(item.TenantID), serviceStringField("policy_id", item.ID.String()), serviceStringField("action", action))
	}
}

func decodePolicyFile(value string) ([]byte, error) {
	clean := strings.TrimSpace(value)
	if clean == "" {
		return nil, nil
	}
	if comma := strings.Index(clean, ","); comma >= 0 {
		clean = clean[comma+1:]
	}
	content, err := base64.StdEncoding.DecodeString(clean)
	if err != nil {
		return nil, domain.ErrInvalidPolicyFile
	}
	return content, nil
}
