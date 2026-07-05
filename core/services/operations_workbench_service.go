package services

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

const maxOperationsWorkbenchCards = 500

func (s *TenantService) GetOperationsWorkbench(ctx context.Context, query ports.OperationsWorkbenchQuery) (*domain.OperationsWorkbench, error) {
	if query.TenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.log.Warn().Err(err).Str("operation", "get operations workbench").Msg("invalid operations workbench query")
		return nil, err
	}
	limit := query.Limit
	if limit <= 0 {
		limit = 200
	}
	if limit > maxOperationsWorkbenchCards {
		limit = maxOperationsWorkbenchCards
	}
	filter := domain.OperationsWorkbenchFilter{
		TenantID: query.TenantID,
		Lane:     cleanOptionalString(query.Lane),
		Category: cleanOptionalString(query.Category),
		Severity: cleanOptionalString(query.Severity),
		Search:   cleanOptionalString(query.Search),
		Limit:    limit,
		Offset:   query.Offset,
	}
	cards, err := s.operationsWorkbench.ListOperationsWorkbenchCards(ctx, filter)
	if err != nil {
		s.log.Error().Err(err).Str("operation", "get operations workbench").Str("tenant_id", query.TenantID.String()).Msg("operations workbench aggregation failed")
		return nil, err
	}
	generatedAt := time.Now().UTC()
	applyOperationsWorkbenchActions(cards)
	return &domain.OperationsWorkbench{GeneratedAt: generatedAt, Summary: buildOperationsWorkbenchSummary(cards, generatedAt), Cards: cards}, nil
}

func (s *TenantService) ActOperationsWorkbenchCard(ctx context.Context, cmd ports.OperationsWorkbenchActionCommand) (*domain.OperationsWorkbenchActionResult, error) {
	if cmd.TenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate operations workbench action tenant", err)
		return nil, err
	}
	cardKey := strings.TrimSpace(cmd.CardKey)
	action := strings.ToLower(strings.TrimSpace(cmd.Action))
	if cardKey == "" || action == "" {
		err := domain.ErrUnsupportedOperationsWorkbenchAction
		s.logError("validate operations workbench action", err, serviceTenantIDField(cmd.TenantID), serviceStringField("card_key", cardKey), serviceStringField("action", action))
		return nil, err
	}
	cards, err := s.operationsWorkbench.ListOperationsWorkbenchCards(ctx, domain.OperationsWorkbenchFilter{TenantID: cmd.TenantID, Limit: maxOperationsWorkbenchCards})
	if err != nil {
		s.logError("load operations workbench action card", err, serviceTenantIDField(cmd.TenantID), serviceStringField("card_key", cardKey))
		return nil, err
	}
	applyOperationsWorkbenchActions(cards)
	card := findOperationsWorkbenchCard(cards, cardKey)
	if card == nil {
		err := domain.ErrOperationsWorkbenchCardNotFound
		s.logError("operations workbench action card missing", err, serviceTenantIDField(cmd.TenantID), serviceStringField("card_key", cardKey))
		return nil, err
	}
	if !operationsWorkbenchActionAllowed(card.Actions, action) || action == domain.WorkItemActionOpenRecord {
		err := domain.ErrUnsupportedOperationsWorkbenchAction
		s.logError("operations workbench action unsupported", err, serviceTenantIDField(cmd.TenantID), serviceStringField("card_key", cardKey), serviceStringField("action", action), serviceStringField("source_type", card.SourceType))
		return nil, err
	}
	switch card.SourceType {
	case "leave_approval":
		return s.actLeaveApprovalWorkbenchCard(ctx, card, cmd, action)
	default:
		err := domain.ErrUnsupportedOperationsWorkbenchAction
		s.logError("operations workbench source action unsupported", err, serviceTenantIDField(cmd.TenantID), serviceStringField("card_key", cardKey), serviceStringField("action", action), serviceStringField("source_type", card.SourceType))
		return nil, err
	}
}

func buildOperationsWorkbenchSummary(cards []*domain.OperationsWorkbenchCard, now time.Time) domain.OperationsWorkbenchSummary {
	summary := domain.OperationsWorkbenchSummary{ByLane: map[string]int32{}, ByCategory: map[string]int32{}, BySeverity: map[string]int32{}}
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	tomorrow := today.AddDate(0, 0, 1)
	for _, card := range cards {
		if card == nil {
			continue
		}
		summary.Total++
		summary.ByLane[card.Lane]++
		summary.ByCategory[card.Category]++
		summary.BySeverity[card.Severity]++
		if card.Priority <= 25 || card.Severity == "critical" || card.Severity == "high" {
			summary.HighPriority++
		}
		if card.DueAt != nil {
			due := card.DueAt.UTC()
			if due.Before(now) {
				summary.Overdue++
			}
			if !due.Before(today) && due.Before(tomorrow) {
				summary.DueToday++
			}
		}
		switch card.Lane {
		case "approvals":
			summary.Approvals++
		case "exceptions":
			summary.Exceptions++
		case "payroll_blockers":
			summary.PayrollBlockers++
		case "joining":
			summary.Joining++
		case "exit":
			summary.Exits++
		case "compliance":
			summary.Compliance++
		case "ai_recommendations":
			summary.AIRecommendations++
		case "employee_requests":
			summary.EmployeeRequests++
		}
		if card.Category == "documents" {
			summary.Documents++
		}
	}
	return summary
}

func applyOperationsWorkbenchActions(cards []*domain.OperationsWorkbenchCard) {
	for _, card := range cards {
		if card == nil {
			continue
		}
		card.Actions = operationsWorkbenchActions(card)
	}
}

func operationsWorkbenchActions(card *domain.OperationsWorkbenchCard) []domain.WorkItemAction {
	openRecord := domain.WorkItemAction{Key: domain.WorkItemActionOpenRecord, Label: "Open full record", Tone: "neutral", Inline: false}
	switch card.SourceType {
	case "leave_approval":
		return []domain.WorkItemAction{
			{Key: domain.WorkItemActionApprove, Label: "Approve", Tone: "positive", Primary: true, Inline: true, RemarksPlaceholder: "Add an approval note... ", CompletionBadge: "Approved"},
			{Key: domain.WorkItemActionReject, Label: "Reject", Tone: "danger", Inline: true, RequiresRemarks: true, RemarksPlaceholder: "Reason for rejection", CompletionBadge: "Rejected"},
			openRecord,
		}
	default:
		return []domain.WorkItemAction{openRecord}
	}
}

func findOperationsWorkbenchCard(cards []*domain.OperationsWorkbenchCard, cardKey string) *domain.OperationsWorkbenchCard {
	for _, card := range cards {
		if card != nil && card.CardKey == cardKey {
			return card
		}
	}
	return nil
}

func operationsWorkbenchActionAllowed(actions []domain.WorkItemAction, action string) bool {
	for _, item := range actions {
		if item.Key == action {
			return true
		}
	}
	return false
}

func (s *TenantService) actLeaveApprovalWorkbenchCard(ctx context.Context, card *domain.OperationsWorkbenchCard, cmd ports.OperationsWorkbenchActionCommand, action string) (*domain.OperationsWorkbenchActionResult, error) {
	approverID := operationsWorkbenchMetadataUUID(card.Metadata, "approver_id")
	if approverID == uuid.Nil && cmd.ActorID != nil {
		approverID = *cmd.ActorID
	}
	var source any
	var err error
	switch action {
	case domain.WorkItemActionApprove:
		source, err = s.ApproveLeave(ctx, ports.ApproveLeaveCommand{TenantID: cmd.TenantID, ApprovalID: card.SourceID, ApproverID: approverID, Remarks: cleanOptionalString(cmd.Remarks), ActorID: cmd.ActorID})
	case domain.WorkItemActionReject:
		source, err = s.RejectLeave(ctx, ports.RejectLeaveCommand{TenantID: cmd.TenantID, ApprovalID: card.SourceID, ApproverID: approverID, Remarks: cleanOptionalString(cmd.Remarks), ActorID: cmd.ActorID})
	default:
		err = domain.ErrUnsupportedOperationsWorkbenchAction
	}
	if err != nil {
		s.logError("act leave approval workbench card", err, serviceTenantIDField(cmd.TenantID), serviceStringField("card_key", card.CardKey), serviceStringField("action", action))
		return nil, err
	}
	sourceJSON, err := json.Marshal(source)
	if err != nil {
		sourceJSON = json.RawMessage(`{}`)
	}
	return &domain.OperationsWorkbenchActionResult{CardKey: card.CardKey, Action: action, Status: "completed", Badge: operationsWorkbenchCompletionBadge(card.Actions, action), Card: card, Source: sourceJSON}, nil
}

func operationsWorkbenchCompletionBadge(actions []domain.WorkItemAction, action string) string {
	for _, item := range actions {
		if item.Key == action {
			return item.CompletionBadge
		}
	}
	return "Completed"
}

func operationsWorkbenchMetadataUUID(metadata json.RawMessage, key string) uuid.UUID {
	var values map[string]any
	if len(metadata) == 0 || json.Unmarshal(metadata, &values) != nil {
		return uuid.Nil
	}
	raw, ok := values[key].(string)
	if !ok {
		return uuid.Nil
	}
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil
	}
	return id
}

func cleanOptionalString(value *string) *string {
	if value == nil {
		return nil
	}
	clean := strings.TrimSpace(*value)
	if clean == "" {
		return nil
	}
	return &clean
}
