package domain

const (
	RoleSuperAdmin = "superadmin"
	RoleTenant     = "tenant"
	RoleHR         = "hr"
	RoleManager    = "manager"
	RoleEmployee   = "employee"
	RoleApplicant  = "applicant"
)

const (
	LeaveStatusPending  = "pending"
	LeaveStatusApproved = "approved"
	LeaveStatusRejected = "rejected"
	LeaveStatusCanceled = "canceled"
)

const (
	LeaveDayFullDay    = "fullday"
	LeaveDayFirstHalf  = "firsthalf"
	LeaveDaySecondHalf = "secondhalf"
)

const LeaveTypeShortEarnLeave = "earnleave"

const (
	AttendanceCheckin  = "checkin"
	AttendanceCheckout = "checkout"
)

const (
	AttendanceStatusPresent    = "present"
	AttendanceStatusLeave      = "leave"
	AttendanceStatusAbsent     = "absent"
	AttendanceStatusHoliday    = "holiday"
	AttendanceStatusWeekoff    = "weekoff"
	AttendanceStatusIncomplete = "incomplete"
)

const (
	AllocationFixed   = "fixed"
	AllocationMonthly = "monthly"
)

const (
	NotifLeaveApplied    = "leaveapplied"
	NotifLeaveApproved   = "leaveapproved"
	NotifLeaveClarify    = "leaveclarification"
	NotifLeaveRejected   = "leaverejected"
	NotifCompanyPolicy   = "companypolicy"
	NotifGeneralNotif    = "generalnotification"
	NotifUserCelebration = "usercelebration"
)

const (
	NotifChannelPush     = "Push"
	NotifChannelEmail    = "Email"
	NotifChannelSMS      = "SMS"
	NotifChannelWhatsApp = "WhatsApp"
)

const (
	NotifStatusPending    = "Pending"
	NotifStatusSent       = "Sent"
	NotifStatusFailed     = "Failed"
	NotifStatusSuppressed = "Suppressed"
)

const (
	RefTableUserLeaves      = "userleaves"
	RefTableLeaveMessages   = "leavemessages"
	RefTableCompanyPolicy   = "companypolicy"
	RefTableUserCelebration = "usercelebration"
	RefTableHoliday         = "holiday"
)

const (
	TransactionDebit  = "debit"
	TransactionCredit = "credit"
)

const (
	RangeKeyHalfday = "halfday"
	RangeKeyAbsent  = "absent"
)

const (
	OtpForPasswordReset = "passwordreset"
	OtpForLogin         = "login"
)

const (
	SalaryItemEarning   = "earning"
	SalaryItemDeduction = "deduction"
)

const (
	SalaryCodeBasic = "basic"
	SalaryCodeHRA   = "hra"
	SalaryCodeLWP   = "LWP"
)

const (
	DashToolLeaveRequest  = "leaverequest"
	DashToolAttendance    = "attendancerecord"
	DashToolForgotPunch   = "forgottopunch"
	DashToolLeaveApproval = "leaveapproval"
	DashToolCelebration   = "celebration"
	DashToolHolidays      = "holidays"
	DashToolPolicies      = "policies"
)

const (
	AppStatusNew       = "New"
	AppStatusScreening = "Screening"
	AppStatusInterview = "Interview"
	AppStatusOffered   = "Offered"
	AppStatusHired     = "Hired"
	AppStatusRejected  = "Rejected"
	AppStatusWithdrawn = "Withdrawn"
)

const (
	OfferStatusGenerated = "Generated"
	OfferStatusSent      = "Sent"
	OfferStatusAccepted  = "Accepted"
	OfferStatusDeclined  = "Declined"
	OfferStatusRevoked   = "Revoked"
)

const (
	ReqStatusDraft    = "Draft"
	ReqStatusPending  = "Pending"
	ReqStatusApproved = "Approved"
	ReqStatusRejected = "Rejected"
	ReqStatusClosed   = "Closed"
)

const (
	OnboardStatusPending    = "Pending"
	OnboardStatusInProgress = "InProgress"
	OnboardStatusCompleted  = "Completed"
)
