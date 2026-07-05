package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (h *Handler) CreateNotificationMaster(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "create notification master", err, "tenant context is required")
		return
	}
	var cmd ports.NotificationMasterCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode notification master create request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateNotificationMaster(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "create notification master", err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) SendNotification(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "send notification", err, "tenant context is required")
		return
	}
	var cmd ports.SendNotificationCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode send notification request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	result, err := h.svc.SendNotification(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "send notification", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, result)
}

func (h *Handler) ListNotificationInbox(w http.ResponseWriter, r *http.Request) {
	tenantID, userID, ok := h.notificationInboxTenantUser(w, r, "list notification inbox")
	if !ok {
		return
	}
	items, err := h.svc.ListNotificationInbox(r.Context(), tenantID, userID, optionalBoolQuery(r, "is_read"), optionalUUIDQuery(r, "notification_master_id"), queryInt32(r, "limit", 25), queryInt32(r, "offset", 0))
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list notification inbox", err, "failed to list notification inbox")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) ListNotificationInboxPage(w http.ResponseWriter, r *http.Request) {
	tenantID, userID, ok := h.notificationInboxTenantUser(w, r, "list notification inbox page")
	if !ok {
		return
	}
	page, err := h.svc.ListNotificationInboxPage(r.Context(), h.notificationInboxFilter(r, tenantID, userID))
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list notification inbox page", err, "failed to list notification inbox")
		return
	}
	respondJSON(w, http.StatusOK, page)
}

func (h *Handler) CountUnreadNotifications(w http.ResponseWriter, r *http.Request) {
	tenantID, userID, ok := h.notificationInboxTenantUser(w, r, "count unread notifications")
	if !ok {
		return
	}
	count, err := h.svc.CountUnreadNotifications(r.Context(), tenantID, userID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "count unread notifications", err, "failed to count unread notifications")
		return
	}
	respondJSON(w, http.StatusOK, map[string]int64{"unread": count})
}

func (h *Handler) MarkNotificationRead(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "mark notification read", err, "tenant context is required")
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "notificationID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "mark notification read", err, "invalid notification id")
		return
	}
	if err := h.svc.MarkNotificationRead(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "mark notification read", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) MarkNotificationUnread(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "mark notification unread", err, "tenant context is required")
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "notificationID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "mark notification unread", err, "invalid notification id")
		return
	}
	if err := h.svc.MarkNotificationUnread(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "mark notification unread", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) MarkAllNotificationsRead(w http.ResponseWriter, r *http.Request) {
	tenantID, userID, ok := h.notificationInboxTenantUser(w, r, "mark all notifications read")
	if !ok {
		return
	}
	if err := h.svc.MarkAllNotificationsRead(r.Context(), tenantID, userID, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "mark all notifications read", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeleteNotificationInboxItem(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "archive notification inbox item", err, "tenant context is required")
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "notificationID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "archive notification inbox item", err, "invalid notification id")
		return
	}
	if err := h.svc.DeleteNotificationInboxItem(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "archive notification inbox item", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListNotificationLogs(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list notification logs", err, "tenant context is required")
		return
	}
	items, err := h.svc.ListNotificationLogs(r.Context(), tenantID, optionalUUIDQuery(r, "user_id"), optionalStringQuery(r, "status"), optionalStringQuery(r, "channel"), optionalUUIDQuery(r, "bulk_id"), queryInt32(r, "limit", 25), queryInt32(r, "offset", 0))
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list notification logs", err, "failed to list notification logs")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) ListNotificationMasters(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list notification masters", err, "tenant context is required")
		return
	}
	items, err := h.svc.ListNotificationMasters(r.Context(), tenantID, h.actorIDFromRequest(r))
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list notification masters", err, "failed to list notification masters")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) UpdateNotificationMaster(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.lookupRequestIDs(w, r, "notificationMasterID", "update notification master")
	if !ok {
		return
	}
	var cmd ports.NotificationMasterCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode notification master update request", err, "invalid request body")
		return
	}
	cmd.ID = id
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateNotificationMaster(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "update notification master", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) DeleteNotificationMaster(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.lookupRequestIDs(w, r, "notificationMasterID", "delete notification master")
	if !ok {
		return
	}
	if err := h.svc.DeleteNotificationMaster(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete notification master", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListNotificationPreferences(w http.ResponseWriter, r *http.Request) {
	tenantID, userID, ok := h.notificationPreferenceTenantUser(w, r, "list notification preferences")
	if !ok {
		return
	}
	items, err := h.svc.ListEffectiveNotificationPreferences(r.Context(), tenantID, userID, h.actorIDFromRequest(r))
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list notification preferences", err, "failed to list notification preferences")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) UpsertNotificationPreference(w http.ResponseWriter, r *http.Request) {
	tenantID, userID, ok := h.notificationPreferenceTenantUser(w, r, "upsert notification preference")
	if !ok {
		return
	}
	var cmd ports.NotificationPreferenceCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode notification preference request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.UserID = userID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpsertNotificationPreference(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "upsert notification preference", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) DeleteNotificationPreference(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "delete notification preference", err, "tenant context is required")
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "notificationPreferenceID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete notification preference", err, "invalid notification preference id")
		return
	}
	if err := h.svc.DeleteNotificationPreference(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete notification preference", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) UpsertDeviceToken(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "upsert device token", err, "tenant context is required")
		return
	}
	actorID := h.actorIDFromRequest(r)
	if actorID == nil || *actorID == uuid.Nil {
		h.respondError(w, r, http.StatusUnauthorized, "upsert device token", errors.New("user context is required"), "user context is required")
		return
	}
	var cmd ports.DeviceTokenCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode device token request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.UserID = *actorID
	cmd.ActorID = actorID
	item, err := h.svc.UpsertDeviceToken(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "upsert device token", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) ListDeviceTokens(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list device tokens", err, "tenant context is required")
		return
	}
	actorID := h.actorIDFromRequest(r)
	if actorID == nil || *actorID == uuid.Nil {
		h.respondError(w, r, http.StatusUnauthorized, "list device tokens", errors.New("user context is required"), "user context is required")
		return
	}
	items, err := h.svc.ListDeviceTokens(r.Context(), tenantID, *actorID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list device tokens", err, "failed to list device tokens")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) DeleteDeviceToken(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "delete device token", err, "tenant context is required")
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "deviceTokenID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete device token", err, "invalid device token id")
		return
	}
	if err := h.svc.DeleteDeviceToken(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete device token", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CreateTenantNotificationMaster(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "create tenant notification master")
	if !ok {
		return
	}
	var cmd ports.NotificationMasterCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode tenant notification master create request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateNotificationMaster(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "create tenant notification master", err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) ListTenantNotificationMasters(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant notification masters")
	if !ok {
		return
	}
	items, err := h.svc.ListNotificationMasters(r.Context(), tenantID, h.actorIDFromRequest(r))
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list tenant notification masters", err, "failed to list notification masters")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) UpdateTenantNotificationMaster(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLookupRequestIDs(w, r, "notificationMasterID", "update tenant notification master")
	if !ok {
		return
	}
	var cmd ports.NotificationMasterCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode tenant notification master update request", err, "invalid request body")
		return
	}
	cmd.ID = id
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateNotificationMaster(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "update tenant notification master", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) DeleteTenantNotificationMaster(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLookupRequestIDs(w, r, "notificationMasterID", "delete tenant notification master")
	if !ok {
		return
	}
	if err := h.svc.DeleteNotificationMaster(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant notification master", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListTenantNotificationPreferences(w http.ResponseWriter, r *http.Request) {
	tenantID, userID, ok := h.tenantNotificationPreferenceTenantUser(w, r, "list tenant notification preferences")
	if !ok {
		return
	}
	items, err := h.svc.ListEffectiveNotificationPreferences(r.Context(), tenantID, userID, h.actorIDFromRequest(r))
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list tenant notification preferences", err, "failed to list notification preferences")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) SendTenantNotification(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "send tenant notification")
	if !ok {
		return
	}
	var cmd ports.SendNotificationCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode tenant send notification request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	result, err := h.svc.SendNotification(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "send tenant notification", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, result)
}

func (h *Handler) ListTenantNotificationInbox(w http.ResponseWriter, r *http.Request) {
	tenantID, userID, ok := h.tenantNotificationPreferenceTenantUser(w, r, "list tenant notification inbox")
	if !ok {
		return
	}
	items, err := h.svc.ListNotificationInbox(r.Context(), tenantID, userID, optionalBoolQuery(r, "is_read"), optionalUUIDQuery(r, "notification_master_id"), queryInt32(r, "limit", 25), queryInt32(r, "offset", 0))
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list tenant notification inbox", err, "failed to list notification inbox")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) ListTenantNotificationInboxPage(w http.ResponseWriter, r *http.Request) {
	tenantID, userID, ok := h.tenantNotificationPreferenceTenantUser(w, r, "list tenant notification inbox page")
	if !ok {
		return
	}
	page, err := h.svc.ListNotificationInboxPage(r.Context(), h.notificationInboxFilter(r, tenantID, userID))
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list tenant notification inbox page", err, "failed to list notification inbox")
		return
	}
	respondJSON(w, http.StatusOK, page)
}

func (h *Handler) MarkTenantNotificationRead(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "mark tenant notification read")
	if !ok {
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "notificationID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "mark tenant notification read", err, "invalid notification id")
		return
	}
	if err := h.svc.MarkNotificationRead(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "mark tenant notification read", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) MarkTenantNotificationUnread(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "mark tenant notification unread")
	if !ok {
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "notificationID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "mark tenant notification unread", err, "invalid notification id")
		return
	}
	if err := h.svc.MarkNotificationUnread(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "mark tenant notification unread", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) MarkTenantAllNotificationsRead(w http.ResponseWriter, r *http.Request) {
	tenantID, userID, ok := h.tenantNotificationPreferenceTenantUser(w, r, "mark tenant all notifications read")
	if !ok {
		return
	}
	if err := h.svc.MarkAllNotificationsRead(r.Context(), tenantID, userID, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "mark tenant all notifications read", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeleteTenantNotificationInboxItem(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "archive tenant notification inbox item")
	if !ok {
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "notificationID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "archive tenant notification inbox item", err, "invalid notification id")
		return
	}
	if err := h.svc.DeleteNotificationInboxItem(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "archive tenant notification inbox item", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListTenantNotificationLogs(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant notification logs")
	if !ok {
		return
	}
	items, err := h.svc.ListNotificationLogs(r.Context(), tenantID, optionalUUIDQuery(r, "user_id"), optionalStringQuery(r, "status"), optionalStringQuery(r, "channel"), optionalUUIDQuery(r, "bulk_id"), queryInt32(r, "limit", 25), queryInt32(r, "offset", 0))
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list tenant notification logs", err, "failed to list notification logs")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) UpsertTenantNotificationPreference(w http.ResponseWriter, r *http.Request) {
	tenantID, userID, ok := h.tenantNotificationPreferenceTenantUser(w, r, "upsert tenant notification preference")
	if !ok {
		return
	}
	var cmd ports.NotificationPreferenceCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode tenant notification preference request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.UserID = userID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpsertNotificationPreference(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "upsert tenant notification preference", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) UpsertTenantDeviceToken(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "upsert tenant device token")
	if !ok {
		return
	}
	var cmd ports.DeviceTokenCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode tenant device token request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpsertDeviceToken(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "upsert tenant device token", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) ListTenantDeviceTokens(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant device tokens")
	if !ok {
		return
	}
	if deviceID := r.URL.Query().Get("device_id"); deviceID != "" {
		items, err := h.svc.ListDeviceTokensByDeviceID(r.Context(), tenantID, deviceID)
		if err != nil {
			h.respondError(w, r, http.StatusInternalServerError, "list tenant device tokens by device", err, "failed to list device tokens")
			return
		}
		respondJSON(w, http.StatusOK, items)
		return
	}
	userID := optionalUUIDQuery(r, "user_id")
	if userID == nil {
		h.respondError(w, r, http.StatusBadRequest, "list tenant device tokens", errors.New("user_id is required"), "user_id is required")
		return
	}
	items, err := h.svc.ListDeviceTokens(r.Context(), tenantID, *userID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list tenant device tokens", err, "failed to list device tokens")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) DeleteTenantDeviceToken(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "delete tenant device token")
	if !ok {
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "deviceTokenID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant device token", err, "invalid device token id")
		return
	}
	if err := h.svc.DeleteDeviceToken(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant device token", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) notificationPreferenceTenantUser(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	userID := h.actorIDFromRequest(r)
	if raw := r.URL.Query().Get("user_id"); raw != "" {
		parsed, err := uuid.Parse(raw)
		if err != nil {
			h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid user id")
			return uuid.Nil, uuid.Nil, false
		}
		userID = &parsed
	}
	if userID == nil || *userID == uuid.Nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, errors.New("user context is required"), "user context is required")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, *userID, true
}

func (h *Handler) notificationInboxTenantUser(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	return h.notificationPreferenceTenantUser(w, r, operation)
}

func (h *Handler) tenantNotificationPreferenceTenantUser(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, ok := h.superAdminTenantID(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	raw := r.URL.Query().Get("user_id")
	if raw == "" {
		h.respondError(w, r, http.StatusBadRequest, operation, errors.New("user_id is required"), "user_id is required")
		return uuid.Nil, uuid.Nil, false
	}
	userID, err := uuid.Parse(raw)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid user id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, userID, true
}

func optionalBoolQuery(r *http.Request, key string) *bool {
	value := r.URL.Query().Get(key)
	if value == "" {
		return nil
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return nil
	}
	return &parsed
}

func optionalUUIDQuery(r *http.Request, key string) *uuid.UUID {
	value := r.URL.Query().Get(key)
	if value == "" {
		return nil
	}
	parsed, err := uuid.Parse(value)
	if err != nil || parsed == uuid.Nil {
		return nil
	}
	return &parsed
}

func optionalStringQuery(r *http.Request, key string) *string {
	value := r.URL.Query().Get(key)
	if value == "" {
		return nil
	}
	return &value
}

func (h *Handler) notificationInboxFilter(r *http.Request, tenantID uuid.UUID, userID uuid.UUID) domain.NotificationInboxFilter {
	return domain.NotificationInboxFilter{
		TenantID:             tenantID,
		UserID:               userID,
		IsRead:               optionalBoolQuery(r, "is_read"),
		NotificationMasterID: optionalUUIDQuery(r, "notification_master_id"),
		NotificationCode:     optionalStringQuery(r, "notification_code"),
		Search:               optionalStringQuery(r, "search"),
		Limit:                queryInt32(r, "limit", 25),
		Offset:               queryInt32(r, "offset", 0),
	}
}
