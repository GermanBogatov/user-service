package http

import (
	"github.com/GermanBogatov/user-service/pkg/logging"
	"net/http"
)

// live - хэндлер для проверки приложения
func live(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("Healthy"))
	if err != nil {
		logging.Errorf("error write response: %s", err)
	}
}

// readiness - хэндлер для проверки подключаемых бд и так далее
func readiness(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("Healthy"))
	if err != nil {
		logging.Errorf("error write response: %s", err)
	}
}
