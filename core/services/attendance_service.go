package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) PunchAttendance(ctx context.Context, cmd ports.AttendancePunchCommand) (*domain.AttendancePunch, error) {
	if cmd.TenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate attendance punch tenant", err)
		return nil, err
	}
	if cmd.UserID == uuid.Nil {
		err := domain.ErrInvalidLeaveUser
		s.logError("validate attendance punch user", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	if _, err := s.employees.GetEmployeeByUserID(ctx, cmd.TenantID, cmd.UserID); err != nil {
		s.logError("validate attendance punch employee", err, serviceTenantIDField(cmd.TenantID), serviceStringField("user_id", cmd.UserID.String()))
		return nil, err
	}
	attendanceRequired, err := s.employees.GetEmployeeAttendanceRequired(ctx, cmd.TenantID, cmd.UserID)
	if err != nil {
		s.logError("validate attendance punch requirement", err, serviceTenantIDField(cmd.TenantID), serviceStringField("user_id", cmd.UserID.String()))
		return nil, err
	}
	if !attendanceRequired {
		err := domain.ErrAttendanceNotRequired
		s.logError("validate attendance punch requirement", err, serviceTenantIDField(cmd.TenantID), serviceStringField("user_id", cmd.UserID.String()))
		return nil, err
	}
	punchTime, err := parseAttendanceTime(cmd.Time)
	if err != nil {
		s.logError("parse attendance punch time", err, serviceTenantIDField(cmd.TenantID), serviceStringField("user_id", cmd.UserID.String()))
		return nil, err
	}
	punchDate, err := parseAttendanceDate(cmd.Date, punchTime)
	if err != nil {
		s.logError("parse attendance punch date", err, serviceTenantIDField(cmd.TenantID), serviceStringField("user_id", cmd.UserID.String()))
		return nil, err
	}
	if cmd.Source == nil || *cmd.Source == "" {
		cmd.Source = strPtr(domain.AttendanceSourceWeb)
	}
	if cmd.WorkMode == nil || *cmd.WorkMode == "" {
		cmd.WorkMode = strPtr(domain.AttendanceWorkModeOffice)
	}
	if attendanceSourceRequiresLocation(cmd.Source) && (cmd.Latitude == nil || cmd.Longitude == nil) {
		err := domain.ErrAttendanceLocationRequired
		s.logError("validate attendance punch location", err, serviceTenantIDField(cmd.TenantID), serviceStringField("user_id", cmd.UserID.String()))
		return nil, err
	}
	attendance, err := domain.NewAttendance(cmd.TenantID, cmd.UserID, punchDate, punchTime, cmd.Action, cmd.Source, cmd.Latitude, cmd.Longitude, cmd.WorkMode, cmd.Remarks)
	if err != nil {
		s.logError("validate attendance punch", err, serviceTenantIDField(cmd.TenantID), serviceStringField("user_id", cmd.UserID.String()))
		return nil, err
	}
	deviceLog, err := domain.NewDeviceLog(cmd.TenantID, cmd.UserID, cmd.DeviceID, cmd.DeviceType, cmd.IPAddress, cmd.Action)
	if err != nil {
		s.logError("validate attendance device log", err, serviceTenantIDField(cmd.TenantID), serviceStringField("user_id", cmd.UserID.String()))
		return nil, err
	}
	dateKey := attendance.Date.Format("2006-01-02")
	existing, err := s.attendances.ListAttendancesByUserDate(ctx, cmd.TenantID, cmd.UserID, dateKey)
	if err != nil {
		s.logError("list attendance punch day events", err, serviceTenantIDField(cmd.TenantID), serviceStringField("user_id", cmd.UserID.String()))
		return nil, err
	}
	if err := validateAttendancePunchSequence(*attendance.Type, existing); err != nil {
		s.logError("validate attendance punch sequence", err, serviceTenantIDField(cmd.TenantID), serviceStringField("user_id", cmd.UserID.String()))
		return nil, err
	}
	var result *domain.AttendancePunch
	err = s.system.RunAsSystem(ctx, func(txCtx context.Context) error {
		createdAttendance, err := s.attendances.CreateAttendance(txCtx, attendance, cmd.ActorID)
		if err != nil {
			return err
		}
		createdDeviceLog, err := s.attendances.CreateDeviceLog(txCtx, deviceLog, cmd.ActorID)
		if err != nil {
			return err
		}
		result = &domain.AttendancePunch{Attendance: createdAttendance, DeviceLog: createdDeviceLog}
		return nil
	})
	if err != nil {
		s.logError("punch attendance", err, serviceTenantIDField(cmd.TenantID), serviceStringField("user_id", cmd.UserID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) ListAttendancesByUser(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID, startDate string, endDate string) ([]*domain.Attendance, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate attendance list tenant", err)
		return nil, err
	}
	if userID == uuid.Nil {
		err := domain.ErrInvalidLeaveUser
		s.logError("validate attendance list user", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	start, end, err := parseDateRangeOrToday(startDate, endDate)
	if err != nil {
		s.logError("validate attendance list date range", err, serviceTenantIDField(tenantID), serviceStringField("user_id", userID.String()))
		return nil, err
	}
	items, err := s.attendances.ListAttendancesByUser(ctx, tenantID, userID, start.Format("2006-01-02"), end.Format("2006-01-02"))
	if err != nil {
		s.logError("list attendances by user", err, serviceTenantIDField(tenantID), serviceStringField("user_id", userID.String()))
		return nil, err
	}
	return items, nil
}

func (s *TenantService) ListAttendancesByUserDate(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID, date string) ([]*domain.Attendance, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate attendance day tenant", err)
		return nil, err
	}
	if userID == uuid.Nil {
		err := domain.ErrInvalidLeaveUser
		s.logError("validate attendance day user", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	if date == "" {
		date = time.Now().UTC().Format("2006-01-02")
	}
	if _, err := time.Parse("2006-01-02", date); err != nil {
		s.logError("validate attendance day date", domain.ErrInvalidAttendanceDate, serviceTenantIDField(tenantID), serviceStringField("user_id", userID.String()))
		return nil, domain.ErrInvalidAttendanceDate
	}
	items, err := s.attendances.ListAttendancesByUserDate(ctx, tenantID, userID, date)
	if err != nil {
		s.logError("list attendances by user date", err, serviceTenantIDField(tenantID), serviceStringField("user_id", userID.String()))
		return nil, err
	}
	return items, nil
}

func (s *TenantService) CreateAttendanceWorkdaySegment(ctx context.Context, cmd ports.AttendanceSegmentCommand) (*domain.AttendanceWorkdaySegment, error) {
	if cmd.TenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate attendance segment tenant", err)
		return nil, err
	}
	if cmd.UserID == uuid.Nil {
		err := domain.ErrInvalidLeaveUser
		s.logError("validate attendance segment user", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	if _, err := s.employees.GetEmployeeByUserID(ctx, cmd.TenantID, cmd.UserID); err != nil {
		s.logError("validate attendance segment employee", err, serviceTenantIDField(cmd.TenantID), serviceStringField("user_id", cmd.UserID.String()))
		return nil, err
	}
	attendanceRequired, err := s.employees.GetEmployeeAttendanceRequired(ctx, cmd.TenantID, cmd.UserID)
	if err != nil {
		s.logError("validate attendance segment requirement", err, serviceTenantIDField(cmd.TenantID), serviceStringField("user_id", cmd.UserID.String()))
		return nil, err
	}
	if !attendanceRequired {
		err := domain.ErrAttendanceNotRequired
		s.logError("validate attendance segment requirement", err, serviceTenantIDField(cmd.TenantID), serviceStringField("user_id", cmd.UserID.String()))
		return nil, err
	}
	eventTime, err := parseAttendanceTime(cmd.EventTime)
	if err != nil {
		s.logError("parse attendance segment event time", err, serviceTenantIDField(cmd.TenantID), serviceStringField("user_id", cmd.UserID.String()))
		return nil, err
	}
	segmentDate, err := parseAttendanceDate(cmd.Date, eventTime)
	if err != nil {
		s.logError("parse attendance segment date", err, serviceTenantIDField(cmd.TenantID), serviceStringField("user_id", cmd.UserID.String()))
		return nil, err
	}
	if cmd.Source == nil || *cmd.Source == "" {
		cmd.Source = strPtr(domain.AttendanceSourceWeb)
	}
	if cmd.WorkMode == nil || *cmd.WorkMode == "" {
		cmd.WorkMode = strPtr(domain.AttendanceWorkModeField)
	}
	item, err := domain.NewAttendanceWorkdaySegment(domain.AttendanceWorkdaySegment{
		TenantID:                   cmd.TenantID,
		UserID:                     cmd.UserID,
		Date:                       segmentDate,
		EventTime:                  eventTime,
		SegmentType:                cmd.SegmentType,
		Action:                     cmd.Action,
		WorkMode:                   cmd.WorkMode,
		Source:                     cmd.Source,
		AttendanceLocationID:       cmd.AttendanceLocationID,
		ReferenceType:              cmd.ReferenceType,
		ReferenceID:                cmd.ReferenceID,
		ReferenceLabel:             cmd.ReferenceLabel,
		Latitude:                   cmd.Latitude,
		Longitude:                  cmd.Longitude,
		LocationAccuracyMeters:     cmd.LocationAccuracyMeters,
		LocationVerificationStatus: cmd.LocationVerificationStatus,
		Remarks:                    cmd.Remarks,
		Metadata:                   cmd.Metadata,
	})
	if err != nil {
		s.logError("validate attendance segment", err, serviceTenantIDField(cmd.TenantID), serviceStringField("user_id", cmd.UserID.String()))
		return nil, err
	}
	var result *domain.AttendanceWorkdaySegment
	err = s.system.RunAsSystem(ctx, func(txCtx context.Context) error {
		created, err := s.attendances.CreateAttendanceWorkdaySegment(txCtx, item, cmd.ActorID)
		if err != nil {
			return err
		}
		result = created
		return nil
	})
	if err != nil {
		s.logError("create attendance segment", err, serviceTenantIDField(cmd.TenantID), serviceStringField("user_id", cmd.UserID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) ListAttendanceWorkdaySegmentsByUser(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID, startDate string, endDate string) ([]*domain.AttendanceWorkdaySegment, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate attendance segment list tenant", err)
		return nil, err
	}
	if userID == uuid.Nil {
		err := domain.ErrInvalidLeaveUser
		s.logError("validate attendance segment list user", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	start, end, err := parseDateRangeOrToday(startDate, endDate)
	if err != nil {
		s.logError("validate attendance segment list date range", err, serviceTenantIDField(tenantID), serviceStringField("user_id", userID.String()))
		return nil, err
	}
	items, err := s.attendances.ListAttendanceWorkdaySegmentsByUser(ctx, tenantID, userID, start.Format("2006-01-02"), end.Format("2006-01-02"))
	if err != nil {
		s.logError("list attendance segments by user", err, serviceTenantIDField(tenantID), serviceStringField("user_id", userID.String()))
		return nil, err
	}
	return items, nil
}

func (s *TenantService) ListAttendanceWorkdaySegmentsByUserDate(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID, date string) ([]*domain.AttendanceWorkdaySegment, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate attendance segment day tenant", err)
		return nil, err
	}
	if userID == uuid.Nil {
		err := domain.ErrInvalidLeaveUser
		s.logError("validate attendance segment day user", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	if date == "" {
		date = time.Now().UTC().Format("2006-01-02")
	}
	if _, err := time.Parse("2006-01-02", date); err != nil {
		s.logError("validate attendance segment day date", domain.ErrInvalidAttendanceDate, serviceTenantIDField(tenantID), serviceStringField("user_id", userID.String()))
		return nil, domain.ErrInvalidAttendanceDate
	}
	items, err := s.attendances.ListAttendanceWorkdaySegmentsByUserDate(ctx, tenantID, userID, date)
	if err != nil {
		s.logError("list attendance segments by user date", err, serviceTenantIDField(tenantID), serviceStringField("user_id", userID.String()))
		return nil, err
	}
	return items, nil
}

func parseAttendanceTime(value string) (time.Time, error) {
	if value == "" {
		return time.Now().UTC(), nil
	}
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return time.Time{}, domain.ErrInvalidAttendanceDate
	}
	return parsed.UTC(), nil
}

func parseAttendanceDate(value string, fallback time.Time) (time.Time, error) {
	if value == "" {
		return fallback.UTC(), nil
	}
	parsed, err := time.Parse("2006-01-02", value)
	if err != nil {
		return time.Time{}, domain.ErrInvalidAttendanceDate
	}
	return parsed, nil
}

func parseDateRangeOrToday(startDate string, endDate string) (time.Time, time.Time, error) {
	if startDate == "" {
		startDate = time.Now().UTC().Format("2006-01-02")
	}
	if endDate == "" {
		endDate = startDate
	}
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return time.Time{}, time.Time{}, domain.ErrInvalidDateRange
	}
	end, err := time.Parse("2006-01-02", endDate)
	if err != nil || end.Before(start) {
		return time.Time{}, time.Time{}, domain.ErrInvalidDateRange
	}
	return start, end, nil
}

func validateAttendancePunchSequence(action string, items []*domain.Attendance) error {
	checkins := 0
	checkouts := 0
	for _, item := range items {
		if item == nil || item.Type == nil {
			continue
		}
		switch *item.Type {
		case domain.AttendanceCheckin:
			checkins++
		case domain.AttendanceCheckout:
			checkouts++
		}
	}
	switch action {
	case domain.AttendanceCheckin:
		if checkins > checkouts {
			return domain.ErrAttendanceAlreadyCheckedIn
		}
	case domain.AttendanceCheckout:
		if checkins == 0 {
			return domain.ErrAttendanceNotCheckedIn
		}
		if checkouts >= checkins {
			return domain.ErrAttendanceAlreadyCheckedOut
		}
	default:
		return domain.ErrInvalidAttendanceAction
	}
	return nil
}

func attendanceSourceRequiresLocation(source *string) bool {
	if source == nil {
		return true
	}
	switch *source {
	case domain.AttendanceSourceWeb, domain.AttendanceSourceMobile, domain.AttendanceSourceKiosk:
		return true
	default:
		return false
	}
}
