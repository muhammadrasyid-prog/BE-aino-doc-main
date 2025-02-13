package models

import (
	"database/sql"
	"time"
)

type Tour struct {
    TourID        uint          `json:"tour_id" db:"tour_id"`         // Mengganti ID menjadi tour_id
    TourUUID      string        `json:"tour_uuid" db:"tour_uuid"`     // Menambahkan kolom tour_uuid
    TourName      string        `json:"tour_name" db:"tour_name"`
    PageURL       string        `json:"page_url" db:"page_url"`
    PageName      string        `json:"page_name" db:"page_name"` // ini nanti dibuat data json
    StepOrder     int           `json:"step_order" db:"step_order"`
    ElementID     string        `json:"element_id" db:"element_id"` // jadi element id ini hanya yang akan disesuaikan dengan page name
    Title         string        `json:"title" db:"title"`
    Content       string        `json:"content" db:"content"`
    Placement     string        `json:"placement" db:"placement"`
    Buttons       string        `json:"buttons" db:"buttons"`         // Menambahkan kolom buttons untuk menangani tombol input
    Completed     bool          `json:"completed" db:"completed"`
    CreatedAt     time.Time     `json:"created_at" db:"created_at"`
    UpdatedAt     time.Time     `json:"updated_at" db:"updated_at"`
    DeletedAt     sql.NullTime  `json:"deleted_at" db:"deleted_at"`
}

type Step struct {
    TourUUID   string `json:"tour_uuid" db:"tour_uuid"`
    StepOrder  int    `json:"step_order" db:"step_order"`
    ElementID  string `json:"element_id" db:"element_id"`
    Title      string `json:"title" db:"title"`
    Content    string `json:"content" db:"content"`
    Placement  string `json:"placement" db:"placement"`
    Buttons    string `json:"buttons" db:"buttons"`
}

type TourWithSteps struct {
    TourName  string `json:"tour_name" db:"tour_name"`
    PageName  string `json:"page_name" db:"page_name"`
    Steps     []Step `json:"steps" db:"steps"`
}
