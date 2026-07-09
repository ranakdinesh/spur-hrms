package services

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) GetLeave(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID uuid.UUID, canManage bool) (*domain.Leave, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate get leave tenant", err)
		return nil, err
	}
	if id == uuid.Nil {
		err := domain.ErrInvalidLeaveID
		s.logError("validate get leave id", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	leave, err := s.leaveRequests.GetLeave(ctx, tenantID, id)
	if err != nil {
		s.logError("get leave", err, serviceTenantIDField(tenantID), serviceStringField("leave_id", id.String()))
		return nil, err
	}
	if canManage || leave.UserID == actorID || s.isLeaveApprover(ctx, tenantID, leave.ID, actorID) {
		return leave, nil
	}
	err = domain.ErrLeaveMessageUnauthorized
	s.logError("authorize get leave", err, serviceTenantIDField(tenantID), serviceStringField("leave_id", id.String()), serviceStringField("actor_id", actorID.String()))
	return nil, err
}

func (s *TenantService) ListLeaveRequestMessages(ctx context.Context, tenantID uuid.UUID, leaveID uuid.UUID, actorID uuid.UUID, canManage bool) ([]*domain.LeaveRequestMessage, error) {
	if _, err := s.GetLeave(ctx, tenantID, leaveID, actorID, canManage); err != nil {
		return nil, err
	}
	items, err := s.leaveRequests.ListLeaveRequestMessages(ctx, tenantID, leaveID)
	if err != nil {
		s.logError("list leave request messages", err, serviceTenantIDField(tenantID), serviceStringField("leave_id", leaveID.String()))
		return nil, err
	}
	return items, nil
}

func (s *TenantService) CreateLeaveRequestMessage(ctx context.Context, cmd ports.LeaveRequestMessageCommand) (*domain.LeaveRequestMessage, error) {
	if cmd.ActorID == nil || *cmd.ActorID == uuid.Nil {
		err := domain.ErrInvalidLeaveUser
		s.logError("validate leave message actor", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	if cmd.SenderUserID == uuid.Nil {
		cmd.SenderUserID = *cmd.ActorID
	}
	body := strings.TrimSpace(cmd.Body)
	if body == "" {
		err := domain.ErrInvalidLeaveMessage
		s.logError("validate leave message body", err, serviceTenantIDField(cmd.TenantID), serviceStringField("leave_id", cmd.LeaveID.String()))
		return nil, err
	}
	leave, err := s.GetLeave(ctx, cmd.TenantID, cmd.LeaveID, cmd.SenderUserID, cmd.CanManage)
	if err != nil {
		return nil, err
	}
	messageType := strings.TrimSpace(cmd.MessageType)
	if messageType == "" {
		messageType = domain.LeaveMessageClarificationRequest
		if cmd.SenderUserID == leave.UserID {
			messageType = domain.LeaveMessageEmployeeReply
		}
	}
	if !validLeaveMessageType(messageType) {
		err := domain.ErrInvalidLeaveMessageType
		s.logError("validate leave message type", err, serviceTenantIDField(cmd.TenantID), serviceStringField("leave_id", cmd.LeaveID.String()), serviceStringField("message_type", messageType))
		return nil, err
	}
	item := &domain.LeaveRequestMessage{
		TenantID:        cmd.TenantID,
		LeaveID:         cmd.LeaveID,
		SenderUserID:    cmd.SenderUserID,
		RecipientUserID: cmd.RecipientUserID,
		MessageType:     messageType,
		Body:            body,
	}
	created, err := s.leaveRequests.CreateLeaveRequestMessage(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create leave request message", err, serviceTenantIDField(cmd.TenantID), serviceStringField("leave_id", cmd.LeaveID.String()))
		return nil, err
	}
	s.notifyLeaveMessage(ctx, leave, created, cmd.ActorID)
	return created, nil
}

func validLeaveMessageType(value string) bool {
	switch value {
	case domain.LeaveMessageClarificationRequest, domain.LeaveMessageEmployeeReply, domain.LeaveMessageComment:
		return true
	default:
		return false
	}
}

func (s *TenantService) isLeaveApprover(ctx context.Context, tenantID uuid.UUID, leaveID uuid.UUID, actorID uuid.UUID) bool {
	if actorID == uuid.Nil {
		return false
	}
	approvals, err := s.approvalWorkflows.ListLeaveApprovalsByLeave(ctx, tenantID, leaveID)
	if err != nil {
		s.logError("list leave approvers for authorization", err, serviceTenantIDField(tenantID), serviceStringField("leave_id", leaveID.String()))
		return false
	}
	for _, approval := range approvals {
		if approval != nil && approval.ApproverID == actorID {
			return true
		}
	}
	return false
}

func (s *TenantService) leaveMessageRecipients(ctx context.Context, leave *domain.Leave, message *domain.LeaveRequestMessage) []uuid.UUID {
	if leave == nil || message == nil {
		return nil
	}
	if message.RecipientUserID != nil && *message.RecipientUserID != uuid.Nil && *message.RecipientUserID != message.SenderUserID {
		return []uuid.UUID{*message.RecipientUserID}
	}
	recipients := []uuid.UUID{}
	if message.SenderUserID != leave.UserID {
		recipients = append(recipients, leave.UserID)
	} else {
		approvals, err := s.approvalWorkflows.ListLeaveApprovalsByLeave(ctx, leave.TenantID, leave.ID)
		if err != nil {
			s.logError("list leave approvers for message recipients", err, serviceTenantIDField(leave.TenantID), serviceStringField("leave_id", leave.ID.String()))
			return recipients
		}
		for _, approval := range approvals {
			if approval != nil && approval.ApproverID != uuid.Nil && approval.ApproverID != message.SenderUserID {
				recipients = append(recipients, approval.ApproverID)
			}
		}
	}
	return uniqueUUIDs(recipients)
}

func (s *TenantService) notifyLeaveMessage(ctx context.Context, leave *domain.Leave, message *domain.LeaveRequestMessage, actorID *uuid.UUID) {
	recipients := s.leaveMessageRecipients(ctx, leave, message)
	if leave == nil || message == nil || len(recipients) == 0 {
		return
	}
	referenceTable := "hrms.leave_request_messages"
	referenceID := message.ID
	title := "Leave request message"
	body := "A leave request has a new clarification message."
	if message.MessageType == domain.LeaveMessageEmployeeReply {
		body = "An employee replied on a leave request."
	}
	if message.MessageType == domain.LeaveMessageClarificationRequest {
		body = "Your manager requested clarification on a leave request."
	}
	if _, err := s.SendNotification(ctx, ports.SendNotificationCommand{
		TenantID:         leave.TenantID,
		NotificationCode: domain.NotifLeaveClarify,
		UserIDs:          recipients,
		Title:            title,
		Message:          body,
		ReferenceTable:   &referenceTable,
		ReferenceID:      &referenceID,
		Channels:         []string{"in_app", domain.NotifChannelEmail},
		ActorID:          actorID,
	}); err != nil {
		s.logError("send leave message notification", err, serviceTenantIDField(leave.TenantID), serviceStringField("leave_id", leave.ID.String()), serviceStringField("message_id", message.ID.String()))
	}
}
