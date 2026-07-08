package domain

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	AttendanceSegmentOffice      = "office"
	AttendanceSegmentSite        = "site"
	AttendanceSegmentClientSite  = "client_site"
	AttendanceSegmentProjectSite = "project_site"
	AttendanceSegmentTravel      = "travel"
	AttendanceSegmentBreak       = "break"
	AttendanceSegmentRemote      = "remote"
	AttendanceSegmentOther       = "other"

	AttendanceSegmentActionStart    = "start"
	AttendanceSegmentActionEnd      = "end"
	AttendanceSegmentActionCheckin  = "checkin"
	AttendanceSegmentActionCheckout = "checkout"
	AttendanceSegmentActionArrive   = "arrive"
	AttendanceSegmentActionDepart   = "depart"
	AttendanceSegmentActionNote     = "note"

	AttendanceReferenceClient  = "client"
	AttendanceReferenceProject = "project"
	AttendanceReferenceSite    = "site"
	AttendanceReferenceRoute   = "route"
	AttendanceReferenceTicket  = "ticket"
	AttendanceReferenceTask    = "task"
	AttendanceReferenceOther   = "other"
)

var (
	ErrInvalidAttendanceSegmentType   = errors.New("attendance segment type is invalid")
	ErrInvalidAttendanceSegmentAction = errors.New("attendance segment action is invalid")
	ErrInvalidAttendanceReferenceType = errors.New("attendance reference type is invalid")
)

type AttendanceWorkdaySegment struct {
	ID                         uuid.UUID       `json:"id"`
	TenantID                   uuid.UUID       `json:"tenant_id"`
	UserID                     uuid.UUID       `json:"user_id"`
	Date                       time.Time       `json:"date"`
	EventTime                  time.Time       `json:"event_time"`
	SegmentType                string          `json:"segment_type"`
	Action                     string          `json:"action"`
	WorkMode                   *string         `json:"work_mode,omitempty"`
	Source                     *string         `json:"source,omitempty"`
	AttendanceLocationID       *uuid.UUID      `json:"attendance_location_id,omitempty"`
	ReferenceType              *string         `json:"reference_type,omitempty"`
	ReferenceID                *uuid.UUID      `json:"reference_id,omitempty"`
	ReferenceLabel             *string         `json:"reference_label,omitempty"`
	Latitude                   *float64        `json:"latitude,omitempty"`
	Longitude                  *float64        `json:"longitude,omitempty"`
	LocationAccuracyMeters     *float64        `json:"location_accuracy_meters,omitempty"`
	LocationVerificationStatus string          `json:"location_verification_status"`
	Remarks                    *string         `json:"remarks,omitempty"`
	Metadata                   json.RawMessage `json:"metadata,omitempty"`
	Inactive                   bool            `json:"inactive"`
	CreatedAt                  time.Time       `json:"created_at"`
	CreatedBy                  *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt                  time.Time       `json:"updated_at"`
	UpdatedBy                  *uuid.UUID      `json:"updated_by,omitempty"`
}

func NewAttendanceWorkdaySegment(item AttendanceWorkdaySegment) (*AttendanceWorkdaySegment, error) {
	if item.TenantID == uuid.Nil {
		return nil, ErrInvalidTenantID
	}
	if item.UserID == uuid.Nil {
		return nil, ErrInvalidLeaveUser
	}
	if item.Date.IsZero() {
		return nil, ErrInvalidAttendanceDate
	}
	if item.EventTime.IsZero() {
		return nil, ErrInvalidAttendanceDate
	}
	segmentType := strings.ToLower(strings.TrimSpace(item.SegmentType))
	if !validAttendanceSegmentType(segmentType) {
		return nil, ErrInvalidAttendanceSegmentType
	}
	action := strings.ToLower(strings.TrimSpace(item.Action))
	if !validAttendanceSegmentAction(action) {
		return nil, ErrInvalidAttendanceSegmentAction
	}
	if item.WorkMode != nil {
		clean := strings.ToLower(strings.TrimSpace(*item.WorkMode))
		if !validAttendanceSegmentWorkMode(clean) {
			return nil, ErrInvalidAttendanceWorkMode
		}
		item.WorkMode = &clean
	}
	if item.Source != nil {
		clean := strings.ToLower(strings.TrimSpace(*item.Source))
		if !validAttendanceSource(clean) {
			return nil, ErrInvalidAttendanceSource
		}
		item.Source = &clean
	}
	if item.ReferenceType != nil {
		clean := strings.ToLower(strings.TrimSpace(*item.ReferenceType))
		if !validAttendanceReferenceType(clean) {
			return nil, ErrInvalidAttendanceReferenceType
		}
		item.ReferenceType = &clean
	}
	if item.Latitude != nil && (*item.Latitude < -90 || *item.Latitude > 90) {
		return nil, ErrInvalidAttendanceLocation
	}
	if item.Longitude != nil && (*item.Longitude < -180 || *item.Longitude > 180) {
		return nil, ErrInvalidAttendanceLocation
	}
	if len(item.Metadata) == 0 {
		item.Metadata = json.RawMessage(`{}`)
	}
	if item.LocationVerificationStatus == "" {
		item.LocationVerificationStatus = "not_checked"
	}
	now := time.Now().UTC()
	item.Date = dateOnly(item.Date)
	item.EventTime = item.EventTime.UTC()
	item.SegmentType = segmentType
	item.Action = action
	item.ReferenceLabel = cleanOptional(item.ReferenceLabel)
	item.Remarks = cleanOptional(item.Remarks)
	item.CreatedAt = now
	item.UpdatedAt = now
	return &item, nil
}

func validAttendanceSegmentType(value string) bool {
	switch value {
	case AttendanceSegmentOffice, AttendanceSegmentSite, AttendanceSegmentClientSite, AttendanceSegmentProjectSite, AttendanceSegmentTravel, AttendanceSegmentBreak, AttendanceSegmentRemote, AttendanceSegmentOther:
		return true
	default:
		return false
	}
}

func validAttendanceSegmentAction(value string) bool {
	switch value {
	case AttendanceSegmentActionStart, AttendanceSegmentActionEnd, AttendanceSegmentActionCheckin, AttendanceSegmentActionCheckout, AttendanceSegmentActionArrive, AttendanceSegmentActionDepart, AttendanceSegmentActionNote:
		return true
	default:
		return false
	}
}

func validAttendanceSegmentWorkMode(value string) bool {
	switch value {
	case AttendanceWorkModeOffice, AttendanceWorkModeRemote, AttendanceWorkModeField, AttendanceWorkModeHybrid, AttendanceSegmentClientSite, AttendanceSegmentProjectSite:
		return true
	default:
		return false
	}
}

func validAttendanceReferenceType(value string) bool {
	switch value {
	case AttendanceReferenceClient, AttendanceReferenceProject, AttendanceReferenceSite, AttendanceReferenceRoute, AttendanceReferenceTicket, AttendanceReferenceTask, AttendanceReferenceOther:
		return true
	default:
		return false
	}
}
