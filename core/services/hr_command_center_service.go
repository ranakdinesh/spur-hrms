package services

import (
	"context"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

const maxHRCommandCenterCards = 350

type commandCenterSectionDefinition struct {
	key          string
	title        string
	description  string
	routeSection string
	actionLabel  string
	match        func(*domain.OperationsWorkbenchCard, time.Time) bool
}

func (s *TenantService) GetHRCommandCenter(ctx context.Context, query ports.HRCommandCenterQuery) (*domain.HRCommandCenter, error) {
	if query.TenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.log.Warn().Err(err).Str("operation", "get hr command center").Msg("invalid hr command center query")
		return nil, err
	}
	limit := query.Limit
	if limit <= 0 {
		limit = maxHRCommandCenterCards
	}
	if limit > maxHRCommandCenterCards {
		limit = maxHRCommandCenterCards
	}
	workbench, err := s.GetOperationsWorkbench(ctx, ports.OperationsWorkbenchQuery{TenantID: query.TenantID, Limit: limit})
	if err != nil {
		s.log.Error().Err(err).Str("operation", "get hr command center").Str("tenant_id", query.TenantID.String()).Msg("hr command center aggregation failed")
		return nil, err
	}
	return buildHRCommandCenter(workbench), nil
}

func buildHRCommandCenter(workbench *domain.OperationsWorkbench) *domain.HRCommandCenter {
	generatedAt := time.Now().UTC()
	if workbench != nil && !workbench.GeneratedAt.IsZero() {
		generatedAt = workbench.GeneratedAt
	}
	summary := domain.OperationsWorkbenchSummary{ByLane: map[string]int32{}, ByCategory: map[string]int32{}, BySeverity: map[string]int32{}}
	var cards []*domain.OperationsWorkbenchCard
	if workbench != nil {
		summary = workbench.Summary
		cards = append(cards, workbench.Cards...)
	}
	sortCommandCenterCards(cards, generatedAt)
	sections := make([]*domain.HRCommandCenterSection, 0, len(commandCenterSectionDefinitions()))
	for _, def := range commandCenterSectionDefinitions() {
		sectionCards := filterCommandCenterCards(cards, generatedAt, def.match)
		sections = append(sections, &domain.HRCommandCenterSection{
			Key:         def.key,
			Title:       def.title,
			Description: def.description,
			Count:       int32(len(sectionCards)),
			HighPriority: countCommandCenterCards(sectionCards, generatedAt, func(card *domain.OperationsWorkbenchCard, _ time.Time) bool {
				return isHighPriorityCommandCenterCard(card)
			}),
			Overdue:      countCommandCenterCards(sectionCards, generatedAt, isOverdueCommandCenterCard),
			DueToday:     countCommandCenterCards(sectionCards, generatedAt, isDueTodayCommandCenterCard),
			RouteSection: def.routeSection,
			ActionLabel:  def.actionLabel,
			Cards:        topCommandCenterCards(sectionCards, 5),
		})
	}
	return &domain.HRCommandCenter{
		GeneratedAt:  generatedAt,
		Summary:      summary,
		Sections:     sections,
		QuickActions: buildHRCommandCenterActions(sections),
	}
}

func commandCenterSectionDefinitions() []commandCenterSectionDefinition {
	return []commandCenterSectionDefinition{
		{key: "today", title: "Today", description: "Time-sensitive HR work due now or carrying high operational risk.", routeSection: "operations-workbench", actionLabel: "Open Today", match: isTodayCommandCenterCard},
		{key: "approvals", title: "Approvals", description: "Leave, payroll, benefits, access, and workflow approvals waiting for action.", routeSection: "workflow-inbox", actionLabel: "Open Approvals", match: func(card *domain.OperationsWorkbenchCard, _ time.Time) bool {
			return card != nil && (card.Lane == "approvals" || strings.Contains(card.RouteSection, "approvals") || strings.Contains(strings.ToLower(card.ActionLabel), "approve"))
		}},
		{key: "exceptions", title: "Exceptions", description: "Attendance, compliance, employee relations, and other items that need HR review.", routeSection: "operations-workbench", actionLabel: "Review Exceptions", match: func(card *domain.OperationsWorkbenchCard, _ time.Time) bool {
			return card != nil && card.Lane == "exceptions"
		}},
		{key: "payroll_close", title: "Payroll Close", description: "Payroll blockers, salary sheet issues, pay run readiness, and compensation handoff items.", routeSection: "payroll-operations", actionLabel: "Open Payroll", match: isPayrollCloseCommandCenterCard},
		{key: "people_movement", title: "People Movement", description: "Joining, onboarding, exit, transfer, and workforce change work.", routeSection: "workforce-hub", actionLabel: "Open People Movement", match: isPeopleMovementCommandCenterCard},
		{key: "employee_requests", title: "Employee Requests", description: "Employee initiated cases, documents, claims, and service requests.", routeSection: "hr-helpdesk", actionLabel: "Open Requests", match: func(card *domain.OperationsWorkbenchCard, _ time.Time) bool {
			return card != nil && card.Lane == "employee_requests"
		}},
		{key: "insights", title: "Insights", description: "AI recommendations and analytics nudges surfaced inside the HR operating flow.", routeSection: "insights", actionLabel: "Open Insights", match: isInsightCommandCenterCard},
	}
}

func filterCommandCenterCards(cards []*domain.OperationsWorkbenchCard, now time.Time, match func(*domain.OperationsWorkbenchCard, time.Time) bool) []*domain.OperationsWorkbenchCard {
	filtered := make([]*domain.OperationsWorkbenchCard, 0)
	for _, card := range cards {
		if match(card, now) {
			filtered = append(filtered, card)
		}
	}
	return filtered
}

func countCommandCenterCards(cards []*domain.OperationsWorkbenchCard, now time.Time, match func(*domain.OperationsWorkbenchCard, time.Time) bool) int32 {
	var count int32
	for _, card := range cards {
		if match(card, now) {
			count++
		}
	}
	return count
}

func topCommandCenterCards(cards []*domain.OperationsWorkbenchCard, limit int) []*domain.OperationsWorkbenchCard {
	if len(cards) <= limit {
		return cards
	}
	return cards[:limit]
}

func sortCommandCenterCards(cards []*domain.OperationsWorkbenchCard, now time.Time) {
	sort.SliceStable(cards, func(i, j int) bool {
		left, right := cards[i], cards[j]
		if left == nil {
			return false
		}
		if right == nil {
			return true
		}
		leftOverdue, rightOverdue := isOverdueCommandCenterCard(left, now), isOverdueCommandCenterCard(right, now)
		if leftOverdue != rightOverdue {
			return leftOverdue
		}
		if left.Priority != right.Priority {
			return left.Priority < right.Priority
		}
		if left.DueAt != nil && right.DueAt != nil && !left.DueAt.Equal(*right.DueAt) {
			return left.DueAt.Before(*right.DueAt)
		}
		if left.DueAt != nil && right.DueAt == nil {
			return true
		}
		if left.DueAt == nil && right.DueAt != nil {
			return false
		}
		return left.DetectedAt.After(right.DetectedAt)
	})
}

func buildHRCommandCenterActions(sections []*domain.HRCommandCenterSection) []*domain.HRCommandCenterAction {
	actions := make([]*domain.HRCommandCenterAction, 0, len(sections))
	for index, section := range sections {
		if section == nil {
			continue
		}
		actions = append(actions, &domain.HRCommandCenterAction{
			Key:          section.Key,
			Title:        section.ActionLabel,
			Description:  section.Description,
			RouteSection: section.RouteSection,
			BadgeCount:   section.Count,
			Priority:     int32(index + 1),
		})
	}
	return actions
}

func isTodayCommandCenterCard(card *domain.OperationsWorkbenchCard, now time.Time) bool {
	return isDueTodayCommandCenterCard(card, now) || isOverdueCommandCenterCard(card, now) || isHighPriorityCommandCenterCard(card)
}

func isHighPriorityCommandCenterCard(card *domain.OperationsWorkbenchCard) bool {
	return card != nil && (card.Priority <= 25 || card.Severity == "critical" || card.Severity == "high")
}

func isOverdueCommandCenterCard(card *domain.OperationsWorkbenchCard, now time.Time) bool {
	return card != nil && card.DueAt != nil && card.DueAt.UTC().Before(now)
}

func isDueTodayCommandCenterCard(card *domain.OperationsWorkbenchCard, now time.Time) bool {
	if card == nil || card.DueAt == nil {
		return false
	}
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	tomorrow := today.AddDate(0, 0, 1)
	due := card.DueAt.UTC()
	return !due.Before(today) && due.Before(tomorrow)
}

func isPayrollCloseCommandCenterCard(card *domain.OperationsWorkbenchCard, _ time.Time) bool {
	if card == nil {
		return false
	}
	return card.Lane == "payroll_blockers" || containsAnyCommandCenterValue(card.RouteSection, "payroll", "salary", "payslip", "compensation") || containsAnyCommandCenterValue(card.SourceModule, "payroll", "salary", "compensation")
}

func isPeopleMovementCommandCenterCard(card *domain.OperationsWorkbenchCard, _ time.Time) bool {
	if card == nil {
		return false
	}
	return card.Lane == "joining" || card.Lane == "exit" || containsAnyCommandCenterValue(card.RouteSection, "onboarding", "employee-exits", "workforce-hub") || containsAnyCommandCenterValue(card.Category, "joining", "exit", "onboarding", "transfer")
}

func isInsightCommandCenterCard(card *domain.OperationsWorkbenchCard, _ time.Time) bool {
	if card == nil {
		return false
	}
	return card.Lane == "ai_recommendations" || containsAnyCommandCenterValue(card.RouteSection, "insights", "people-analytics") || containsAnyCommandCenterValue(card.Category, "ai", "insight", "analytics")
}

func containsAnyCommandCenterValue(value string, needles ...string) bool {
	clean := strings.ToLower(value)
	for _, needle := range needles {
		if strings.Contains(clean, needle) {
			return true
		}
	}
	return false
}
