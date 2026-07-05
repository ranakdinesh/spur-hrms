package postgres

import (
	"encoding/json"

	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func mapAttendanceLocation(row sqlc.HrmsAttendanceLocation) *domain.AttendanceLocation {
	return &domain.AttendanceLocation{ID: row.ID, TenantID: row.TenantID, BranchID: ptrFromUUID(row.BranchID), Code: row.Code, Name: row.Name, LocationType: row.LocationType, Latitude: ptrFromNumeric(row.Latitude), Longitude: ptrFromNumeric(row.Longitude), RadiusMeters: row.RadiusMeters, Address: ptrFromText(row.Address), City: ptrFromText(row.City), State: ptrFromText(row.State), Country: ptrFromText(row.Country), Pincode: ptrFromText(row.Pincode), EffectiveFrom: ptrFromDate(row.EffectiveFrom), EffectiveTo: ptrFromDate(row.EffectiveTo), Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)}
}

func mapAttendanceLocations(rows []sqlc.HrmsAttendanceLocation) []*domain.AttendanceLocation {
	items := make([]*domain.AttendanceLocation, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapAttendanceLocation(row))
	}
	return items
}

func mapAttendanceLocationAssignment(row sqlc.HrmsAttendanceLocationAssignment) *domain.AttendanceLocationAssignment {
	return &domain.AttendanceLocationAssignment{ID: row.ID, TenantID: row.TenantID, LocationID: row.LocationID, BranchID: ptrFromUUID(row.BranchID), DepartmentID: ptrFromUUID(row.DepartmentID), UserID: ptrFromUUID(row.UserID), EffectiveFrom: ptrFromDate(row.EffectiveFrom), EffectiveTo: ptrFromDate(row.EffectiveTo), Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)}
}

func mapAttendanceLocationAssignments(rows []sqlc.HrmsAttendanceLocationAssignment) []*domain.AttendanceLocationAssignment {
	items := make([]*domain.AttendanceLocationAssignment, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapAttendanceLocationAssignment(row))
	}
	return items
}

func mapAttendanceDevice(row sqlc.HrmsAttendanceDevice) *domain.AttendanceDevice {
	config := json.RawMessage(row.Config)
	if len(config) == 0 {
		config = json.RawMessage(`{}`)
	}
	return &domain.AttendanceDevice{ID: row.ID, TenantID: row.TenantID, AttendanceLocationID: ptrFromUUID(row.AttendanceLocationID), BranchID: ptrFromUUID(row.BranchID), Code: row.Code, Name: row.Name, Vendor: ptrFromText(row.Vendor), Model: ptrFromText(row.Model), SerialNumber: ptrFromText(row.SerialNumber), DeviceIdentifier: ptrFromText(row.DeviceIdentifier), IntegrationType: row.IntegrationType, DirectionMode: row.DirectionMode, Timezone: row.Timezone, Status: row.Status, LastSeenAt: ptrFromTimestamptz(row.LastSeenAt), Config: config, Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)}
}

func mapAttendanceDevices(rows []sqlc.HrmsAttendanceDevice) []*domain.AttendanceDevice {
	items := make([]*domain.AttendanceDevice, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapAttendanceDevice(row))
	}
	return items
}

func mapEmployeeAttendanceDevice(row sqlc.HrmsEmployeeAttendanceDevice) *domain.EmployeeAttendanceDevice {
	return &domain.EmployeeAttendanceDevice{ID: row.ID, TenantID: row.TenantID, UserID: row.UserID, DeviceID: row.DeviceID, DeviceUserID: row.DeviceUserID, CredentialType: row.CredentialType, CardNumber: ptrFromText(row.CardNumber), EffectiveFrom: ptrFromDate(row.EffectiveFrom), EffectiveTo: ptrFromDate(row.EffectiveTo), Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)}
}

func mapEmployeeAttendanceDevices(rows []sqlc.HrmsEmployeeAttendanceDevice) []*domain.EmployeeAttendanceDevice {
	items := make([]*domain.EmployeeAttendanceDevice, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapEmployeeAttendanceDevice(row))
	}
	return items
}

func mapRawAttendanceEvent(row sqlc.HrmsRawAttendanceEvent) *domain.RawAttendanceEvent {
	payload := json.RawMessage(row.RawPayload)
	if len(payload) == 0 {
		payload = json.RawMessage(`{}`)
	}
	return &domain.RawAttendanceEvent{ID: row.ID, TenantID: row.TenantID, DeviceID: row.DeviceID, EmployeeDeviceMappingID: ptrFromUUID(row.EmployeeDeviceMappingID), AttendanceID: ptrFromUUID(row.AttendanceID), ExternalEventID: ptrFromText(row.ExternalEventID), DeviceUserID: ptrFromText(row.DeviceUserID), EventTime: timeFromTimestamptz(row.EventTime), EventType: ptrFromText(row.EventType), AttendanceType: ptrFromText(row.AttendanceType), ImportBatchID: ptrFromText(row.ImportBatchID), ProcessingStatus: row.ProcessingStatus, ProcessingError: ptrFromText(row.ProcessingError), RawPayload: payload, ProcessedAt: ptrFromTimestamptz(row.ProcessedAt), Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)}
}

func mapRawAttendanceEvents(rows []sqlc.HrmsRawAttendanceEvent) []*domain.RawAttendanceEvent {
	items := make([]*domain.RawAttendanceEvent, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapRawAttendanceEvent(row))
	}
	return items
}
