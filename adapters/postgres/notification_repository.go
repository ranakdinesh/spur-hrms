package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *Store) CreateNotificationMaster(ctx context.Context, item *domain.NotificationMaster, actorID *uuid.UUID) (*domain.NotificationMaster, error) {
	row, err := s.getQueries(ctx).CreateNotificationMaster(ctx, sqlc.CreateNotificationMasterParams{TenantID: item.TenantID, Code: item.Code, Name: textFromPtr(item.Name), Description: textFromPtr(item.Description), IsInAppEnabled: item.IsInAppEnabled, IsEmailEnabled: item.IsEmailEnabled, IsPushEnabled: item.IsPushEnabled, EmailSubjectTemplate: textFromPtr(item.EmailSubjectTemplate), EmailTextTemplate: textFromPtr(item.EmailTextTemplate), EmailHtmlTemplate: textFromPtr(item.EmailHTMLTemplate), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create notification master", fmt.Errorf("hrms: create notification master: %w", err), tenantIDField(item.TenantID), stringField("notification_code", item.Code))
	}
	return mapNotificationMaster(row), nil
}

func (s *Store) ListNotificationMasters(ctx context.Context, tenantID uuid.UUID) ([]*domain.NotificationMaster, error) {
	rows, err := s.getQueries(ctx).ListNotificationMasters(ctx, tenantID)
	if err != nil {
		return nil, s.logDBError(ctx, "list notification masters", err, tenantIDField(tenantID))
	}
	return mapNotificationMasters(rows), nil
}

func (s *Store) GetNotificationMaster(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.NotificationMaster, error) {
	row, err := s.getQueries(ctx).GetNotificationMaster(ctx, sqlc.GetNotificationMasterParams{TenantID: tenantID, ID: id})
	if err != nil {
		return nil, s.logDBError(ctx, "get notification master", fmt.Errorf("hrms: get notification master: %w", err), tenantIDField(tenantID), stringField("notification_master_id", id.String()))
	}
	return mapNotificationMaster(row), nil
}

func (s *Store) GetNotificationMasterByCode(ctx context.Context, tenantID uuid.UUID, code string) (*domain.NotificationMaster, error) {
	row, err := s.getQueries(ctx).GetNotificationMasterByCode(ctx, sqlc.GetNotificationMasterByCodeParams{TenantID: tenantID, Code: code})
	if err != nil {
		return nil, s.logDBError(ctx, "get notification master by code", fmt.Errorf("hrms: get notification master by code: %w", err), tenantIDField(tenantID), stringField("notification_code", code))
	}
	return mapNotificationMaster(row), nil
}

func (s *Store) UpdateNotificationMaster(ctx context.Context, item *domain.NotificationMaster, actorID *uuid.UUID) (*domain.NotificationMaster, error) {
	row, err := s.getQueries(ctx).UpdateNotificationMaster(ctx, sqlc.UpdateNotificationMasterParams{TenantID: item.TenantID, ID: item.ID, Code: item.Code, Name: textFromPtr(item.Name), Description: textFromPtr(item.Description), IsInAppEnabled: item.IsInAppEnabled, IsEmailEnabled: item.IsEmailEnabled, IsPushEnabled: item.IsPushEnabled, EmailSubjectTemplate: textFromPtr(item.EmailSubjectTemplate), EmailTextTemplate: textFromPtr(item.EmailTextTemplate), EmailHtmlTemplate: textFromPtr(item.EmailHTMLTemplate), UpdatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "update notification master", fmt.Errorf("hrms: update notification master: %w", err), tenantIDField(item.TenantID), stringField("notification_master_id", item.ID.String()))
	}
	return mapNotificationMaster(row), nil
}

func (s *Store) DeleteNotificationMaster(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteNotificationMaster(ctx, sqlc.SoftDeleteNotificationMasterParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete notification master", fmt.Errorf("hrms: delete notification master: %w", err), tenantIDField(tenantID), stringField("notification_master_id", id.String()))
	}
	return nil
}

func (s *Store) ListNotificationPreferencesByUser(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID) ([]*domain.NotificationPreference, error) {
	rows, err := s.getQueries(ctx).ListNotificationPreferencesByUser(ctx, sqlc.ListNotificationPreferencesByUserParams{TenantID: tenantID, UserID: userID})
	if err != nil {
		return nil, s.logDBError(ctx, "list notification preferences", err, tenantIDField(tenantID), stringField("user_id", userID.String()))
	}
	return mapNotificationPreferences(rows), nil
}

func (s *Store) UpsertNotificationPreference(ctx context.Context, item *domain.NotificationPreference, actorID *uuid.UUID) (*domain.NotificationPreference, error) {
	row, err := s.getQueries(ctx).UpsertNotificationPreference(ctx, sqlc.UpsertNotificationPreferenceParams{TenantID: item.TenantID, UserID: item.UserID, NotificationMasterID: item.NotificationMasterID, IsInAppEnabled: item.IsInAppEnabled, IsEmailEnabled: item.IsEmailEnabled, IsPushEnabled: item.IsPushEnabled, DigestFrequency: item.DigestFrequency, QuietHoursStart: optionalClockTimeFromPtr(item.QuietHoursStart), QuietHoursEnd: optionalClockTimeFromPtr(item.QuietHoursEnd), Timezone: textFromPtr(item.Timezone), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "upsert notification preference", fmt.Errorf("hrms: upsert notification preference: %w", err), tenantIDField(item.TenantID), stringField("user_id", item.UserID.String()), stringField("notification_master_id", item.NotificationMasterID.String()))
	}
	return mapNotificationPreference(row), nil
}

func (s *Store) DeleteNotificationPreference(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteNotificationPreference(ctx, sqlc.SoftDeleteNotificationPreferenceParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete notification preference", fmt.Errorf("hrms: delete notification preference: %w", err), tenantIDField(tenantID), stringField("notification_preference_id", id.String()))
	}
	return nil
}

func (s *Store) CreateNotificationInboxItem(ctx context.Context, item *domain.NotificationInboxItem, actorID *uuid.UUID) (*domain.NotificationInboxItem, error) {
	row, err := s.getQueries(ctx).CreateNotificationInboxItem(ctx, sqlc.CreateNotificationInboxItemParams{TenantID: item.TenantID, UserID: item.UserID, NotificationMasterID: item.NotificationMasterID, Title: item.Title, Message: item.Message, ReferenceTable: textFromPtr(item.ReferenceTable), ReferenceID: uuidFromPtr(item.ReferenceID), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create notification inbox item", fmt.Errorf("hrms: create notification inbox item: %w", err), tenantIDField(item.TenantID), stringField("user_id", item.UserID.String()))
	}
	return mapNotificationInbox(row), nil
}

func (s *Store) ListNotificationInboxByUser(ctx context.Context, filter domain.NotificationInboxFilter) ([]*domain.NotificationInboxItem, error) {
	rows, err := s.getQueries(ctx).ListNotificationInboxByUser(ctx, sqlc.ListNotificationInboxByUserParams{TenantID: filter.TenantID, UserID: filter.UserID, IsRead: boolFromPtr(filter.IsRead), NotificationMasterID: uuidFromPtr(filter.NotificationMasterID), NotificationCode: textFromPtr(filter.NotificationCode), Search: textFromPtr(filter.Search), Offset: filter.Offset, Limit: filter.Limit})
	if err != nil {
		return nil, s.logDBError(ctx, "list notification inbox", err, tenantIDField(filter.TenantID), stringField("user_id", filter.UserID.String()))
	}
	return mapNotificationInboxItems(rows), nil
}

func (s *Store) CountNotificationInboxByUser(ctx context.Context, filter domain.NotificationInboxFilter) (int64, error) {
	count, err := s.getQueries(ctx).CountNotificationInboxByUser(ctx, sqlc.CountNotificationInboxByUserParams{TenantID: filter.TenantID, UserID: filter.UserID, IsRead: boolFromPtr(filter.IsRead), NotificationMasterID: uuidFromPtr(filter.NotificationMasterID), NotificationCode: textFromPtr(filter.NotificationCode), Search: textFromPtr(filter.Search)})
	if err != nil {
		return 0, s.logDBError(ctx, "count notification inbox", err, tenantIDField(filter.TenantID), stringField("user_id", filter.UserID.String()))
	}
	return count, nil
}

func (s *Store) CountUnreadNotificationsByUser(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID) (int64, error) {
	count, err := s.getQueries(ctx).CountUnreadNotificationsByUser(ctx, sqlc.CountUnreadNotificationsByUserParams{TenantID: tenantID, UserID: userID})
	if err != nil {
		return 0, s.logDBError(ctx, "count unread notifications", err, tenantIDField(tenantID), stringField("user_id", userID.String()))
	}
	return count, nil
}

func (s *Store) MarkNotificationInboxItemRead(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).MarkNotificationInboxItemRead(ctx, sqlc.MarkNotificationInboxItemReadParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "mark notification read", err, tenantIDField(tenantID), stringField("notification_id", id.String()))
	}
	return nil
}

func (s *Store) MarkNotificationInboxItemUnread(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).MarkNotificationInboxItemUnread(ctx, sqlc.MarkNotificationInboxItemUnreadParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "mark notification unread", err, tenantIDField(tenantID), stringField("notification_id", id.String()))
	}
	return nil
}

func (s *Store) MarkNotificationInboxAllRead(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).MarkNotificationInboxAllRead(ctx, sqlc.MarkNotificationInboxAllReadParams{TenantID: tenantID, UserID: userID, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "mark all notifications read", err, tenantIDField(tenantID), stringField("user_id", userID.String()))
	}
	return nil
}

func (s *Store) DeleteNotificationInboxItem(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteNotificationInboxItem(ctx, sqlc.SoftDeleteNotificationInboxItemParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "archive notification inbox item", err, tenantIDField(tenantID), stringField("notification_id", id.String()))
	}
	return nil
}

func (s *Store) CreateNotificationLog(ctx context.Context, item *domain.NotificationLog, actorID *uuid.UUID) (*domain.NotificationLog, error) {
	row, err := s.getQueries(ctx).CreateNotificationLog(ctx, sqlc.CreateNotificationLogParams{TenantID: item.TenantID, NotificationMasterID: uuidFromPtr(item.NotificationMasterID), UserID: uuidFromPtr(item.UserID), Channel: item.Channel, TargetAddress: item.TargetAddress, Subject: textFromPtr(item.Subject), Body: textFromPtr(item.Body), Status: item.Status, ErrorMessage: textFromPtr(item.ErrorMessage), ExternalReferenceID: textFromPtr(item.ExternalReferenceID), IdempotencyKey: textFromPtr(item.IdempotencyKey), BulkID: uuidFromPtr(item.BulkID), CreatedBy: uuidFromPtr(actorID), SentDateStatus: item.Status})
	if err != nil {
		return nil, s.logDBError(ctx, "create notification log", fmt.Errorf("hrms: create notification log: %w", err), tenantIDField(item.TenantID), stringField("channel", item.Channel), stringField("status", item.Status))
	}
	return mapNotificationLog(row), nil
}

func (s *Store) ListNotificationLogs(ctx context.Context, tenantID uuid.UUID, userID *uuid.UUID, status *string, channel *string, bulkID *uuid.UUID, limit int32, offset int32) ([]*domain.NotificationLog, error) {
	rows, err := s.getQueries(ctx).ListNotificationLogs(ctx, sqlc.ListNotificationLogsParams{TenantID: tenantID, UserID: uuidFromPtr(userID), Status: textFromPtr(status), Channel: textFromPtr(channel), BulkID: uuidFromPtr(bulkID), Limit: limit, Offset: offset})
	if err != nil {
		return nil, s.logDBError(ctx, "list notification logs", err, tenantIDField(tenantID))
	}
	return mapNotificationLogs(rows), nil
}

func (s *Store) GetEmailProviderSettings(ctx context.Context, tenantID uuid.UUID) (*domain.EmailProviderSettings, error) {
	row, err := s.getQueries(ctx).GetEmailProviderSettings(ctx, tenantID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrEmailProviderSettingsNotFound
		}
		return nil, s.logDBError(ctx, "get email provider settings", err, tenantIDField(tenantID))
	}
	return mapEmailProviderSettings(row), nil
}

func (s *Store) UpsertEmailProviderSettings(ctx context.Context, item *domain.EmailProviderSettings, actorID *uuid.UUID) (*domain.EmailProviderSettings, error) {
	row, err := s.getQueries(ctx).UpsertEmailProviderSettings(ctx, sqlc.UpsertEmailProviderSettingsParams{
		TenantID: item.TenantID, Provider: item.Provider, IsEnabled: item.IsEnabled, FromName: textFromPtr(item.FromName), FromEmail: item.FromEmail, ReplyToEmail: textFromPtr(item.ReplyToEmail),
		SmtpHost: textFromPtr(item.SMTPHost), SmtpPort: int4FromPtr(item.SMTPPort), SmtpUsername: textFromPtr(item.SMTPUsername), SmtpPassword: textFromPtr(item.SMTPPassword), SmtpEncryption: item.SMTPEncryption,
		SendgridApiKey: textFromPtr(item.SendGridAPIKey), SendgridSandboxMode: item.SendGridSandboxMode, WebhookSigningSecret: textFromPtr(item.WebhookSigningSecret),
		SpfStatus: textFromPtr(item.SPFStatus), DkimStatus: textFromPtr(item.DKIMStatus), DmarcStatus: textFromPtr(item.DMARCStatus), Column18: []byte(item.Metadata), CreatedBy: uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "upsert email provider settings", err, tenantIDField(item.TenantID), stringField("provider", item.Provider))
	}
	return mapEmailProviderSettings(row), nil
}

func (s *Store) UpdateEmailProviderTestResult(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, message *string, actorID *uuid.UUID) (*domain.EmailProviderSettings, error) {
	row, err := s.getQueries(ctx).UpdateEmailProviderTestResult(ctx, sqlc.UpdateEmailProviderTestResultParams{TenantID: tenantID, ID: id, LastTestStatus: textFromPtr(&status), LastTestMessage: textFromPtr(message), UpdatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "update email provider test result", err, tenantIDField(tenantID), stringField("email_provider_settings_id", id.String()))
	}
	return mapEmailProviderSettings(row), nil
}

func (s *Store) DeleteEmailProviderSettings(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteEmailProviderSettings(ctx, sqlc.SoftDeleteEmailProviderSettingsParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete email provider settings", err, tenantIDField(tenantID), stringField("email_provider_settings_id", id.String()))
	}
	return nil
}

func (s *Store) GetCommunicationProviderSettings(ctx context.Context, tenantID uuid.UUID) (*domain.CommunicationProviderSettings, error) {
	row, err := s.getQueries(ctx).GetCommunicationProviderSettings(ctx, tenantID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrCommunicationProviderSettingsNotFound
		}
		return nil, s.logDBError(ctx, "get communication provider settings", err, tenantIDField(tenantID))
	}
	return mapCommunicationProviderSettings(row), nil
}

func (s *Store) UpsertCommunicationProviderSettings(ctx context.Context, item *domain.CommunicationProviderSettings, actorID *uuid.UUID) (*domain.CommunicationProviderSettings, error) {
	row, err := s.getQueries(ctx).UpsertCommunicationProviderSettings(ctx, sqlc.UpsertCommunicationProviderSettingsParams{
		TenantID: item.TenantID, SmsProvider: item.SMSProvider, SmsEnabled: item.SMSEnabled, SmsSenderID: textFromPtr(item.SMSSenderID), SmsAuthKey: textFromPtr(item.SMSAuthKey), SmsTemplateID: textFromPtr(item.SMSTemplateID), SmsRoute: textFromPtr(item.SMSRoute), SmsCountryCode: textFromPtr(item.SMSCountryCode), SmsBaseUrl: textFromPtr(item.SMSBaseURL),
		WhatsappProvider: item.WhatsAppProvider, WhatsappEnabled: item.WhatsAppEnabled, WhatsappAuthKey: textFromPtr(item.WhatsAppAuthKey), WhatsappAppName: textFromPtr(item.WhatsAppAppName), WhatsappSourceNumber: textFromPtr(item.WhatsAppSourceNumber), WhatsappTemplateID: textFromPtr(item.WhatsAppTemplateID), WhatsappTemplateName: textFromPtr(item.WhatsAppTemplateName), WhatsappNamespace: textFromPtr(item.WhatsAppNamespace), WhatsappBaseUrl: textFromPtr(item.WhatsAppBaseURL),
		WebhookSigningSecret: textFromPtr(item.WebhookSigningSecret), Column20: []byte(item.Metadata), CreatedBy: uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "upsert communication provider settings", err, tenantIDField(item.TenantID), stringField("sms_provider", item.SMSProvider), stringField("whatsapp_provider", item.WhatsAppProvider))
	}
	return mapCommunicationProviderSettings(row), nil
}

func (s *Store) UpdateCommunicationProviderTestResult(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, channel string, status string, message *string, actorID *uuid.UUID) (*domain.CommunicationProviderSettings, error) {
	row, err := s.getQueries(ctx).UpdateCommunicationProviderTestResult(ctx, sqlc.UpdateCommunicationProviderTestResultParams{TenantID: tenantID, ID: id, Column3: channel, SmsLastTestStatus: textFromPtr(&status), SmsLastTestMessage: textFromPtr(message), UpdatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "update communication provider test result", err, tenantIDField(tenantID), stringField("communication_provider_settings_id", id.String()), stringField("channel", channel))
	}
	return mapCommunicationProviderSettings(row), nil
}

func (s *Store) DeleteCommunicationProviderSettings(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteCommunicationProviderSettings(ctx, sqlc.SoftDeleteCommunicationProviderSettingsParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete communication provider settings", err, tenantIDField(tenantID), stringField("communication_provider_settings_id", id.String()))
	}
	return nil
}

func (s *Store) GetStorageProviderSettings(ctx context.Context, tenantID uuid.UUID) (*domain.StorageProviderSettings, error) {
	row, err := s.getQueries(ctx).GetStorageProviderSettings(ctx, tenantID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrStorageProviderSettingsNotFound
		}
		return nil, s.logDBError(ctx, "get storage provider settings", err, tenantIDField(tenantID))
	}
	return mapStorageProviderSettings(row), nil
}

func (s *Store) UpsertStorageProviderSettings(ctx context.Context, item *domain.StorageProviderSettings, actorID *uuid.UUID) (*domain.StorageProviderSettings, error) {
	row, err := s.getQueries(ctx).UpsertStorageProviderSettings(ctx, sqlc.UpsertStorageProviderSettingsParams{
		TenantID: item.TenantID, Provider: item.Provider, IsEnabled: item.IsEnabled, Bucket: item.Bucket, Region: textFromPtr(item.Region), Endpoint: textFromPtr(item.Endpoint), AccessKeyID: textFromPtr(item.AccessKeyID), SecretAccessKey: textFromPtr(item.SecretAccessKey), UseSsl: item.UseSSL, ForcePathStyle: item.ForcePathStyle, ObjectPrefix: textFromPtr(item.ObjectPrefix), PublicBaseUrl: textFromPtr(item.PublicBaseURL), MaxFileSizeBytes: item.MaxFileSizeBytes, AllowedContentTypes: textFromPtr(item.AllowedContentTypes), Column15: []byte(item.Metadata), CreatedBy: uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "upsert storage provider settings", err, tenantIDField(item.TenantID), stringField("provider", item.Provider), stringField("bucket", item.Bucket))
	}
	return mapStorageProviderSettings(row), nil
}

func (s *Store) UpdateStorageProviderTestResult(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, message *string, actorID *uuid.UUID) (*domain.StorageProviderSettings, error) {
	row, err := s.getQueries(ctx).UpdateStorageProviderTestResult(ctx, sqlc.UpdateStorageProviderTestResultParams{TenantID: tenantID, ID: id, LastTestStatus: textFromPtr(&status), LastTestMessage: textFromPtr(message), UpdatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "update storage provider test result", err, tenantIDField(tenantID), stringField("storage_provider_settings_id", id.String()))
	}
	return mapStorageProviderSettings(row), nil
}

func (s *Store) DeleteStorageProviderSettings(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteStorageProviderSettings(ctx, sqlc.SoftDeleteStorageProviderSettingsParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete storage provider settings", err, tenantIDField(tenantID), stringField("storage_provider_settings_id", id.String()))
	}
	return nil
}

func (s *Store) GetPushProviderSettings(ctx context.Context, tenantID uuid.UUID) (*domain.PushProviderSettings, error) {
	row, err := s.getQueries(ctx).GetPushProviderSettings(ctx, tenantID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrPushProviderSettingsNotFound
		}
		return nil, s.logDBError(ctx, "get push provider settings", err, tenantIDField(tenantID))
	}
	return mapPushProviderSettings(row), nil
}

func (s *Store) UpsertPushProviderSettings(ctx context.Context, item *domain.PushProviderSettings, actorID *uuid.UUID) (*domain.PushProviderSettings, error) {
	row, err := s.getQueries(ctx).UpsertPushProviderSettings(ctx, sqlc.UpsertPushProviderSettingsParams{
		TenantID: item.TenantID, Provider: item.Provider, IsEnabled: item.IsEnabled, ProjectID: textFromPtr(item.ProjectID), ClientEmail: textFromPtr(item.ClientEmail), PrivateKey: textFromPtr(item.PrivateKey), PrivateKeyID: textFromPtr(item.PrivateKeyID), AuthUri: textFromPtr(item.AuthURI), TokenUri: textFromPtr(item.TokenURI),
		AndroidEnabled: item.AndroidEnabled, IosEnabled: item.IOSEnabled, WebEnabled: item.WebEnabled, DefaultClickAction: textFromPtr(item.DefaultClickAction), DefaultImageUrl: textFromPtr(item.DefaultImageURL), TtlSeconds: item.TTLSeconds, CollapseKey: textFromPtr(item.CollapseKey), Column17: []byte(item.Metadata), CreatedBy: uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "upsert push provider settings", err, tenantIDField(item.TenantID), stringField("provider", item.Provider))
	}
	return mapPushProviderSettings(row), nil
}

func (s *Store) UpdatePushProviderTestResult(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, message *string, actorID *uuid.UUID) (*domain.PushProviderSettings, error) {
	row, err := s.getQueries(ctx).UpdatePushProviderTestResult(ctx, sqlc.UpdatePushProviderTestResultParams{TenantID: tenantID, ID: id, LastTestStatus: textFromPtr(&status), LastTestMessage: textFromPtr(message), UpdatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "update push provider test result", err, tenantIDField(tenantID), stringField("push_provider_settings_id", id.String()))
	}
	return mapPushProviderSettings(row), nil
}

func (s *Store) DeletePushProviderSettings(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeletePushProviderSettings(ctx, sqlc.SoftDeletePushProviderSettingsParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete push provider settings", err, tenantIDField(tenantID), stringField("push_provider_settings_id", id.String()))
	}
	return nil
}

func (s *Store) UpdateNotificationLogDelivery(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, update ports.NotificationDeliveryUpdate, actorID *uuid.UUID) (*domain.NotificationLog, error) {
	row, err := s.getQueries(ctx).UpdateNotificationLogDelivery(ctx, sqlc.UpdateNotificationLogDeliveryParams{TenantID: tenantID, ID: id, Status: update.Status, ErrorMessage: textFromPtr(update.ErrorMessage), ExternalReferenceID: textFromPtr(update.ExternalReferenceID), Provider: textFromPtr(update.Provider), ProviderMessageID: textFromPtr(update.ProviderMessageID), ProviderEventStatus: textFromPtr(update.ProviderEventStatus), ProviderEventAt: timestamptzFromPtr(update.ProviderEventAt), UpdatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "update notification log delivery", err, tenantIDField(tenantID), stringField("notification_log_id", id.String()))
	}
	return mapNotificationLog(row), nil
}

func boolFromPtr(value *bool) pgtype.Bool {
	if value == nil {
		return pgtype.Bool{Valid: false}
	}
	return pgtype.Bool{Bool: *value, Valid: true}
}

func (s *Store) UpsertDeviceToken(ctx context.Context, item *domain.DeviceToken, actorID *uuid.UUID) (*domain.DeviceToken, error) {
	row, err := s.getQueries(ctx).UpsertDeviceToken(ctx, sqlc.UpsertDeviceTokenParams{TenantID: item.TenantID, UserID: item.UserID, DeviceToken: item.DeviceToken, DeviceType: textFromPtr(item.DeviceType), DeviceID: textFromPtr(item.DeviceID), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "upsert device token", fmt.Errorf("hrms: upsert device token: %w", err), tenantIDField(item.TenantID), stringField("user_id", item.UserID.String()))
	}
	return mapDeviceToken(row), nil
}

func (s *Store) DeactivateRotatedDeviceTokens(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID, activeToken string, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).DeactivateRotatedDeviceTokens(ctx, sqlc.DeactivateRotatedDeviceTokensParams{TenantID: tenantID, UserID: userID, DeviceToken: activeToken, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "deactivate rotated device tokens", fmt.Errorf("hrms: deactivate rotated device tokens: %w", err), tenantIDField(tenantID), stringField("user_id", userID.String()))
	}
	return nil
}

func (s *Store) ListDeviceTokensByUser(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID) ([]*domain.DeviceToken, error) {
	rows, err := s.getQueries(ctx).ListDeviceTokensByUser(ctx, sqlc.ListDeviceTokensByUserParams{TenantID: tenantID, UserID: userID})
	if err != nil {
		return nil, s.logDBError(ctx, "list device tokens by user", err, tenantIDField(tenantID), stringField("user_id", userID.String()))
	}
	return mapDeviceTokens(rows), nil
}

func (s *Store) ListDeviceTokensByDeviceID(ctx context.Context, tenantID uuid.UUID, deviceID string) ([]*domain.DeviceToken, error) {
	rows, err := s.getQueries(ctx).ListDeviceTokensByDeviceID(ctx, sqlc.ListDeviceTokensByDeviceIDParams{TenantID: tenantID, DeviceID: textFromString(deviceID)})
	if err != nil {
		return nil, s.logDBError(ctx, "list device tokens by device id", err, tenantIDField(tenantID), stringField("device_id", deviceID))
	}
	return mapDeviceTokens(rows), nil
}

func (s *Store) DeleteDeviceToken(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteDeviceToken(ctx, sqlc.SoftDeleteDeviceTokenParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete device token", fmt.Errorf("hrms: delete device token: %w", err), tenantIDField(tenantID), stringField("device_token_id", id.String()))
	}
	return nil
}
