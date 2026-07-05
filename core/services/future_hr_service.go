package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) ListBoundedAIAgents(ctx context.Context, tenantID uuid.UUID) ([]domain.BoundedAIAgentDefinition, error) {
	if tenantID == uuid.Nil {
		return nil, domain.ErrInvalidBoundedAIAgentRun
	}
	return boundedAIAgentDefinitions(), nil
}

func (s *TenantService) RunBoundedAIAgents(ctx context.Context, cmd ports.BoundedAIAgentRunCommand) (*domain.BoundedAIAgentRunResult, error) {
	if cmd.TenantID == uuid.Nil {
		s.logError("validate bounded ai agent run", domain.ErrInvalidBoundedAIAgentRun)
		return nil, domain.ErrInvalidBoundedAIAgentRun
	}
	definitions := boundedAIAgentDefinitions()
	selected := map[string]bool{}
	for _, key := range cmd.Agents {
		key = strings.TrimSpace(key)
		if key != "" {
			selected[key] = true
		}
	}
	runAll := len(selected) == 0
	actions := make([]*domain.AIAgentActionLog, 0, len(definitions))
	runDate := time.Now().UTC().Format("20060102")
	modelVersion := "deterministic-bounded-v1"
	for _, agent := range definitions {
		if !runAll && !selected[agent.Key] {
			continue
		}
		action, err := s.CreateAIAgentAction(ctx, ports.AIActionCommand{
			TenantID:            cmd.TenantID,
			ActionKey:           fmt.Sprintf("%s-%s", agent.Key, runDate),
			AgentKey:            agent.Key,
			AgentName:           agent.Name,
			ActionType:          "workflow_recommendation",
			Status:              domain.AIActionStatusProposed,
			Severity:            agent.Severity,
			Title:               agent.Name,
			Summary:             boundedAgentSummary(agent),
			VisibilityScope:     agent.VisibilityScope,
			ProposedAction:      rawJSON(map[string]any{"workflow": agent.Workflow, "steps": boundedAgentSteps(agent.Key)}),
			InputSnapshot:       rawJSON(map[string]any{"signals": agent.Signals, "run_date": runDate}),
			OutputSnapshot:      rawJSON(map[string]any{"mode": "deterministic_fallback", "sidecar_required": false}),
			Explainability:      rawJSON(map[string]any{"guardrails": agent.Guardrails, "reason": "Generated from configured HRMS workflow signals with human review required."}),
			ConfidenceScore:     72,
			ModelVersion:        &modelVersion,
			RequiresHumanReview: true,
			ActorID:             cmd.ActorID,
		})
		if err != nil {
			s.logError("run bounded ai agent", err, serviceTenantIDField(cmd.TenantID), serviceStringField("agent_key", agent.Key))
			return nil, err
		}
		actions = append(actions, action)
	}
	return &domain.BoundedAIAgentRunResult{Agents: definitions, Actions: actions}, nil
}

func (s *TenantService) GetPeopleAnalyticsWorkspace(ctx context.Context, tenantID uuid.UUID) (*domain.PeopleAnalyticsWorkspace, error) {
	workspace, err := s.peopleAnalytics.GetPeopleAnalyticsWorkspace(ctx, tenantID)
	if err != nil {
		s.logError("get people analytics workspace", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	return workspace, nil
}

func (s *TenantService) ListPrivacyEcosystemWorkspace(ctx context.Context, filter domain.PrivacyEcosystemFilter) (*domain.PrivacyEcosystemWorkspace, error) {
	filter.Limit = limitAI(filter.Limit)
	consents, err := s.privacyEcosystem.ListPrivacyConsents(ctx, filter)
	if err != nil {
		s.logError("list privacy consents", err, serviceTenantIDField(filter.TenantID))
		return nil, err
	}
	erasure, err := s.privacyEcosystem.ListDataErasureRequests(ctx, filter)
	if err != nil {
		s.logError("list erasure requests", err, serviceTenantIDField(filter.TenantID))
		return nil, err
	}
	integrations, err := s.privacyEcosystem.ListEcosystemIntegrationHooks(ctx, filter)
	if err != nil {
		s.logError("list ecosystem integrations", err, serviceTenantIDField(filter.TenantID))
		return nil, err
	}
	mobile, err := s.privacyEcosystem.ListMobileAPIConstraints(ctx, filter)
	if err != nil {
		s.logError("list mobile api constraints", err, serviceTenantIDField(filter.TenantID))
		return nil, err
	}
	return &domain.PrivacyEcosystemWorkspace{Consents: consents, Erasure: erasure, Integrations: integrations, Mobile: mobile, Summary: privacySummary(consents, erasure, integrations, mobile)}, nil
}

func (s *TenantService) UpsertPrivacyConsent(ctx context.Context, cmd ports.PrivacyConsentCommand) (*domain.PrivacyConsent, error) {
	grantedAt, err := parseOptionalAITime(cmd.GrantedAt)
	if err != nil {
		return nil, err
	}
	revokedAt, err := parseOptionalAITime(cmd.RevokedAt)
	if err != nil {
		return nil, err
	}
	expiresAt, err := parseOptionalAITime(cmd.ExpiresAt)
	if err != nil {
		return nil, err
	}
	item, err := domain.NewPrivacyConsent(domain.PrivacyConsent{TenantID: cmd.TenantID, EmployeeUserID: cmd.EmployeeUserID, WorkerProfileID: cmd.WorkerProfileID, ConsentKey: cmd.ConsentKey, ConsentArea: cmd.ConsentArea, Status: cmd.Status, LawfulBasis: cmd.LawfulBasis, Channel: cmd.Channel, Source: cmd.Source, Purpose: cmd.Purpose, GrantedAt: grantedAt, RevokedAt: revokedAt, ExpiresAt: expiresAt, Evidence: cmd.Evidence, Metadata: cmd.Metadata})
	if err != nil {
		s.logError("validate privacy consent", err, serviceTenantIDField(cmd.TenantID), serviceStringField("consent_key", cmd.ConsentKey))
		return nil, err
	}
	return s.privacyEcosystem.UpsertPrivacyConsent(ctx, item, cmd.ActorID)
}

func (s *TenantService) CreateDataErasureRequest(ctx context.Context, cmd ports.DataErasureRequestCommand) (*domain.DataErasureRequest, error) {
	dueAt, err := parseOptionalAITime(cmd.DueAt)
	if err != nil {
		return nil, err
	}
	item, err := domain.NewDataErasureRequest(domain.DataErasureRequest{TenantID: cmd.TenantID, RequestKey: cmd.RequestKey, SubjectUserID: cmd.SubjectUserID, WorkerProfileID: cmd.WorkerProfileID, RequestType: cmd.RequestType, Status: cmd.Status, Priority: cmd.Priority, RequestedBy: cmd.RequestedBy, Reason: cmd.Reason, Scope: cmd.Scope, RetainedReason: cmd.RetainedReason, DueAt: dueAt, AuditSummary: cmd.AuditSummary})
	if err != nil {
		s.logError("validate erasure request", err, serviceTenantIDField(cmd.TenantID), serviceStringField("request_key", cmd.RequestKey))
		return nil, err
	}
	return s.privacyEcosystem.CreateDataErasureRequest(ctx, item, cmd.ActorID)
}

func (s *TenantService) UpdateDataErasureRequestStatus(ctx context.Context, cmd ports.DataErasureStatusCommand) (*domain.DataErasureRequest, error) {
	status := strings.TrimSpace(cmd.Status)
	if status == "" || cmd.TenantID == uuid.Nil || cmd.ID == uuid.Nil {
		err := domain.ErrInvalidDataErasureRequest
		s.logError("validate erasure request status", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	return s.privacyEcosystem.UpdateDataErasureRequestStatus(ctx, cmd.TenantID, cmd.ID, status, cmd.RetainedReason, cmd.AuditSummary, cmd.ActorID)
}

func (s *TenantService) UpsertEcosystemIntegrationHook(ctx context.Context, cmd ports.EcosystemIntegrationHookCommand) (*domain.EcosystemIntegrationHook, error) {
	item, err := domain.NewEcosystemIntegrationHook(domain.EcosystemIntegrationHook{TenantID: cmd.TenantID, HookKey: cmd.HookKey, Provider: cmd.Provider, Channel: cmd.Channel, Direction: cmd.Direction, Status: cmd.Status, DisplayName: cmd.DisplayName, EndpointURL: cmd.EndpointURL, EventTypes: cmd.EventTypes, SecretRef: cmd.SecretRef, ConsentRequired: cmd.ConsentRequired, MobileSafe: cmd.MobileSafe, Config: cmd.Config})
	if err != nil {
		s.logError("validate ecosystem integration", err, serviceTenantIDField(cmd.TenantID), serviceStringField("hook_key", cmd.HookKey))
		return nil, err
	}
	return s.privacyEcosystem.UpsertEcosystemIntegrationHook(ctx, item, cmd.ActorID)
}

func (s *TenantService) UpsertMobileAPIConstraint(ctx context.Context, cmd ports.MobileAPIConstraintCommand) (*domain.MobileAPIConstraint, error) {
	item, err := domain.NewMobileAPIConstraint(domain.MobileAPIConstraint{TenantID: cmd.TenantID, ConstraintKey: cmd.ConstraintKey, Workflow: cmd.Workflow, MinAndroidVersion: cmd.MinAndroidVersion, MinIOSVersion: cmd.MinIOSVersion, OfflineSupported: cmd.OfflineSupported, LowBandwidthMode: cmd.LowBandwidthMode, RequiresLocation: cmd.RequiresLocation, RequiresDeviceBinding: cmd.RequiresDeviceBinding, MaxPayloadKB: cmd.MaxPayloadKB, Status: cmd.Status, Notes: cmd.Notes, Config: cmd.Config})
	if err != nil {
		s.logError("validate mobile api constraint", err, serviceTenantIDField(cmd.TenantID), serviceStringField("constraint_key", cmd.ConstraintKey))
		return nil, err
	}
	return s.privacyEcosystem.UpsertMobileAPIConstraint(ctx, item, cmd.ActorID)
}

func boundedAIAgentDefinitions() []domain.BoundedAIAgentDefinition {
	return []domain.BoundedAIAgentDefinition{
		{Key: "leave_intelligence", Name: "Leave Intelligence Agent", Workflow: "leave_approval", Severity: "medium", VisibilityScope: domain.AIVisibilityHR, Signals: []string{"leave balance", "overlap", "team absence", "policy exception"}, Guardrails: []string{"Never approve leave automatically.", "Show policy reasons and balance impact.", "Require manager/HR final action."}},
		{Key: "payroll_anomaly", Name: "Payroll Anomaly Agent", Workflow: "payroll_run", Severity: "high", VisibilityScope: domain.AIVisibilityHR, Signals: []string{"LOP mismatch", "variable import variance", "statutory missing", "duplicate slip"}, Guardrails: []string{"Never change payroll amounts automatically.", "Show source rows and reconciliation reason.", "Require payroll reviewer approval."}},
		{Key: "attrition_risk", Name: "Attrition Risk Agent", Workflow: "retention", Severity: "high", VisibilityScope: domain.AIVisibilityHR, Signals: []string{"low engagement", "high absence", "feedback trend", "wellbeing risk"}, Guardrails: []string{"Do not expose individual risk to managers without HR review.", "Use aggregate manager views.", "Require human retention plan."}},
		{Key: "engagement_health", Name: "Engagement Health Agent", Workflow: "engagement", Severity: "medium", VisibilityScope: domain.AIVisibilityManagerAggregate, Signals: []string{"project load", "pulse trend", "check-in mood", "work log burden"}, Guardrails: []string{"Prefer aggregate team signals.", "Suppress small groups.", "Avoid medical or sensitive inference."}},
		{Key: "compliance_guard", Name: "Compliance Guard Agent", Workflow: "compliance", Severity: "high", VisibilityScope: domain.AIVisibilityHR, Signals: []string{"missing documents", "expired certificates", "statutory checklist", "policy acknowledgement"}, Guardrails: []string{"Create checklist recommendation only.", "Keep evidence trail.", "Do not waive compliance automatically."}},
		{Key: "onboarding_autopilot", Name: "Onboarding Autopilot Agent", Workflow: "onboarding", Severity: "medium", VisibilityScope: domain.AIVisibilityHR, Signals: []string{"offer accepted", "joining date", "document pending", "task delay"}, Guardrails: []string{"Create tasks and reminders only after HR review.", "Respect tenant communication consent.", "Keep candidate/employee audit trail."}},
	}
}

func boundedAgentSummary(agent domain.BoundedAIAgentDefinition) string {
	return fmt.Sprintf("%s reviewed %s signals and proposed a human-approved next action.", agent.Name, agent.Workflow)
}

func boundedAgentSteps(agentKey string) []string {
	switch agentKey {
	case "leave_intelligence":
		return []string{"Review balance and overlap", "Check team capacity", "Send approve/reject recommendation"}
	case "payroll_anomaly":
		return []string{"Compare payroll inputs", "Flag variance rows", "Hold affected group for reviewer"}
	case "attrition_risk":
		return []string{"Review HR-only risk signals", "Prepare retention conversation", "Track human outcome"}
	case "engagement_health":
		return []string{"Aggregate team trend", "Suppress small groups", "Recommend manager action"}
	case "compliance_guard":
		return []string{"Scan pending checklist", "Prioritize critical evidence", "Create follow-up task"}
	case "onboarding_autopilot":
		return []string{"Check joining readiness", "Find pending tasks", "Recommend reminders"}
	default:
		return []string{"Review signal", "Propose action", "Wait for human approval"}
	}
}

func privacySummary(consents []*domain.PrivacyConsent, erasure []*domain.DataErasureRequest, integrations []*domain.EcosystemIntegrationHook, mobile []*domain.MobileAPIConstraint) map[string]int32 {
	summary := map[string]int32{"consents": int32(len(consents)), "erasure_open": 0, "integrations_active": 0, "mobile_constraints": int32(len(mobile))}
	for _, item := range erasure {
		if item.Status != "completed" && item.Status != "rejected" && item.Status != "cancelled" {
			summary["erasure_open"]++
		}
	}
	for _, item := range integrations {
		if item.Status == "active" {
			summary["integrations_active"]++
		}
	}
	return summary
}

func parseOptionalAITime(value string) (*time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, nil
	}
	parsed, err := parseAITime(value)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}
