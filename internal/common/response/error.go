package response

import (
	"github.com/GermanBogatov/user-service/internal/common/apperror"
	"github.com/GermanBogatov/user-service/pkg/logging"
	"github.com/pkg/errors"
	"net/http"
)

// The RespondError function helps to convert the business error to standardized JSON response
func RespondError(w http.ResponseWriter, r *http.Request, err error) {
	w.Header().Set("Content-Type", "application/json")
	var appErr *apperror.AppError
	if errors.As(err, &appErr) {
		logging.Errorf("code [%d] method [%s] path [%s]: %v", appErr.StatusCode, r.Method, r.URL.Path, err)
		w.WriteHeader(appErr.StatusCode)
		_, errWrite := w.Write(NewAppErrorResponse(appErr).Marshal())
		if errWrite != nil {
			logging.Errorf("error write response: %s", errWrite)
		}

		return
	}

	logging.Errorf("code [%d] method [%s] path [%s]: %v", http.StatusInternalServerError, r.Method, r.URL.Path, err)
	w.WriteHeader(http.StatusInternalServerError)
	_, errWrite := w.Write(NewAppErrorResponse(apperror.InternalServerError(err)).Marshal())
	if errWrite != nil {
		logging.Errorf("error write response: %s", errWrite)
	}
}

// NewAppErrorResponse - json ответ с кастомной ошибкой и статусом
func NewAppErrorResponse(err *apperror.AppError) *ViewResponse {
	return &ViewResponse{
		Code:      err.StatusCode,
		Error:     err.Error(),
		ErrorType: err.ErrType,
	}
}
