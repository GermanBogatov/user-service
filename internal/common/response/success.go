package response

import (
	"github.com/GermanBogatov/user-service/pkg/logging"
	"net/http"
)

// RespondSuccessCreate - метод по возврату ответа с успешным созданием
func RespondSuccessCreate(w http.ResponseWriter, resp ViewResponse) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err := w.Write(resp.Marshal())
	if err != nil {
		logging.Errorf("error write response: %s", err)
	}

	return nil
}

// RespondSuccess - метод по возврату успешного ответа
func RespondSuccess(w http.ResponseWriter, resp ViewResponse) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(resp.Marshal())
	if err != nil {
		logging.Errorf("error write response: %s", err)
	}

	return nil
}
