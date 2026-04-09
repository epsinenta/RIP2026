package handler

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

const bearerPrefix = "Bearer"

func corsAllowedOriginsFromEnv() []string {
	raw := strings.TrimSpace(os.Getenv("CORS_ALLOWED_ORIGINS"))
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func isAllowedCORSOrigin(origin string, extra []string) bool {
	if origin == "" {
		return true
	}
	static := []string{
		"http://localhost:3000",
		"http://127.0.0.1:3000",
		"http://localhost:8000",
		"http://127.0.0.1:8000",
		"http://localhost:8080",
		"http://127.0.0.1:8080",
		"tauri://localhost",
		"http://tauri.localhost",
		"https://tauri.localhost",
	}
	for _, o := range static {
		if origin == o {
			return true
		}
	}
	for _, o := range extra {
		if origin == o {
			return true
		}
	}
	if strings.HasPrefix(origin, "https://") && strings.HasSuffix(origin, ".github.io") {
		return true
	}
	return false
}

func CORSMiddleware() gin.HandlerFunc {
	extra := corsAllowedOriginsFromEnv()
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept, Origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, Authorization")

		if origin != "" && isAllowedCORSOrigin(origin, extra) {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func (h *Handler) ModeratorMiddleware(allowedRole bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := extractTokenFromHeader(c.Request)

		if tokenString == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		jwtKey := os.Getenv("JWT_KEY")
		if jwtKey == "" {
			jwtKey = "default-secret-key-change-in-production"
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, nil
			}
			return []byte(jwtKey), nil
		})

		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		blacklisted, err := h.Repository.IsTokenBlacklisted(context.Background(), tokenString)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if blacklisted {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		jwtIsModerator, ok := claims["is_moderator"].(bool)
		if !ok {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if allowedRole && !jwtIsModerator {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Set("user_id", userIDStr)
		c.Set("is_moderator", jwtIsModerator)
		c.Next()
	}
}

func (h *Handler) WithOptionalAuthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		prefix := "Bearer "

		if tokenString == "" || !strings.HasPrefix(tokenString, prefix) {
			c.Set("user_id", "")
			c.Next()
			return
		}

		tokenString = strings.TrimPrefix(tokenString, prefix)

		blacklisted, err := h.Repository.IsTokenBlacklisted(context.Background(), tokenString)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if blacklisted {
			c.Set("user_id", "")
			c.Next()
			return
		}

		jwtKey := os.Getenv("JWT_KEY")
		if jwtKey == "" {
			jwtKey = "default-secret-key-change-in-production"
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtKey), nil
		})

		if err != nil || token == nil || !token.Valid {
			c.Set("user_id", "")
			c.Next()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.Set("user_id", "")
			c.Next()
			return
		}

		userIDValue, exists := claims["user_id"]
		if !exists || userIDValue == nil {
			c.Set("user_id", "")
			c.Next()
			return
		}

		var userIDStr string
		switch v := userIDValue.(type) {
		case string:
			userIDStr = v
		case float64:
			userIDStr = fmt.Sprintf("%.0f", v)
		default:
			userIDStr = fmt.Sprintf("%v", v)
		}

		c.Set("user_id", userIDStr)
		c.Next()
	}
}

func extractTokenFromHeader(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")
	if bearerToken == "" {
		return ""
	}
	parts := strings.SplitN(bearerToken, " ", 2)
	if len(parts) != 2 || parts[0] != bearerPrefix {
		return ""
	}
	return parts[1]
}

func getUserID(ctx *gin.Context) (uint, error) {
	userIDVal, exists := ctx.Get("user_id")
	if !exists || userIDVal == nil {
		return 0, fmt.Errorf("user_id not found")
	}
	userIDStr, ok := userIDVal.(string)
	if !ok || userIDStr == "" {
		return 0, fmt.Errorf("user_id not found")
	}
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(userID), nil
}
