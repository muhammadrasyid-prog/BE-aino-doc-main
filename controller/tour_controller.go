// controller/tour_guide.go

package controller

import (
	"database/sql"
	"document/models"
	"document/service"
	"encoding/json"
	"fmt"

	// "go/token"
	"log"
	"net/http"

	// "regexp"
	// "strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

// AddTourGuide adds a new tour guide entry
func AddTourGuide(c echo.Context) error {
	// Ambil token dari header Authorization
	tokenString := c.Request().Header.Get("Authorization")
	secretKey := "secretJwToken" // Jika Anda menggunakan secretKey, pastikan ini adalah kunci yang tepat

	if tokenString == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak ditemukan!",
			"status":  false,
		})
	}

	// Periksa apakah tokenString mengandung "Bearer "
	if !strings.HasPrefix(tokenString, "Bearer ") {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}

	// Hapus "Bearer " dari tokenString
	tokenOnly := strings.TrimPrefix(tokenString, "Bearer ")

	// Dekripsi token JWE
	decrypted, err := DecryptJWE(tokenOnly, secretKey)
	if err != nil {
		fmt.Println("Gagal mendekripsi token:", err)
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}

	// Token yang sudah dideskripsi
	fmt.Println("Token yang sudah dideskripsi:", decrypted)

	var tourGuide models.Tour
	if err := c.Bind(&tourGuide); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}
	
	// Save the tour guide to the database
	err = service.AddTourGuide(tourGuide)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to add tour guide"})
	}
	return c.JSON(http.StatusOK, tourGuide)
}

func ShowTourById(c echo.Context) error {
	id := c.Param("id")

	var getTour models.Tour

	getTour, err := service.ShowTourById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Print(err)
			response := models.Response{
				Code:	404,
				Message: "Tour tidak ditemukan!",
				Status:	false,
			}
			return c.JSON(http.StatusNotFound, response)
		} else {
			log.Print(err)
			return c.JSON(http.StatusInternalServerError, &models.Response{
				Code:	500,
				Message: "Terjadi Kesalahan Internal pada Server. Mohon coba beberapa saat lagi!",
				Status:	false,
			}) 
		}
	}
	return c.JSON(http.StatusOK, getTour)
}

// UpdateTourGuide updates an existing tour guide entry
func UpdateTourGuide(c echo.Context) error {
    id := c.Param("id")

    // Mendapatkan data sebelumnya berdasarkan ID
    previousContent, errGet := service.ShowTourById(id)
    if errGet != nil {
        log.Print(errGet)
        return c.JSON(http.StatusNotFound, &models.Response{
            Code:    404,
            Message: "Tour tidak ditemukan!",
            Status:  false,
        })
    }

    // Mendapatkan token dari header Authorization
    tokenString := c.Request().Header.Get("Authorization")
    secretKey := "secretJwToken"

    if tokenString == "" {
        return c.JSON(http.StatusUnauthorized, map[string]interface{}{
            "code":    401,
            "message": "Token tidak ditemukan!",
            "status":  false,
        })
    }

    // Periksa apakah token diawali dengan "Bearer "
    if !strings.HasPrefix(tokenString, "Bearer ") {
        return c.JSON(http.StatusUnauthorized, map[string]interface{}{
            "code":    401,
            "message": "Token tidak valid!",
            "status":  false,
        })
    }

    // Menghilangkan prefix "Bearer " dari tokenString
    tokenOnly := strings.TrimPrefix(tokenString, "Bearer ")

	// Dekripsi token JWE
    decrypted, err := DecryptJWE(tokenOnly, secretKey)
    if err != nil {
        fmt.Println("Gagal mendekripsi token:", err)
        return c.JSON(http.StatusUnauthorized, map[string]interface{}{
            "code":    401,
            "message": "Token tidak valid!",
            "status":  false,
        })
    }

    // Klaim JWT
    var claims JwtCustomClaims

    // Mengurai klaim dari token hasil dekripsi
    errJ := json.Unmarshal([]byte(decrypted), &claims)
    if errJ != nil {
        fmt.Println("Gagal mengurai klaim:", errJ)
        return c.JSON(http.StatusUnauthorized, map[string]interface{}{
            "code":    401,
            "message": "Token tidak valid!",
            "status":  false,
        })
    }
	userName := c.Get("user_name").(string)

    // Token valid, bisa lanjutkan logika lain seperti mengecek expiration, issuer, dll.
    fmt.Println("Token berhasil diurai:", claims)

	// Bind the tour guide data from the request body
	if userName == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Invalid token atau token tidak ditemukan!",
			"status":  false,
		})
	}

	var editTour models.Tour
	if err := c.Bind(&editTour); err != nil {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Code:    400,
			Message: "Invalid input!",
			Status:  false,
		})
	}

	err = c.Validate(&editTour)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Code:    422,
			Message: "Data tidak boleh kosong!",
			Status:  false,
		})
	}

	_, errService := service.UpdateTourGuide(editTour, id, userName)
	if errService != nil {
		log.Print(errService)
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Code:    500,
			Message: "Gagal memperbarui tour guide!",
			Status:  false,
		})
	}
		
    // Berhasil melakukan pembaruan
	log.Println(previousContent)
    return c.JSON(http.StatusOK, &models.Response{
        Code:    200,
        Message: "Product berhasil diperbarui!",
        Status:  true,
    })
}


// DeleteTourGuide deletes a tour guide entry
// DeleteTourGuide deletes a tour guide entry based on tour_name
func DeleteTourGuideByName(c echo.Context) error {
	tokenString := c.Request().Header.Get("Authorization")
	secretKey := "secretJwToken"

	// Cek apakah token ada
	if tokenString == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak ditemukan",
			"status":  false,
		})
	}

	// Cek apakah token valid
	if !strings.HasPrefix(tokenString, "Bearer ") {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid",
			"status":  false,
		})
	}

	tokenOnly := strings.TrimPrefix(tokenString, "Bearer ")

	// Dekripsi token
	decrypted, err := DecryptJWE(tokenOnly, secretKey)
	if err != nil {
		fmt.Println("Gagal mendekripsi token:", err)
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid",
			"status":  false,
		})
	}

	// Parsing claims dari token
	var claims JwtCustomClaims
	errJ := json.Unmarshal([]byte(decrypted), &claims)
	if errJ != nil {
		fmt.Println("Gagal mengurai klaim:", errJ)
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid",
			"status":  false,
		})
	}

	userName := c.Get("user_name").(string) // Ambil username dari context
	tourName := c.Param("tour_name")        // Ambil parameter tour_name dari URL

	// Periksa apakah tour dengan nama tersebut ada
	_, errGet := service.ShowTourByName(tourName) // Fungsi ini harus diperbarui di service
	if errGet != nil {
		log.Println("Kesalahan saat penghapusan:", errGet)
		return c.JSON(http.StatusNotFound, &models.Response{
			Code:    404,
			Message: "Gagal menghapus tour. Tour tidak ditemukan!",
			Status:  false,
		})
	}

	// Hapus tour berdasarkan tour_name
	errService := service.DeleteTourGuideByName(tourName, userName) // Fungsi ini harus diperbarui di service
	if errService != nil {
		log.Println("Kesalahan saat penghapusan:", errService)
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Code:    500,
			Message: "Terjadi kesalahan internal pada server. Mohon coba beberapa saat lagi!",
			Status:  false,
		})
	}

	// Jika berhasil, kirimkan response sukses
	return c.JSON(http.StatusOK, &models.Response{
		Code:    200,
		Message: "Tour berhasil dihapus!",
		Status:  true,
	})
}

func DeleteTourStep(c echo.Context) error {
	tokenString := c.Request().Header.Get("Authorization")
	secretKey := "secretJwToken"

	// Cek apakah token ada
	if tokenString == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak ditemukan",
			"status":  false,
		})
	}

	// Cek apakah token valid
	if !strings.HasPrefix(tokenString, "Bearer ") {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid",
			"status":  false,
		})
	}

	tokenOnly := strings.TrimPrefix(tokenString, "Bearer ")

	// Dekripsi token
	decrypted, err := DecryptJWE(tokenOnly, secretKey)
	if err != nil {
		fmt.Println("Gagal mendekripsi token:", err)
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid",
			"status":  false,
		})
	}

	// Parsing claims dari token
	var claims JwtCustomClaims
	errJ := json.Unmarshal([]byte(decrypted), &claims)
	if errJ != nil {
		fmt.Println("Gagal mengurai klaim:", errJ)
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid",
			"status":  false,
		})
	}

	userName := c.Get("user_name").(string) // Ambil username dari context	
	id := c.Param("id") // Ambil parameter id dari URL

	// Periksa apakah tour dengan uuid tersebut ada
	_, errGet := service.ShowTourById(id) // Fungsi ini harus diperbarui di service
	if errGet != nil {
		log.Println("Kesalahan saat penghapusan:", errGet)
		return c.JSON(http.StatusNotFound, &models.Response{
			Code:    404,
			Message: "Gagal menghapus tour. Tour tidak ditemukan!",
			Status:  false,
		})
	}

	// Hapus tour berdasarkan uuid
	errService := service.DeleteTourStep(id, userName) // Fungsi ini harus diperbarui di service
	if errService != nil {
		log.Println("Kesalahan saat penghapusan:", errService)
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Code:    500,
			Message: "Terjadi kesalahan internal pada server. Mohon coba beberapa saat lagi!",
			Status:  false,
		})
	}

	// jika berhasil, kirimkan response sukses
	return c.JSON(http.StatusOK, &models.Response{
		Code:    200,
		Message: "Tour berhasil dihapus!",
		Status:  true,
	})
   
}

// GetTourGuideByID retrieves a tour guide by ID
func GetTourGuideByID(c echo.Context) error {
	id := c.Param("id")

	var getTour models.Tour

	getTour, err := service.ShowTourById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Print(err)
			response := models.Response{
				Code:    404,
				Message: "Tour tidak ditemukan!",
				Status:  false,
			}
			return c.JSON(http.StatusNotFound, response)
		} else {
			log.Print(err)
			return c.JSON(http.StatusInternalServerError, &models.Response{
				Code:    500,
				Message: "Terjadi kesalahan internal pada server. Mohon coba beberapa saat lagi!",
				Status:  false,
			})
		}
	}

	return c.JSON(http.StatusOK, getTour)

}

// GetAllTourGuides retrieves all tour guides
func GetAllTourGuides(c echo.Context) error {
	tourGuides, err := service.GetAllTourGuides()
	if err != nil {
		log.Print(err)
		response := models.Response{
			Code:    500,
			Message: "Gagal mengambil data tour guide!",
			Status:  false,
		}
		return c.JSON(http.StatusInternalServerError, response)
	}
	return c.JSON(http.StatusOK, tourGuides)
}

// GetSpecAllTourGuides retrieves all tour steps grouped by tour_name and page_name
func GetTourWithStepsByName(c echo.Context) error {
    tourName := c.Param("tour_name")
    log.Println("Received tour_name:", tourName)

    // Panggil service untuk mendapatkan data
	tourWithSteps, err := service.GetTourWithStepsByName(tourName)
    if err != nil {
        // Jika tidak ditemukan, kembalikan 404
        if err == sql.ErrNoRows {
            return c.JSON(http.StatusNotFound, map[string]interface{}{
                "message": "Tour not found",
            })
        }
        // Jika ada error lainnya
        return c.JSON(http.StatusInternalServerError, map[string]interface{}{
            "message": "Gagal mendapatkan data tour",
            "error":   err.Error(),
        })
    }

    return c.JSON(http.StatusOK, tourWithSteps)
}

// Function GET sorting Ascending dan Descending untuk multikolom
func GetTourGuideSort(c echo.Context) error {
	sortBy := c.QueryParam("sortBy") // Kolom yang digunakan untuk sorting
	orderBy := c.QueryParam("orderBy") // Urutan sorting (asc/desc)

	// Daftar kolom yang diperbolehkan untuk sorting
	validColumns := map[string]bool{
		"tour_name": true,
		"created_at": true,
		"updated_at": true,
		"step_order": true,
	}

	// Validasi apakah kolom yang diterima sesuai dengan kolom yang diperbolehkan
	if _, valid := validColumns[sortBy]; !valid {
		return c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "Invalid sort column",
			Status:  false,
		})
	}

	// Validasi untuk order (asc/desc)
	if orderBy != "asc" && orderBy != "desc" {
		return c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "Invalid order parameter. Allowed values: asc, desc",
			Status:  false,
		})
	}

	// Panggil service untuk mengambil data berdasarkan sorting
	tourGuides, err := service.GetTourGuideSort(sortBy, orderBy)
	if err != nil {
		log.Print(err)
		response := models.Response{
			Code:    500,
			Message: "Gagal mengambil data tour guide!",
			Status:  false,
		}
		return c.JSON(http.StatusInternalServerError, response)
	}
	return c.JSON(http.StatusOK, tourGuides)
}

