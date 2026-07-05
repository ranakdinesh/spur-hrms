package postgres

import (
	"math"
	"math/big"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func mapAttendance(row sqlc.HrmsAttendance) *domain.Attendance {
	return &domain.Attendance{
		ID:                         row.ID,
		TenantID:                   row.TenantID,
		UserID:                     row.UserID,
		Date:                       timeFromDate(row.Date),
		Time:                       ptrFromTimestamptz(row.Time),
		Type:                       ptrFromText(row.Type),
		Status:                     ptrFromText(row.Status),
		Source:                     ptrFromText(row.Source),
		Latitude:                   ptrFromNumeric(row.Latitude),
		Longitude:                  ptrFromNumeric(row.Longitude),
		WorkMode:                   ptrFromText(row.WorkMode),
		Remarks:                    ptrFromText(row.Remarks),
		AttendanceLocationID:       ptrFromUUID(row.AttendanceLocationID),
		AttendanceDeviceID:         ptrFromUUID(row.AttendanceDeviceID),
		RawAttendanceEventID:       ptrFromUUID(row.RawAttendanceEventID),
		LocationAccuracyMeters:     ptrFromNumeric(row.LocationAccuracyMeters),
		LocationVerificationStatus: row.LocationVerificationStatus,
		Inactive:                   row.Inactive,
		CreatedAt:                  timeFromTimestamptz(row.CreatedAt),
		CreatedBy:                  ptrFromUUID(row.CreatedBy),
		UpdatedAt:                  timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:                  ptrFromUUID(row.UpdatedBy),
	}
}

func mapAttendances(rows []sqlc.HrmsAttendance) []*domain.Attendance {
	items := make([]*domain.Attendance, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapAttendance(row))
	}
	return items
}

func mapDeviceLog(row sqlc.HrmsDeviceLog) *domain.DeviceLog {
	return &domain.DeviceLog{
		ID:         row.ID,
		TenantID:   row.TenantID,
		UserID:     row.UserID,
		DeviceID:   ptrFromText(row.DeviceID),
		DeviceType: ptrFromText(row.DeviceType),
		IPAddress:  ptrFromText(row.IpAddress),
		Action:     ptrFromText(row.Action),
		Inactive:   row.Inactive,
		LoggedAt:   timeFromTimestamptz(row.LoggedAt),
		CreatedAt:  timeFromTimestamptz(row.CreatedAt),
		CreatedBy:  ptrFromUUID(row.CreatedBy),
		UpdatedAt:  timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:  ptrFromUUID(row.UpdatedBy),
	}
}

func mapDeviceLogs(rows []sqlc.HrmsDeviceLog) []*domain.DeviceLog {
	items := make([]*domain.DeviceLog, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapDeviceLog(row))
	}
	return items
}

func numericFromCoordinatePtr(value *float64) pgtype.Numeric {
	if value == nil {
		return pgtype.Numeric{Valid: false}
	}
	scaled := int64(math.Round(*value * 10000000))
	return pgtype.Numeric{Int: big.NewInt(scaled), Exp: -7, Valid: true}
}

func ptrFromNumeric(value pgtype.Numeric) *float64 {
	if !value.Valid {
		return nil
	}
	floatValue, err := value.Float64Value()
	if err != nil || !floatValue.Valid {
		return nil
	}
	return &floatValue.Float64
}
