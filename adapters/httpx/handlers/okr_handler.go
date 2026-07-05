package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (h *Handler) ListOKRCycles(w http.ResponseWriter, r *http.Request) {
	if tenantID, err := h.tenantIDFromRequest(r); err == nil {
		h.listOKRCyclesForTenant(w, r, tenantID, "list okr cycles")
	} else {
		h.respondError(w, r, http.StatusUnauthorized, "list okr cycles", err, "tenant context is required")
	}
}

func (h *Handler) CreateOKRCycle(w http.ResponseWriter, r *http.Request) {
	if tenantID, err := h.tenantIDFromRequest(r); err == nil {
		h.createOKRCycleForTenant(w, r, tenantID, "create okr cycle")
	} else {
		h.respondError(w, r, http.StatusUnauthorized, "create okr cycle", err, "tenant context is required")
	}
}

func (h *Handler) GetOKRCycle(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.okrCycleRequestIDs(w, r, "get okr cycle")
	if !ok {
		return
	}
	h.getOKRCycleForTenant(w, r, tenantID, id, "get okr cycle")
}

func (h *Handler) UpdateOKRCycle(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.okrCycleRequestIDs(w, r, "update okr cycle")
	if !ok {
		return
	}
	h.updateOKRCycleForTenant(w, r, tenantID, id, "update okr cycle")
}

func (h *Handler) DeleteOKRCycle(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.okrCycleRequestIDs(w, r, "delete okr cycle")
	if !ok {
		return
	}
	if err := h.svc.DeleteOKRCycle(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete okr cycle", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) UpdateOKRCycleStatus(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.okrCycleRequestIDs(w, r, "update okr cycle status")
	if !ok {
		return
	}
	h.updateOKRStatusForTenant(w, r, tenantID, id, "update okr cycle status", true)
}

func (h *Handler) ListObjectives(w http.ResponseWriter, r *http.Request) {
	if tenantID, err := h.tenantIDFromRequest(r); err == nil {
		h.listObjectivesForTenant(w, r, tenantID, "list objectives")
	} else {
		h.respondError(w, r, http.StatusUnauthorized, "list objectives", err, "tenant context is required")
	}
}

func (h *Handler) CreateObjective(w http.ResponseWriter, r *http.Request) {
	if tenantID, err := h.tenantIDFromRequest(r); err == nil {
		h.createObjectiveForTenant(w, r, tenantID, "create objective")
	} else {
		h.respondError(w, r, http.StatusUnauthorized, "create objective", err, "tenant context is required")
	}
}

func (h *Handler) GetObjective(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.objectiveRequestIDs(w, r, "get objective")
	if !ok {
		return
	}
	h.getObjectiveForTenant(w, r, tenantID, id, "get objective")
}

func (h *Handler) UpdateObjective(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.objectiveRequestIDs(w, r, "update objective")
	if !ok {
		return
	}
	h.updateObjectiveForTenant(w, r, tenantID, id, "update objective")
}

func (h *Handler) DeleteObjective(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.objectiveRequestIDs(w, r, "delete objective")
	if !ok {
		return
	}
	if err := h.svc.DeleteObjective(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete objective", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) UpdateObjectiveStatus(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.objectiveRequestIDs(w, r, "update objective status")
	if !ok {
		return
	}
	h.updateOKRStatusForTenant(w, r, tenantID, id, "update objective status", false)
}

func (h *Handler) ListKeyResults(w http.ResponseWriter, r *http.Request) {
	if tenantID, err := h.tenantIDFromRequest(r); err == nil {
		h.listKeyResultsForTenant(w, r, tenantID, "list key results")
	} else {
		h.respondError(w, r, http.StatusUnauthorized, "list key results", err, "tenant context is required")
	}
}

func (h *Handler) CreateKeyResult(w http.ResponseWriter, r *http.Request) {
	if tenantID, err := h.tenantIDFromRequest(r); err == nil {
		h.createKeyResultForTenant(w, r, tenantID, "create key result")
	} else {
		h.respondError(w, r, http.StatusUnauthorized, "create key result", err, "tenant context is required")
	}
}

func (h *Handler) GetKeyResult(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.keyResultRequestIDs(w, r, "get key result")
	if !ok {
		return
	}
	h.getKeyResultForTenant(w, r, tenantID, id, "get key result")
}

func (h *Handler) UpdateKeyResult(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.keyResultRequestIDs(w, r, "update key result")
	if !ok {
		return
	}
	h.updateKeyResultForTenant(w, r, tenantID, id, "update key result")
}

func (h *Handler) DeleteKeyResult(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.keyResultRequestIDs(w, r, "delete key result")
	if !ok {
		return
	}
	if err := h.svc.DeleteKeyResult(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete key result", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CreateKeyResultCheckIn(w http.ResponseWriter, r *http.Request) {
	if tenantID, err := h.tenantIDFromRequest(r); err == nil {
		h.createKeyResultCheckInForTenant(w, r, tenantID, "create key result check-in")
	} else {
		h.respondError(w, r, http.StatusUnauthorized, "create key result check-in", err, "tenant context is required")
	}
}

func (h *Handler) ListKeyResultCheckIns(w http.ResponseWriter, r *http.Request) {
	if tenantID, err := h.tenantIDFromRequest(r); err == nil {
		h.listKeyResultCheckInsForTenant(w, r, tenantID, "list key result check-ins")
	} else {
		h.respondError(w, r, http.StatusUnauthorized, "list key result check-ins", err, "tenant context is required")
	}
}

func (h *Handler) GetOKRSummary(w http.ResponseWriter, r *http.Request) {
	if tenantID, err := h.tenantIDFromRequest(r); err == nil {
		h.getOKRSummaryForTenant(w, r, tenantID, "get okr summary")
	} else {
		h.respondError(w, r, http.StatusUnauthorized, "get okr summary", err, "tenant context is required")
	}
}

func (h *Handler) ListTenantOKRCycles(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "list tenant okr cycles"); ok {
		h.listOKRCyclesForTenant(w, r, tenantID, "list tenant okr cycles")
	}
}

func (h *Handler) CreateTenantOKRCycle(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "create tenant okr cycle"); ok {
		h.createOKRCycleForTenant(w, r, tenantID, "create tenant okr cycle")
	}
}

func (h *Handler) GetTenantOKRCycle(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminOKRCycleRequestIDs(w, r, "get tenant okr cycle")
	if ok {
		h.getOKRCycleForTenant(w, r, tenantID, id, "get tenant okr cycle")
	}
}

func (h *Handler) UpdateTenantOKRCycle(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminOKRCycleRequestIDs(w, r, "update tenant okr cycle")
	if ok {
		h.updateOKRCycleForTenant(w, r, tenantID, id, "update tenant okr cycle")
	}
}

func (h *Handler) DeleteTenantOKRCycle(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminOKRCycleRequestIDs(w, r, "delete tenant okr cycle")
	if !ok {
		return
	}
	if err := h.svc.DeleteOKRCycle(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant okr cycle", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) UpdateTenantOKRCycleStatus(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminOKRCycleRequestIDs(w, r, "update tenant okr cycle status")
	if ok {
		h.updateOKRStatusForTenant(w, r, tenantID, id, "update tenant okr cycle status", true)
	}
}

func (h *Handler) ListTenantObjectives(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "list tenant objectives"); ok {
		h.listObjectivesForTenant(w, r, tenantID, "list tenant objectives")
	}
}

func (h *Handler) CreateTenantObjective(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "create tenant objective"); ok {
		h.createObjectiveForTenant(w, r, tenantID, "create tenant objective")
	}
}

func (h *Handler) GetTenantObjective(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminObjectiveRequestIDs(w, r, "get tenant objective")
	if ok {
		h.getObjectiveForTenant(w, r, tenantID, id, "get tenant objective")
	}
}

func (h *Handler) UpdateTenantObjective(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminObjectiveRequestIDs(w, r, "update tenant objective")
	if ok {
		h.updateObjectiveForTenant(w, r, tenantID, id, "update tenant objective")
	}
}

func (h *Handler) DeleteTenantObjective(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminObjectiveRequestIDs(w, r, "delete tenant objective")
	if !ok {
		return
	}
	if err := h.svc.DeleteObjective(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant objective", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) UpdateTenantObjectiveStatus(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminObjectiveRequestIDs(w, r, "update tenant objective status")
	if ok {
		h.updateOKRStatusForTenant(w, r, tenantID, id, "update tenant objective status", false)
	}
}

func (h *Handler) ListTenantKeyResults(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "list tenant key results"); ok {
		h.listKeyResultsForTenant(w, r, tenantID, "list tenant key results")
	}
}

func (h *Handler) CreateTenantKeyResult(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "create tenant key result"); ok {
		h.createKeyResultForTenant(w, r, tenantID, "create tenant key result")
	}
}

func (h *Handler) GetTenantKeyResult(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminKeyResultRequestIDs(w, r, "get tenant key result")
	if ok {
		h.getKeyResultForTenant(w, r, tenantID, id, "get tenant key result")
	}
}

func (h *Handler) UpdateTenantKeyResult(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminKeyResultRequestIDs(w, r, "update tenant key result")
	if ok {
		h.updateKeyResultForTenant(w, r, tenantID, id, "update tenant key result")
	}
}

func (h *Handler) DeleteTenantKeyResult(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminKeyResultRequestIDs(w, r, "delete tenant key result")
	if !ok {
		return
	}
	if err := h.svc.DeleteKeyResult(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant key result", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CreateTenantKeyResultCheckIn(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "create tenant key result check-in"); ok {
		h.createKeyResultCheckInForTenant(w, r, tenantID, "create tenant key result check-in")
	}
}

func (h *Handler) ListTenantKeyResultCheckIns(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "list tenant key result check-ins"); ok {
		h.listKeyResultCheckInsForTenant(w, r, tenantID, "list tenant key result check-ins")
	}
}

func (h *Handler) GetTenantOKRSummary(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "get tenant okr summary"); ok {
		h.getOKRSummaryForTenant(w, r, tenantID, "get tenant okr summary")
	}
}

func (h *Handler) listOKRCyclesForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	items, err := h.svc.ListOKRCycles(r.Context(), domain.OKRCycleFilter{TenantID: tenantID, Status: optionalStringQuery(r, "status"), Search: optionalStringQuery(r, "search")})
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list okr cycles")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) createOKRCycleForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.OKRCycleCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateOKRCycle(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) getOKRCycleForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	item, err := h.svc.GetOKRCycle(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) updateOKRCycleForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.OKRCycleCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.ID = id
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateOKRCycle(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) updateOKRStatusForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string, cycle bool) {
	var cmd ports.OKRStatusCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.ID = id
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	if cycle {
		item, err := h.svc.UpdateOKRCycleStatus(r.Context(), cmd)
		if err != nil {
			h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
			return
		}
		respondJSON(w, http.StatusOK, item)
		return
	}
	item, err := h.svc.UpdateObjectiveStatus(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) listObjectivesForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	cycleID, ok := h.optionalUUIDQuery(w, r, "cycle_id", operation)
	if !ok {
		return
	}
	items, err := h.svc.ListObjectives(r.Context(), domain.ObjectiveFilter{TenantID: tenantID, CycleID: cycleID, OwnerType: optionalStringQuery(r, "owner_type"), Status: optionalStringQuery(r, "status"), Search: optionalStringQuery(r, "search")})
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list objectives")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) createObjectiveForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.ObjectiveCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateObjective(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) getObjectiveForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	item, err := h.svc.GetObjective(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) updateObjectiveForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.ObjectiveCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.ID = id
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateObjective(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) listKeyResultsForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	objectiveID, ok := h.optionalUUIDQuery(w, r, "objective_id", operation)
	if !ok {
		return
	}
	items, err := h.svc.ListKeyResults(r.Context(), domain.KeyResultFilter{TenantID: tenantID, ObjectiveID: objectiveID, Status: optionalStringQuery(r, "status")})
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list key results")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) createKeyResultForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.KeyResultCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateKeyResult(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) getKeyResultForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	item, err := h.svc.GetKeyResult(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) updateKeyResultForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.KeyResultCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.ID = id
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateKeyResult(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) createKeyResultCheckInForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.KeyResultCheckInCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateKeyResultCheckIn(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) listKeyResultCheckInsForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	keyResultID, ok := h.optionalUUIDQuery(w, r, "key_result_id", operation)
	if !ok {
		return
	}
	objectiveID, ok := h.optionalUUIDQuery(w, r, "objective_id", operation)
	if !ok {
		return
	}
	items, err := h.svc.ListKeyResultCheckIns(r.Context(), domain.KeyResultCheckInFilter{TenantID: tenantID, KeyResultID: keyResultID, ObjectiveID: objectiveID})
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list key result check-ins")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) getOKRSummaryForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	cycleID, ok := h.optionalUUIDQuery(w, r, "cycle_id", operation)
	if !ok {
		return
	}
	items, err := h.svc.GetOKRSummary(r.Context(), tenantID, cycleID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to get okr summary")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) okrCycleRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	id, ok := h.uuidURLParam(w, r, "cycleID", operation, "invalid cycle id")
	return tenantID, id, ok
}

func (h *Handler) objectiveRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	id, ok := h.uuidURLParam(w, r, "objectiveID", operation, "invalid objective id")
	return tenantID, id, ok
}

func (h *Handler) keyResultRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	id, ok := h.uuidURLParam(w, r, "keyResultID", operation, "invalid key result id")
	return tenantID, id, ok
}

func (h *Handler) superAdminOKRCycleRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, ok := h.superAdminTenantID(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	id, ok := h.uuidURLParam(w, r, "cycleID", operation, "invalid cycle id")
	return tenantID, id, ok
}

func (h *Handler) superAdminObjectiveRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, ok := h.superAdminTenantID(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	id, ok := h.uuidURLParam(w, r, "objectiveID", operation, "invalid objective id")
	return tenantID, id, ok
}

func (h *Handler) superAdminKeyResultRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, ok := h.superAdminTenantID(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	id, ok := h.uuidURLParam(w, r, "keyResultID", operation, "invalid key result id")
	return tenantID, id, ok
}
