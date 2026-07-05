package domain

import (
	"errors"
	"testing"
)

func TestEnumValidatorsNormalizeLowercaseEnums(t *testing.T) {
	tests := []struct {
		name      string
		validate  func(string) (string, error)
		input     string
		expected  string
		wantError bool
	}{
		{"role", ValidateRole, " HR ", RoleHR, false},
		{"leave status", ValidateLeaveStatus, "APPROVED", LeaveStatusApproved, false},
		{"leave day", ValidateLeaveDayType, "FirstHalf", LeaveDayFirstHalf, false},
		{"attendance type", ValidateAttendanceType, "CHECKOUT", AttendanceCheckout, false},
		{"attendance status", ValidateAttendanceStatus, "WeekOff", AttendanceStatusWeekoff, false},
		{"allocation", ValidateAllocationType, "Monthly", AllocationMonthly, false},
		{"notification code", ValidateNotificationCode, "LeaveApproved", NotifLeaveApproved, false},
		{"transaction", ValidateTransactionType, "Credit", TransactionCredit, false},
		{"regularisation", ValidateRegularisationKey, "HalfDay", RangeKeyHalfday, false},
		{"otp", ValidateOTPFor, "Login", OtpForLogin, false},
		{"salary item", ValidateSalaryItemType, "Deduction", SalaryItemDeduction, false},
		{"invalid", ValidateRole, "owner", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.validate(tt.input)
			if tt.wantError {
				if !errors.Is(err, ErrInvalidEnumValue) {
					t.Fatalf("expected ErrInvalidEnumValue, got %v", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.expected {
				t.Fatalf("got %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestEnumValidatorsKeepExactCaseEnums(t *testing.T) {
	tests := []struct {
		name     string
		validate func(string) (string, error)
		input    string
		expected string
	}{
		{"notification channel", ValidateNotificationChannel, " Push ", NotifChannelPush},
		{"notification status", ValidateNotificationStatus, "Sent", NotifStatusSent},
		{"application status", ValidateApplicationStatus, "Screening", AppStatusScreening},
		{"offer status", ValidateOfferStatus, "Accepted", OfferStatusAccepted},
		{"requisition status", ValidateRequisitionStatus, "Closed", ReqStatusClosed},
		{"onboarding status", ValidateOnboardingStatus, "InProgress", OnboardStatusInProgress},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.validate(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.expected {
				t.Fatalf("got %q, want %q", got, tt.expected)
			}
		})
	}

	if _, err := ValidateNotificationChannel("push"); !errors.Is(err, ErrInvalidEnumValue) {
		t.Fatalf("expected exact-case validation error, got %v", err)
	}
}
