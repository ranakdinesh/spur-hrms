package postgres

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func mapNotificationMaster(row sqlc.HrmsNotificationMaster) *domain.NotificationMaster {
	return &domain.NotificationMaster{
		ID:                   row.ID,
		TenantID:             row.TenantID,
		Code:                 row.Code,
		Name:                 ptrFromText(row.Name),
		Description:          ptrFromText(row.Description),
		IsInAppEnabled:       row.IsInAppEnabled,
		IsEmailEnabled:       row.IsEmailEnabled,
		IsPushEnabled:        row.IsPushEnabled,
		EmailSubjectTemplate: ptrFromText(row.EmailSubjectTemplate),
		EmailTextTemplate:    ptrFromText(row.EmailTextTemplate),
		EmailHTMLTemplate:    ptrFromText(row.EmailHtmlTemplate),
		Inactive:             row.Inactive,
		CreatedAt:            timeFromTimestamptz(row.CreatedAt),
		CreatedBy:            ptrFromUUID(row.CreatedBy),
		UpdatedAt:            timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:            ptrFromUUID(row.UpdatedBy),
	}
}

func mapNotificationMasters(rows []sqlc.HrmsNotificationMaster) []*domain.NotificationMaster {
	items := make([]*domain.NotificationMaster, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapNotificationMaster(row))
	}
	return items
}

func mapNotificationPreference(row sqlc.HrmsNotificationPreference) *domain.NotificationPreference {
	return &domain.NotificationPreference{
		ID:                   row.ID,
		TenantID:             row.TenantID,
		UserID:               row.UserID,
		NotificationMasterID: row.NotificationMasterID,
		IsInAppEnabled:       row.IsInAppEnabled,
		IsEmailEnabled:       row.IsEmailEnabled,
		IsPushEnabled:        row.IsPushEnabled,
		DigestFrequency:      row.DigestFrequency,
		QuietHoursStart:      ptrFromOptionalClockTime(row.QuietHoursStart),
		QuietHoursEnd:        ptrFromOptionalClockTime(row.QuietHoursEnd),
		Timezone:             ptrFromText(row.Timezone),
		Inactive:             row.Inactive,
		CreatedAt:            timeFromTimestamptz(row.CreatedAt),
		CreatedBy:            ptrFromUUID(row.CreatedBy),
		UpdatedAt:            timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:            ptrFromUUID(row.UpdatedBy),
	}
}

func mapNotificationPreferences(rows []sqlc.HrmsNotificationPreference) []*domain.NotificationPreference {
	items := make([]*domain.NotificationPreference, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapNotificationPreference(row))
	}
	return items
}

func optionalClockTimeFromPtr(value *string) pgtype.Time {
	if value == nil || *value == "" {
		return pgtype.Time{Valid: false}
	}
	return timeFromClockString(*value)
}

func ptrFromOptionalClockTime(value pgtype.Time) *string {
	if !value.Valid {
		return nil
	}
	minutes := int(value.Microseconds / int64(60*1000*1000))
	formatted := fmt.Sprintf("%02d:%02d", minutes/60, minutes%60)
	return &formatted
}

func mapNotificationInbox(row sqlc.HrmsNotificationInbox) *domain.NotificationInboxItem {
	return &domain.NotificationInboxItem{
		ID:                   row.ID,
		TenantID:             row.TenantID,
		UserID:               row.UserID,
		NotificationMasterID: row.NotificationMasterID,
		Title:                row.Title,
		Message:              row.Message,
		ReferenceTable:       ptrFromText(row.ReferenceTable),
		ReferenceID:          ptrFromUUID(row.ReferenceID),
		IsRead:               row.IsRead,
		ReadDate:             ptrFromTimestamptz(row.ReadDate),
		Inactive:             row.Inactive,
		CreatedAt:            timeFromTimestamptz(row.CreatedAt),
		CreatedBy:            ptrFromUUID(row.CreatedBy),
		UpdatedAt:            timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:            ptrFromUUID(row.UpdatedBy),
	}
}

func mapNotificationInboxListRow(row sqlc.ListNotificationInboxByUserRow) *domain.NotificationInboxItem {
	notificationCode := row.NotificationCode
	return &domain.NotificationInboxItem{
		ID:                   row.ID,
		TenantID:             row.TenantID,
		UserID:               row.UserID,
		NotificationMasterID: row.NotificationMasterID,
		NotificationCode:     &notificationCode,
		NotificationName:     ptrFromText(row.NotificationName),
		Title:                row.Title,
		Message:              row.Message,
		ReferenceTable:       ptrFromText(row.ReferenceTable),
		ReferenceID:          ptrFromUUID(row.ReferenceID),
		IsRead:               row.IsRead,
		ReadDate:             ptrFromTimestamptz(row.ReadDate),
		Inactive:             row.Inactive,
		CreatedAt:            timeFromTimestamptz(row.CreatedAt),
		CreatedBy:            ptrFromUUID(row.CreatedBy),
		UpdatedAt:            timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:            ptrFromUUID(row.UpdatedBy),
	}
}

func mapNotificationInboxItems(rows []sqlc.ListNotificationInboxByUserRow) []*domain.NotificationInboxItem {
	items := make([]*domain.NotificationInboxItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapNotificationInboxListRow(row))
	}
	return items
}

func mapNotificationLog(row sqlc.HrmsNotificationLog) *domain.NotificationLog {
	return &domain.NotificationLog{
		ID:                   row.ID,
		TenantID:             row.TenantID,
		NotificationMasterID: ptrFromUUID(row.NotificationMasterID),
		UserID:               ptrFromUUID(row.UserID),
		Channel:              row.Channel,
		TargetAddress:        row.TargetAddress,
		Subject:              ptrFromText(row.Subject),
		Body:                 ptrFromText(row.Body),
		Status:               row.Status,
		SentDate:             ptrFromTimestamptz(row.SentDate),
		ErrorMessage:         ptrFromText(row.ErrorMessage),
		ExternalReferenceID:  ptrFromText(row.ExternalReferenceID),
		IdempotencyKey:       ptrFromText(row.IdempotencyKey),
		BulkID:               ptrFromUUID(row.BulkID),
		Provider:             ptrFromText(row.Provider),
		ProviderMessageID:    ptrFromText(row.ProviderMessageID),
		ProviderEventStatus:  ptrFromText(row.ProviderEventStatus),
		ProviderEventAt:      ptrFromTimestamptz(row.ProviderEventAt),
		AttemptCount:         row.AttemptCount,
		LastAttemptAt:        ptrFromTimestamptz(row.LastAttemptAt),
		Inactive:             row.Inactive,
		CreatedAt:            timeFromTimestamptz(row.CreatedAt),
		CreatedBy:            ptrFromUUID(row.CreatedBy),
		UpdatedAt:            timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:            ptrFromUUID(row.UpdatedBy),
	}
}

func mapEmailProviderSettings(row sqlc.HrmsEmailProviderSetting) *domain.EmailProviderSettings {
	item := &domain.EmailProviderSettings{
		ID: row.ID, TenantID: row.TenantID, Provider: row.Provider, IsEnabled: row.IsEnabled, FromName: ptrFromText(row.FromName), FromEmail: row.FromEmail, ReplyToEmail: ptrFromText(row.ReplyToEmail),
		SMTPHost: ptrFromText(row.SmtpHost), SMTPPort: int32PtrFromPg(row.SmtpPort), SMTPUsername: ptrFromText(row.SmtpUsername), SMTPPassword: ptrFromText(row.SmtpPassword), SMTPEncryption: row.SmtpEncryption,
		SendGridAPIKey: ptrFromText(row.SendgridApiKey), SendGridSandboxMode: row.SendgridSandboxMode, WebhookSigningSecret: ptrFromText(row.WebhookSigningSecret),
		SPFStatus: ptrFromText(row.SpfStatus), DKIMStatus: ptrFromText(row.DkimStatus), DMARCStatus: ptrFromText(row.DmarcStatus), LastTestAt: ptrFromTimestamptz(row.LastTestAt), LastTestStatus: ptrFromText(row.LastTestStatus), LastTestMessage: ptrFromText(row.LastTestMessage),
		Metadata: row.Metadata, Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy),
	}
	if len(item.Metadata) == 0 {
		item.Metadata = []byte(`{}`)
	}
	return item
}

func mapCommunicationProviderSettings(row sqlc.HrmsCommunicationProviderSetting) *domain.CommunicationProviderSettings {
	item := &domain.CommunicationProviderSettings{
		ID: row.ID, TenantID: row.TenantID, SMSProvider: row.SmsProvider, SMSEnabled: row.SmsEnabled, SMSSenderID: ptrFromText(row.SmsSenderID), SMSAuthKey: ptrFromText(row.SmsAuthKey), SMSTemplateID: ptrFromText(row.SmsTemplateID), SMSRoute: ptrFromText(row.SmsRoute), SMSCountryCode: ptrFromText(row.SmsCountryCode), SMSBaseURL: ptrFromText(row.SmsBaseUrl),
		WhatsAppProvider: row.WhatsappProvider, WhatsAppEnabled: row.WhatsappEnabled, WhatsAppAuthKey: ptrFromText(row.WhatsappAuthKey), WhatsAppAppName: ptrFromText(row.WhatsappAppName), WhatsAppSourceNumber: ptrFromText(row.WhatsappSourceNumber), WhatsAppTemplateID: ptrFromText(row.WhatsappTemplateID), WhatsAppTemplateName: ptrFromText(row.WhatsappTemplateName), WhatsAppNamespace: ptrFromText(row.WhatsappNamespace), WhatsAppBaseURL: ptrFromText(row.WhatsappBaseUrl), WebhookSigningSecret: ptrFromText(row.WebhookSigningSecret),
		SMSLastTestAt: ptrFromTimestamptz(row.SmsLastTestAt), SMSLastTestStatus: ptrFromText(row.SmsLastTestStatus), SMSLastTestMessage: ptrFromText(row.SmsLastTestMessage), WhatsAppLastTestAt: ptrFromTimestamptz(row.WhatsappLastTestAt), WhatsAppLastTestStatus: ptrFromText(row.WhatsappLastTestStatus), WhatsAppLastTestMessage: ptrFromText(row.WhatsappLastTestMessage),
		Metadata: row.Metadata, Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy),
	}
	if len(item.Metadata) == 0 {
		item.Metadata = []byte(`{}`)
	}
	return item
}

func mapStorageProviderSettings(row sqlc.HrmsStorageProviderSetting) *domain.StorageProviderSettings {
	item := &domain.StorageProviderSettings{
		ID: row.ID, TenantID: row.TenantID, Provider: row.Provider, IsEnabled: row.IsEnabled, Bucket: row.Bucket, Region: ptrFromText(row.Region), Endpoint: ptrFromText(row.Endpoint), AccessKeyID: ptrFromText(row.AccessKeyID), SecretAccessKey: ptrFromText(row.SecretAccessKey), UseSSL: row.UseSsl, ForcePathStyle: row.ForcePathStyle, ObjectPrefix: ptrFromText(row.ObjectPrefix), PublicBaseURL: ptrFromText(row.PublicBaseUrl), MaxFileSizeBytes: row.MaxFileSizeBytes, AllowedContentTypes: ptrFromText(row.AllowedContentTypes),
		LastTestAt: ptrFromTimestamptz(row.LastTestAt), LastTestStatus: ptrFromText(row.LastTestStatus), LastTestMessage: ptrFromText(row.LastTestMessage), Metadata: row.Metadata, Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy),
	}
	if len(item.Metadata) == 0 {
		item.Metadata = []byte(`{}`)
	}
	return item
}

func mapPushProviderSettings(row sqlc.HrmsPushProviderSetting) *domain.PushProviderSettings {
	item := &domain.PushProviderSettings{
		ID: row.ID, TenantID: row.TenantID, Provider: row.Provider, IsEnabled: row.IsEnabled, ProjectID: ptrFromText(row.ProjectID), ClientEmail: ptrFromText(row.ClientEmail), PrivateKey: ptrFromText(row.PrivateKey), PrivateKeyID: ptrFromText(row.PrivateKeyID), AuthURI: ptrFromText(row.AuthUri), TokenURI: ptrFromText(row.TokenUri),
		AndroidEnabled: row.AndroidEnabled, IOSEnabled: row.IosEnabled, WebEnabled: row.WebEnabled, DefaultClickAction: ptrFromText(row.DefaultClickAction), DefaultImageURL: ptrFromText(row.DefaultImageUrl), TTLSeconds: row.TtlSeconds, CollapseKey: ptrFromText(row.CollapseKey),
		LastTestAt: ptrFromTimestamptz(row.LastTestAt), LastTestStatus: ptrFromText(row.LastTestStatus), LastTestMessage: ptrFromText(row.LastTestMessage), Metadata: row.Metadata, Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy),
	}
	if len(item.Metadata) == 0 {
		item.Metadata = []byte(`{}`)
	}
	return item
}

func int32PtrFromPg(value pgtype.Int4) *int32 {
	if !value.Valid {
		return nil
	}
	return &value.Int32
}

func mapNotificationLogs(rows []sqlc.HrmsNotificationLog) []*domain.NotificationLog {
	items := make([]*domain.NotificationLog, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapNotificationLog(row))
	}
	return items
}

func mapDeviceToken(row sqlc.HrmsDeviceToken) *domain.DeviceToken {
	return &domain.DeviceToken{
		ID:          row.ID,
		TenantID:    row.TenantID,
		UserID:      row.UserID,
		DeviceToken: row.DeviceToken,
		DeviceType:  ptrFromText(row.DeviceType),
		DeviceID:    ptrFromText(row.DeviceID),
		Inactive:    row.Inactive,
		CreatedAt:   timeFromTimestamptz(row.CreatedAt),
		CreatedBy:   ptrFromUUID(row.CreatedBy),
		UpdatedAt:   timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:   ptrFromUUID(row.UpdatedBy),
	}
}

func mapDeviceTokens(rows []sqlc.HrmsDeviceToken) []*domain.DeviceToken {
	items := make([]*domain.DeviceToken, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapDeviceToken(row))
	}
	return items
}
