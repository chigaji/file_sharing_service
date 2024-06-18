package middleware

import (
	"net/http"

	"github.com/chigaji/file_sharing_service/pkg/handlers"
	"github.com/labstack/echo/v4"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")

		if authHeader == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Missing or invalid token")
		}

		tokenStr := authHeader[len("Bearer "):]

		claims, err := handlers.ValidateJwtToken(tokenStr)

		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid Token")
		}

		c.Set("username", claims.Username)
		return next(c)
	}
}
