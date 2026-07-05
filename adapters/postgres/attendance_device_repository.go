package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (s *Store) CreateAttendanceLocation(ctx context.Context, item *domain.AttendanceLocation, actorID *uuid.UUID) (*domain.AttendanceLocation, error) {
	row, err := s.getQueries(ctx).CreateAttendanceLocation(ctx, sqlc.CreateAttendanceLocationParams{TenantID: item.TenantID, BranchID: uuidFromPtr(item.BranchID), Code: item.Code, Name: item.Name, LocationType: item.LocationType, Latitude: numericFromCoordinatePtr(item.Latitude), Longitude: numericFromCoordinatePtr(item.Longitude), RadiusMeters: item.RadiusMeters, Address: textFromPtr(item.Address), City: textFromPtr(item.City), State: textFromPtr(item.State), Country: textFromPtr(item.Country), Pincode: textFromPtr(item.Pincode), EffectiveFrom: dateFromPtr(item.EffectiveFrom), EffectiveTo: dateFromPtr(item.EffectiveTo), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create attendance location", err, tenantIDField(item.TenantID), stringField("code", item.Code))
	}
	return mapAttendanceLocation(row), nil
}

func (s *Store) ListAttendanceLocations(ctx context.Context, tenantID uuid.UUID) ([]*domain.AttendanceLocation, error) {
	rows, err := s.getQueries(ctx).ListAttendanceLocations(ctx, tenantID)
	if err != nil {
		return nil, s.logDBError(ctx, "list attendance locations", err, tenantIDField(tenantID))
	}
	return mapAttendanceLocations(rows), nil
}

func (s *Store) GetAttendanceLocation(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.AttendanceLocation, error) {
	row, err := s.getQueries(ctx).GetAttendanceLocation(ctx, sqlc.GetAttendanceLocationParams{TenantID: tenantID, ID: id})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrAttendanceLocationNotFound
		}
		return nil, s.logDBError(ctx, "get attendance location", err, tenantIDField(tenantID), stringField("location_id", id.String()))
	}
	return mapAttendanceLocation(row), nil
}

func (s *Store) UpdateAttendanceLocation(ctx context.Context, item *domain.AttendanceLocation, actorID *uuid.UUID) (*domain.AttendanceLocation, error) {
	row, err := s.getQueries(ctx).UpdateAttendanceLocation(ctx, sqlc.UpdateAttendanceLocationParams{TenantID: item.TenantID, ID: item.ID, BranchID: uuidFromPtr(item.BranchID), Code: item.Code, Name: item.Name, LocationType: item.LocationType, Latitude: numericFromCoordinatePtr(item.Latitude), Longitude: numericFromCoordinatePtr(item.Longitude), RadiusMeters: item.RadiusMeters, Address: textFromPtr(item.Address), City: textFromPtr(item.City), State: textFromPtr(item.State), Country: textFromPtr(item.Country), Pincode: textFromPtr(item.Pincode), EffectiveFrom: dateFromPtr(item.EffectiveFrom), EffectiveTo: dateFromPtr(item.EffectiveTo), UpdatedBy: uuidFromPtr(actorID)})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrAttendanceLocationNotFound
		}
		return nil, s.logDBError(ctx, "update attendance location", err, tenantIDField(item.TenantID), stringField("location_id", item.ID.String()))
	}
	return mapAttendanceLocation(row), nil
}

func (s *Store) DeleteAttendanceLocation(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteAttendanceLocation(ctx, sqlc.SoftDeleteAttendanceLocationParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete attendance location", err, tenantIDField(tenantID), stringField("location_id", id.String()))
	}
	return nil
}

func (s *Store) CreateAttendanceLocationAssignment(ctx context.Context, item *domain.AttendanceLocationAssignment, actorID *uuid.UUID) (*domain.AttendanceLocationAssignment, error) {
	row, err := s.getQueries(ctx).CreateAttendanceLocationAssignment(ctx, sqlc.CreateAttendanceLocationAssignmentParams{TenantID: item.TenantID, LocationID: item.LocationID, BranchID: uuidFromPtr(item.BranchID), DepartmentID: uuidFromPtr(item.DepartmentID), UserID: uuidFromPtr(item.UserID), EffectiveFrom: dateFromPtr(item.EffectiveFrom), EffectiveTo: dateFromPtr(item.EffectiveTo), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create attendance location assignment", err, tenantIDField(item.TenantID), stringField("location_id", item.LocationID.String()))
	}
	return mapAttendanceLocationAssignment(row), nil
}

func (s *Store) ListAttendanceLocationAssignments(ctx context.Context, tenantID uuid.UUID, locationID *uuid.UUID) ([]*domain.AttendanceLocationAssignment, error) {
	if locationID != nil && *locationID != uuid.Nil {
		rows, err := s.getQueries(ctx).ListAttendanceLocationAssignmentsByLocation(ctx, sqlc.ListAttendanceLocationAssignmentsByLocationParams{TenantID: tenantID, LocationID: *locationID})
		if err != nil {
			return nil, s.logDBError(ctx, "list attendance location assignments by location", err, tenantIDField(tenantID), stringField("location_id", locationID.String()))
		}
		return mapAttendanceLocationAssignments(rows), nil
	}
	rows, err := s.getQueries(ctx).ListAttendanceLocationAssignments(ctx, tenantID)
	if err != nil {
		return nil, s.logDBError(ctx, "list attendance location assignments", err, tenantIDField(tenantID))
	}
	return mapAttendanceLocationAssignments(rows), nil
}

func (s *Store) UpdateAttendanceLocationAssignment(ctx context.Context, item *domain.AttendanceLocationAssignment, actorID *uuid.UUID) (*domain.AttendanceLocationAssignment, error) {
	row, err := s.getQueries(ctx).UpdateAttendanceLocationAssignment(ctx, sqlc.UpdateAttendanceLocationAssignmentParams{TenantID: item.TenantID, ID: item.ID, LocationID: item.LocationID, BranchID: uuidFromPtr(item.BranchID), DepartmentID: uuidFromPtr(item.DepartmentID), UserID: uuidFromPtr(item.UserID), EffectiveFrom: dateFromPtr(item.EffectiveFrom), EffectiveTo: dateFromPtr(item.EffectiveTo), UpdatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "update attendance location assignment", err, tenantIDField(item.TenantID), stringField("assignment_id", item.ID.String()))
	}
	return mapAttendanceLocationAssignment(row), nil
}

func (s *Store) DeleteAttendanceLocationAssignment(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteAttendanceLocationAssignment(ctx, sqlc.SoftDeleteAttendanceLocationAssignmentParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete attendance location assignment", err, tenantIDField(tenantID), stringField("assignment_id", id.String()))
	}
	return nil
}

func (s *Store) CreateAttendanceDevice(ctx context.Context, item *domain.AttendanceDevice, actorID *uuid.UUID) (*domain.AttendanceDevice, error) {
	row, err := s.getQueries(ctx).CreateAttendanceDevice(ctx, sqlc.CreateAttendanceDeviceParams{TenantID: item.TenantID, AttendanceLocationID: uuidFromPtr(item.AttendanceLocationID), BranchID: uuidFromPtr(item.BranchID), Code: item.Code, Name: item.Name, Vendor: textFromPtr(item.Vendor), Model: textFromPtr(item.Model), SerialNumber: textFromPtr(item.SerialNumber), DeviceIdentifier: textFromPtr(item.DeviceIdentifier), IntegrationType: item.IntegrationType, DirectionMode: item.DirectionMode, Timezone: item.Timezone, Status: item.Status, Config: []byte(item.Config), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create attendance device", err, tenantIDField(item.TenantID), stringField("code", item.Code))
	}
	return mapAttendanceDevice(row), nil
}

func (s *Store) ListAttendanceDevices(ctx context.Context, tenantID uuid.UUID) ([]*domain.AttendanceDevice, error) {
	rows, err := s.getQueries(ctx).ListAttendanceDevices(ctx, tenantID)
	if err != nil {
		return nil, s.logDBError(ctx, "list attendance devices", err, tenantIDField(tenantID))
	}
	return mapAttendanceDevices(rows), nil
}

func (s *Store) GetAttendanceDevice(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.AttendanceDevice, error) {
	row, err := s.getQueries(ctx).GetAttendanceDevice(ctx, sqlc.GetAttendanceDeviceParams{TenantID: tenantID, ID: id})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrAttendanceDeviceNotFound
		}
		return nil, s.logDBError(ctx, "get attendance device", err, tenantIDField(tenantID), stringField("device_id", id.String()))
	}
	return mapAttendanceDevice(row), nil
}

func (s *Store) UpdateAttendanceDevice(ctx context.Context, item *domain.AttendanceDevice, actorID *uuid.UUID) (*domain.AttendanceDevice, error) {
	row, err := s.getQueries(ctx).UpdateAttendanceDevice(ctx, sqlc.UpdateAttendanceDeviceParams{TenantID: item.TenantID, ID: item.ID, AttendanceLocationID: uuidFromPtr(item.AttendanceLocationID), BranchID: uuidFromPtr(item.BranchID), Code: item.Code, Name: item.Name, Vendor: textFromPtr(item.Vendor), Model: textFromPtr(item.Model), SerialNumber: textFromPtr(item.SerialNumber), DeviceIdentifier: textFromPtr(item.DeviceIdentifier), IntegrationType: item.IntegrationType, DirectionMode: item.DirectionMode, Timezone: item.Timezone, Status: item.Status, Config: []byte(item.Config), UpdatedBy: uuidFromPtr(actorID)})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrAttendanceDeviceNotFound
		}
		return nil, s.logDBError(ctx, "update attendance device", err, tenantIDField(item.TenantID), stringField("device_id", item.ID.String()))
	}
	return mapAttendanceDevice(row), nil
}

func (s *Store) DeleteAttendanceDevice(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteAttendanceDevice(ctx, sqlc.SoftDeleteAttendanceDeviceParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete attendance device", err, tenantIDField(tenantID), stringField("device_id", id.String()))
	}
	return nil
}

func (s *Store) CreateEmployeeAttendanceDevice(ctx context.Context, item *domain.EmployeeAttendanceDevice, actorID *uuid.UUID) (*domain.EmployeeAttendanceDevice, error) {
	row, err := s.getQueries(ctx).CreateEmployeeAttendanceDevice(ctx, sqlc.CreateEmployeeAttendanceDeviceParams{TenantID: item.TenantID, UserID: item.UserID, DeviceID: item.DeviceID, DeviceUserID: item.DeviceUserID, CredentialType: item.CredentialType, CardNumber: textFromPtr(item.CardNumber), EffectiveFrom: dateFromPtr(item.EffectiveFrom), EffectiveTo: dateFromPtr(item.EffectiveTo), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create employee attendance device", err, tenantIDField(item.TenantID), stringField("user_id", item.UserID.String()))
	}
	return mapEmployeeAttendanceDevice(row), nil
}

func (s *Store) ListEmployeeAttendanceDevices(ctx context.Context, tenantID uuid.UUID, userID *uuid.UUID) ([]*domain.EmployeeAttendanceDevice, error) {
	if userID != nil && *userID != uuid.Nil {
		rows, err := s.getQueries(ctx).ListEmployeeAttendanceDevicesByUser(ctx, sqlc.ListEmployeeAttendanceDevicesByUserParams{TenantID: tenantID, UserID: *userID})
		if err != nil {
			return nil, s.logDBError(ctx, "list employee attendance devices by user", err, tenantIDField(tenantID), stringField("user_id", userID.String()))
		}
		return mapEmployeeAttendanceDevices(rows), nil
	}
	rows, err := s.getQueries(ctx).ListEmployeeAttendanceDevices(ctx, tenantID)
	if err != nil {
		return nil, s.logDBError(ctx, "list employee attendance devices", err, tenantIDField(tenantID))
	}
	return mapEmployeeAttendanceDevices(rows), nil
}

func (s *Store) GetEmployeeAttendanceDevice(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.EmployeeAttendanceDevice, error) {
	row, err := s.getQueries(ctx).GetEmployeeAttendanceDevice(ctx, sqlc.GetEmployeeAttendanceDeviceParams{TenantID: tenantID, ID: id})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrEmployeeAttendanceDeviceNotFound
		}
		return nil, s.logDBError(ctx, "get employee attendance device", err, tenantIDField(tenantID), stringField("mapping_id", id.String()))
	}
	return mapEmployeeAttendanceDevice(row), nil
}

func (s *Store) UpdateEmployeeAttendanceDevice(ctx context.Context, item *domain.EmployeeAttendanceDevice, actorID *uuid.UUID) (*domain.EmployeeAttendanceDevice, error) {
	row, err := s.getQueries(ctx).UpdateEmployeeAttendanceDevice(ctx, sqlc.UpdateEmployeeAttendanceDeviceParams{TenantID: item.TenantID, ID: item.ID, UserID: item.UserID, DeviceID: item.DeviceID, DeviceUserID: item.DeviceUserID, CredentialType: item.CredentialType, CardNumber: textFromPtr(item.CardNumber), EffectiveFrom: dateFromPtr(item.EffectiveFrom), EffectiveTo: dateFromPtr(item.EffectiveTo), UpdatedBy: uuidFromPtr(actorID)})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrEmployeeAttendanceDeviceNotFound
		}
		return nil, s.logDBError(ctx, "update employee attendance device", err, tenantIDField(item.TenantID), stringField("mapping_id", item.ID.String()))
	}
	return mapEmployeeAttendanceDevice(row), nil
}

func (s *Store) DeleteEmployeeAttendanceDevice(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteEmployeeAttendanceDevice(ctx, sqlc.SoftDeleteEmployeeAttendanceDeviceParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete employee attendance device", err, tenantIDField(tenantID), stringField("mapping_id", id.String()))
	}
	return nil
}

func (s *Store) CreateRawAttendanceEvent(ctx context.Context, item *domain.RawAttendanceEvent, actorID *uuid.UUID) (*domain.RawAttendanceEvent, error) {
	row, err := s.getQueries(ctx).CreateRawAttendanceEvent(ctx, sqlc.CreateRawAttendanceEventParams{TenantID: item.TenantID, DeviceID: item.DeviceID, EmployeeDeviceMappingID: uuidFromPtr(item.EmployeeDeviceMappingID), ExternalEventID: textFromPtr(item.ExternalEventID), DeviceUserID: textFromPtr(item.DeviceUserID), EventTime: timestamptzFromTime(item.EventTime), EventType: textFromPtr(item.EventType), AttendanceType: textFromPtr(item.AttendanceType), ImportBatchID: textFromPtr(item.ImportBatchID), RawPayload: []byte(item.RawPayload), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create raw attendance event", err, tenantIDField(item.TenantID), stringField("device_id", item.DeviceID.String()))
	}
	return mapRawAttendanceEvent(row), nil
}

func (s *Store) ListRawAttendanceEvents(ctx context.Context, tenantID uuid.UUID, limit int32) ([]*domain.RawAttendanceEvent, error) {
	rows, err := s.getQueries(ctx).ListRawAttendanceEvents(ctx, sqlc.ListRawAttendanceEventsParams{TenantID: tenantID, Limit: limit})
	if err != nil {
		return nil, s.logDBError(ctx, "list raw attendance events", err, tenantIDField(tenantID))
	}
	return mapRawAttendanceEvents(rows), nil
}
