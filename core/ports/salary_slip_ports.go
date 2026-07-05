package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

type SalarySlipRepo interface {
	CreateSalarySlip(ctx context.Context, item *domain.SalarySlip, actorID *uuid.UUID) (*domain.SalarySlip, error)
	UpdateSalarySlip(ctx context.Context, item *domain.SalarySlip, actorID *uuid.UUID) (*domain.SalarySlip, error)
	ListSalarySlipsByUser(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID) ([]*domain.SalarySlip, error)
	ListRecentSalarySlipsByUser(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID, limit int32) ([]*domain.SalarySlip, error)
	ListSalarySlipsByTenantPeriod(ctx context.Context, tenantID uuid.UUID, month int32, year int32) ([]*domain.SalarySlip, error)
	GetSalarySlip(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.SalarySlip, error)
	GetSalarySlipByPeriod(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID, month int32, year int32) (*domain.SalarySlip, error)
	CreateSalarySlipItem(ctx context.Context, item *domain.SalarySlipItem, actorID *uuid.UUID) (*domain.SalarySlipItem, error)
	CreateSalarySlipLeave(ctx context.Context, item *domain.SalarySlipLeave, actorID *uuid.UUID) (*domain.SalarySlipLeave, error)
	ListSalarySlipItems(ctx context.Context, tenantID uuid.UUID, slipID uuid.UUID) ([]*domain.SalarySlipItem, error)
	ListSalarySlipLeaves(ctx context.Context, tenantID uuid.UUID, slipID uuid.UUID) ([]*domain.SalarySlipLeave, error)
	DeleteSalarySlipItemsBySlip(ctx context.Context, tenantID uuid.UUID, slipID uuid.UUID, actorID *uuid.UUID) error
	DeleteSalarySlipLeavesBySlip(ctx context.Context, tenantID uuid.UUID, slipID uuid.UUID, actorID *uuid.UUID) error
	GetSalarySlipFormat(ctx context.Context, tenantID uuid.UUID) (*domain.SalarySlipFormat, error)
	UpsertSalarySlipFormat(ctx context.Context, item *domain.SalarySlipFormat, actorID *uuid.UUID) (*domain.SalarySlipFormat, error)
}

type SalarySlipPDFRenderer interface {
	RenderSalarySlipPDF(ctx context.Context, doc SalarySlipDocument) ([]byte, error)
}

type SalarySlipStorage interface {
	StoreSalarySlipPDF(ctx context.Context, input StoreSalarySlipPDFInput) (string, error)
}

type StoreSalarySlipPDFInput struct {
	TenantID    uuid.UUID
	UserID      uuid.UUID
	SlipID      uuid.UUID
	Month       int32
	Year        int32
	FileName    string
	ContentType string
	Content     []byte
}

type SalarySlipDocument struct {
	Format   *domain.SalarySlipFormat `json:"format"`
	Slip     *domain.SalarySlip       `json:"slip"`
	Employee *domain.Employee         `json:"employee,omitempty"`
}

type SalarySlipFormatCommand struct {
	TenantID                  uuid.UUID      `json:"tenant_id"`
	Title                     string         `json:"title"`
	Subtitle                  *string        `json:"subtitle,omitempty"`
	LogoPath                  *string        `json:"logo_path,omitempty"`
	PrimaryColor              string         `json:"primary_color"`
	AccentColor               string         `json:"accent_color"`
	ShowLeaveBalance          bool           `json:"show_leave_balance"`
	ShowYTDSummary            bool           `json:"show_ytd_summary"`
	ShowEmployeeBank          bool           `json:"show_employee_bank"`
	ShowEmployerContributions bool           `json:"show_employer_contributions"`
	FooterText                *string        `json:"footer_text,omitempty"`
	CustomFields              map[string]any `json:"custom_fields"`
	ActorID                   *uuid.UUID     `json:"-"`
}

type GenerateSalarySlipCommand struct {
	TenantID    uuid.UUID  `json:"tenant_id"`
	UserID      uuid.UUID  `json:"user_id"`
	FYID        uuid.UUID  `json:"fy_id"`
	Month       int32      `json:"month"`
	Year        int32      `json:"year"`
	PresentDays *int       `json:"present_days,omitempty"`
	AbsentDays  *int       `json:"absent_days,omitempty"`
	TotalDays   *int       `json:"total_days,omitempty"`
	IsSpecial   bool       `json:"is_special"`
	Regenerate  bool       `json:"regenerate"`
	ActorID     *uuid.UUID `json:"-"`
}

type BulkGenerateSalarySlipsCommand struct {
	TenantID   uuid.UUID  `json:"tenant_id"`
	FYID       uuid.UUID  `json:"fy_id"`
	Month      int32      `json:"month"`
	Year       int32      `json:"year"`
	Regenerate bool       `json:"regenerate"`
	ActorID    *uuid.UUID `json:"-"`
}
