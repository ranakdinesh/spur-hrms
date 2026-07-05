package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) CreateNotificationMaster(ctx context.Context, cmd ports.NotificationMasterCommand) (*domain.NotificationMaster, error) {
	item, err := domain.NewNotificationMaster(domain.NotificationMasterInput{TenantID: cmd.TenantID, Code: cmd.Code, Name: cmd.Name, Description: cmd.Description, IsInAppEnabled: cmd.IsInAppEnabled, IsEmailEnabled: cmd.IsEmailEnabled, IsPushEnabled: cmd.IsPushEnabled, EmailSubjectTemplate: cmd.EmailSubjectTemplate, EmailTextTemplate: cmd.EmailTextTemplate, EmailHTMLTemplate: cmd.EmailHTMLTemplate})
	if err != nil {
		s.logError("validate notification master create", err, serviceTenantIDField(cmd.TenantID), serviceStringField("notification_code", cmd.Code))
		return nil, err
	}
	result, err := s.notifications.CreateNotificationMaster(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create notification master", err, serviceTenantIDField(cmd.TenantID), serviceStringField("notification_code", item.Code))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) ListNotificationMasters(ctx context.Context, tenantID uuid.UUID, actorID *uuid.UUID) ([]*domain.NotificationMaster, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate notification master list tenant", err)
		return nil, err
	}
	if _, err := s.seedMissingNotificationMasters(ctx, tenantID, actorID); err != nil {
		return nil, err
	}
	items, err := s.notifications.ListNotificationMasters(ctx, tenantID)
	if err != nil {
		s.logError("list notification masters", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	return items, nil
}

func (s *TenantService) GetNotificationMaster(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.NotificationMaster, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate notification master get tenant", err)
		return nil, err
	}
	if id == uuid.Nil {
		err := domain.ErrInvalidNotificationMasterID
		s.logError("validate notification master get id", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	item, err := s.notifications.GetNotificationMaster(ctx, tenantID, id)
	if err != nil {
		s.logError("get notification master", err, serviceTenantIDField(tenantID), serviceStringField("notification_master_id", id.String()))
		return nil, err
	}
	return item, nil
}

func (s *TenantService) UpdateNotificationMaster(ctx context.Context, cmd ports.NotificationMasterCommand) (*domain.NotificationMaster, error) {
	if cmd.ID == uuid.Nil {
		err := domain.ErrInvalidNotificationMasterID
		s.logError("validate notification master update id", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	item, err := domain.NewNotificationMaster(domain.NotificationMasterInput{TenantID: cmd.TenantID, Code: cmd.Code, Name: cmd.Name, Description: cmd.Description, IsInAppEnabled: cmd.IsInAppEnabled, IsEmailEnabled: cmd.IsEmailEnabled, IsPushEnabled: cmd.IsPushEnabled, EmailSubjectTemplate: cmd.EmailSubjectTemplate, EmailTextTemplate: cmd.EmailTextTemplate, EmailHTMLTemplate: cmd.EmailHTMLTemplate})
	if err != nil {
		s.logError("validate notification master update", err, serviceTenantIDField(cmd.TenantID), serviceStringField("notification_master_id", cmd.ID.String()))
		return nil, err
	}
	item.ID = cmd.ID
	result, err := s.notifications.UpdateNotificationMaster(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("update notification master", err, serviceTenantIDField(cmd.TenantID), serviceStringField("notification_master_id", cmd.ID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) DeleteNotificationMaster(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate notification master delete tenant", err)
		return err
	}
	if id == uuid.Nil {
		err := domain.ErrInvalidNotificationMasterID
		s.logError("validate notification master delete id", err, serviceTenantIDField(tenantID))
		return err
	}
	if err := s.notifications.DeleteNotificationMaster(ctx, tenantID, id, actorID); err != nil {
		s.logError("delete notification master", err, serviceTenantIDField(tenantID), serviceStringField("notification_master_id", id.String()))
		return err
	}
	return nil
}

func (s *TenantService) ListEffectiveNotificationPreferences(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID, actorID *uuid.UUID) ([]*domain.EffectiveNotificationPreference, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate notification preference tenant", err)
		return nil, err
	}
	if userID == uuid.Nil {
		err := domain.ErrInvalidNotificationUser
		s.logError("validate notification preference user", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	masters, err := s.ListNotificationMasters(ctx, tenantID, actorID)
	if err != nil {
		return nil, err
	}
	preferences, err := s.notifications.ListNotificationPreferencesByUser(ctx, tenantID, userID)
	if err != nil {
		s.logError("list notification preferences", err, serviceTenantIDField(tenantID), serviceStringField("user_id", userID.String()))
		return nil, err
	}
	emailReady := false
	if settings, err := s.activeEmailProviderSettings(ctx, tenantID); err == nil && settings != nil && settings.IsEnabled {
		emailReady = true
	}
	pushReady := false
	if settings, err := s.activePushProviderSettings(ctx, tenantID); err == nil && settings != nil && settings.IsEnabled {
		pushReady = true
	}
	byMasterID := map[uuid.UUID]*domain.NotificationPreference{}
	for _, pref := range preferences {
		if pref != nil {
			byMasterID[pref.NotificationMasterID] = pref
		}
	}
	result := make([]*domain.EffectiveNotificationPreference, 0, len(masters))
	for _, master := range masters {
		if master == nil {
			continue
		}
		pref := byMasterID[master.ID]
		effective := &domain.EffectiveNotificationPreference{
			Master:              master,
			Preference:          pref,
			IsInAppEnabled:      master.IsInAppEnabled,
			IsEmailEnabled:      master.IsEmailEnabled,
			IsPushEnabled:       master.IsPushEnabled,
			DigestFrequency:     domain.DigestFrequencyInstant,
			UsesTenantDefaults:  pref == nil,
			ProviderReady:       map[string]bool{"in_app": true, "email": emailReady, "push": pushReady},
			FutureReadyFeatures: []string{"category controls", "channel opt-outs", "quiet hours", "digest frequency", "email provider webhooks", "flutter device token rotation"},
		}
		if pref != nil {
			effective.IsInAppEnabled = master.IsInAppEnabled && pref.IsInAppEnabled
			effective.IsEmailEnabled = master.IsEmailEnabled && pref.IsEmailEnabled
			effective.IsPushEnabled = master.IsPushEnabled && pref.IsPushEnabled
			effective.DigestFrequency = pref.DigestFrequency
			effective.QuietHoursStart = pref.QuietHoursStart
			effective.QuietHoursEnd = pref.QuietHoursEnd
			effective.Timezone = pref.Timezone
		}
		result = append(result, effective)
	}
	return result, nil
}

func (s *TenantService) UpsertNotificationPreference(ctx context.Context, cmd ports.NotificationPreferenceCommand) (*domain.NotificationPreference, error) {
	item, err := domain.NewNotificationPreference(domain.NotificationPreferenceInput{TenantID: cmd.TenantID, UserID: cmd.UserID, NotificationMasterID: cmd.NotificationMasterID, IsInAppEnabled: cmd.IsInAppEnabled, IsEmailEnabled: cmd.IsEmailEnabled, IsPushEnabled: cmd.IsPushEnabled, DigestFrequency: cmd.DigestFrequency, QuietHoursStart: cmd.QuietHoursStart, QuietHoursEnd: cmd.QuietHoursEnd, Timezone: cmd.Timezone})
	if err != nil {
		s.logError("validate notification preference upsert", err, serviceTenantIDField(cmd.TenantID), serviceStringField("user_id", cmd.UserID.String()))
		return nil, err
	}
	if _, err := s.GetNotificationMaster(ctx, cmd.TenantID, cmd.NotificationMasterID); err != nil {
		return nil, err
	}
	result, err := s.notifications.UpsertNotificationPreference(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("upsert notification preference", err, serviceTenantIDField(cmd.TenantID), serviceStringField("user_id", cmd.UserID.String()), serviceStringField("notification_master_id", cmd.NotificationMasterID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) DeleteNotificationPreference(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate notification preference delete tenant", err)
		return err
	}
	if id == uuid.Nil {
		err := domain.ErrInvalidNotificationMasterID
		s.logError("validate notification preference delete id", err, serviceTenantIDField(tenantID))
		return err
	}
	if err := s.notifications.DeleteNotificationPreference(ctx, tenantID, id, actorID); err != nil {
		s.logError("delete notification preference", err, serviceTenantIDField(tenantID), serviceStringField("notification_preference_id", id.String()))
		return err
	}
	return nil
}

func (s *TenantService) seedMissingNotificationMasters(ctx context.Context, tenantID uuid.UUID, actorID *uuid.UUID) (int, error) {
	existing, err := s.notifications.ListNotificationMasters(ctx, tenantID)
	if err != nil {
		s.logError("list notification masters for seed", err, serviceTenantIDField(tenantID))
		return 0, err
	}
	byCode := map[string]bool{}
	for _, item := range existing {
		if item != nil {
			byCode[item.Code] = true
		}
	}
	count := 0
	for _, input := range domain.DefaultNotificationMasterInputs(tenantID) {
		code, err := domain.ValidateNotificationCode(input.Code)
		if err != nil {
			s.logError("validate default notification master", err, serviceTenantIDField(tenantID), serviceStringField("notification_code", input.Code))
			return count, err
		}
		if byCode[code] {
			continue
		}
		s.log.Warn().Str("tenant_id", tenantID.String()).Str("notification_code", code).Msg("hrms: notification master missing, seeding default")
		item, err := domain.NewNotificationMaster(input)
		if err != nil {
			s.logError("create default notification master domain", err, serviceTenantIDField(tenantID), serviceStringField("notification_code", code))
			return count, err
		}
		if _, err = s.notifications.CreateNotificationMaster(ctx, item, actorID); err != nil {
			s.logError("seed notification master", err, serviceTenantIDField(tenantID), serviceStringField("notification_code", code))
			return count, err
		}
		count++
	}
	return count, nil
}

func (s *TenantService) SendNotification(ctx context.Context, cmd ports.SendNotificationCommand) (*domain.NotificationSendResult, error) {
	input := domain.NotificationSendInput{TenantID: cmd.TenantID, NotificationCode: cmd.NotificationCode, UserIDs: cmd.UserIDs, Title: cmd.Title, Message: cmd.Message, ReferenceTable: cmd.ReferenceTable, ReferenceID: cmd.ReferenceID, Channels: cmd.Channels, IdempotencyKey: cmd.IdempotencyKey, ActorID: cmd.ActorID}
	result, err := s.sendNotification(ctx, input)
	if err != nil {
		s.logError("send notification", err, serviceTenantIDField(cmd.TenantID), serviceStringField("notification_code", cmd.NotificationCode))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) ListNotificationInbox(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID, isRead *bool, notificationMasterID *uuid.UUID, limit int32, offset int32) ([]*domain.NotificationInboxItem, error) {
	page, err := s.ListNotificationInboxPage(ctx, domain.NotificationInboxFilter{TenantID: tenantID, UserID: userID, IsRead: isRead, NotificationMasterID: notificationMasterID, Limit: limit, Offset: offset})
	if err != nil {
		return nil, err
	}
	return page.Items, nil
}

func (s *TenantService) ListNotificationInboxPage(ctx context.Context, filter domain.NotificationInboxFilter) (*domain.NotificationInboxPage, error) {
	tenantID := filter.TenantID
	userID := filter.UserID
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate notification inbox tenant", err)
		return nil, err
	}
	if userID == uuid.Nil {
		err := domain.ErrInvalidNotificationUser
		s.logError("validate notification inbox user", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	filter.TenantID = tenantID
	filter.UserID = userID
	filter.NotificationCode = cleanStringPtr(filter.NotificationCode)
	filter.Search = cleanStringPtr(filter.Search)
	limit, offset := normalizeListWindow(filter.Limit, filter.Offset)
	filter.Limit = limit
	filter.Offset = offset
	items, err := s.notifications.ListNotificationInboxByUser(ctx, filter)
	if err != nil {
		s.logError("list notification inbox", err, serviceTenantIDField(tenantID), serviceStringField("user_id", userID.String()))
		return nil, err
	}
	total, err := s.notifications.CountNotificationInboxByUser(ctx, filter)
	if err != nil {
		s.logError("count notification inbox", err, serviceTenantIDField(tenantID), serviceStringField("user_id", userID.String()))
		return nil, err
	}
	unread, err := s.notifications.CountUnreadNotificationsByUser(ctx, tenantID, userID)
	if err != nil {
		s.logError("count unread notifications", err, serviceTenantIDField(tenantID), serviceStringField("user_id", userID.String()))
		return nil, err
	}
	page := &domain.NotificationInboxPage{Items: items, Total: total, Unread: unread, Limit: limit, Offset: offset}
	if int64(offset)+int64(len(items)) < total {
		next := offset + limit
		page.NextOffset = &next
	}
	return page, nil
}

func (s *TenantService) CountUnreadNotifications(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID) (int64, error) {
	count, err := s.notifications.CountUnreadNotificationsByUser(ctx, tenantID, userID)
	if err != nil {
		s.logError("count unread notifications", err, serviceTenantIDField(tenantID), serviceStringField("user_id", userID.String()))
		return 0, err
	}
	return count, nil
}

func (s *TenantService) MarkNotificationRead(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.notifications.MarkNotificationInboxItemRead(ctx, tenantID, id, actorID); err != nil {
		s.logError("mark notification read", err, serviceTenantIDField(tenantID), serviceStringField("notification_id", id.String()))
		return err
	}
	return nil
}

func (s *TenantService) MarkNotificationUnread(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.notifications.MarkNotificationInboxItemUnread(ctx, tenantID, id, actorID); err != nil {
		s.logError("mark notification unread", err, serviceTenantIDField(tenantID), serviceStringField("notification_id", id.String()))
		return err
	}
	return nil
}

func (s *TenantService) MarkAllNotificationsRead(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID, actorID *uuid.UUID) error {
	if err := s.notifications.MarkNotificationInboxAllRead(ctx, tenantID, userID, actorID); err != nil {
		s.logError("mark all notifications read", err, serviceTenantIDField(tenantID), serviceStringField("user_id", userID.String()))
		return err
	}
	return nil
}

func (s *TenantService) DeleteNotificationInboxItem(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate notification archive tenant", err)
		return err
	}
	if id == uuid.Nil {
		err := domain.ErrInvalidNotificationMasterID
		s.logError("validate notification archive id", err, serviceTenantIDField(tenantID))
		return err
	}
	if err := s.notifications.DeleteNotificationInboxItem(ctx, tenantID, id, actorID); err != nil {
		s.logError("archive notification inbox item", err, serviceTenantIDField(tenantID), serviceStringField("notification_id", id.String()))
		return err
	}
	return nil
}

func (s *TenantService) ListNotificationLogs(ctx context.Context, tenantID uuid.UUID, userID *uuid.UUID, status *string, channel *string, bulkID *uuid.UUID, limit int32, offset int32) ([]*domain.NotificationLog, error) {
	limit, offset = normalizeListWindow(limit, offset)
	items, err := s.notifications.ListNotificationLogs(ctx, tenantID, userID, cleanStringPtr(status), cleanStringPtr(channel), bulkID, limit, offset)
	if err != nil {
		s.logError("list notification logs", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	return items, nil
}

func (s *TenantService) sendNotification(ctx context.Context, input domain.NotificationSendInput) (*domain.NotificationSendResult, error) {
	if input.TenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	code, err := domain.ValidateNotificationCode(input.NotificationCode)
	if err != nil {
		return nil, domain.ErrInvalidNotificationCode
	}
	if strings.TrimSpace(input.Title) == "" {
		return nil, domain.ErrInvalidNotificationTitle
	}
	if strings.TrimSpace(input.Message) == "" {
		return nil, domain.ErrInvalidNotificationMessage
	}
	bulkID := uuid.New()
	if input.BulkID != nil && *input.BulkID != uuid.Nil {
		bulkID = *input.BulkID
	}
	result := &domain.NotificationSendResult{BulkID: bulkID}
	masters, err := s.ListNotificationMasters(ctx, input.TenantID, input.ActorID)
	if err != nil {
		return nil, err
	}
	var master *domain.NotificationMaster
	for _, item := range masters {
		if item != nil && item.Code == code {
			master = item
			break
		}
	}
	if master == nil {
		result.MissingMaster = true
		s.log.Warn().Str("tenant_id", input.TenantID.String()).Str("notification_code", code).Msg("hrms: notification master missing, skipping send")
		return result, nil
	}
	employees, err := s.ListEmployees(ctx, input.TenantID)
	if err != nil {
		return nil, err
	}
	employeeByUser := map[uuid.UUID]*domain.EmployeeListItem{}
	for _, employee := range employees {
		if employee != nil {
			employeeByUser[employee.UserID] = employee
		}
	}
	channels := normalizeNotificationChannels(input.Channels)
	for _, userID := range uniqueUUIDs(input.UserIDs) {
		if userID == uuid.Nil {
			continue
		}
		result.Processed++
		effectiveRows, err := s.ListEffectiveNotificationPreferences(ctx, input.TenantID, userID, input.ActorID)
		if err != nil {
			return nil, err
		}
		effective := effectivePreferenceForMaster(effectiveRows, master.ID)
		if effective == nil {
			effective = &domain.EffectiveNotificationPreference{Master: master, IsInAppEnabled: master.IsInAppEnabled, IsEmailEnabled: master.IsEmailEnabled, IsPushEnabled: master.IsPushEnabled, DigestFrequency: domain.DigestFrequencyInstant}
		}
		if containsString(channels, "in_app") {
			if effective.IsInAppEnabled && effective.DigestFrequency != domain.DigestFrequencyMuted {
				item := &domain.NotificationInboxItem{TenantID: input.TenantID, UserID: userID, NotificationMasterID: master.ID, Title: strings.TrimSpace(input.Title), Message: strings.TrimSpace(input.Message), ReferenceTable: input.ReferenceTable, ReferenceID: input.ReferenceID}
				if _, err := s.notifications.CreateNotificationInboxItem(ctx, item, input.ActorID); err != nil {
					return nil, err
				}
				result.InboxCount++
			} else {
				result.Suppressed++
			}
		}
		employee := employeeByUser[userID]
		deviceTokens := []*domain.DeviceToken{}
		if containsString(channels, domain.NotifChannelPush) {
			deviceTokens, err = s.notifications.ListDeviceTokensByUser(ctx, input.TenantID, userID)
			if err != nil {
				return nil, err
			}
		}
		for _, channel := range []string{domain.NotifChannelEmail, domain.NotifChannelPush} {
			if !containsString(channels, channel) {
				continue
			}
			status := domain.NotifStatusPending
			target := userID.String()
			var message *string
			if channel == domain.NotifChannelEmail {
				if !effective.IsEmailEnabled || effective.DigestFrequency == domain.DigestFrequencyMuted {
					status = domain.NotifStatusSuppressed
					msg := "email disabled by notification preferences"
					message = &msg
				}
				if employee != nil && employee.Email != nil && strings.TrimSpace(*employee.Email) != "" {
					target = strings.TrimSpace(*employee.Email)
				} else if status != domain.NotifStatusSuppressed {
					status = domain.NotifStatusSuppressed
					msg := "employee email address is missing"
					message = &msg
				}
			}
			if channel == domain.NotifChannelPush && (!effective.IsPushEnabled || effective.DigestFrequency == domain.DigestFrequencyMuted) {
				status = domain.NotifStatusSuppressed
				msg := "push disabled by notification preferences"
				message = &msg
			}
			if channel == domain.NotifChannelPush && status != domain.NotifStatusSuppressed && len(deviceTokens) == 0 {
				status = domain.NotifStatusSuppressed
				msg := "no registered push device tokens"
				message = &msg
			}
			idempotency := notificationIdempotencyKey(input.IdempotencyKey, bulkID, userID, channel)
			log := &domain.NotificationLog{TenantID: input.TenantID, NotificationMasterID: &master.ID, UserID: &userID, Channel: channel, TargetAddress: target, Subject: &input.Title, Body: &input.Message, Status: status, ErrorMessage: message, IdempotencyKey: &idempotency, BulkID: &bulkID}
			created, err := s.notifications.CreateNotificationLog(ctx, log, input.ActorID)
			if err != nil {
				return nil, err
			}
			if channel == domain.NotifChannelEmail && status != domain.NotifStatusSuppressed {
				created = s.deliverEmailNotification(ctx, created, employee, input.ActorID)
				if created.Status == domain.NotifStatusFailed {
					result.Failed++
				} else if created.Status == domain.NotifStatusSent {
					result.Sent++
				} else if created.Status == domain.NotifStatusSuppressed {
					result.Suppressed++
				}
			}
			if channel == domain.NotifChannelPush && status != domain.NotifStatusSuppressed {
				created = s.deliverPushNotification(ctx, created, deviceTokens, input.ReferenceID, input.ActorID)
				if created.Status == domain.NotifStatusFailed {
					result.Failed++
				} else if created.Status == domain.NotifStatusSent {
					result.Sent++
				} else if created.Status == domain.NotifStatusSuppressed {
					result.Suppressed++
				}
			}
			result.Logs = append(result.Logs, created)
			result.LogCount++
			if status == domain.NotifStatusSuppressed {
				result.Suppressed++
			}
		}
	}
	return result, nil
}

func (s *TenantService) deliverEmailNotification(ctx context.Context, log *domain.NotificationLog, employee *domain.EmployeeListItem, actorID *uuid.UUID) *domain.NotificationLog {
	if log == nil {
		return nil
	}
	fail := func(message string) *domain.NotificationLog {
		updated, err := s.emailProviders.UpdateNotificationLogDelivery(ctx, log.TenantID, log.ID, ports.NotificationDeliveryUpdate{Status: domain.NotifStatusFailed, ErrorMessage: &message}, actorID)
		if err != nil {
			s.logError("update failed email notification log", err, serviceTenantIDField(log.TenantID), serviceStringField("notification_log_id", log.ID.String()))
			return log
		}
		return updated
	}
	if s.emailDelivery == nil {
		msg := "email delivery sender is not configured"
		s.log.Warn().Str("tenant_id", log.TenantID.String()).Msg("hrms: email delivery skipped because sender is not configured")
		return fail(msg)
	}
	settings, err := s.activeEmailProviderSettings(ctx, log.TenantID)
	if err != nil {
		msg := "email provider settings are not configured"
		s.log.Warn().Err(err).Str("tenant_id", log.TenantID.String()).Msg("hrms: email delivery skipped because provider settings are missing")
		return fail(msg)
	}
	if !settings.IsEnabled {
		msg := "email provider is disabled"
		s.log.Warn().Str("tenant_id", log.TenantID.String()).Str("provider", settings.Provider).Msg("hrms: email delivery skipped because provider is disabled")
		updated, updateErr := s.emailProviders.UpdateNotificationLogDelivery(ctx, log.TenantID, log.ID, ports.NotificationDeliveryUpdate{Status: domain.NotifStatusSuppressed, ErrorMessage: &msg, Provider: &settings.Provider}, actorID)
		if updateErr != nil {
			s.logError("update suppressed email notification log", updateErr, serviceTenantIDField(log.TenantID), serviceStringField("notification_log_id", log.ID.String()))
			return log
		}
		return updated
	}
	subject := ""
	if log.Subject != nil {
		subject = *log.Subject
	}
	body := ""
	if log.Body != nil {
		body = *log.Body
	}
	rendered := s.renderNotificationEmail(ctx, log, employee)
	if rendered.Subject != "" {
		subject = rendered.Subject
	}
	if rendered.TextBody != "" {
		body = rendered.TextBody
	}
	result, err := s.emailDelivery.SendEmail(ctx, settings, ports.EmailMessage{TenantID: log.TenantID, To: log.TargetAddress, Subject: subject, TextBody: body, HTMLBody: rendered.HTMLBody, IdempotencyKey: log.IdempotencyKey})
	if err != nil {
		s.logError("deliver email notification", err, serviceTenantIDField(log.TenantID), serviceStringField("provider", settings.Provider), serviceStringField("notification_log_id", log.ID.String()))
		return fail(err.Error())
	}
	status := domain.NotifStatusSent
	provider := result.Provider
	externalRef := cleanStringPtr(&result.ExternalReference)
	messageID := cleanStringPtr(&result.MessageID)
	eventStatus := cleanStringPtr(&result.EventStatus)
	updated, err := s.emailProviders.UpdateNotificationLogDelivery(ctx, log.TenantID, log.ID, ports.NotificationDeliveryUpdate{Status: status, ExternalReferenceID: externalRef, Provider: &provider, ProviderMessageID: messageID, ProviderEventStatus: eventStatus}, actorID)
	if err != nil {
		s.logError("update sent email notification log", err, serviceTenantIDField(log.TenantID), serviceStringField("notification_log_id", log.ID.String()))
		return log
	}
	return updated
}

func (s *TenantService) deliverPushNotification(ctx context.Context, log *domain.NotificationLog, tokens []*domain.DeviceToken, referenceID *uuid.UUID, actorID *uuid.UUID) *domain.NotificationLog {
	if log == nil {
		return nil
	}
	fail := func(message string) *domain.NotificationLog {
		updated, err := s.pushProviders.UpdateNotificationLogDelivery(ctx, log.TenantID, log.ID, ports.NotificationDeliveryUpdate{Status: domain.NotifStatusFailed, ErrorMessage: &message}, actorID)
		if err != nil {
			s.logError("update failed push notification log", err, serviceTenantIDField(log.TenantID), serviceStringField("notification_log_id", log.ID.String()))
			return log
		}
		return updated
	}
	if s.pushDelivery == nil {
		msg := "push delivery sender is not configured"
		s.log.Warn().Str("tenant_id", log.TenantID.String()).Msg("hrms: push delivery skipped because sender is not configured")
		return fail(msg)
	}
	settings, err := s.activePushProviderSettings(ctx, log.TenantID)
	if err != nil {
		msg := "push provider settings are not configured"
		s.log.Warn().Err(err).Str("tenant_id", log.TenantID.String()).Msg("hrms: push delivery skipped because provider settings are missing")
		return fail(msg)
	}
	if !settings.IsEnabled {
		msg := "push provider is disabled"
		updated, updateErr := s.pushProviders.UpdateNotificationLogDelivery(ctx, log.TenantID, log.ID, ports.NotificationDeliveryUpdate{Status: domain.NotifStatusSuppressed, ErrorMessage: &msg, Provider: &settings.Provider}, actorID)
		if updateErr != nil {
			s.logError("update suppressed push notification log", updateErr, serviceTenantIDField(log.TenantID), serviceStringField("notification_log_id", log.ID.String()))
			return log
		}
		return updated
	}
	title, body := "", ""
	if log.Subject != nil {
		title = *log.Subject
	}
	if log.Body != nil {
		body = *log.Body
	}
	var lastResult *ports.PushDeliveryResult
	failures := []string{}
	for _, token := range tokens {
		if token == nil || strings.TrimSpace(token.DeviceToken) == "" {
			continue
		}
		result, err := s.pushDelivery.SendPush(ctx, settings, ports.PushMessage{TenantID: log.TenantID, Token: token.DeviceToken, Title: title, Body: body, IdempotencyKey: log.IdempotencyKey, ReferenceID: referenceID, Data: pushDataForLog(log, referenceID)})
		if err != nil {
			failures = append(failures, err.Error())
			continue
		}
		lastResult = result
	}
	if lastResult == nil {
		if len(failures) == 0 {
			failures = append(failures, "no valid push device tokens")
		}
		return fail(strings.Join(failures, "; "))
	}
	provider := lastResult.Provider
	externalRef := cleanStringPtr(&lastResult.ExternalReference)
	messageID := cleanStringPtr(&lastResult.MessageID)
	eventStatus := cleanStringPtr(&lastResult.EventStatus)
	updated, err := s.pushProviders.UpdateNotificationLogDelivery(ctx, log.TenantID, log.ID, ports.NotificationDeliveryUpdate{Status: domain.NotifStatusSent, ExternalReferenceID: externalRef, Provider: &provider, ProviderMessageID: messageID, ProviderEventStatus: eventStatus}, actorID)
	if err != nil {
		s.logError("update sent push notification log", err, serviceTenantIDField(log.TenantID), serviceStringField("notification_log_id", log.ID.String()))
		return log
	}
	return updated
}

func pushDataForLog(log *domain.NotificationLog, referenceID *uuid.UUID) map[string]string {
	data := map[string]string{"notification_log_id": log.ID.String(), "tenant_id": log.TenantID.String()}
	if log.NotificationMasterID != nil {
		data["notification_master_id"] = log.NotificationMasterID.String()
	}
	if referenceID != nil {
		data["reference_id"] = referenceID.String()
	}
	return data
}

func normalizeNotificationChannels(channels []string) []string {
	if len(channels) == 0 {
		return []string{"in_app", domain.NotifChannelEmail, domain.NotifChannelPush}
	}
	result := make([]string, 0, len(channels))
	for _, channel := range channels {
		clean := strings.TrimSpace(channel)
		if strings.EqualFold(clean, "in-app") || strings.EqualFold(clean, "in_app") {
			clean = "in_app"
		}
		if strings.EqualFold(clean, "email") {
			clean = domain.NotifChannelEmail
		}
		if strings.EqualFold(clean, "push") {
			clean = domain.NotifChannelPush
		}
		if clean != "" && !containsString(result, clean) {
			result = append(result, clean)
		}
	}
	return result
}

func effectivePreferenceForMaster(items []*domain.EffectiveNotificationPreference, masterID uuid.UUID) *domain.EffectiveNotificationPreference {
	for _, item := range items {
		if item != nil && item.Master != nil && item.Master.ID == masterID {
			return item
		}
	}
	return nil
}

func uniqueUUIDs(values []uuid.UUID) []uuid.UUID {
	seen := map[uuid.UUID]bool{}
	result := make([]uuid.UUID, 0, len(values))
	for _, value := range values {
		if value == uuid.Nil || seen[value] {
			continue
		}
		seen[value] = true
		result = append(result, value)
	}
	return result
}

func containsString(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func notificationIdempotencyKey(base *string, bulkID uuid.UUID, userID uuid.UUID, channel string) string {
	if base != nil && strings.TrimSpace(*base) != "" {
		return fmt.Sprintf("%s:%s:%s", strings.TrimSpace(*base), userID.String(), channel)
	}
	return fmt.Sprintf("%s:%s:%s", bulkID.String(), userID.String(), channel)
}

func normalizeListWindow(limit int32, offset int32) (int32, int32) {
	if limit <= 0 || limit > 100 {
		limit = 25
	}
	if offset < 0 {
		offset = 0
	}
	return limit, offset
}

func cleanStringPtr(value *string) *string {
	if value == nil {
		return nil
	}
	clean := strings.TrimSpace(*value)
	if clean == "" {
		return nil
	}
	return &clean
}

func (s *TenantService) UpsertDeviceToken(ctx context.Context, cmd ports.DeviceTokenCommand) (*domain.DeviceToken, error) {
	item, err := domain.NewDeviceToken(domain.DeviceTokenInput{TenantID: cmd.TenantID, UserID: cmd.UserID, DeviceToken: cmd.DeviceToken, DeviceType: cmd.DeviceType, DeviceID: cmd.DeviceID})
	if err != nil {
		s.logError("validate device token upsert", err, serviceTenantIDField(cmd.TenantID), serviceStringField("user_id", cmd.UserID.String()))
		return nil, err
	}
	result, err := s.notifications.UpsertDeviceToken(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("upsert device token", err, serviceTenantIDField(cmd.TenantID), serviceStringField("user_id", cmd.UserID.String()))
		return nil, err
	}
	if err := s.notifications.DeactivateRotatedDeviceTokens(ctx, cmd.TenantID, cmd.UserID, item.DeviceToken, cmd.ActorID); err != nil {
		s.logError("deactivate rotated device tokens", err, serviceTenantIDField(cmd.TenantID), serviceStringField("user_id", cmd.UserID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) ListDeviceTokens(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID) ([]*domain.DeviceToken, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate device token tenant", err)
		return nil, err
	}
	if userID == uuid.Nil {
		err := domain.ErrInvalidNotificationUser
		s.logError("validate device token user", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	items, err := s.notifications.ListDeviceTokensByUser(ctx, tenantID, userID)
	if err != nil {
		s.logError("list device tokens", err, serviceTenantIDField(tenantID), serviceStringField("user_id", userID.String()))
		return nil, err
	}
	return items, nil
}

func (s *TenantService) ListDeviceTokensByDeviceID(ctx context.Context, tenantID uuid.UUID, deviceID string) ([]*domain.DeviceToken, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate device token device tenant", err)
		return nil, err
	}
	clean := strings.TrimSpace(deviceID)
	if clean == "" {
		err := domain.ErrInvalidDeviceID
		s.logError("validate device token device id", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	items, err := s.notifications.ListDeviceTokensByDeviceID(ctx, tenantID, clean)
	if err != nil {
		s.logError("list device tokens by device id", err, serviceTenantIDField(tenantID), serviceStringField("device_id", clean))
		return nil, err
	}
	return items, nil
}

func (s *TenantService) DeleteDeviceToken(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate device token delete tenant", err)
		return err
	}
	if id == uuid.Nil {
		err := domain.ErrInvalidDeviceTokenID
		s.logError("validate device token delete id", err, serviceTenantIDField(tenantID))
		return err
	}
	if err := s.notifications.DeleteDeviceToken(ctx, tenantID, id, actorID); err != nil {
		s.logError("delete device token", err, serviceTenantIDField(tenantID), serviceStringField("device_token_id", id.String()))
		return err
	}
	return nil
}
