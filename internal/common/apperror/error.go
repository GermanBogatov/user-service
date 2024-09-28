package apperror

import (
	"github.com/pkg/errors"
	"net/http"
)

// AppError - структура ошибки приложения
type AppError struct {
	StatusCode int
	ErrType    string
	Err        error
}

// Error - вывод ошибки в строку
func (a *AppError) Error() string {
	return a.Err.Error()
}

// NewAppErr - создание ошибки приложения через ошибку типа `error`
func NewAppErr(statusCode int, errType string, err error) *AppError {
	return &AppError{
		StatusCode: statusCode,
		Err:        err,
		ErrType:    errType,
	}
}

// BadRequestError - ошибка c кодом 400
func BadRequestError(err error) *AppError {
	return NewAppErr(http.StatusBadRequest, ErrType400, err)
}

// ApplicationError - ошибка приложения
func ApplicationError(err error) *AppError {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}

	return InternalServerError(err)
}

// InternalServerError - ошибка c кодом 500
func InternalServerError(err error) *AppError {
	if errors.Is(err, ErrUserNotFound) {
		return NotFoundError(err)
	}

	if errors.Is(err, ErrUserIsExistWithEmail) {
		return ConflictError(err)
	}

	if errors.Is(err, ErrRefreshTokenNotFound) {
		return ConflictError(err)
	}

	return NewAppErr(http.StatusInternalServerError, ErrType500, err)

}

// UnauthorizedError - ошибка c кодом 401
func UnauthorizedError(err error) *AppError {
	return NewAppErr(http.StatusUnauthorized, ErrType401, err)
}

// ConflictError - ошибка c кодом 409
func ConflictError(err error) *AppError {
	return NewAppErr(http.StatusConflict, ErrType409, err)
}

// NotFoundError - ошибка c кодом 404
func NotFoundError(err error) *AppError {
	return NewAppErr(http.StatusNotFound, ErrType404, err)
}
