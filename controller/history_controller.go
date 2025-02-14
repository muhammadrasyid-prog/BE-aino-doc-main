package controller

import (
	// "document/models"
	"document/service"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

// func GetTimelineHistory(c echo.Context) error {
// 	// Ambil token dan periksa role
// 	tokenString := c.Request().Header.Get("Authorization")
// 	if tokenString == "" {
// 		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
// 			"code":    401,
// 			"message": "Token tidak ditemukan!",
// 			"status":  false,
// 		})
// 	}

// 	decrypted, err := DecryptJWE(strings.TrimPrefix(tokenString, "Bearer "), "secretJwToken")
// 	if err != nil {
// 		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
// 			"code":    401,
// 			"message": "Token tidak valid!",
// 			"status":  false,
// 		})
// 	}

// 	var claims JwtCustomClaims
// 	err = json.Unmarshal([]byte(decrypted), &claims)
// 	if err != nil {
// 		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
// 			"code":    401,
// 			"message": "Token tidak valid!",
// 			"status":  false,
// 		})
// 	}

// 	// Hanya superadmin (SA) dan admin (A) yang boleh mengakses
// 	if claims.RoleCode != "SA" && claims.RoleCode != "A" {
// 		return c.JSON(http.StatusForbidden, map[string]interface{}{
// 			"code":    403,
// 			"message": "Akses ditolak!",
// 			"status":  false,
// 		})
// 	}

// 	// Ambil data history dari service
// 	history, err := service.GetTimelineHistory(db.DB)
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
// 			"code":    500,
// 			"message": "Gagal mengambil data history",
// 			"status":  false,
// 		})
// 	}

// 	return c.JSON(http.StatusOK, history)
// }

func GetRecentTimelineHistory(c echo.Context) error {
	// Ambil token dan periksa role
	tokenString := c.Request().Header.Get("Authorization")
	if tokenString == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak ditemukan!",
			"status":  false,
		})
	}

	decrypted, err := DecryptJWE(strings.TrimPrefix(tokenString, "Bearer "), "secretJwToken")
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}

	var claims JwtCustomClaims
	err = json.Unmarshal([]byte(decrypted), &claims)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}

	// Hanya superadmin (SA) dan admin (A) yang boleh mengakses
	if claims.RoleCode != "SA" && claims.RoleCode != "A" {
		return c.JSON(http.StatusForbidden, map[string]interface{}{
			"code":    403,
			"message": "Akses ditolak!",
			"status":  false,
		})
	}

	// Ambil data recent timeline dari service
	history, err := service.GetRecentTimelineHistory(db.DB)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "Gagal mengambil data history terbaru",
			"status":  false,
		})
	}

	return c.JSON(http.StatusOK, history)
}

func GetOlderTimelineHistory(c echo.Context) error {
	// Ambil token dan periksa role
	tokenString := c.Request().Header.Get("Authorization")
	if tokenString == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak ditemukan!",
			"status":  false,
		})
	}

	decrypted, err := DecryptJWE(strings.TrimPrefix(tokenString, "Bearer "), "secretJwToken")
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}

	var claims JwtCustomClaims
	err = json.Unmarshal([]byte(decrypted), &claims)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}

	// Hanya superadmin (SA) dan admin (A) yang boleh mengakses
	if claims.RoleCode != "SA" && claims.RoleCode != "A" {
		return c.JSON(http.StatusForbidden, map[string]interface{}{
			"code":    403,
			"message": "Akses ditolak!",
			"status":  false,
		})
	}

	// Ambil query parameter untuk pagination
	lastCreatedAtStr := c.QueryParam("last_created_at")
	if lastCreatedAtStr == "" {
		lastCreatedAtStr = time.Now().Format(time.RFC3339) // Default ke waktu sekarang
	}

	lastCreatedAt, err := time.Parse(time.RFC3339, lastCreatedAtStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "Format last_created_at harus dalam bentuk ISO 8601 (YYYY-MM-DDTHH:MM:SSZ)",
			"status":  false,
		})
	}

	// Ambil limit
	limitStr := c.QueryParam("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10 // Default limit
	}

	// Ambil older timeline dari service
	history, err := service.GetOlderTimelineHistory(db.DB, lastCreatedAt, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "Gagal mengambil data history lama",
			"status":  false,
		})
	}

	return c.JSON(http.StatusOK, history)
}
