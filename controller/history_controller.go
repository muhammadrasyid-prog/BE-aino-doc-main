package controller

import (
	// "document/models"
	"document/service"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func GetTimelineHistory(c echo.Context) error {
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

	// Ambil data history dari service
	history, err := service.GetTimelineHistory(db.DB)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "Gagal mengambil data history",
			"status":  false,
		})
	}

	return c.JSON(http.StatusOK, history)
}
