package service

import (
	"database/sql"
	"document/models"
	"log"
	"time"
	// "time"
	// "github.com/google/uuid"
)

// GetTimelineHistory mengambil semua data history dokumen dari tabel form_ms
// func GetTimelineHistory(db *sql.DB) ([]models.TimelineHistory, error) {
// 	rows, err := db.Query(`
// 		SELECT
// 			f.form_uuid, f.form_number, f.form_ticket, f.form_status,
// 			f.document_id, f.user_id, f.project_id,
// 			d.document_name, p.project_name,
// 			f.created_by, f.created_at, f.updated_by, f.updated_at
// 		FROM form_ms f
// 		LEFT JOIN document_ms d ON f.document_id = d.document_id
// 		LEFT JOIN project_ms p ON f.project_id = p.project_id
// 		WHERE f.deleted_at IS NULL
// 		ORDER BY f.created_at DESC
// 	`)
// 	if err != nil {
// 		log.Println("Error executing query:", err)
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	// Slice untuk menampung data history
// 	var historyList []models.TimelineHistory

// 	// Iterasi hasil query
// 	for rows.Next() {
// 		var history models.TimelineHistory
// 		var updatedBy sql.NullString
// 		var updatedAt sql.NullTime
// 		var formUUID string

// 		// Scan hasil query ke variabel
// 		err := rows.Scan(
// 			&formUUID, &history.FormNumber, &history.FormTicket, &history.FormStatus,
// 			&history.DocumentID, &history.UserID, &history.ProjectID,
// 			&history.DocumentName, &history.ProjectName,
// 			&history.CreatedBy, &history.CreatedAt, &updatedBy, &updatedAt,
// 		)
// 		if err != nil {
// 			log.Println("Error scanning row:", err)
// 			return nil, err
// 		}

// 		// Set FormUUID
// 		history.FormUUID = formUUID

// 		// Handle NULL untuk kolom UpdatedBy
// 		if updatedBy.Valid {
// 			history.UpdatedBy = sql.NullString{
// 				String: updatedBy.String, // Assign nilai string
// 				Valid:  true,             // Tandai sebagai valid
// 			}
// 		} else {
// 			history.UpdatedBy = sql.NullString{
// 				String: "",    // Nilai default string
// 				Valid:  false, // Tandai sebagai tidak valid
// 			}
// 		}

// 		// Handle NULL untuk kolom UpdatedAt
// 		if updatedAt.Valid {
// 			history.UpdatedAt = sql.NullTime{
// 				Time:  updatedAt.Time, // Assign nilai time.Time
// 				Valid: true,           // Tandai sebagai valid
// 			}
// 		} else {
// 			history.UpdatedAt = sql.NullTime{
// 				Time:  time.Time{}, // Nilai default time.Time
// 				Valid: false,       // Tandai sebagai tidak valid
// 			}
// 		}

// 		// Tambahkan history ke list
// 		historyList = append(historyList, history)
// 	}

// 	// Cek error setelah iterasi
// 	if err = rows.Err(); err != nil {
// 		log.Println("Error in row iteration:", err)
// 		return nil, err
// 	}

// 	return historyList, nil

func GetRecentTimelineHistory(db *sql.DB) ([]models.TimelineHistory, error) {
	rows, err := db.Query(`
		SELECT 
			f.form_uuid, f.form_number, f.form_ticket, f.form_status,
			d.document_uuid, d.document_name,
			p.project_uuid, p.project_name,
			f.created_by, f.created_at, f.updated_by, f.updated_at
		FROM form_ms f
		LEFT JOIN document_ms d ON f.document_id = d.document_id
		LEFT JOIN project_ms p ON f.project_id = p.project_id
		WHERE f.deleted_at IS NULL
		ORDER BY f.created_at DESC
		LIMIT 3; -- Ambil 3 data terbaru
	`)
	if err != nil {
		log.Println("Error executing recent query:", err)
		return nil, err
	}
	defer rows.Close()

	return scanTimelineHistory(rows)
}

func GetOlderTimelineHistory(db *sql.DB, lastCreatedAt time.Time, limit int) ([]models.TimelineHistory, error) {
	log.Printf("Fetching older timeline starting from: %v with limit %d\n", lastCreatedAt, limit)

	rows, err := db.Query(`
		WITH RecentForms AS (
			SELECT form_uuid
			FROM form_ms
			WHERE deleted_at IS NULL
			ORDER BY created_at DESC
			LIMIT 3
		)
		SELECT 
			f.form_uuid, f.form_number, f.form_ticket, f.form_status,
			d.document_uuid, d.document_name,
			p.project_uuid, p.project_name,
			f.created_by, f.created_at, f.updated_by, f.updated_at
		FROM form_ms f
		LEFT JOIN document_ms d ON f.document_id = d.document_id
		LEFT JOIN project_ms p ON f.project_id = p.project_id
		WHERE f.deleted_at IS NULL 
		AND f.form_uuid NOT IN (SELECT form_uuid FROM RecentForms)
		ORDER BY f.created_at DESC
		LIMIT $1
	`, limit)
	if err != nil {
		log.Println("Error executing older query:", err)
		return nil, err
	}
	defer rows.Close()

	return scanTimelineHistory(rows)
}

func scanTimelineHistory(rows *sql.Rows) ([]models.TimelineHistory, error) {
	var historyList []models.TimelineHistory

	for rows.Next() {
		var history models.TimelineHistory
		var updatedBy sql.NullString
		var updatedAt sql.NullTime
		var formUUID string

		err := rows.Scan(
			&formUUID, &history.FormNumber, &history.FormTicket, &history.FormStatus,
			&history.DocumentUUID, &history.DocumentName,
			&history.ProjectUUID, &history.ProjectName,
			&history.CreatedBy, &history.CreatedAt, &updatedBy, &updatedAt,
		)
		if err != nil {
			log.Println("Error scanning row:", err)
			return nil, err
		}

		history.FormUUID = formUUID
		history.UpdatedBy = updatedBy
		history.UpdatedAt = updatedAt

		historyList = append(historyList, history)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error in row iteration:", err)
		return nil, err
	}

	return historyList, nil
}
