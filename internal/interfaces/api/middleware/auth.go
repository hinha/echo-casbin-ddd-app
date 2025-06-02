package middleware

import (
	"net/http"
	"strings"

	"github.com/hinha/echo-casbin-ddd-app/internal/config"
	"github.com/hinha/echo-casbin-ddd-app/internal/infrastructure/auth"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

// JWTMiddleware creates a JWT middleware
func JWTMiddleware(config *config.Config) echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		SigningKey:  []byte(config.JWT.Secret),
		TokenLookup: "header:Authorization",
	})
}

// APIKeyMiddleware creates an API key middleware
func APIKeyMiddleware(config *config.Config, apiKeyService *auth.APIKeyService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			apiKey := c.Request().Header.Get(config.APIKey.HeaderName)
			if apiKey == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "API key is required"})
			}

			valid, _, err := apiKeyService.ValidateAPIKey(c.Request().Context(), apiKey)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
			}

			if !valid {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid API key"})
			}

			return next(c)
		}
	}
}

// CasbinMiddleware creates a Casbin middleware
func CasbinMiddleware(casbinService *auth.CasbinService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get user from JWT token
			user := c.Get("user")
			if user == nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
			}

			// Extract username from JWT claims
			claims, ok := user.(*auth.Claims)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
			}

			// Get request path and method
			path := c.Request().URL.Path
			method := c.Request().Method

			// Check if user has permission
			allowed, err := casbinService.Enforce(claims.Username, "default", path, method)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
			}

			if !allowed {
				return c.JSON(http.StatusForbidden, map[string]string{"error": "Forbidden"})
			}

			return next(c)
		}
	}
}

// ExtractTokenFromHeader extracts the token from the Authorization header
func ExtractTokenFromHeader(header string) string {
	parts := strings.Split(header, " ")
	if len(parts) == 2 && parts[0] == "Bearer" {
		return parts[1]
	}
	return ""
}
