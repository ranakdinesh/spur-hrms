package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (h *Handler) CreateAttendanceLocation(w http.ResponseWriter, r *http.Request) {
	if tenantID, err := h.tenantIDFromRequest(r); err == nil {
		h.createAttendanceLocationForTenant(w, r, tenantID, "create attendance location")
	} else {
		h.respondError(w, r, http.StatusUnauthorized, "create attendance location", err, "tenant context is required")
	}
}
func (h *Handler) ListAttendanceLocations(w http.ResponseWriter, r *http.Request) {
	if tenantID, err := h.tenantIDFromRequest(r); err == nil {
		h.listAttendanceLocationsForTenant(w, r, tenantID, "list attendance locations")
	} else {
		h.respondError(w, r, http.StatusUnauthorized, "list attendance locations", err, "tenant context is required")
	}
}
func (h *Handler) UpdateAttendanceLocation(w http.ResponseWriter, r *http.Request) {
	if tenantID, err := h.tenantIDFromRequest(r); err == nil {
		h.updateAttendanceLocationForTenant(w, r, tenantID, "update attendance location")
	} else {
		h.respondError(w, r, http.StatusUnauthorized, "update attendance location", err, "tenant context is required")
	}
}
func (h *Handler) DeleteAttendanceLocation(w http.ResponseWriter, r *http.Request) {
	if tenantID, err := h.tenantIDFromRequest(r); err == nil {
		h.deleteAttendanceLocationForTenant(w, r, tenantID, "delete attendance location")
	} else {
		h.respondError(w, r, http.StatusUnauthorized, "delete attendance location", err, "tenant context is required")
	}
}

func (h *Handler) CreateTenantAttendanceLocation(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "create tenant attendance location"); ok {
		h.createAttendanceLocationForTenant(w, r, tenantID, "create tenant attendance location")
	}
}
func (h *Handler) ListTenantAttendanceLocations(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "list tenant attendance locations"); ok {
		h.listAttendanceLocationsForTenant(w, r, tenantID, "list tenant attendance locations")
	}
}
func (h *Handler) UpdateTenantAttendanceLocation(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "update tenant attendance location"); ok {
		h.updateAttendanceLocationForTenant(w, r, tenantID, "update tenant attendance location")
	}
}
func (h *Handler) DeleteTenantAttendanceLocation(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "delete tenant attendance location"); ok {
		h.deleteAttendanceLocationForTenant(w, r, tenantID, "delete tenant attendance location")
	}
}

func (h *Handler) CreateAttendanceLocationAssignment(w http.ResponseWriter, r *http.Request) {
	if tenantID, err := h.tenantIDFromRequest(r); err == nil {
		h.createAttendanceLocationAssignmentForTenant(w, r, tenantID, "create attendance location assignment")
	} else {
		h.respondError(w, r, http.StatusUnauthorized, "create attendance location assignment", err, "tenant context is required")
	}
}
func (h *Handler) ListAttendanceLocationAssignments(w http.ResponseWriter, r *http.Request) {
	if tenantID, err := h.tenantIDFromRequest(r); err == nil {
		h.listAttendanceLocationAssignmentsForTenant(w, r, tenantID, "list attendance location assignments")
	} else {
		h.respondError(w, r, http.StatusUnauthorized, "list attendance location assignments", err, "tenant context is required")
	}
}
func (h *Handler) UpdateAttendanceLocationAssignment(w http.ResponseWriter, r *http.Request) {
	if tenantID, err := h.tenantIDFromRequest(r); err == nil {
		h.updateAttendanceLocationAssignmentForTenant(w, r, tenantID, "update attendance location assignment")
	} else {
		h.respondError(w, r, http.StatusUnauthorized, "update attendance location assignment", err, "tenant context is required")
	}
}
func (h *Handler) DeleteAttendanceLocationAssignment(w http.ResponseWriter, r *http.Request) {
	if tenantID, err := h.tenantIDFromRequest(r); err == nil {
		h.deleteAttendanceLocationAssignmentForTenant(w, r, tenantID, "delete attendance location assignment")
	} else {
		h.respondError(w, r, http.StatusUnauthorized, "delete attendance location assignment", err, "tenant context is required")
	}
}
func (h *Handler) CreateTenantAttendanceLocationAssignment(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "create tenant attendance location assignment"); ok {
		h.createAttendanceLocationAssignmentForTenant(w, r, tenantID, "create tenant attendance location assignment")
	}
}
func (h *Handler) ListTenantAttendanceLocationAssignments(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "list tenant attendance location assignments"); ok {
		h.listAttendanceLocationAssignmentsForTenant(w, r, tenantID, "list tenant attendance location assignments")
	}
}
func (h *Handler) UpdateTenantAttendanceLocationAssignment(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "update tenant attendance location assignment"); ok {
		h.updateAttendanceLocationAssignmentForTenant(w, r, tenantID, "update tenant attendance location assignment")
	}
}
func (h *Handler) DeleteTenantAttendanceLocationAssignment(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "delete tenant attendance location assignment"); ok {
		h.deleteAttendanceLocationAssignmentForTenant(w, r, tenantID, "delete tenant attendance location assignment")
	}
}

func (h *Handler) CreateAttendanceDevice(w http.ResponseWriter, r *http.Request) {
	if tenantID, err := h.tenantIDFromRequest(r); err == nil {
		h.createAttendanceDeviceForTenant(w, r, tenantID, "create attendance device")
	} else {
		h.respondError(w, r, http.StatusUnauthorized, "create attendance device", err, "tenant context is required")
	}
}
func (h *Handler) ListAttendanceDevices(w http.ResponseWriter, r *http.Request) {
	if tenantID, err := h.tenantIDFromRequest(r); err == nil {
		h.listAttendanceDevicesForTenant(w, r, tenantID, "list attendance devices")
	} else {
		h.respondError(w, r, http.StatusUnauthorized, "list attendance devices", err, "tenant context is required")
	}
}
func (h *Handler) UpdateAttendanceDevice(w http.ResponseWriter, r *http.Request) {
	if tenantID, err := h.tenantIDFromRequest(r); err == nil {
		h.updateAttendanceDeviceForTenant(w, r, tenantID, "update attendance device")
	} else {
		h.respondError(w, r, http.StatusUnauthorized, "update attendance device", err, "tenant context is required")
	}
}
func (h *Handler) DeleteAttendanceDevice(w http.ResponseWriter, r *http.Request) {
	if tenantID, err := h.tenantIDFromRequest(r); err == nil {
		h.deleteAttendanceDeviceForTenant(w, r, tenantID, "delete attendance device")
	} else {
		h.respondError(w, r, http.StatusUnauthorized, "delete attendance device", err, "tenant context is required")
	}
}
func (h *Handler) CreateTenantAttendanceDevice(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "create tenant attendance device"); ok {
		h.createAttendanceDeviceForTenant(w, r, tenantID, "create tenant attendance device")
	}
}
func (h *Handler) ListTenantAttendanceDevices(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "list tenant attendance devices"); ok {
		h.listAttendanceDevicesForTenant(w, r, tenantID, "list tenant attendance devices")
	}
}
func (h *Handler) UpdateTenantAttendanceDevice(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "update tenant attendance device"); ok {
		h.updateAttendanceDeviceForTenant(w, r, tenantID, "update tenant attendance device")
	}
}
func (h *Handler) DeleteTenantAttendanceDevice(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "delete tenant attendance device"); ok {
		h.deleteAttendanceDeviceForTenant(w, r, tenantID, "delete tenant attendance device")
	}
}

func (h *Handler) CreateEmployeeAttendanceDevice(w http.ResponseWriter, r *http.Request) {
	if tenantID, err := h.tenantIDFromRequest(r); err == nil {
		h.createEmployeeAttendanceDeviceForTenant(w, r, tenantID, "create employee attendance device")
	} else {
		h.respondError(w, r, http.StatusUnauthorized, "create employee attendance device", err, "tenant context is required")
	}
}
func (h *Handler) ListEmployeeAttendanceDevices(w http.ResponseWriter, r *http.Request) {
	if tenantID, err := h.tenantIDFromRequest(r); err == nil {
		h.listEmployeeAttendanceDevicesForTenant(w, r, tenantID, "list employee attendance devices")
	} else {
		h.respondError(w, r, http.StatusUnauthorized, "list employee attendance devices", err, "tenant context is required")
	}
}
func (h *Handler) UpdateEmployeeAttendanceDevice(w http.ResponseWriter, r *http.Request) {
	if tenantID, err := h.tenantIDFromRequest(r); err == nil {
		h.updateEmployeeAttendanceDeviceForTenant(w, r, tenantID, "update employee attendance device")
	} else {
		h.respondError(w, r, http.StatusUnauthorized, "update employee attendance device", err, "tenant context is required")
	}
}
func (h *Handler) DeleteEmployeeAttendanceDevice(w http.ResponseWriter, r *http.Request) {
	if tenantID, err := h.tenantIDFromRequest(r); err == nil {
		h.deleteEmployeeAttendanceDeviceForTenant(w, r, tenantID, "delete employee attendance device")
	} else {
		h.respondError(w, r, http.StatusUnauthorized, "delete employee attendance device", err, "tenant context is required")
	}
}
func (h *Handler) CreateTenantEmployeeAttendanceDevice(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "create tenant employee attendance device"); ok {
		h.createEmployeeAttendanceDeviceForTenant(w, r, tenantID, "create tenant employee attendance device")
	}
}
func (h *Handler) ListTenantEmployeeAttendanceDevices(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "list tenant employee attendance devices"); ok {
		h.listEmployeeAttendanceDevicesForTenant(w, r, tenantID, "list tenant employee attendance devices")
	}
}
func (h *Handler) UpdateTenantEmployeeAttendanceDevice(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "update tenant employee attendance device"); ok {
		h.updateEmployeeAttendanceDeviceForTenant(w, r, tenantID, "update tenant employee attendance device")
	}
}
func (h *Handler) DeleteTenantEmployeeAttendanceDevice(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "delete tenant employee attendance device"); ok {
		h.deleteEmployeeAttendanceDeviceForTenant(w, r, tenantID, "delete tenant employee attendance device")
	}
}

func (h *Handler) createAttendanceLocationForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.AttendanceLocationCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateAttendanceLocation(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}
func (h *Handler) listAttendanceLocationsForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	items, err := h.svc.ListAttendanceLocations(r.Context(), tenantID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list attendance locations")
		return
	}
	respondJSON(w, http.StatusOK, items)
}
func (h *Handler) updateAttendanceLocationForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.AttendanceLocationCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "locationID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid location id")
		return
	}
	cmd.ID = id
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateAttendanceLocation(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}
func (h *Handler) deleteAttendanceLocationForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	id, err := uuid.Parse(chi.URLParam(r, "locationID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid location id")
		return
	}
	if err := h.svc.DeleteAttendanceLocation(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) createAttendanceLocationAssignmentForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.AttendanceLocationAssignmentCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateAttendanceLocationAssignment(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}
func (h *Handler) listAttendanceLocationAssignmentsForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	locationID := optionalUUIDQuery(r, "location_id")
	items, err := h.svc.ListAttendanceLocationAssignments(r.Context(), tenantID, locationID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list attendance location assignments")
		return
	}
	respondJSON(w, http.StatusOK, items)
}
func (h *Handler) updateAttendanceLocationAssignmentForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.AttendanceLocationAssignmentCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "assignmentID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid assignment id")
		return
	}
	cmd.ID = id
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateAttendanceLocationAssignment(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}
func (h *Handler) deleteAttendanceLocationAssignmentForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	id, err := uuid.Parse(chi.URLParam(r, "assignmentID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid assignment id")
		return
	}
	if err := h.svc.DeleteAttendanceLocationAssignment(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) createAttendanceDeviceForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.AttendanceDeviceCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateAttendanceDevice(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}
func (h *Handler) listAttendanceDevicesForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	items, err := h.svc.ListAttendanceDevices(r.Context(), tenantID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list attendance devices")
		return
	}
	respondJSON(w, http.StatusOK, items)
}
func (h *Handler) updateAttendanceDeviceForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.AttendanceDeviceCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "deviceID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid device id")
		return
	}
	cmd.ID = id
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateAttendanceDevice(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}
func (h *Handler) deleteAttendanceDeviceForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	id, err := uuid.Parse(chi.URLParam(r, "deviceID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid device id")
		return
	}
	if err := h.svc.DeleteAttendanceDevice(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) createEmployeeAttendanceDeviceForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.EmployeeAttendanceDeviceCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateEmployeeAttendanceDevice(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}
func (h *Handler) listEmployeeAttendanceDevicesForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	userID := optionalUUIDQuery(r, "user_id")
	items, err := h.svc.ListEmployeeAttendanceDevices(r.Context(), tenantID, userID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list employee attendance devices")
		return
	}
	respondJSON(w, http.StatusOK, items)
}
func (h *Handler) updateEmployeeAttendanceDeviceForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.EmployeeAttendanceDeviceCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "mappingID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid mapping id")
		return
	}
	cmd.ID = id
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateEmployeeAttendanceDevice(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}
func (h *Handler) deleteEmployeeAttendanceDeviceForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	id, err := uuid.Parse(chi.URLParam(r, "mappingID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid mapping id")
		return
	}
	if err := h.svc.DeleteEmployeeAttendanceDevice(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
