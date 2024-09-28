package helpers

import (
	"crypto/sha256"
	"fmt"
	"github.com/GermanBogatov/user-service/internal/config"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"net/http"
)

// GetUuidFromPath - получение uuid-значения из пути запроса
func GetUuidFromPath(r *http.Request, key string) (uuid.UUID, error) {
	var uuidByte uuid.UUID
	var err error

	uuidString := chi.URLParam(r, key)
	if uuidString == "" {
		return uuidByte, fmt.Errorf("value not found for key=[%s]", key)
	}

	uuidByte, err = uuid.Parse(uuidString)
	if err != nil {
		return uuidByte, errors.Wrap(err, fmt.Sprintf("invalid parse uuid=[%s] from path", uuidString))
	}

	return uuidByte, nil
}

// GetStringFromPath - получение строки из пути
func GetStringFromPath(r *http.Request, key string) (string, error) {
	value := chi.URLParam(r, key)
	if value == "" {
		return "", fmt.Errorf("value not found for key=[%s]", key)
	}

	return value, nil
}

// GetStringWithDefaultFromQuery - получение строкового значения из query. Если его нет, то заменять дефолтным
func GetStringWithDefaultFromQuery(r *http.Request, key, defaultParam string) string {
	param := r.URL.Query().Get(key)
	if len(param) == 0 {
		return defaultParam
	}

	return param
}

// GetOptionalParamFromQuery - получение значения из query, либо nil при его отсутствии
func GetOptionalParamFromQuery(r *http.Request, key string) *string {
	param := r.URL.Query().Get(key)
	if len(param) == 0 {
		return nil
	}

	return &param
}

func GeneratePasswordHash(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(config.JWTTokenSalt)))
}
