package helpers

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"net/http"
)

// GetUuidFromHeader - получение uuid-значения из хедера
func GetUuidFromHeader(r *http.Request, key string) (uuid.UUID, error) {
	var (
		uuidByte uuid.UUID
		err      error
	)

	uuidString := r.Header.Get(key)
	if uuidString == "" {
		return uuidByte, fmt.Errorf("не было найдено значение по ключу=[%s]", key)
	}

	uuidByte, err = uuid.Parse(uuidString)
	if err != nil {
		return uuidByte, errors.Wrap(err, fmt.Sprintf("не удалось распарсить uuid=[%s] из хедера", uuidString))
	}

	return uuidByte, nil
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
