package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (s *Store) CreateAttendance(ctx context.Context, item *domain.Attendance, actorID *uuid.UUID) (*domain.Attendance, error) {
	row, err := s.getQueries(ctx).CreateAttendance(ctx, sqlc.CreateAttendanceParams{TenantID: item.TenantID, UserID: item.UserID, Date: dateFromTime(item.Date), Time: timestamptzFromPtr(item.Time), Type: textFromPtr(item.Type), Status: textFromPtr(item.Status), Source: textFromPtr(item.Source), Latitude: numericFromCoordinatePtr(item.Latitude), Longitude: numericFromCoordinatePtr(item.Longitude), WorkMode: textFromPtr(item.WorkMode), Remarks: textFromPtr(item.Remarks), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create attendance", err, tenantIDField(item.TenantID), stringField("user_id", item.UserID.String()))
	}
	return mapAttendance(row), nil
}

func (s *Store) UpdateAttendance(ctx context.Context, item *domain.Attendance, actorID *uuid.UUID) (*domain.Attendance, error) {
	row, err := s.getQueries(ctx).UpdateAttendance(ctx, sqlc.UpdateAttendanceParams{TenantID: item.TenantID, ID: item.ID, Time: timestamptzFromPtr(item.Time), Type: textFromPtr(item.Type), Status: textFromPtr(item.Status), Source: textFromPtr(item.Source), Latitude: numericFromCoordinatePtr(item.Latitude), Longitude: numericFromCoordinatePtr(item.Longitude), WorkMode: textFromPtr(item.WorkMode), Remarks: textFromPtr(item.Remarks), UpdatedBy: uuidFromPtr(actorID)})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrAttendanceNotFound
		}
		return nil, s.logDBError(ctx, "update attendance", err, tenantIDField(item.TenantID), stringField("attendance_id", item.ID.String()))
	}
	return mapAttendance(row), nil
}

func (s *Store) ListAttendancesByUser(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID, startDate string, endDate string) ([]*domain.Attendance, error) {
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, err
	}
	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, err
	}
	rows, err := s.getQueries(ctx).ListAttendancesByUser(ctx, sqlc.ListAttendancesByUserParams{TenantID: tenantID, UserID: userID, Date: dateFromTime(start), Date_2: dateFromTime(end)})
	if err != nil {
		return nil, s.logDBError(ctx, "list attendances by user", err, tenantIDField(tenantID), stringField("user_id", userID.String()))
	}
	return mapAttendances(rows), nil
}

func (s *Store) ListAttendancesByUserDate(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID, date string) ([]*domain.Attendance, error) {
	parsed, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, err
	}
	rows, err := s.getQueries(ctx).ListAttendancesByUserDate(ctx, sqlc.ListAttendancesByUserDateParams{TenantID: tenantID, UserID: userID, Date: dateFromTime(parsed)})
	if err != nil {
		return nil, s.logDBError(ctx, "list attendances by user date", err, tenantIDField(tenantID), stringField("user_id", userID.String()))
	}
	return mapAttendances(rows), nil
}

func (s *Store) ListAttendancesByDate(ctx context.Context, tenantID uuid.UUID, date string) ([]*domain.Attendance, error) {
	parsed, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, err
	}
	rows, err := s.getQueries(ctx).ListAttendancesByDate(ctx, sqlc.ListAttendancesByDateParams{TenantID: tenantID, Date: dateFromTime(parsed)})
	if err != nil {
		return nil, s.logDBError(ctx, "list attendances by date", err, tenantIDField(tenantID))
	}
	return mapAttendances(rows), nil
}

func (s *Store) GetAttendance(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.Attendance, error) {
	row, err := s.getQueries(ctx).GetAttendance(ctx, sqlc.GetAttendanceParams{TenantID: tenantID, ID: id})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrAttendanceNotFound
		}
		return nil, s.logDBError(ctx, "get attendance", err, tenantIDField(tenantID), stringField("attendance_id", id.String()))
	}
	return mapAttendance(row), nil
}

func (s *Store) CreateAttendanceWorkdaySegment(ctx context.Context, item *domain.AttendanceWorkdaySegment, actorID *uuid.UUID) (*domain.AttendanceWorkdaySegment, error) {
	row, err := s.getQueries(ctx).CreateAttendanceWorkdaySegment(ctx, sqlc.CreateAttendanceWorkdaySegmentParams{
		TenantID:                   item.TenantID,
		UserID:                     item.UserID,
		Date:                       dateFromTime(item.Date),
		EventTime:                  timestamptzFromPtr(&item.EventTime),
		SegmentType:                item.SegmentType,
		Action:                     item.Action,
		WorkMode:                   textFromPtr(item.WorkMode),
		Source:                     textFromPtr(item.Source),
		AttendanceLocationID:       uuidFromPtr(item.AttendanceLocationID),
		ReferenceType:              textFromPtr(item.ReferenceType),
		ReferenceID:                uuidFromPtr(item.ReferenceID),
		ReferenceLabel:             textFromPtr(item.ReferenceLabel),
		Latitude:                   numericFromCoordinatePtr(item.Latitude),
		Longitude:                  numericFromCoordinatePtr(item.Longitude),
		LocationAccuracyMeters:     numericFromFloatPtr(item.LocationAccuracyMeters),
		LocationVerificationStatus: item.LocationVerificationStatus,
		Remarks:                    textFromPtr(item.Remarks),
		Metadata:                   jsonBytesFromRaw(item.Metadata),
		CreatedBy:                  uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "create attendance workday segment", err, tenantIDField(item.TenantID), stringField("user_id", item.UserID.String()))
	}
	return mapAttendanceWorkdaySegment(row), nil
}

func (s *Store) ListAttendanceWorkdaySegmentsByUser(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID, startDate string, endDate string) ([]*domain.AttendanceWorkdaySegment, error) {
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, err
	}
	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, err
	}
	rows, err := s.getQueries(ctx).ListAttendanceWorkdaySegmentsByUser(ctx, sqlc.ListAttendanceWorkdaySegmentsByUserParams{TenantID: tenantID, UserID: userID, Date: dateFromTime(start), Date_2: dateFromTime(end)})
	if err != nil {
		return nil, s.logDBError(ctx, "list attendance workday segments by user", err, tenantIDField(tenantID), stringField("user_id", userID.String()))
	}
	return mapAttendanceWorkdaySegments(rows), nil
}

func (s *Store) ListAttendanceWorkdaySegmentsByUserDate(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID, date string) ([]*domain.AttendanceWorkdaySegment, error) {
	parsed, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, err
	}
	rows, err := s.getQueries(ctx).ListAttendanceWorkdaySegmentsByUserDate(ctx, sqlc.ListAttendanceWorkdaySegmentsByUserDateParams{TenantID: tenantID, UserID: userID, Date: dateFromTime(parsed)})
	if err != nil {
		return nil, s.logDBError(ctx, "list attendance workday segments by user date", err, tenantIDField(tenantID), stringField("user_id", userID.String()))
	}
	return mapAttendanceWorkdaySegments(rows), nil
}

func (s *Store) CreateDeviceLog(ctx context.Context, item *domain.DeviceLog, actorID *uuid.UUID) (*domain.DeviceLog, error) {
	row, err := s.getQueries(ctx).CreateDeviceLog(ctx, sqlc.CreateDeviceLogParams{TenantID: item.TenantID, UserID: item.UserID, DeviceID: textFromPtr(item.DeviceID), DeviceType: textFromPtr(item.DeviceType), IpAddress: textFromPtr(item.IPAddress), Action: textFromPtr(item.Action), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create device log", err, tenantIDField(item.TenantID), stringField("user_id", item.UserID.String()))
	}
	return mapDeviceLog(row), nil
}

func (s *Store) ListDeviceLogsByUser(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID) ([]*domain.DeviceLog, error) {
	rows, err := s.getQueries(ctx).ListDeviceLogsByUser(ctx, sqlc.ListDeviceLogsByUserParams{TenantID: tenantID, UserID: userID})
	if err != nil {
		return nil, s.logDBError(ctx, "list device logs by user", err, tenantIDField(tenantID), stringField("user_id", userID.String()))
	}
	return mapDeviceLogs(rows), nil
}
