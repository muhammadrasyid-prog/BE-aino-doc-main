package models

import (
	"database/sql"
	"time"
)

// Struktur untuk merepresentasikan data history timeline
type TimelineHistory struct {
	// FormID       uint           `json:"form_id" db:"form_id"`
	FormUUID     string         `json:"form_uuid" db:"form_uuid"`
	FormNumber   string         `json:"form_number,omitempty" db:"form_number"`
	FormTicket   string         `json:"form_ticket,omitempty" db:"form_ticket"`
	FormStatus   string         `json:"form_status" db:"form_status"`
	ProjectUUID  string         `json:"project_uuid" db:"project_uuid"`
	ProjectName  string         `json:"project_name,omitempty" db:"project_name"` // Tambah ini
	DocumentUUID string         `json:"document_uuid" db:"document_uuid"`
	DocumentName string         `json:"document_name,omitempty" db:"document_name"` // Tambah ini
	CreatedBy    string         `json:"created_by" db:"created_by"`
	CreatedAt    time.Time      `json:"created_at" db:"created_at"`
	UpdatedBy    sql.NullString `json:"updated_by,omitempty" db:"updated_by"`
	UpdatedAt    sql.NullTime   `json:"updated_at,omitempty" db:"updated_at"`
}
