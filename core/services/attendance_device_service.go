package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) CreateAttendanceLocation(ctx context.Context, cmd ports.AttendanceLocationCommand) (*domain.AttendanceLocation, error) {
	item, err := s.buildAttendanceLocation(ctx, cmd)
	if err != nil {
		return nil, err
	}
	result, err := s.attendanceLocations.CreateAttendanceLocation(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create attendance location", err, serviceTenantIDField(cmd.TenantID), serviceStringField("code", item.Code))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) ListAttendanceLocations(ctx context.Context, tenantID uuid.UUID) ([]*domain.AttendanceLocation, error) {
	if tenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	items, err := s.attendanceLocations.ListAttendanceLocations(ctx, tenantID)
	if err != nil {
		s.logError("list attendance locations", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	return items, nil
}

func (s *TenantService) UpdateAttendanceLocation(ctx context.Context, cmd ports.AttendanceLocationCommand) (*domain.AttendanceLocation, error) {
	if cmd.ID == uuid.Nil {
		return nil, domain.ErrInvalidAttendanceLocationID
	}
	item, err := s.buildAttendanceLocation(ctx, cmd)
	if err != nil {
		return nil, err
	}
	item.ID = cmd.ID
	result, err := s.attendanceLocations.UpdateAttendanceLocation(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("update attendance location", err, serviceTenantIDField(cmd.TenantID), serviceStringField("location_id", cmd.ID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) DeleteAttendanceLocation(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil {
		return domain.ErrInvalidTenantID
	}
	if id == uuid.Nil {
		return domain.ErrInvalidAttendanceLocationID
	}
	if err := s.attendanceLocations.DeleteAttendanceLocation(ctx, tenantID, id, actorID); err != nil {
		s.logError("delete attendance location", err, serviceTenantIDField(tenantID), serviceStringField("location_id", id.String()))
		return err
	}
	return nil
}

func (s *TenantService) CreateAttendanceLocationAssignment(ctx context.Context, cmd ports.AttendanceLocationAssignmentCommand) (*domain.AttendanceLocationAssignment, error) {
	item, err := s.buildAttendanceLocationAssignment(ctx, cmd)
	if err != nil {
		return nil, err
	}
	result, err := s.attendanceLocations.CreateAttendanceLocationAssignment(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create attendance location assignment", err, serviceTenantIDField(cmd.TenantID), serviceStringField("location_id", cmd.LocationID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) ListAttendanceLocationAssignments(ctx context.Context, tenantID uuid.UUID, locationID *uuid.UUID) ([]*domain.AttendanceLocationAssignment, error) {
	if tenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	items, err := s.attendanceLocations.ListAttendanceLocationAssignments(ctx, tenantID, locationID)
	if err != nil {
		s.logError("list attendance location assignments", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	return items, nil
}

func (s *TenantService) UpdateAttendanceLocationAssignment(ctx context.Context, cmd ports.AttendanceLocationAssignmentCommand) (*domain.AttendanceLocationAssignment, error) {
	if cmd.ID == uuid.Nil {
		return nil, domain.ErrInvalidAttendanceLocationID
	}
	item, err := s.buildAttendanceLocationAssignment(ctx, cmd)
	if err != nil {
		return nil, err
	}
	item.ID = cmd.ID
	result, err := s.attendanceLocations.UpdateAttendanceLocationAssignment(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("update attendance location assignment", err, serviceTenantIDField(cmd.TenantID), serviceStringField("assignment_id", cmd.ID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) DeleteAttendanceLocationAssignment(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil {
		return domain.ErrInvalidTenantID
	}
	if id == uuid.Nil {
		return domain.ErrInvalidAttendanceLocationID
	}
	if err := s.attendanceLocations.DeleteAttendanceLocationAssignment(ctx, tenantID, id, actorID); err != nil {
		s.logError("delete attendance location assignment", err, serviceTenantIDField(tenantID), serviceStringField("assignment_id", id.String()))
		return err
	}
	return nil
}

func (s *TenantService) CreateAttendanceDevice(ctx context.Context, cmd ports.AttendanceDeviceCommand) (*domain.AttendanceDevice, error) {
	item, err := s.buildAttendanceDevice(ctx, cmd)
	if err != nil {
		return nil, err
	}
	result, err := s.attendanceDevices.CreateAttendanceDevice(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create attendance device", err, serviceTenantIDField(cmd.TenantID), serviceStringField("code", item.Code))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) ListAttendanceDevices(ctx context.Context, tenantID uuid.UUID) ([]*domain.AttendanceDevice, error) {
	if tenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	items, err := s.attendanceDevices.ListAttendanceDevices(ctx, tenantID)
	if err != nil {
		s.logError("list attendance devices", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	return items, nil
}

func (s *TenantService) UpdateAttendanceDevice(ctx context.Context, cmd ports.AttendanceDeviceCommand) (*domain.AttendanceDevice, error) {
	if cmd.ID == uuid.Nil {
		return nil, domain.ErrInvalidAttendanceDeviceID
	}
	item, err := s.buildAttendanceDevice(ctx, cmd)
	if err != nil {
		return nil, err
	}
	item.ID = cmd.ID
	result, err := s.attendanceDevices.UpdateAttendanceDevice(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("update attendance device", err, serviceTenantIDField(cmd.TenantID), serviceStringField("device_id", cmd.ID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) DeleteAttendanceDevice(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil {
		return domain.ErrInvalidTenantID
	}
	if id == uuid.Nil {
		return domain.ErrInvalidAttendanceDeviceID
	}
	if err := s.attendanceDevices.DeleteAttendanceDevice(ctx, tenantID, id, actorID); err != nil {
		s.logError("delete attendance device", err, serviceTenantIDField(tenantID), serviceStringField("device_id", id.String()))
		return err
	}
	return nil
}

func (s *TenantService) CreateEmployeeAttendanceDevice(ctx context.Context, cmd ports.EmployeeAttendanceDeviceCommand) (*domain.EmployeeAttendanceDevice, error) {
	item, err := s.buildEmployeeAttendanceDevice(ctx, cmd)
	if err != nil {
		return nil, err
	}
	result, err := s.attendanceDevices.CreateEmployeeAttendanceDevice(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create employee attendance device", err, serviceTenantIDField(cmd.TenantID), serviceStringField("user_id", cmd.UserID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) ListEmployeeAttendanceDevices(ctx context.Context, tenantID uuid.UUID, userID *uuid.UUID) ([]*domain.EmployeeAttendanceDevice, error) {
	if tenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	items, err := s.attendanceDevices.ListEmployeeAttendanceDevices(ctx, tenantID, userID)
	if err != nil {
		s.logError("list employee attendance devices", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	return items, nil
}

func (s *TenantService) UpdateEmployeeAttendanceDevice(ctx context.Context, cmd ports.EmployeeAttendanceDeviceCommand) (*domain.EmployeeAttendanceDevice, error) {
	if cmd.ID == uuid.Nil {
		return nil, domain.ErrInvalidAttendanceDeviceID
	}
	item, err := s.buildEmployeeAttendanceDevice(ctx, cmd)
	if err != nil {
		return nil, err
	}
	item.ID = cmd.ID
	result, err := s.attendanceDevices.UpdateEmployeeAttendanceDevice(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("update employee attendance device", err, serviceTenantIDField(cmd.TenantID), serviceStringField("mapping_id", cmd.ID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) DeleteEmployeeAttendanceDevice(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil {
		return domain.ErrInvalidTenantID
	}
	if id == uuid.Nil {
		return domain.ErrInvalidAttendanceDeviceID
	}
	if err := s.attendanceDevices.DeleteEmployeeAttendanceDevice(ctx, tenantID, id, actorID); err != nil {
		s.logError("delete employee attendance device", err, serviceTenantIDField(tenantID), serviceStringField("mapping_id", id.String()))
		return err
	}
	return nil
}

func (s *TenantService) buildAttendanceLocation(ctx context.Context, cmd ports.AttendanceLocationCommand) (*domain.AttendanceLocation, error) {
	effectiveFrom, err := parseAttendancePolicyOptionalDate(cmd.EffectiveFrom)
	if err != nil {
		return nil, err
	}
	effectiveTo, err := parseAttendancePolicyOptionalDate(cmd.EffectiveTo)
	if err != nil {
		return nil, err
	}
	if cmd.BranchID != nil && *cmd.BranchID != uuid.Nil {
		if _, err := s.branches.GetBranch(ctx, cmd.TenantID, *cmd.BranchID); err != nil {
			return nil, err
		}
	}
	return domain.NewAttendanceLocation(domain.AttendanceLocationInput{TenantID: cmd.TenantID, BranchID: cmd.BranchID, Code: cmd.Code, Name: cmd.Name, LocationType: cmd.LocationType, Latitude: cmd.Latitude, Longitude: cmd.Longitude, RadiusMeters: cmd.RadiusMeters, Address: cmd.Address, City: cmd.City, State: cmd.State, Country: cmd.Country, Pincode: cmd.Pincode, EffectiveFrom: effectiveFrom, EffectiveTo: effectiveTo})
}

func (s *TenantService) buildAttendanceLocationAssignment(ctx context.Context, cmd ports.AttendanceLocationAssignmentCommand) (*domain.AttendanceLocationAssignment, error) {
	effectiveFrom, err := parseAttendancePolicyOptionalDate(cmd.EffectiveFrom)
	if err != nil {
		return nil, err
	}
	effectiveTo, err := parseAttendancePolicyOptionalDate(cmd.EffectiveTo)
	if err != nil {
		return nil, err
	}
	if _, err := s.attendanceLocations.GetAttendanceLocation(ctx, cmd.TenantID, cmd.LocationID); err != nil {
		return nil, err
	}
	if cmd.BranchID != nil && *cmd.BranchID != uuid.Nil {
		if _, err := s.branches.GetBranch(ctx, cmd.TenantID, *cmd.BranchID); err != nil {
			return nil, err
		}
	}
	if cmd.DepartmentID != nil && *cmd.DepartmentID != uuid.Nil {
		if _, err := s.departments.GetDepartment(ctx, cmd.TenantID, *cmd.DepartmentID); err != nil {
			return nil, err
		}
	}
	if cmd.UserID != nil && *cmd.UserID != uuid.Nil {
		if _, err := s.employees.GetEmployeeByUserID(ctx, cmd.TenantID, *cmd.UserID); err != nil {
			return nil, err
		}
	}
	return domain.NewAttendanceLocationAssignment(domain.AttendanceLocationAssignmentInput{TenantID: cmd.TenantID, LocationID: cmd.LocationID, BranchID: cmd.BranchID, DepartmentID: cmd.DepartmentID, UserID: cmd.UserID, EffectiveFrom: effectiveFrom, EffectiveTo: effectiveTo})
}

func (s *TenantService) buildAttendanceDevice(ctx context.Context, cmd ports.AttendanceDeviceCommand) (*domain.AttendanceDevice, error) {
	if cmd.AttendanceLocationID != nil && *cmd.AttendanceLocationID != uuid.Nil {
		if _, err := s.attendanceLocations.GetAttendanceLocation(ctx, cmd.TenantID, *cmd.AttendanceLocationID); err != nil {
			return nil, err
		}
	}
	if cmd.BranchID != nil && *cmd.BranchID != uuid.Nil {
		if _, err := s.branches.GetBranch(ctx, cmd.TenantID, *cmd.BranchID); err != nil {
			return nil, err
		}
	}
	return domain.NewAttendanceDevice(domain.AttendanceDeviceInput{TenantID: cmd.TenantID, AttendanceLocationID: cmd.AttendanceLocationID, BranchID: cmd.BranchID, Code: cmd.Code, Name: cmd.Name, Vendor: cmd.Vendor, Model: cmd.Model, SerialNumber: cmd.SerialNumber, DeviceIdentifier: cmd.DeviceIdentifier, IntegrationType: cmd.IntegrationType, DirectionMode: cmd.DirectionMode, Timezone: cmd.Timezone, Status: cmd.Status, Config: cmd.Config})
}

func (s *TenantService) buildEmployeeAttendanceDevice(ctx context.Context, cmd ports.EmployeeAttendanceDeviceCommand) (*domain.EmployeeAttendanceDevice, error) {
	effectiveFrom, err := parseAttendancePolicyOptionalDate(cmd.EffectiveFrom)
	if err != nil {
		return nil, err
	}
	effectiveTo, err := parseAttendancePolicyOptionalDate(cmd.EffectiveTo)
	if err != nil {
		return nil, err
	}
	if _, err := s.employees.GetEmployeeByUserID(ctx, cmd.TenantID, cmd.UserID); err != nil {
		return nil, err
	}
	if _, err := s.attendanceDevices.GetAttendanceDevice(ctx, cmd.TenantID, cmd.DeviceID); err != nil {
		return nil, err
	}
	return domain.NewEmployeeAttendanceDevice(domain.EmployeeAttendanceDeviceInput{TenantID: cmd.TenantID, UserID: cmd.UserID, DeviceID: cmd.DeviceID, DeviceUserID: cmd.DeviceUserID, CredentialType: cmd.CredentialType, CardNumber: cmd.CardNumber, EffectiveFrom: effectiveFrom, EffectiveTo: effectiveTo})
}
