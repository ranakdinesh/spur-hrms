package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidNotificationMasterID = errors.New("notification_master_id is required")
	ErrInvalidNotificationCode     = errors.New("notification code is invalid")
	ErrInvalidNotificationName     = errors.New("notification name is required")
	ErrNotificationMasterNotFound  = errors.New("notification master not found")
	ErrInvalidNotificationUser     = errors.New("notification user_id is required")
	ErrInvalidDigestFrequency      = errors.New("digest frequency is invalid")
	ErrInvalidQuietHours           = errors.New("quiet hours must include both start and end")
	ErrInvalidNotificationTitle    = errors.New("notification title is required")
	ErrInvalidNotificationMessage  = errors.New("notification message is required")
	ErrInvalidDeviceToken          = errors.New("device token is required")
	ErrInvalidDeviceID             = errors.New("device_id is required")
	ErrInvalidDeviceTokenID        = errors.New("device token id is required")
)

const (
	DigestFrequencyInstant = "instant"
	DigestFrequencyDaily   = "daily"
	DigestFrequencyWeekly  = "weekly"
	DigestFrequencyMuted   = "muted"
)

type NotificationMaster struct {
	ID                   uuid.UUID  `json:"id"`
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
	Inactive             bool       `json:"inactive"`
	CreatedAt            time.Time  `json:"created_at"`
	CreatedBy            *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt            time.Time  `json:"updated_at"`
	UpdatedBy            *uuid.UUID `json:"updated_by,omitempty"`
}

type NotificationMasterInput struct {
	TenantID             uuid.UUID
	Code                 string
	Name                 *string
	Description          *string
	IsInAppEnabled       bool
	IsEmailEnabled       bool
	IsPushEnabled        bool
	EmailSubjectTemplate *string
	EmailTextTemplate    *string
	EmailHTMLTemplate    *string
}

func NewNotificationMaster(input NotificationMasterInput) (*NotificationMaster, error) {
	if input.TenantID == uuid.Nil {
		return nil, ErrInvalidTenantID
	}
	code, err := ValidateNotificationCode(input.Code)
	if err != nil {
		return nil, ErrInvalidNotificationCode
	}
	name := cleanOptional(input.Name)
	if name == nil || strings.TrimSpace(*name) == "" {
		return nil, ErrInvalidNotificationName
	}
	now := time.Now().UTC()
	return &NotificationMaster{TenantID: input.TenantID, Code: code, Name: name, Description: cleanOptional(input.Description), IsInAppEnabled: input.IsInAppEnabled, IsEmailEnabled: input.IsEmailEnabled, IsPushEnabled: input.IsPushEnabled, EmailSubjectTemplate: cleanOptional(input.EmailSubjectTemplate), EmailTextTemplate: cleanOptional(input.EmailTextTemplate), EmailHTMLTemplate: cleanOptional(input.EmailHTMLTemplate), CreatedAt: now, UpdatedAt: now}, nil
}

func DefaultNotificationMasterInputs(tenantID uuid.UUID) []NotificationMasterInput {
	return []NotificationMasterInput{
		{TenantID: tenantID, Code: NotifLeaveApplied, Name: stringPtr("Leave Applied"), Description: stringPtr("Leave request submissions and approval workflow updates."), IsInAppEnabled: true, IsEmailEnabled: true, IsPushEnabled: true, EmailSubjectTemplate: stringPtr("{{tenant_name}} leave request submitted"), EmailTextTemplate: defaultEmailTextTemplate("Your leave request has been submitted.")},
		{TenantID: tenantID, Code: NotifLeaveApproved, Name: stringPtr("Leave Approved"), Description: stringPtr("Approved leave request notifications."), IsInAppEnabled: true, IsEmailEnabled: true, IsPushEnabled: true, EmailSubjectTemplate: stringPtr("{{tenant_name}} leave request approved"), EmailTextTemplate: defaultEmailTextTemplate("Your leave request has been approved.")},
		{TenantID: tenantID, Code: NotifLeaveRejected, Name: stringPtr("Leave Rejected"), Description: stringPtr("Rejected leave request notifications."), IsInAppEnabled: true, IsEmailEnabled: true, IsPushEnabled: true, EmailSubjectTemplate: stringPtr("{{tenant_name}} leave request update"), EmailTextTemplate: defaultEmailTextTemplate("Your leave request has been reviewed.")},
		{TenantID: tenantID, Code: NotifCompanyPolicy, Name: stringPtr("Company Policy"), Description: stringPtr("Policy publish and acknowledgement reminders."), IsInAppEnabled: true, IsEmailEnabled: true, IsPushEnabled: false, EmailSubjectTemplate: stringPtr("{{tenant_name}} policy update"), EmailTextTemplate: defaultEmailTextTemplate("A company policy update is available.")},
		{TenantID: tenantID, Code: NotifUserCelebration, Name: stringPtr("Celebrations"), Description: stringPtr("Birthday, work anniversary, festival, and team event reminders."), IsInAppEnabled: true, IsEmailEnabled: true, IsPushEnabled: true, EmailSubjectTemplate: stringPtr("{{tenant_name}} celebration reminder"), EmailTextTemplate: defaultEmailTextTemplate("There is a team celebration update.")},
		{TenantID: tenantID, Code: NotifGeneralNotif, Name: stringPtr("General Notifications"), Description: stringPtr("Tenant announcements and general HR reminders."), IsInAppEnabled: true, IsEmailEnabled: false, IsPushEnabled: true, EmailSubjectTemplate: stringPtr("{{tenant_name}} notification"), EmailTextTemplate: defaultEmailTextTemplate("You have a new notification.")},
	}
}

func defaultEmailTextTemplate(intro string) *string {
	return stringPtr("Hello {{employee_name}},\n\n" + intro + "\n\n{{notification_message}}\n\nRegards,\n{{tenant_name}}")
}

type NotificationPreference struct {
	ID                   uuid.UUID  `json:"id"`
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
	Inactive             bool       `json:"inactive"`
	CreatedAt            time.Time  `json:"created_at"`
	CreatedBy            *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt            time.Time  `json:"updated_at"`
	UpdatedBy            *uuid.UUID `json:"updated_by,omitempty"`
}

type NotificationPreferenceInput struct {
	TenantID             uuid.UUID
	UserID               uuid.UUID
	NotificationMasterID uuid.UUID
	IsInAppEnabled       bool
	IsEmailEnabled       bool
	IsPushEnabled        bool
	DigestFrequency      string
	QuietHoursStart      *string
	QuietHoursEnd        *string
	Timezone             *string
}

func NewNotificationPreference(input NotificationPreferenceInput) (*NotificationPreference, error) {
	if input.TenantID == uuid.Nil {
		return nil, ErrInvalidTenantID
	}
	if input.UserID == uuid.Nil {
		return nil, ErrInvalidNotificationUser
	}
	if input.NotificationMasterID == uuid.Nil {
		return nil, ErrInvalidNotificationMasterID
	}
	frequency, err := ValidateDigestFrequency(input.DigestFrequency)
	if err != nil {
		return nil, err
	}
	if (input.QuietHoursStart == nil) != (input.QuietHoursEnd == nil) {
		return nil, ErrInvalidQuietHours
	}
	now := time.Now().UTC()
	return &NotificationPreference{TenantID: input.TenantID, UserID: input.UserID, NotificationMasterID: input.NotificationMasterID, IsInAppEnabled: input.IsInAppEnabled, IsEmailEnabled: input.IsEmailEnabled, IsPushEnabled: input.IsPushEnabled, DigestFrequency: frequency, QuietHoursStart: cleanOptional(input.QuietHoursStart), QuietHoursEnd: cleanOptional(input.QuietHoursEnd), Timezone: cleanOptional(input.Timezone), CreatedAt: now, UpdatedAt: now}, nil
}

func ValidateDigestFrequency(value string) (string, error) {
	frequency := strings.ToLower(strings.TrimSpace(value))
	if frequency == "" {
		frequency = DigestFrequencyInstant
	}
	switch frequency {
	case DigestFrequencyInstant, DigestFrequencyDaily, DigestFrequencyWeekly, DigestFrequencyMuted:
		return frequency, nil
	default:
		return "", ErrInvalidDigestFrequency
	}
}

type EffectiveNotificationPreference struct {
	Master              *NotificationMaster     `json:"master"`
	Preference          *NotificationPreference `json:"preference,omitempty"`
	IsInAppEnabled      bool                    `json:"is_in_app_enabled"`
	IsEmailEnabled      bool                    `json:"is_email_enabled"`
	IsPushEnabled       bool                    `json:"is_push_enabled"`
	DigestFrequency     string                  `json:"digest_frequency"`
	QuietHoursStart     *string                 `json:"quiet_hours_start,omitempty"`
	QuietHoursEnd       *string                 `json:"quiet_hours_end,omitempty"`
	Timezone            *string                 `json:"timezone,omitempty"`
	UsesTenantDefaults  bool                    `json:"uses_tenant_defaults"`
	ProviderReady       map[string]bool         `json:"provider_ready"`
	FutureReadyFeatures []string                `json:"future_ready_features"`
}

type NotificationInboxItem struct {
	ID                   uuid.UUID  `json:"id"`
	TenantID             uuid.UUID  `json:"tenant_id"`
	UserID               uuid.UUID  `json:"user_id"`
	NotificationMasterID uuid.UUID  `json:"notification_master_id"`
	NotificationCode     *string    `json:"notification_code,omitempty"`
	NotificationName     *string    `json:"notification_name,omitempty"`
	Title                string     `json:"title"`
	Message              string     `json:"message"`
	ReferenceTable       *string    `json:"reference_table,omitempty"`
	ReferenceID          *uuid.UUID `json:"reference_id,omitempty"`
	IsRead               bool       `json:"is_read"`
	ReadDate             *time.Time `json:"read_date,omitempty"`
	Inactive             bool       `json:"inactive"`
	CreatedAt            time.Time  `json:"created_at"`
	CreatedBy            *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt            time.Time  `json:"updated_at"`
	UpdatedBy            *uuid.UUID `json:"updated_by,omitempty"`
}

type NotificationInboxFilter struct {
	TenantID             uuid.UUID
	UserID               uuid.UUID
	IsRead               *bool
	NotificationMasterID *uuid.UUID
	NotificationCode     *string
	Search               *string
	Limit                int32
	Offset               int32
}

type NotificationInboxPage struct {
	Items      []*NotificationInboxItem `json:"items"`
	Total      int64                    `json:"total"`
	Unread     int64                    `json:"unread"`
	Limit      int32                    `json:"limit"`
	Offset     int32                    `json:"offset"`
	NextOffset *int32                   `json:"next_offset,omitempty"`
}

type NotificationLog struct {
	ID                   uuid.UUID  `json:"id"`
	TenantID             uuid.UUID  `json:"tenant_id"`
	NotificationMasterID *uuid.UUID `json:"notification_master_id,omitempty"`
	UserID               *uuid.UUID `json:"user_id,omitempty"`
	Channel              string     `json:"channel"`
	TargetAddress        string     `json:"target_address"`
	Subject              *string    `json:"subject,omitempty"`
	Body                 *string    `json:"body,omitempty"`
	Status               string     `json:"status"`
	SentDate             *time.Time `json:"sent_date,omitempty"`
	ErrorMessage         *string    `json:"error_message,omitempty"`
	ExternalReferenceID  *string    `json:"external_reference_id,omitempty"`
	IdempotencyKey       *string    `json:"idempotency_key,omitempty"`
	BulkID               *uuid.UUID `json:"bulk_id,omitempty"`
	Provider             *string    `json:"provider,omitempty"`
	ProviderMessageID    *string    `json:"provider_message_id,omitempty"`
	ProviderEventStatus  *string    `json:"provider_event_status,omitempty"`
	ProviderEventAt      *time.Time `json:"provider_event_at,omitempty"`
	AttemptCount         int32      `json:"attempt_count"`
	LastAttemptAt        *time.Time `json:"last_attempt_at,omitempty"`
	Inactive             bool       `json:"inactive"`
	CreatedAt            time.Time  `json:"created_at"`
	CreatedBy            *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt            time.Time  `json:"updated_at"`
	UpdatedBy            *uuid.UUID `json:"updated_by,omitempty"`
}

type NotificationSendInput struct {
	TenantID         uuid.UUID
	NotificationCode string
	UserIDs          []uuid.UUID
	Title            string
	Message          string
	ReferenceTable   *string
	ReferenceID      *uuid.UUID
	Channels         []string
	IdempotencyKey   *string
	BulkID           *uuid.UUID
	ActorID          *uuid.UUID
}

type NotificationSendResult struct {
	BulkID        uuid.UUID          `json:"bulk_id"`
	Processed     int32              `json:"processed"`
	InboxCount    int32              `json:"inbox_count"`
	LogCount      int32              `json:"log_count"`
	Sent          int32              `json:"sent"`
	Suppressed    int32              `json:"suppressed"`
	Failed        int32              `json:"failed"`
	MissingMaster bool               `json:"missing_master"`
	Logs          []*NotificationLog `json:"logs"`
}

type DeviceToken struct {
	ID          uuid.UUID  `json:"id"`
	TenantID    uuid.UUID  `json:"tenant_id"`
	UserID      uuid.UUID  `json:"user_id"`
	DeviceToken string     `json:"device_token"`
	DeviceType  *string    `json:"device_type,omitempty"`
	DeviceID    *string    `json:"device_id,omitempty"`
	Inactive    bool       `json:"inactive"`
	CreatedAt   time.Time  `json:"created_at"`
	CreatedBy   *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt   time.Time  `json:"updated_at"`
	UpdatedBy   *uuid.UUID `json:"updated_by,omitempty"`
}

type DeviceTokenInput struct {
	TenantID    uuid.UUID
	UserID      uuid.UUID
	DeviceToken string
	DeviceType  *string
	DeviceID    *string
}

func NewDeviceToken(input DeviceTokenInput) (*DeviceToken, error) {
	if input.TenantID == uuid.Nil {
		return nil, ErrInvalidTenantID
	}
	if input.UserID == uuid.Nil {
		return nil, ErrInvalidNotificationUser
	}
	token := strings.TrimSpace(input.DeviceToken)
	if token == "" {
		return nil, ErrInvalidDeviceToken
	}
	deviceID := cleanOptional(input.DeviceID)
	if deviceID == nil {
		return nil, ErrInvalidDeviceID
	}
	now := time.Now().UTC()
	return &DeviceToken{TenantID: input.TenantID, UserID: input.UserID, DeviceToken: token, DeviceType: cleanOptional(input.DeviceType), DeviceID: deviceID, CreatedAt: now, UpdatedAt: now}, nil
}
