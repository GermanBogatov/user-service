package apperror

import "github.com/pkg/errors"

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrRequiredXUserID = errors.New("не задан X-USER-ID в хедере")
)

const (
	ErrType500 = "INTERNAL_SERVER_ERROR"
	ErrType400 = "INVALID_CONTENT_FIELD"
	ErrType404 = "NOT_FOUND"
	ErrType401 = "UNAUTHORIZED"
	ErrType409 = "CONFLICT"
)
