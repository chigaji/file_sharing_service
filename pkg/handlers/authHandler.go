package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/chigaji/file_sharing_service/pkg/models"
	"github.com/chigaji/file_sharing_service/pkg/storage"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
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

func RegisterHandler(c echo.Context) error {
	var user models.User
	if err := c.Bind(&user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not hash password")
	}

	user.Password = string(hashedPassword)

	if err := storage.SaveUser(user); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not save user")
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "User successfully registered",
	})

}
func LoginHandser(c echo.Context) error {
	var creds Credentials

	if err := json.NewDecoder(c.Request().Body).Decode(&creds); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	fmt.Println(creds.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not hash password")
	}

	creds.Password = string(hashedPassword)

	storedPassword, err := storage.GetUserPassword(creds.Username)

	fmt.Println("pwd", creds.Password)
	fmt.Println("hashed pwd", storedPassword)

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
