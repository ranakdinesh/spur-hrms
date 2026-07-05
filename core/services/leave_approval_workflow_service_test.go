package services

import (
	"testing"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func TestLeaveApprovalStepComplete(t *testing.T) {
	leaveID := uuid.New()
	currentID := uuid.New()
	otherID := uuid.New()
	tests := []struct {
		name     string
		current  *domain.LeaveApproval
		items    []*domain.LeaveApproval
		wantDone bool
	}{
		{"nil current", nil, nil, false},
		{
			name:    "single all approval completes",
			current: &domain.LeaveApproval{ID: currentID, LeaveID: leaveID, StepOrder: 1, DecisionRule: domain.LeaveApprovalDecisionAll, RequiredApprovals: 1},
			items: []*domain.LeaveApproval{
				{ID: currentID, LeaveID: leaveID, StepOrder: 1, Status: domain.LeaveStatusPending},
			},
			wantDone: true,
		},
		{
			name:    "all requires every same-step approval",
			current: &domain.LeaveApproval{ID: currentID, LeaveID: leaveID, StepOrder: 1, DecisionRule: domain.LeaveApprovalDecisionAll, RequiredApprovals: 1},
			items: []*domain.LeaveApproval{
				{ID: currentID, LeaveID: leaveID, StepOrder: 1, Status: domain.LeaveStatusPending},
				{ID: otherID, LeaveID: leaveID, StepOrder: 1, Status: domain.LeaveStatusPending},
			},
			wantDone: false,
		},
		{
			name:    "all completes when other same-step approval already approved",
			current: &domain.LeaveApproval{ID: currentID, LeaveID: leaveID, StepOrder: 1, DecisionRule: domain.LeaveApprovalDecisionAll, RequiredApprovals: 1},
			items: []*domain.LeaveApproval{
				{ID: currentID, LeaveID: leaveID, StepOrder: 1, Status: domain.LeaveStatusPending},
				{ID: otherID, LeaveID: leaveID, StepOrder: 1, Status: domain.LeaveStatusApproved},
			},
			wantDone: true,
		},
		{
			name:    "any completes after required count",
			current: &domain.LeaveApproval{ID: currentID, LeaveID: leaveID, StepOrder: 1, DecisionRule: domain.LeaveApprovalDecisionAny, RequiredApprovals: 1},
			items: []*domain.LeaveApproval{
				{ID: currentID, LeaveID: leaveID, StepOrder: 1, Status: domain.LeaveStatusPending},
				{ID: otherID, LeaveID: leaveID, StepOrder: 1, Status: domain.LeaveStatusPending},
			},
			wantDone: true,
		},
		{
			name:    "ignores other leave and step approvals",
			current: &domain.LeaveApproval{ID: currentID, LeaveID: leaveID, StepOrder: 2, DecisionRule: domain.LeaveApprovalDecisionAll, RequiredApprovals: 1},
			items: []*domain.LeaveApproval{
				{ID: currentID, LeaveID: leaveID, StepOrder: 2, Status: domain.LeaveStatusPending},
				{ID: uuid.New(), LeaveID: uuid.New(), StepOrder: 2, Status: domain.LeaveStatusPending},
				{ID: uuid.New(), LeaveID: leaveID, StepOrder: 1, Status: domain.LeaveStatusPending},
			},
			wantDone: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := leaveApprovalStepComplete(tt.items, tt.current); got != tt.wantDone {
				t.Fatalf("got %v, want %v", got, tt.wantDone)
			}
		})
	}
}

func TestLeaveApprovalStepAlreadyCreated(t *testing.T) {
	items := []*domain.LeaveApproval{{StepOrder: 1}, {StepOrder: 3}}
	if !leaveApprovalStepAlreadyCreated(items, 3) {
		t.Fatal("expected existing step order")
	}
	if leaveApprovalStepAlreadyCreated(items, 2) {
		t.Fatal("did not expect missing step order")
	}
}
