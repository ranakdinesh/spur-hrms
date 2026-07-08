package handlers

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
	"github.com/ranakdinesh/spur-hrms/pkg/permissions"
	"github.com/ranakdinesh/spur-identity/adapters/http/httputil"
)

type selfServiceTenantService struct {
	ports.TenantService
	listLeavesCalled       bool
	previewLeaveCalled     bool
	cancelLeaveCalled      bool
	listApprovalsCalled    bool
	createSegmentCalled    bool
	listSegmentsCalled     bool
	attendanceStatusCalled bool
	createCompOffCalled    bool
	listCompOffCalled      bool
	reviewCompOffCalled    bool
	gotUserID              uuid.UUID
	gotApproverID          uuid.UUID
	gotAttendanceUserID    *uuid.UUID
}

func (s *selfServiceTenantService) ListLeavesByUser(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID) ([]*domain.Leave, error) {
	s.listLeavesCalled = true
	s.gotUserID = userID
	return []*domain.Leave{}, nil
}

func (s *selfServiceTenantService) PreviewLeave(ctx context.Context, cmd ports.ApplyLeaveCommand) (*domain.LeavePreview, error) {
	s.previewLeaveCalled = true
	s.gotUserID = cmd.UserID
	return &domain.LeavePreview{TenantID: cmd.TenantID, UserID: cmd.UserID, LeaveTypeID: cmd.LeaveTypeID, FYID: cmd.FYID, Allowed: true}, nil
}

func (s *selfServiceTenantService) CancelLeave(ctx context.Context, cmd ports.CancelLeaveCommand) (*domain.LeaveApplication, error) {
	s.cancelLeaveCalled = true
	s.gotUserID = cmd.UserID
	return &domain.LeaveApplication{Leave: &domain.Leave{ID: cmd.LeaveID, TenantID: cmd.TenantID, UserID: cmd.UserID, Status: domain.LeaveStatusCanceled}}, nil
}

func (s *selfServiceTenantService) ListPendingApprovalsByApprover(ctx context.Context, tenantID uuid.UUID, approverID uuid.UUID) ([]*domain.LeaveApproval, error) {
	s.listApprovalsCalled = true
	s.gotApproverID = approverID
	return []*domain.LeaveApproval{}, nil
}

func (s *selfServiceTenantService) CreateAttendanceWorkdaySegment(ctx context.Context, cmd ports.AttendanceSegmentCommand) (*domain.AttendanceWorkdaySegment, error) {
	s.createSegmentCalled = true
	s.gotUserID = cmd.UserID
	return &domain.AttendanceWorkdaySegment{ID: uuid.New(), TenantID: cmd.TenantID, UserID: cmd.UserID, SegmentType: cmd.SegmentType, Action: cmd.Action}, nil
}

func (s *selfServiceTenantService) ListAttendanceWorkdaySegmentsByUser(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID, startDate string, endDate string) ([]*domain.AttendanceWorkdaySegment, error) {
	s.listSegmentsCalled = true
	s.gotUserID = userID
	return []*domain.AttendanceWorkdaySegment{}, nil
}

func (s *selfServiceTenantService) ListAttendanceWorkdaySegmentsByUserDate(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID, date string) ([]*domain.AttendanceWorkdaySegment, error) {
	s.listSegmentsCalled = true
	s.gotUserID = userID
	return []*domain.AttendanceWorkdaySegment{}, nil
}

func (s *selfServiceTenantService) ListAttendanceDailyStatuses(ctx context.Context, query ports.AttendanceStatusQuery) ([]*domain.AttendanceDailyStatus, error) {
	s.attendanceStatusCalled = true
	s.gotAttendanceUserID = query.UserID
	return []*domain.AttendanceDailyStatus{}, nil
}

func (s *selfServiceTenantService) CreateCompOffRequest(ctx context.Context, cmd ports.CompOffRequestCommand) (*domain.CompOffRequest, error) {
	s.createCompOffCalled = true
	s.gotUserID = cmd.UserID
	return &domain.CompOffRequest{ID: uuid.New(), TenantID: cmd.TenantID, UserID: cmd.UserID, Status: domain.CompOffStatusPending, RequestedDays: cmd.RequestedDays}, nil
}

func (s *selfServiceTenantService) ListCompOffRequestsByUser(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID) ([]*domain.CompOffRequest, error) {
	s.listCompOffCalled = true
	s.gotUserID = userID
	return []*domain.CompOffRequest{}, nil
}

func (s *selfServiceTenantService) ReviewCompOffRequest(ctx context.Context, cmd ports.CompOffReviewCommand) (*domain.CompOffRequest, error) {
	s.reviewCompOffCalled = true
	return &domain.CompOffRequest{ID: cmd.RequestID, TenantID: cmd.TenantID, Status: cmd.Status}, nil
}

func TestEmployeeCannotListOtherUserLeaves(t *testing.T) {
	tenantID := uuid.New()
	actorID := uuid.New()
	otherID := uuid.New()
	svc := &selfServiceTenantService{}
	handler := New(
		svc,
		func(context.Context) string { return tenantID.String() },
		func(context.Context) uuid.UUID { return actorID },
		func(context.Context) bool { return false },
		nil,
	)
	request := httptest.NewRequest(http.MethodGet, "/hrms/leaves?user_id="+otherID.String(), nil)
	request = request.WithContext(httputil.SetPermissions(request.Context(), []string{permissions.LeaveSelfView, permissions.LeavesList}))
	recorder := httptest.NewRecorder()

	handler.ListLeaves(recorder, request)

	if recorder.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d; body=%s", recorder.Code, http.StatusForbidden, recorder.Body.String())
	}
	if svc.listLeavesCalled {
		t.Fatal("service should not be called for other user's leaves")
	}
}

func TestEmployeeCanListOwnLeaves(t *testing.T) {
	tenantID := uuid.New()
	actorID := uuid.New()
	svc := &selfServiceTenantService{}
	handler := New(
		svc,
		func(context.Context) string { return tenantID.String() },
		func(context.Context) uuid.UUID { return actorID },
		func(context.Context) bool { return false },
		nil,
	)
	request := httptest.NewRequest(http.MethodGet, "/hrms/leaves", nil)
	request = request.WithContext(httputil.SetPermissions(request.Context(), []string{permissions.LeaveSelfView}))
	recorder := httptest.NewRecorder()

	handler.ListLeaves(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", recorder.Code, http.StatusOK, recorder.Body.String())
	}
	if !svc.listLeavesCalled || svc.gotUserID != actorID {
		t.Fatalf("service user = %s called=%t, want actor %s", svc.gotUserID, svc.listLeavesCalled, actorID)
	}
}

func TestEmployeeLeavePreviewDefaultsToActor(t *testing.T) {
	tenantID := uuid.New()
	actorID := uuid.New()
	leaveTypeID := uuid.New()
	svc := &selfServiceTenantService{}
	handler := New(
		svc,
		func(context.Context) string { return tenantID.String() },
		func(context.Context) uuid.UUID { return actorID },
		func(context.Context) bool { return false },
		nil,
	)
	body := []byte(`{"leave_type_id":"` + leaveTypeID.String() + `","start_date":"2026-07-10","end_date":"2026-07-10"}`)
	request := httptest.NewRequest(http.MethodPost, "/hrms/leaves/preview", bytes.NewReader(body))
	request = request.WithContext(httputil.SetPermissions(request.Context(), []string{permissions.LeaveSelfApply}))
	recorder := httptest.NewRecorder()

	handler.PreviewLeave(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", recorder.Code, http.StatusOK, recorder.Body.String())
	}
	if !svc.previewLeaveCalled || svc.gotUserID != actorID {
		t.Fatalf("preview user = %s called=%t, want actor %s", svc.gotUserID, svc.previewLeaveCalled, actorID)
	}
}

func TestEmployeeCannotCancelOtherUserLeave(t *testing.T) {
	tenantID := uuid.New()
	actorID := uuid.New()
	otherID := uuid.New()
	leaveID := uuid.New()
	svc := &selfServiceTenantService{}
	handler := New(
		svc,
		func(context.Context) string { return tenantID.String() },
		func(context.Context) uuid.UUID { return actorID },
		func(context.Context) bool { return false },
		nil,
	)
	body := []byte(`{"user_id":"` + otherID.String() + `"}`)
	request := httptest.NewRequest(http.MethodPost, "/hrms/leaves/"+leaveID.String()+"/cancel", bytes.NewReader(body))
	request = request.WithContext(httputil.SetPermissions(request.Context(), []string{permissions.LeaveSelfApply, permissions.LeavesApply}))
	recorder := httptest.NewRecorder()

	handler.cancelLeaveForTenant(recorder, request, tenantID, leaveID, "cancel leave")

	if recorder.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d; body=%s", recorder.Code, http.StatusForbidden, recorder.Body.String())
	}
	if svc.cancelLeaveCalled {
		t.Fatal("service should not be called for other user's leave cancel")
	}
}

func TestEmployeeCanCancelOwnLeave(t *testing.T) {
	tenantID := uuid.New()
	actorID := uuid.New()
	leaveID := uuid.New()
	svc := &selfServiceTenantService{}
	handler := New(
		svc,
		func(context.Context) string { return tenantID.String() },
		func(context.Context) uuid.UUID { return actorID },
		func(context.Context) bool { return false },
		nil,
	)
	request := httptest.NewRequest(http.MethodPost, "/hrms/leaves/"+leaveID.String()+"/cancel", bytes.NewReader([]byte(`{}`)))
	request = request.WithContext(httputil.SetPermissions(request.Context(), []string{permissions.LeaveSelfApply}))
	recorder := httptest.NewRecorder()

	handler.cancelLeaveForTenant(recorder, request, tenantID, leaveID, "cancel leave")

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", recorder.Code, http.StatusOK, recorder.Body.String())
	}
	if !svc.cancelLeaveCalled || svc.gotUserID != actorID {
		t.Fatalf("cancel user = %s called=%t, want actor %s", svc.gotUserID, svc.cancelLeaveCalled, actorID)
	}
}

func TestEmployeeCannotReadOtherApprovalQueue(t *testing.T) {
	tenantID := uuid.New()
	actorID := uuid.New()
	otherID := uuid.New()
	svc := &selfServiceTenantService{}
	handler := New(
		svc,
		func(context.Context) string { return tenantID.String() },
		func(context.Context) uuid.UUID { return actorID },
		func(context.Context) bool { return false },
		nil,
	)
	request := httptest.NewRequest(http.MethodGet, "/hrms/leave-approvals?approver_id="+otherID.String(), nil)
	request = request.WithContext(httputil.SetPermissions(request.Context(), []string{permissions.LeavesApprove}))
	recorder := httptest.NewRecorder()

	handler.ListLeaveApprovals(recorder, request)

	if recorder.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d; body=%s", recorder.Code, http.StatusForbidden, recorder.Body.String())
	}
	if svc.listApprovalsCalled {
		t.Fatal("service should not be called for other user's approval queue")
	}
}

func TestEmployeeAttendanceStatusDefaultsToActor(t *testing.T) {
	tenantID := uuid.New()
	actorID := uuid.New()
	svc := &selfServiceTenantService{}
	handler := New(
		svc,
		func(context.Context) string { return tenantID.String() },
		func(context.Context) uuid.UUID { return actorID },
		func(context.Context) bool { return false },
		nil,
	)
	request := httptest.NewRequest(http.MethodGet, "/hrms/attendances/status?date=2026-07-08", nil)
	request = request.WithContext(httputil.SetPermissions(request.Context(), []string{permissions.AttendanceSelfView}))
	recorder := httptest.NewRecorder()

	handler.ListAttendanceDailyStatuses(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", recorder.Code, http.StatusOK, recorder.Body.String())
	}
	if !svc.attendanceStatusCalled || svc.gotAttendanceUserID == nil || *svc.gotAttendanceUserID != actorID {
		t.Fatalf("attendance status user = %v called=%t, want actor %s", svc.gotAttendanceUserID, svc.attendanceStatusCalled, actorID)
	}
}

func TestEmployeeCannotReadOtherAttendanceStatus(t *testing.T) {
	tenantID := uuid.New()
	actorID := uuid.New()
	otherID := uuid.New()
	svc := &selfServiceTenantService{}
	handler := New(
		svc,
		func(context.Context) string { return tenantID.String() },
		func(context.Context) uuid.UUID { return actorID },
		func(context.Context) bool { return false },
		nil,
	)
	request := httptest.NewRequest(http.MethodGet, "/hrms/attendances/status?date=2026-07-08&user_id="+otherID.String(), nil)
	request = request.WithContext(httputil.SetPermissions(request.Context(), []string{permissions.AttendanceSelfView, permissions.AttendanceList}))
	recorder := httptest.NewRecorder()

	handler.ListAttendanceDailyStatuses(recorder, request)

	if recorder.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d; body=%s", recorder.Code, http.StatusForbidden, recorder.Body.String())
	}
	if svc.attendanceStatusCalled {
		t.Fatal("service should not be called for other user's attendance status")
	}
}

func TestEmployeeCanCreateOwnAttendanceSegment(t *testing.T) {
	tenantID := uuid.New()
	actorID := uuid.New()
	svc := &selfServiceTenantService{}
	handler := New(
		svc,
		func(context.Context) string { return tenantID.String() },
		func(context.Context) uuid.UUID { return actorID },
		func(context.Context) bool { return false },
		nil,
	)
	body := []byte(`{"segment_type":"client_site","action":"arrive","date":"2026-07-08","event_time":"2026-07-08T10:00:00Z"}`)
	request := httptest.NewRequest(http.MethodPost, "/hrms/attendances/segments", bytes.NewReader(body))
	request = request.WithContext(httputil.SetPermissions(request.Context(), []string{permissions.AttendanceSelfPunch}))
	recorder := httptest.NewRecorder()

	handler.CreateAttendanceWorkdaySegment(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d; body=%s", recorder.Code, http.StatusCreated, recorder.Body.String())
	}
	if !svc.createSegmentCalled || svc.gotUserID != actorID {
		t.Fatalf("segment user = %s called=%t, want actor %s", svc.gotUserID, svc.createSegmentCalled, actorID)
	}
}

func TestEmployeeCannotListOtherUserAttendanceSegments(t *testing.T) {
	tenantID := uuid.New()
	actorID := uuid.New()
	otherID := uuid.New()
	svc := &selfServiceTenantService{}
	handler := New(
		svc,
		func(context.Context) string { return tenantID.String() },
		func(context.Context) uuid.UUID { return actorID },
		func(context.Context) bool { return false },
		nil,
	)
	request := httptest.NewRequest(http.MethodGet, "/hrms/attendances/segments?date=2026-07-08&user_id="+otherID.String(), nil)
	request = request.WithContext(httputil.SetPermissions(request.Context(), []string{permissions.AttendanceSelfView}))
	recorder := httptest.NewRecorder()

	handler.ListAttendanceWorkdaySegments(recorder, request)

	if recorder.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d; body=%s", recorder.Code, http.StatusForbidden, recorder.Body.String())
	}
	if svc.listSegmentsCalled {
		t.Fatal("service should not be called for other user's attendance segments")
	}
}

func TestEmployeeCanCreateOwnCompOffRequest(t *testing.T) {
	tenantID := uuid.New()
	actorID := uuid.New()
	svc := &selfServiceTenantService{}
	handler := New(
		svc,
		func(context.Context) string { return tenantID.String() },
		func(context.Context) uuid.UUID { return actorID },
		func(context.Context) bool { return false },
		nil,
	)
	body := []byte(`{"work_date":"2026-07-08","worked_minutes":480,"requested_days":1,"reason":"Worked on weekly off"}`)
	request := httptest.NewRequest(http.MethodPost, "/hrms/comp-off-requests", bytes.NewReader(body))
	request = request.WithContext(httputil.SetPermissions(request.Context(), []string{permissions.LeaveSelfCompOffRequest}))
	recorder := httptest.NewRecorder()

	handler.CreateCompOffRequest(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d; body=%s", recorder.Code, http.StatusCreated, recorder.Body.String())
	}
	if !svc.createCompOffCalled || svc.gotUserID != actorID {
		t.Fatalf("comp-off user = %s called=%t, want actor %s", svc.gotUserID, svc.createCompOffCalled, actorID)
	}
}

func TestEmployeeCannotCreateOtherUserCompOffRequest(t *testing.T) {
	tenantID := uuid.New()
	actorID := uuid.New()
	otherID := uuid.New()
	svc := &selfServiceTenantService{}
	handler := New(
		svc,
		func(context.Context) string { return tenantID.String() },
		func(context.Context) uuid.UUID { return actorID },
		func(context.Context) bool { return false },
		nil,
	)
	body := []byte(`{"user_id":"` + otherID.String() + `","work_date":"2026-07-08","worked_minutes":480,"requested_days":1}`)
	request := httptest.NewRequest(http.MethodPost, "/hrms/comp-off-requests", bytes.NewReader(body))
	request = request.WithContext(httputil.SetPermissions(request.Context(), []string{permissions.LeaveSelfCompOffRequest}))
	recorder := httptest.NewRecorder()

	handler.CreateCompOffRequest(recorder, request)

	if recorder.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d; body=%s", recorder.Code, http.StatusForbidden, recorder.Body.String())
	}
	if svc.createCompOffCalled {
		t.Fatal("service should not be called for other user's comp-off request")
	}
}

func TestEmployeeCannotListOtherUserCompOffRequests(t *testing.T) {
	tenantID := uuid.New()
	actorID := uuid.New()
	otherID := uuid.New()
	svc := &selfServiceTenantService{}
	handler := New(
		svc,
		func(context.Context) string { return tenantID.String() },
		func(context.Context) uuid.UUID { return actorID },
		func(context.Context) bool { return false },
		nil,
	)
	request := httptest.NewRequest(http.MethodGet, "/hrms/comp-off-requests?user_id="+otherID.String(), nil)
	request = request.WithContext(httputil.SetPermissions(request.Context(), []string{permissions.LeaveSelfCompOffRequest}))
	recorder := httptest.NewRecorder()

	handler.ListCompOffRequests(recorder, request)

	if recorder.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d; body=%s", recorder.Code, http.StatusForbidden, recorder.Body.String())
	}
	if svc.listCompOffCalled {
		t.Fatal("service should not be called for other user's comp-off requests")
	}
}

func TestEmployeeCannotApproveCompOffRequest(t *testing.T) {
	tenantID := uuid.New()
	actorID := uuid.New()
	requestID := uuid.New()
	svc := &selfServiceTenantService{}
	handler := New(
		svc,
		func(context.Context) string { return tenantID.String() },
		func(context.Context) uuid.UUID { return actorID },
		func(context.Context) bool { return false },
		nil,
	)
	request := httptest.NewRequest(http.MethodPost, "/hrms/comp-off-requests/"+requestID.String()+"/approve", bytes.NewReader([]byte(`{}`)))
	request = request.WithContext(httputil.SetPermissions(request.Context(), []string{permissions.LeaveSelfCompOffRequest}))
	recorder := httptest.NewRecorder()

	handler.reviewCompOffRequestForTenant(recorder, request, tenantID, requestID, "approve comp-off request", domain.CompOffStatusApproved)

	if recorder.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d; body=%s", recorder.Code, http.StatusForbidden, recorder.Body.String())
	}
	if svc.reviewCompOffCalled {
		t.Fatal("service should not be called for employee comp-off approval")
	}
}
