package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

type NotificationRepo interface {
	CreateNotificationMaster(ctx context.Context, item *domain.NotificationMaster, actorID *uuid.UUID) (*domain.NotificationMaster, error)
	ListNotificationMasters(ctx context.Context, tenantID uuid.UUID) ([]*domain.NotificationMaster, error)
	GetNotificationMaster(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.NotificationMaster, error)
	GetNotificationMasterByCode(ctx context.Context, tenantID uuid.UUID, code string) (*domain.NotificationMaster, error)
	UpdateNotificationMaster(ctx context.Context, item *domain.NotificationMaster, actorID *uuid.UUID) (*domain.NotificationMaster, error)
	DeleteNotificationMaster(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	ListNotificationPreferencesByUser(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID) ([]*domain.NotificationPreference, error)
	UpsertNotificationPreference(ctx context.Context, item *domain.NotificationPreference, actorID *uuid.UUID) (*domain.NotificationPreference, error)
	DeleteNotificationPreference(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	CreateNotificationInboxItem(ctx context.Context, item *domain.NotificationInboxItem, actorID *uuid.UUID) (*domain.NotificationInboxItem, error)
	ListNotificationInboxByUser(ctx context.Context, filter domain.NotificationInboxFilter) ([]*domain.NotificationInboxItem, error)
	CountNotificationInboxByUser(ctx context.Context, filter domain.NotificationInboxFilter) (int64, error)
	CountUnreadNotificationsByUser(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID) (int64, error)
	MarkNotificationInboxItemRead(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	MarkNotificationInboxItemUnread(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	MarkNotificationInboxAllRead(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID, actorID *uuid.UUID) error
	DeleteNotificationInboxItem(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	CreateNotificationLog(ctx context.Context, item *domain.NotificationLog, actorID *uuid.UUID) (*domain.NotificationLog, error)
	ListNotificationLogs(ctx context.Context, tenantID uuid.UUID, userID *uuid.UUID, status *string, channel *string, bulkID *uuid.UUID, limit int32, offset int32) ([]*domain.NotificationLog, error)
	UpsertDeviceToken(ctx context.Context, item *domain.DeviceToken, actorID *uuid.UUID) (*domain.DeviceToken, error)
	DeactivateRotatedDeviceTokens(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID, activeToken string, actorID *uuid.UUID) error
	ListDeviceTokensByUser(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID) ([]*domain.DeviceToken, error)
	ListDeviceTokensByDeviceID(ctx context.Context, tenantID uuid.UUID, deviceID string) ([]*domain.DeviceToken, error)
	DeleteDeviceToken(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
}

type NotificationMasterCommand struct {
	ID                   uuid.UUID  `json:"id,omitempty"`
	TenantID             uuid.UUID  `json:"tenant_id"`
	Code                 string     `json:"code"`
	Name                 *string    `json:"name,omitempty"`
	Description          *string    `json:"description,omitempty"`
	IsInAppEnabled       bool       `json:"is_in_app_enabled"`
	IsEmailEnabled       bool       `json:"is_email_enabled"`
	IsPushEnabled        bool       `json:"is_push_enabled"`
	EmailSubjectTemplate *string    `json:"email_subject_template,omitempty"`
	EmailTextTemplate    *string    `json:"email_text_template,omitempty"`
	EmailHTMLTemplate    *string    `json:"email_html_template,omitempty"`
	ActorID              *uuid.UUID `json:"-"`
}

type NotificationPreferenceCommand struct {
	ID                   uuid.UUID  `json:"id,omitempty"`
	TenantID             uuid.UUID  `json:"tenant_id"`
	UserID               uuid.UUID  `json:"user_id"`
	NotificationMasterID uuid.UUID  `json:"notification_master_id"`
	IsInAppEnabled       bool       `json:"is_in_app_enabled"`
	IsEmailEnabled       bool       `json:"is_email_enabled"`
	IsPushEnabled        bool       `json:"is_push_enabled"`
	DigestFrequency      string     `json:"digest_frequency"`
	QuietHoursStart      *string    `json:"quiet_hours_start,omitempty"`
	QuietHoursEnd        *string    `json:"quiet_hours_end,omitempty"`
	Timezone             *string    `json:"timezone,omitempty"`
	ActorID              *uuid.UUID `json:"-"`
}

type SendNotificationCommand struct {
	TenantID         uuid.UUID   `json:"tenant_id"`
	NotificationCode string      `json:"notification_code"`
	UserIDs          []uuid.UUID `json:"user_ids"`
	Title            string      `json:"title"`
	Message          string      `json:"message"`
	ReferenceTable   *string     `json:"reference_table,omitempty"`
	ReferenceID      *uuid.UUID  `json:"reference_id,omitempty"`
	Channels         []string    `json:"channels,omitempty"`
	IdempotencyKey   *string     `json:"idempotency_key,omitempty"`
	ActorID          *uuid.UUID  `json:"-"`
}

type DeviceTokenCommand struct {
	ID          uuid.UUID  `json:"id,omitempty"`
	TenantID    uuid.UUID  `json:"tenant_id"`
	UserID      uuid.UUID  `json:"user_id"`
	DeviceToken string     `json:"device_token"`
	DeviceType  *string    `json:"device_type,omitempty"`
	DeviceID    *string    `json:"device_id,omitempty"`
	ActorID     *uuid.UUID `json:"-"`
}
