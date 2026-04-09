package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"

	"web_backend/internal/app/repository"
	"web_backend/internal/app/serializer"
)

func (h *Handler) CreateUser(ctx *gin.Context) {
	var j serializer.SignUpRequest
	if err := ctx.BindJSON(&j); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	if j.Login == "" {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("field 'login' is required"))
		return
	}
	if j.Password == "" {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("field 'password' is required"))
		return
	}
	user, err := h.Repository.CreateUser(j)
	if err != nil {
		if errors.Is(err, repository.ErrAlreadyExists) {
			h.errorHandler(ctx, http.StatusConflict, err)
		} else if errors.Is(err, repository.ErrNotAllowed) {
			h.errorHandler(ctx, http.StatusForbidden, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}
	ctx.Header("Location", "/api/users")
	ctx.JSON(http.StatusCreated, serializer.SignUpResponseFromUser(user))
}

func (h *Handler) SignIn(ctx *gin.Context) {
	var j serializer.SignInRequest
	if err := ctx.BindJSON(&j); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	if j.Login == "" {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("field 'login' is required"))
		return
	}
	if j.Password == "" {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("field 'password' is required"))
		return
	}
	token, err := h.Repository.SignIn(j)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) || err.Error() == "неверный логин или пароль" {
			h.errorHandler(ctx, http.StatusUnauthorized, fmt.Errorf("invalid login or password"))
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *Handler) SignOut(ctx *gin.Context) {
	tokenString := extractTokenFromHeader(ctx.Request)
	if tokenString == "" {
		h.errorHandler(ctx, http.StatusUnauthorized, errors.New("no token provided"))
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

	if err != nil || token == nil {
		h.errorHandler(ctx, http.StatusUnauthorized, err)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		h.errorHandler(ctx, http.StatusBadRequest, errors.New("invalid token claims"))
		return
	}

	ttl, err := tokenTTLFromClaims(claims)
	if err != nil {
		ctx.Status(http.StatusNoContent)
		return
	}

	userID, _ := claims["user_id"].(string)
	if userID == "" {
		userID = "unknown"
	}

	err = h.Repository.AddTokenToBlacklist(context.Background(), tokenString, ttl, userID)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func tokenTTLFromClaims(claims jwt.MapClaims) (time.Duration, error) {
	expVal, ok := claims["exp"]
	if !ok {
		return 0, errors.New("exp not present")
	}

	var expUnix int64
	switch v := expVal.(type) {
	case float64:
		expUnix = int64(v)
	case int64:
		expUnix = v
	case json.Number:
		i, err := v.Int64()
		if err != nil {
			return 0, err
		}
		expUnix = i
	default:
		return 0, fmt.Errorf("unsupported exp type %T", v)
	}

	expTime := time.Unix(expUnix, 0)
	ttl := time.Until(expTime)
	if ttl < 0 {
		return 0, errors.New("token already expired")
	}
	return ttl, nil
}
