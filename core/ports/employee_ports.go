package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

type EmployeeRepo interface {
	CreateEmployee(ctx context.Context, employee *domain.Employee, actorID *uuid.UUID) (*domain.Employee, error)
	CountActiveEmployees(ctx context.Context, tenantID uuid.UUID) (int32, error)
	DeactivateEmployee(ctx context.Context, tenantID uuid.UUID, employeeID uuid.UUID, actorID *uuid.UUID) error
	EmployeeCodeExists(ctx context.Context, tenantID uuid.UUID, employeeCode string) (bool, error)
	EmployeeCodeExistsForOtherEmployee(ctx context.Context, tenantID uuid.UUID, employeeID uuid.UUID, employeeCode string) (bool, error)
	GetEmployeeByCode(ctx context.Context, tenantID uuid.UUID, employeeCode string) (*domain.Employee, error)
	GetEmployeeAttendanceRequired(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID) (bool, error)
	GetEmployeeByUserID(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID) (*domain.Employee, error)
	GetEmployeeProfile(ctx context.Context, tenantID uuid.UUID, employeeID uuid.UUID) (*domain.EmployeeProfile, error)
	ListEmployees(ctx context.Context, tenantID uuid.UUID) ([]*domain.EmployeeListItem, error)
	UpdateEmployee(ctx context.Context, employee *domain.Employee, actorID *uuid.UUID) (*domain.Employee, error)
	UpsertEmployeeStatutory(ctx context.Context, item EmployeeStatutoryCommand) (*domain.EmployeeStatutory, error)
	UpsertPrimaryEmployeeBank(ctx context.Context, item EmployeeBankCommand) (*domain.EmployeeBank, error)
}

type EmployeeIdentityPort interface {
	CheckEmployeeIdentityAvailability(ctx context.Context, cmd EmployeeIdentityAvailabilityCommand) error
	CreateEmployeeIdentity(ctx context.Context, cmd CreateEmployeeIdentityCommand) (*EmployeeIdentity, error)
	AssignEmployeeRole(ctx context.Context, cmd AssignEmployeeRoleCommand) error
	DeactivateEmployeeIdentity(ctx context.Context, cmd DeactivateEmployeeIdentityCommand) error
	SendEmployeePasswordReset(ctx context.Context, cmd EmployeeCredentialResetCommand) error
	SetEmployeeTemporaryPassword(ctx context.Context, cmd EmployeeTemporaryPasswordCommand) error
}

type EmployeeIdentityAvailabilityCommand struct {
	TenantID uuid.UUID
	Email    string
	Mobile   string
}

type CreateEmployeeIdentityCommand struct {
	TenantID  uuid.UUID
	FirstName string
	LastName  string
	Email     string
	Mobile    string
	Password  string
	Role      string
	ActorID   *uuid.UUID
}

type AssignEmployeeRoleCommand struct {
	TenantID uuid.UUID
	UserID   uuid.UUID
	Role     string
	ActorID  *uuid.UUID
}

type DeactivateEmployeeIdentityCommand struct {
	TenantID uuid.UUID
	UserID   uuid.UUID
	ActorID  *uuid.UUID
}

type EmployeeIdentity struct {
	UserID uuid.UUID
}

type EmployeeCredentialEventRepo interface {
	CreateEmployeeCredentialEvent(ctx context.Context, event *domain.EmployeeCredentialEvent) (*domain.EmployeeCredentialEvent, error)
	ListEmployeeCredentialEvents(ctx context.Context, tenantID uuid.UUID, employeeID uuid.UUID, limit int32) ([]*domain.EmployeeCredentialEvent, error)
}

type EmployeeCredentialResetCommand struct {
	TenantID     uuid.UUID
	UserID       uuid.UUID
	Email        string
	EmployeeID   uuid.UUID
	Employee     string
	EmployeeCode string
	ActorID      *uuid.UUID
}

type EmployeeTemporaryPasswordCommand struct {
	TenantID          uuid.UUID
	UserID            uuid.UUID
	Email             string
	TemporaryPassword string
	EmployeeID        uuid.UUID
	Employee          string
	EmployeeCode      string
	ActorID           *uuid.UUID
}

type EmployeeCredentialActionCommand struct {
	TenantID          uuid.UUID  `json:"tenant_id"`
	EmployeeID        uuid.UUID  `json:"employee_id"`
	TemporaryPassword string     `json:"temporary_password,omitempty"`
	ActorID           *uuid.UUID `json:"-"`
}

type CreateEmployeeCommand struct {
	TenantID              uuid.UUID                 `json:"tenant_id"`
	UserID                uuid.UUID                 `json:"user_id,omitempty"`
	EmployeeCode          *string                   `json:"employee_code,omitempty"`
	FirstName             string                    `json:"first_name"`
	MiddleName            *string                   `json:"middle_name,omitempty"`
	LastName              *string                   `json:"last_name,omitempty"`
	Email                 *string                   `json:"email,omitempty"`
	Mobile                *string                   `json:"mobile,omitempty"`
	Password              string                    `json:"password,omitempty"`
	DOB                   string                    `json:"dob,omitempty"`
	Gender                *string                   `json:"gender,omitempty"`
	MaritalStatus         *string                   `json:"marital_status,omitempty"`
	BloodGroup            *string                   `json:"blood_group,omitempty"`
	ProfilePhotoPath      *string                   `json:"profile_photo_path,omitempty"`
	Address               *string                   `json:"address,omitempty"`
	City                  *string                   `json:"city,omitempty"`
	State                 *string                   `json:"state,omitempty"`
	Country               *string                   `json:"country,omitempty"`
	Pincode               *string                   `json:"pincode,omitempty"`
	EmergencyContact      *string                   `json:"emergency_contact,omitempty"`
	JoiningDate           string                    `json:"joining_date,omitempty"`
	DepartmentID          *uuid.UUID                `json:"department_id,omitempty"`
	BranchID              *uuid.UUID                `json:"branch_id,omitempty"`
	DesignationID         *uuid.UUID                `json:"designation_id,omitempty"`
	ReportingManagerID    *uuid.UUID                `json:"reporting_manager_id,omitempty"`
	EmploymentTypeID      *uuid.UUID                `json:"employment_type_id,omitempty"`
	Role                  string                    `json:"role,omitempty"`
	Grade                 *string                   `json:"grade,omitempty"`
	ExperienceYear        int32                     `json:"experience_year"`
	ExperienceMonth       int32                     `json:"experience_month"`
	ProbationStatus       string                    `json:"probation_status,omitempty"`
	ProbationStartDate    string                    `json:"probation_start_date,omitempty"`
	ProbationEndDate      string                    `json:"probation_end_date,omitempty"`
	ProbationDurationDays int32                     `json:"probation_duration_days"`
	ProbationConfirmedAt  string                    `json:"probation_confirmed_at,omitempty"`
	IsPayrollStaff        bool                      `json:"is_payroll_staff"`
	Statutory             *EmployeeStatutoryCommand `json:"statutory,omitempty"`
	ActorID               *uuid.UUID                `json:"-"`
}

type UpdateEmployeeCommand struct {
	ID                    uuid.UUID                 `json:"id"`
	TenantID              uuid.UUID                 `json:"tenant_id"`
	EmployeeCode          *string                   `json:"employee_code,omitempty"`
	FirstName             string                    `json:"first_name"`
	MiddleName            *string                   `json:"middle_name,omitempty"`
	LastName              *string                   `json:"last_name,omitempty"`
	Email                 *string                   `json:"email,omitempty"`
	Mobile                *string                   `json:"mobile,omitempty"`
	DOB                   string                    `json:"dob,omitempty"`
	Gender                *string                   `json:"gender,omitempty"`
	MaritalStatus         *string                   `json:"marital_status,omitempty"`
	BloodGroup            *string                   `json:"blood_group,omitempty"`
	ProfilePhotoPath      *string                   `json:"profile_photo_path,omitempty"`
	Address               *string                   `json:"address,omitempty"`
	City                  *string                   `json:"city,omitempty"`
	State                 *string                   `json:"state,omitempty"`
	Country               *string                   `json:"country,omitempty"`
	Pincode               *string                   `json:"pincode,omitempty"`
	EmergencyContact      *string                   `json:"emergency_contact,omitempty"`
	JoiningDate           string                    `json:"joining_date,omitempty"`
	ResignationDate       string                    `json:"resignation_date,omitempty"`
	DepartmentID          *uuid.UUID                `json:"department_id,omitempty"`
	BranchID              *uuid.UUID                `json:"branch_id,omitempty"`
	DesignationID         *uuid.UUID                `json:"designation_id,omitempty"`
	ReportingManagerID    *uuid.UUID                `json:"reporting_manager_id,omitempty"`
	EmploymentTypeID      *uuid.UUID                `json:"employment_type_id,omitempty"`
	Role                  string                    `json:"role,omitempty"`
	Grade                 *string                   `json:"grade,omitempty"`
	ExperienceYear        int32                     `json:"experience_year"`
	ExperienceMonth       int32                     `json:"experience_month"`
	ProbationStatus       string                    `json:"probation_status,omitempty"`
	ProbationStartDate    string                    `json:"probation_start_date,omitempty"`
	ProbationEndDate      string                    `json:"probation_end_date,omitempty"`
	ProbationDurationDays int32                     `json:"probation_duration_days"`
	ProbationConfirmedAt  string                    `json:"probation_confirmed_at,omitempty"`
	IsPayrollStaff        bool                      `json:"is_payroll_staff"`
	Bank                  *EmployeeBankCommand      `json:"bank,omitempty"`
	Statutory             *EmployeeStatutoryCommand `json:"statutory,omitempty"`
	ActorID               *uuid.UUID                `json:"-"`
}

type EmployeeBankCommand struct {
	TenantID      uuid.UUID  `json:"tenant_id"`
	UserID        uuid.UUID  `json:"user_id"`
	BankName      *string    `json:"bank_name,omitempty"`
	AccountNumber *string    `json:"account_number,omitempty"`
	IFSCCode      *string    `json:"ifsc_code,omitempty"`
	AccountType   *string    `json:"account_type,omitempty"`
	BranchName    *string    `json:"branch_name,omitempty"`
	ActorID       *uuid.UUID `json:"-"`
}

type EmployeeStatutoryCommand struct {
	TenantID       uuid.UUID  `json:"tenant_id"`
	UserID         uuid.UUID  `json:"user_id"`
	PFNo           *string    `json:"pf_no,omitempty"`
	UANNo          *string    `json:"uan_no,omitempty"`
	ESICNo         *string    `json:"esic_no,omitempty"`
	PAN            *string    `json:"pan,omitempty"`
	Aadhaar        *string    `json:"aadhaar,omitempty"`
	PTApplicable   bool       `json:"pt_applicable"`
	PFApplicable   bool       `json:"pf_applicable"`
	ESICApplicable bool       `json:"esic_applicable"`
	LWFApplicable  bool       `json:"lwf_applicable"`
	ActorID        *uuid.UUID `json:"-"`
}
