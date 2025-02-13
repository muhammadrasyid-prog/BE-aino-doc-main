package service

import (
	"database/sql"
	"document/models"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

// SaveTourGuide inserts a new tour guide into the database
func AddTourGuide(addTour models.Tour) error {
	currentTimestamp := time.Now().UnixNano() / int64(time.Microsecond)
	uniqueID := uuid.New().ID()

	tour_id := currentTimestamp + int64(uniqueID)

	uuid := uuid.New()
	uuidString := uuid.String()
	_, err := db.NamedExec("INSERT INTO tours (tour_id, tour_uuid, tour_name, page_url, page_name, step_order, element_id, title, content, placement, buttons) VALUES (:tour_id, :tour_uuid, :tour_name, :page_url, :page_name, :step_order, :element_id, :title, :content, :placement, :buttons)", map[string]interface{}{
		"tour_id":    tour_id,
		"tour_uuid":  uuidString,
		"tour_name":  addTour.TourName,
		"page_url":   addTour.PageURL,
		"page_name":  addTour.PageName,
		"step_order": addTour.StepOrder,
		"element_id": addTour.ElementID,
		"title":      addTour.Title,
		"content":    addTour.Content,
		"placement":  addTour.Placement,
		"buttons":    addTour.Buttons,
	})
	if err != nil {
		log.Println("Error inserting tour guide:", err)
		return err
	}
	return nil
}


func ShowTourById(id string) (models.Tour, error) {
	var tour models.Tour

	err := db.Get(&tour, "SELECT tour_id, tour_uuid, tour_name, page_url, page_name, step_order, element_id, title, content, placement, completed, created_at, updated_at, deleted_at FROM tours WHERE tour_uuid = $1", id)
	if err != nil {
		return models.Tour{}, err
	}
	return tour, nil
}

// ShowTourByName retrieves a tour by its name
func ShowTourByName(tourName string) (*models.Tour, error) {
    var tour models.Tour

    // Debugging: cetak tourName yang diterima dari request
    log.Println("Mencari tour dengan nama:", tourName)

    // Gunakan query untuk mencari nama tour dengan mengabaikan kapitalisasi dan spasi
    err := db.Get(&tour, `
        SELECT * 
        FROM tours 
        WHERE LOWER(REPLACE(tour_name, ' ', '')) = LOWER(REPLACE($1, ' ', '')) 
        AND deleted_at IS NULL`, tourName)

    if err != nil {
        // Jika tidak ditemukan, kembalikan error
        log.Println("Error saat mencari tour:", err)
        return nil, err
    }

    // Kembalikan hasil jika ditemukan
    return &tour, nil
}

// GetTourGuideFromDB retrieves a tour guide by ID from the database
func GetTourGuide(id uint) (*models.Tour, error) {
	var tourGuide models.Tour
	row := db.QueryRow(`
		SELECT
			tour_id, tour_uuid, tour_name, page_url, page_name, step_order, element_id, title, content, placement, completed, created_at, updated_at, deleted_at
		FROM tours WHERE tour_id = ?`, id)

	err := row.Scan(
		&tourGuide.TourID,
		&tourGuide.TourUUID,
		&tourGuide.TourName,
		&tourGuide.PageURL,
		&tourGuide.PageName,
		&tourGuide.StepOrder,
		&tourGuide.ElementID,
		&tourGuide.Title,
		&tourGuide.Content,
		&tourGuide.Placement,
		&tourGuide.Completed,
		&tourGuide.CreatedAt,
		&tourGuide.UpdatedAt,
		&tourGuide.DeletedAt,
	)
	if err != nil {
		log.Println("Error retrieving tour guide:", err)
		return nil, err
	}
	return &tourGuide, nil
}

// GetTourGuideByID retrieves a tour guide by ID
func GetTourGuideByID(id uint) (*models.Tour, error) {
	var tourGuide models.Tour

	// Prepare the SQL query to get the tour guide by ID
	query := `
		SELECT
			tour_id, tour_uuid, tour_name, page_url, page_name, step_order, element_id, title, content, placement, completed, created_at, updated_at, deleted_at
		FROM tours
		WHERE tour_id = $1` // Menggunakan $1 untuk PostgreSQL

	// Execute the query
	row := db.QueryRow(query, id)

	// Scan the result into the tourGuide struct
	err := row.Scan(
		&tourGuide.TourID,
		&tourGuide.TourUUID,
		&tourGuide.TourName,
		&tourGuide.PageURL,
		&tourGuide.PageName,
		&tourGuide.StepOrder,
		&tourGuide.ElementID,
		&tourGuide.Title,
		&tourGuide.Content,
		&tourGuide.Placement,
		&tourGuide.Completed,
		&tourGuide.CreatedAt,
		&tourGuide.UpdatedAt,
		&tourGuide.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			// Return nil if no rows were found
			return nil, nil
		}
		// Log and return the error if something else went wrong
		log.Println("Error retrieving tour guide:", err)
		return nil, err
	}

	return &tourGuide, nil
}

// Fungsi untuk cek apakah tour_name sudah ada di database
func GetTourUUIDByName(tourName string) (string, error) {
	var tourUUID string
	err := db.Get(&tourUUID, `SELECT tour_uuid FROM tours WHERE tour_name = $1 LIMIT 1`, tourName)
	if err != nil {
		return "", err // Return error jika tour_name tidak ditemukan
	}
	return tourUUID, nil // Return tour_uuid yang ada
}

// GetAllTourGuidesFromDB retrieves all tour guides from the database
func GetAllTourGuides() ([]models.Tour, error) {
	tourGuide := []models.Tour{}

	rows, errSelect := db.Queryx("SELECT tour_uuid, tour_name, page_url, page_name, step_order, element_id, title, content, placement, completed, created_at, updated_at from tours WHERE deleted_at IS NULL")
	if errSelect != nil {
		return nil, errSelect
	}

	for rows.Next() {
		place := models.Tour{}
		rows.StructScan(&place)
		tourGuide = append(tourGuide, place)
	}

	return tourGuide, nil
}

// GetTourWithStepsByName mengambil semua langkah berdasarkan tour_name
func GetTourWithStepsByName(tourName string) (*models.TourWithSteps, error) {
    var tourWithSteps models.TourWithSteps

    fmt.Println("Tour Name: ", tourName)

    // Ambil data tour berdasarkan tour_name
    err := db.Get(&tourWithSteps, `
        SELECT 
            tour_name,
            page_name
        FROM
            tours
        WHERE
            LOWER(REPLACE(tour_name, ' ', '')) = LOWER(REPLACE($1, ' ', ''))
    `, tourName)

    // Cek error saat mengambil tour
    if err != nil {
        fmt.Println("Error getting tour:", err)
        return nil, err
    }

    // Ambil semua langkah (steps) berdasarkan tour_name
    err = db.Select(&tourWithSteps.Steps, `
        SELECT 
            tour_uuid,
            step_order,
            element_id,
            title,
            content,
            placement,
            buttons
        FROM
            tours
        WHERE
            LOWER(REPLACE(tour_name, ' ', '')) = LOWER(REPLACE($1, ' ', ''))
			AND deleted_at IS NULL  -- Filter langkah yang belum dihapus
        ORDER BY step_order
    `, tourName)

    // Cek error saat mengambil steps
    if err != nil {
        fmt.Println("Error getting steps:", err)
        return nil, err
    }

    return &tourWithSteps, nil
}

// Function sorting Ascending dan Descending untuk multikolom
func GetTourGuideSort(sortBy string, orderBy string) ([]models.Tour, error) {
	tourGuide := []models.Tour{}

	// Menentukan query sorting secara dinamis
	query := fmt.Sprintf("SELECT tour_uuid, tour_name, page_url, page_name, step_order, element_id, title, content, placement, completed, created_at, updated_at from tours WHERE deleted_at IS NULL ORDER BY %s %s", sortBy, orderBy)

	// Menjalankan query dan mengambil hasilnya
	rows, err := db.Queryx(query)
	if err != nil {
		return nil, err
	}

	// Memindai hasil query ke dalam struct Tour
	for rows.Next() {
		place := models.Tour{}
		err := rows.StructScan(&place)
		if err != nil {
			return nil, err
		}
		tourGuide = append(tourGuide, place)
	}

	return tourGuide, nil
}


// func UpdateTourGuide(updateTour models.Tour, id string, username string) (models.Tour, error) {
// 	currentTime := time.Now()

// 	_, err := db.NamedExec("UPDATE tours SET tour_name = :tour_name, page_name = :page_name, step_order = :step_order, element_id = :element_id, title = :title, content = :content, placement = :placement, buttons = :buttons, updated_at = :updated_at WHERE tour_uuid = :id", map[string]interface{}{
// 		"tour_name":  updateTour.TourName,
// 		// "page_url":   updateTour.PageURL,
// 		"page_name":  updateTour.PageName,
// 		"step_order": updateTour.StepOrder,
// 		"element_id": updateTour.ElementID,
// 		"title":      updateTour.Title,
// 		"content":    updateTour.Content,
// 		"placement":  updateTour.Placement,
// 		"buttons":    updateTour.Buttons,
// 		"updated_at": currentTime,
// 		"id":         id,
// 	})
// 	if err != nil {
// 		log.Println("Error updating tour guide", err)
// 		return models.Tour{}, err
// 	}
// 	return updateTour, nil
// }

func UpdateTourGuide(updateTour models.Tour, id string, username string) (models.Tour, error) {
    currentTime := time.Now()

    // Menghilangkan spasi dan merubah tour_name menjadi lowercase
    // tourName := strings.ToLower(strings.Replace(updateTour.TourName, " ", "", -1))

    // Update query untuk memperbarui data berdasarkan id dan tour_name yang sudah diformat
    _, err := db.NamedExec(`
        UPDATE tours
        SET 
            tour_name = :tour_name,
            page_name = :page_name,
            step_order = :step_order,
            element_id = :element_id,
            title = :title,
            content = :content,
            placement = :placement,
            buttons = :buttons,
            updated_at = :updated_at
        WHERE 
            tour_uuid = :id
    `, map[string]interface{}{
        "tour_name":  updateTour.TourName,  // Pastikan tour_name diformat
        "page_name":  updateTour.PageName,
        "step_order": updateTour.StepOrder,
        "element_id": updateTour.ElementID,
        "title":      updateTour.Title,
        "content":    updateTour.Content,
        "placement":  updateTour.Placement,
        "buttons":    updateTour.Buttons,
        "updated_at": currentTime,
        "id":         id,
    })
    
    if err != nil {
        log.Println("Error updating tour guide", err)
        return models.Tour{}, err
    }

    return updateTour, nil
}


// DeleteTourGuide deletes a tour guide entry from the database
func DeleteTourGuideByName(tourName, username string) error {
	currentTime := time.Now()
	var tourID int64
    err := db.Get(&tourID, `
        SELECT tour_id 
        FROM tours 
        WHERE LOWER(REPLACE(tour_name, ' ', '')) = LOWER(REPLACE($1, ' ', ''))
    `, tourName)
    if err != nil {
        log.Println("Error getting tour_id:", err)
        return err
    }

    // Jika tour ditemukan, lakukan soft delete dengan mengisi field `deleted_at`
    result, err := db.NamedExec(`
        UPDATE tours 
        SET deleted_at = :deleted_at 
        WHERE LOWER(REPLACE(tour_name, ' ', '')) = LOWER(REPLACE(:tour_name, ' ', ''))
    `, map[string]interface{}{
        "deleted_at": currentTime,
        "tour_name":  tourName,
    })
    if err != nil {
        log.Println("Error deleting tour guide:", err)
        return err
    }

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func DeleteTourStep(id, username string) error {
	currentTime := time.Now()
	var tourID int64
	err := db.Get(&tourID, `
		SELECT tour_id 
		FROM tours 
		WHERE tour_uuid = $1
	`, id)
	if err != nil {
		log.Println("Error getting tour_id:", err)
		return err
   }

	// Jika tour ditemukan, lakukan soft delete dengan mengisi field `deleted_at`
	result, err := db.NamedExec(`
		UPDATE tours 
		SET deleted_at = :deleted_at 
		WHERE tour_uuid = :tour_uuid
	`, map[string]interface{}{
		"deleted_at": currentTime,
		"tour_uuid":  id,
	})
	if err != nil {
		log.Println("Error deleting tour guide:", err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

