package domain

import "time"

type HRCommandCenter struct {
	GeneratedAt  time.Time                  `json:"generated_at"`
	Summary      OperationsWorkbenchSummary `json:"summary"`
	Sections     []*HRCommandCenterSection  `json:"sections"`
	QuickActions []*HRCommandCenterAction   `json:"quick_actions"`
}

type HRCommandCenterSection struct {
	Key          string                     `json:"key"`
	Title        string                     `json:"title"`
	Description  string                     `json:"description"`
	Count        int32                      `json:"count"`
	HighPriority int32                      `json:"high_priority"`
	Overdue      int32                      `json:"overdue"`
	DueToday     int32                      `json:"due_today"`
	RouteSection string                     `json:"route_section"`
	ActionLabel  string                     `json:"action_label"`
	Cards        []*OperationsWorkbenchCard `json:"cards"`
}

type HRCommandCenterAction struct {
	Key          string `json:"key"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	RouteSection string `json:"route_section"`
	Permission   string `json:"permission,omitempty"`
	BadgeCount   int32  `json:"badge_count"`
	Priority     int32  `json:"priority"`
}
