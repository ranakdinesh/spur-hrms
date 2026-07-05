package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (h *Handler) ListAssetItems(w http.ResponseWriter, r *http.Request) {
	h.withAssetTenant(w, r, "list asset items", func(tenantID uuid.UUID, operation string) { h.listAssetItemsForTenant(w, r, tenantID, operation) })
}

func (h *Handler) CreateAssetItem(w http.ResponseWriter, r *http.Request) {
	h.withAssetTenant(w, r, "create asset item", func(tenantID uuid.UUID, operation string) {
		h.saveAssetItemForTenant(w, r, tenantID, uuid.Nil, operation)
	})
}

func (h *Handler) UpdateAssetItem(w http.ResponseWriter, r *http.Request) {
	h.withAssetTenantID(w, r, "assetID", "update asset item", h.saveAssetItemForTenant)
}

func (h *Handler) UpdateAssetItemStatus(w http.ResponseWriter, r *http.Request) {
	h.withAssetTenantID(w, r, "assetID", "update asset item status", h.updateAssetStatusForTenant("asset"))
}

func (h *Handler) DeleteAssetItem(w http.ResponseWriter, r *http.Request) {
	h.withAssetTenantID(w, r, "assetID", "delete asset item", func(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
		if err := h.svc.DeleteAssetItem(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
			h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})
}

func (h *Handler) ListAccessCatalogItems(w http.ResponseWriter, r *http.Request) {
	h.withAssetTenant(w, r, "list access catalog", func(tenantID uuid.UUID, operation string) { h.listAccessCatalogForTenant(w, r, tenantID, operation) })
}

func (h *Handler) CreateAccessCatalogItem(w http.ResponseWriter, r *http.Request) {
	h.withAssetTenant(w, r, "create access catalog item", func(tenantID uuid.UUID, operation string) {
		h.saveAccessCatalogForTenant(w, r, tenantID, uuid.Nil, operation)
	})
}

func (h *Handler) UpdateAccessCatalogItem(w http.ResponseWriter, r *http.Request) {
	h.withAssetTenantID(w, r, "accessID", "update access catalog item", h.saveAccessCatalogForTenant)
}

func (h *Handler) DeleteAccessCatalogItem(w http.ResponseWriter, r *http.Request) {
	h.withAssetTenantID(w, r, "accessID", "delete access catalog item", func(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
		if err := h.svc.DeleteAccessCatalogItem(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
			h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})
}

func (h *Handler) ListAssetAssignments(w http.ResponseWriter, r *http.Request) {
	h.withAssetTenant(w, r, "list asset assignments", func(tenantID uuid.UUID, operation string) { h.listAssetAssignmentsForTenant(w, r, tenantID, operation) })
}

func (h *Handler) CreateAssetAssignment(w http.ResponseWriter, r *http.Request) {
	h.withAssetTenant(w, r, "create asset assignment", func(tenantID uuid.UUID, operation string) {
		h.saveAssetAssignmentForTenant(w, r, tenantID, uuid.Nil, operation)
	})
}

func (h *Handler) UpdateAssetAssignment(w http.ResponseWriter, r *http.Request) {
	h.withAssetTenantID(w, r, "assignmentID", "update asset assignment", h.saveAssetAssignmentForTenant)
}

func (h *Handler) UpdateAssetAssignmentStatus(w http.ResponseWriter, r *http.Request) {
	h.withAssetTenantID(w, r, "assignmentID", "update asset assignment status", h.updateAssetStatusForTenant("assignment"))
}

func (h *Handler) ListAccessLifecycleTasks(w http.ResponseWriter, r *http.Request) {
	h.withAssetTenant(w, r, "list access lifecycle tasks", func(tenantID uuid.UUID, operation string) { h.listAccessTasksForTenant(w, r, tenantID, operation) })
}

func (h *Handler) CreateAccessLifecycleTask(w http.ResponseWriter, r *http.Request) {
	h.withAssetTenant(w, r, "create access lifecycle task", func(tenantID uuid.UUID, operation string) {
		h.saveAccessTaskForTenant(w, r, tenantID, uuid.Nil, operation)
	})
}

func (h *Handler) UpdateAccessLifecycleTask(w http.ResponseWriter, r *http.Request) {
	h.withAssetTenantID(w, r, "taskID", "update access lifecycle task", h.saveAccessTaskForTenant)
}

func (h *Handler) UpdateAccessLifecycleTaskStatus(w http.ResponseWriter, r *http.Request) {
	h.withAssetTenantID(w, r, "taskID", "update access lifecycle task status", h.updateAssetStatusForTenant("access_task"))
}

func (h *Handler) ListAssetAccessEvents(w http.ResponseWriter, r *http.Request) {
	h.withAssetTenant(w, r, "list asset access events", func(tenantID uuid.UUID, operation string) {
		items, err := h.svc.ListAssetAccessEvents(r.Context(), assetAccessFilterFromRequest(r, tenantID), optionalStringQuery(r, "source_type"), optionalUUIDQuery(r, "source_id"))
		if err != nil {
			h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
			return
		}
		respondJSON(w, http.StatusOK, items)
	})
}

func (h *Handler) GetAssetAccessSummary(w http.ResponseWriter, r *http.Request) {
	h.withAssetTenant(w, r, "get asset access summary", func(tenantID uuid.UUID, operation string) {
		items, err := h.svc.GetAssetAccessSummary(r.Context(), tenantID)
		if err != nil {
			h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
			return
		}
		respondJSON(w, http.StatusOK, items)
	})
}

func (h *Handler) ListTenantAssetItems(w http.ResponseWriter, r *http.Request) {
	h.ListAssetItems(w, r)
}
func (h *Handler) CreateTenantAssetItem(w http.ResponseWriter, r *http.Request) {
	h.CreateAssetItem(w, r)
}
func (h *Handler) UpdateTenantAssetItem(w http.ResponseWriter, r *http.Request) {
	h.UpdateAssetItem(w, r)
}
func (h *Handler) UpdateTenantAssetItemStatus(w http.ResponseWriter, r *http.Request) {
	h.UpdateAssetItemStatus(w, r)
}
func (h *Handler) DeleteTenantAssetItem(w http.ResponseWriter, r *http.Request) {
	h.DeleteAssetItem(w, r)
}
func (h *Handler) ListTenantAccessCatalogItems(w http.ResponseWriter, r *http.Request) {
	h.ListAccessCatalogItems(w, r)
}
func (h *Handler) CreateTenantAccessCatalogItem(w http.ResponseWriter, r *http.Request) {
	h.CreateAccessCatalogItem(w, r)
}
func (h *Handler) UpdateTenantAccessCatalogItem(w http.ResponseWriter, r *http.Request) {
	h.UpdateAccessCatalogItem(w, r)
}
func (h *Handler) DeleteTenantAccessCatalogItem(w http.ResponseWriter, r *http.Request) {
	h.DeleteAccessCatalogItem(w, r)
}
func (h *Handler) ListTenantAssetAssignments(w http.ResponseWriter, r *http.Request) {
	h.ListAssetAssignments(w, r)
}
func (h *Handler) CreateTenantAssetAssignment(w http.ResponseWriter, r *http.Request) {
	h.CreateAssetAssignment(w, r)
}
func (h *Handler) UpdateTenantAssetAssignment(w http.ResponseWriter, r *http.Request) {
	h.UpdateAssetAssignment(w, r)
}
func (h *Handler) UpdateTenantAssetAssignmentStatus(w http.ResponseWriter, r *http.Request) {
	h.UpdateAssetAssignmentStatus(w, r)
}
func (h *Handler) ListTenantAccessLifecycleTasks(w http.ResponseWriter, r *http.Request) {
	h.ListAccessLifecycleTasks(w, r)
}
func (h *Handler) CreateTenantAccessLifecycleTask(w http.ResponseWriter, r *http.Request) {
	h.CreateAccessLifecycleTask(w, r)
}
func (h *Handler) UpdateTenantAccessLifecycleTask(w http.ResponseWriter, r *http.Request) {
	h.UpdateAccessLifecycleTask(w, r)
}
func (h *Handler) UpdateTenantAccessLifecycleTaskStatus(w http.ResponseWriter, r *http.Request) {
	h.UpdateAccessLifecycleTaskStatus(w, r)
}
func (h *Handler) ListTenantAssetAccessEvents(w http.ResponseWriter, r *http.Request) {
	h.ListAssetAccessEvents(w, r)
}
func (h *Handler) GetTenantAssetAccessSummary(w http.ResponseWriter, r *http.Request) {
	h.GetAssetAccessSummary(w, r)
}

func (h *Handler) listAssetItemsForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	items, err := h.svc.ListAssetItems(r.Context(), assetAccessFilterFromRequest(r, tenantID))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) saveAssetItemForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.AssetItemCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID, cmd.ID, cmd.ActorID = tenantID, id, h.actorIDFromRequest(r)
	var item *domain.AssetItem
	var err error
	if id == uuid.Nil {
		item, err = h.svc.CreateAssetItem(r.Context(), cmd)
	} else {
		item, err = h.svc.UpdateAssetItem(r.Context(), cmd)
	}
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) listAccessCatalogForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	items, err := h.svc.ListAccessCatalogItems(r.Context(), assetAccessFilterFromRequest(r, tenantID))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) saveAccessCatalogForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.AccessCatalogItemCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID, cmd.ID, cmd.ActorID = tenantID, id, h.actorIDFromRequest(r)
	var item *domain.AccessCatalogItem
	var err error
	if id == uuid.Nil {
		item, err = h.svc.CreateAccessCatalogItem(r.Context(), cmd)
	} else {
		item, err = h.svc.UpdateAccessCatalogItem(r.Context(), cmd)
	}
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) listAssetAssignmentsForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	items, err := h.svc.ListAssetAssignments(r.Context(), assetAccessFilterFromRequest(r, tenantID))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) saveAssetAssignmentForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.AssetAssignmentCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID, cmd.ID, cmd.ActorID = tenantID, id, h.actorIDFromRequest(r)
	var item *domain.AssetAssignment
	var err error
	if id == uuid.Nil {
		item, err = h.svc.CreateAssetAssignment(r.Context(), cmd)
	} else {
		item, err = h.svc.UpdateAssetAssignment(r.Context(), cmd)
	}
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) listAccessTasksForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	items, err := h.svc.ListAccessLifecycleTasks(r.Context(), assetAccessFilterFromRequest(r, tenantID))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) saveAccessTaskForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.AccessLifecycleTaskCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID, cmd.ID, cmd.ActorID = tenantID, id, h.actorIDFromRequest(r)
	var item *domain.AccessLifecycleTask
	var err error
	if id == uuid.Nil {
		item, err = h.svc.CreateAccessLifecycleTask(r.Context(), cmd)
	} else {
		item, err = h.svc.UpdateAccessLifecycleTask(r.Context(), cmd)
	}
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) updateAssetStatusForTenant(kind string) func(http.ResponseWriter, *http.Request, uuid.UUID, uuid.UUID, string) {
	return func(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
		var cmd ports.AssetAccessStatusCommand
		if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
			h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
			return
		}
		cmd.TenantID, cmd.ID, cmd.ActorID = tenantID, id, h.actorIDFromRequest(r)
		var result any
		var err error
		switch kind {
		case "asset":
			result, err = h.svc.UpdateAssetItemStatus(r.Context(), cmd)
		case "assignment":
			result, err = h.svc.UpdateAssetAssignmentStatus(r.Context(), cmd)
		default:
			result, err = h.svc.UpdateAccessLifecycleTaskStatus(r.Context(), cmd)
		}
		if err != nil {
			h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
			return
		}
		respondJSON(w, http.StatusOK, result)
	}
}

func (h *Handler) withAssetTenant(w http.ResponseWriter, r *http.Request, operation string, fn func(uuid.UUID, string)) {
	tenantID, err := h.assetTenantID(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return
	}
	fn(tenantID, operation)
}

func (h *Handler) withAssetTenantID(w http.ResponseWriter, r *http.Request, idParam string, operation string, fn func(http.ResponseWriter, *http.Request, uuid.UUID, uuid.UUID, string)) {
	tenantID, err := h.assetTenantID(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, idParam))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid id")
		return
	}
	fn(w, r, tenantID, id, operation)
}

func (h *Handler) assetTenantID(r *http.Request) (uuid.UUID, error) {
	if raw := chi.URLParam(r, "tenantID"); raw != "" {
		return uuid.Parse(raw)
	}
	return h.tenantIDFromRequest(r)
}

func assetAccessFilterFromRequest(r *http.Request, tenantID uuid.UUID) domain.AssetAccessFilter {
	return domain.AssetAccessFilter{TenantID: tenantID, WorkerProfileID: optionalUUIDQuery(r, "worker_profile_id"), AssetID: optionalUUIDQuery(r, "asset_id"), AccessItemID: optionalUUIDQuery(r, "access_item_id"), ExitRequestID: optionalUUIDQuery(r, "exit_request_id"), Status: optionalStringQuery(r, "status"), Category: optionalStringQuery(r, "category"), AccessType: optionalStringQuery(r, "access_type"), Search: optionalStringQuery(r, "search"), Limit: queryInt32(r, "limit", 50), Offset: queryInt32(r, "offset", 0)}
}
