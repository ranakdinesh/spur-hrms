package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) notifyLeaveApplied(ctx context.Context, application *domain.LeaveApplication, actorID *uuid.UUID) {
	if application == nil || application.Leave == nil || application.Approval == nil || application.Approval.ApproverID == uuid.Nil {
		return
	}
	referenceTable := "hrms.leaves"
	referenceID := application.Leave.ID
	title := "Leave request pending approval"
	message := fmt.Sprintf("A leave request for %s to %s is waiting for your approval.", application.Leave.StartDate.Format("02 Jan 2006"), application.Leave.EndDate.Format("02 Jan 2006"))
	if _, err := s.SendNotification(ctx, ports.SendNotificationCommand{
		TenantID:         application.Leave.TenantID,
		NotificationCode: domain.NotifLeaveApplied,
		UserIDs:          []uuid.UUID{application.Approval.ApproverID},
		Title:            title,
		Message:          message,
		ReferenceTable:   &referenceTable,
		ReferenceID:      &referenceID,
		Channels:         []string{"in_app", domain.NotifChannelEmail},
		ActorID:          actorID,
	}); err != nil {
		s.logError("send leave applied notification", err, serviceTenantIDField(application.Leave.TenantID), serviceStringField("leave_id", application.Leave.ID.String()))
	}
}

func (s *TenantService) notifyLeaveReviewed(ctx context.Context, application *domain.LeaveApplication, notificationCode string, actorID *uuid.UUID) {
	if application == nil || application.Leave == nil || application.Leave.UserID == uuid.Nil {
		return
	}
	referenceTable := "hrms.leaves"
	referenceID := application.Leave.ID
	status := "reviewed"
	if notificationCode == domain.NotifLeaveApproved {
		status = "approved"
	}
	if notificationCode == domain.NotifLeaveRejected {
		status = "rejected"
	}
	title := fmt.Sprintf("Leave request %s", status)
	message := fmt.Sprintf("Your leave request for %s to %s has been %s.", application.Leave.StartDate.Format("02 Jan 2006"), application.Leave.EndDate.Format("02 Jan 2006"), status)
	if _, err := s.SendNotification(ctx, ports.SendNotificationCommand{
		TenantID:         application.Leave.TenantID,
		NotificationCode: notificationCode,
		UserIDs:          []uuid.UUID{application.Leave.UserID},
		Title:            title,
		Message:          message,
		ReferenceTable:   &referenceTable,
		ReferenceID:      &referenceID,
		Channels:         []string{"in_app", domain.NotifChannelEmail},
		ActorID:          actorID,
	}); err != nil {
		s.logError("send leave reviewed notification", err, serviceTenantIDField(application.Leave.TenantID), serviceStringField("leave_id", application.Leave.ID.String()))
	}
}
