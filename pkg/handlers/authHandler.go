package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/chigaji/file_sharing_service/pkg/storage"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

var jwtKey = []byte("e40adb88-d03b-48ce-a2ea-374616eb")

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func LoginHandser(c echo.Context) error {
	var creds Credentials

	if err := json.NewDecoder(c.Request().Body).Decode(&creds); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	storedPassword, err := storage.GetUserPassword(creds.Username)

	if err != nil || storedPassword != creds.Password {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")
	}

	tokenString, err := generateJwtToken(creds.Username)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not generate token")
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": tokenString,
	})
}

func generateJwtToken(username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		return "", nil
	}
	return tokenString, nil
}

func ValidateJwtToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}
	return claims, nil
}
