package domain

import (
	"errors"
	"strings"
)

var ErrInvalidEnumValue = errors.New("invalid enum value")

func ValidateRole(value string) (string, error) {
	return validateLowerEnum(value, RoleSuperAdmin, RoleTenant, RoleHR, RoleManager, RoleEmployee, RoleApplicant)
}

func ValidateLeaveStatus(value string) (string, error) {
	return validateLowerEnum(value, LeaveStatusPending, LeaveStatusApproved, LeaveStatusRejected, LeaveStatusCanceled)
}

func ValidateLeaveDayType(value string) (string, error) {
	return validateLowerEnum(value, LeaveDayFullDay, LeaveDayFirstHalf, LeaveDaySecondHalf)
}

func ValidateAttendanceType(value string) (string, error) {
	return validateLowerEnum(value, AttendanceCheckin, AttendanceCheckout)
}

func ValidateAttendanceStatus(value string) (string, error) {
	return validateLowerEnum(value, AttendanceStatusPresent, AttendanceStatusLeave, AttendanceStatusAbsent, AttendanceStatusHoliday, AttendanceStatusWeekoff, AttendanceStatusIncomplete)
}

func ValidateAllocationType(value string) (string, error) {
	return validateLowerEnum(value, AllocationFixed, AllocationMonthly)
}

func ValidateNotificationCode(value string) (string, error) {
	return validateLowerEnum(value, NotifLeaveApplied, NotifLeaveApproved, NotifLeaveClarify, NotifLeaveRejected, NotifCompanyPolicy, NotifGeneralNotif, NotifUserCelebration)
}

func ValidateNotificationChannel(value string) (string, error) {
	return validateExactEnum(strings.TrimSpace(value), NotifChannelPush, NotifChannelEmail)
}

func ValidateNotificationStatus(value string) (string, error) {
	return validateExactEnum(strings.TrimSpace(value), NotifStatusPending, NotifStatusSent, NotifStatusFailed)
}

func ValidateReferenceTable(value string) (string, error) {
	return validateLowerEnum(value, RefTableUserLeaves, RefTableLeaveMessages, RefTableCompanyPolicy, RefTableUserCelebration, RefTableHoliday)
}

func ValidateTransactionType(value string) (string, error) {
	return validateLowerEnum(value, TransactionDebit, TransactionCredit)
}

func ValidateRegularisationKey(value string) (string, error) {
	return validateLowerEnum(value, RangeKeyHalfday, RangeKeyAbsent)
}

func ValidateOTPFor(value string) (string, error) {
	return validateLowerEnum(value, OtpForPasswordReset, OtpForLogin)
}

func ValidateSalaryItemType(value string) (string, error) {
	return validateLowerEnum(value, SalaryItemEarning, SalaryItemDeduction)
}

func ValidateApplicationStatus(value string) (string, error) {
	return validateExactEnum(strings.TrimSpace(value), AppStatusNew, AppStatusScreening, AppStatusInterview, AppStatusOffered, AppStatusHired, AppStatusRejected, AppStatusWithdrawn)
}

func ValidateOfferStatus(value string) (string, error) {
	return validateExactEnum(strings.TrimSpace(value), OfferStatusGenerated, OfferStatusSent, OfferStatusAccepted, OfferStatusDeclined, OfferStatusRevoked)
}

func ValidateRequisitionStatus(value string) (string, error) {
	return validateExactEnum(strings.TrimSpace(value), ReqStatusDraft, ReqStatusPending, ReqStatusApproved, ReqStatusRejected, ReqStatusClosed)
}

func ValidateOnboardingStatus(value string) (string, error) {
	return validateExactEnum(strings.TrimSpace(value), OnboardStatusPending, OnboardStatusInProgress, OnboardStatusCompleted)
}

func validateLowerEnum(value string, allowed ...string) (string, error) {
	return validateExactEnum(strings.ToLower(strings.TrimSpace(value)), allowed...)
}

func validateExactEnum(value string, allowed ...string) (string, error) {
	for _, item := range allowed {
		if value == item {
			return value, nil
		}
	}
	return "", ErrInvalidEnumValue
}
