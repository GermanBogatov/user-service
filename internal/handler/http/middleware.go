package http

import (
	"github.com/GermanBogatov/user-service/internal/common/apperror"
	"github.com/GermanBogatov/user-service/internal/common/metrics"
	"github.com/GermanBogatov/user-service/internal/common/response"
	"github.com/GermanBogatov/user-service/internal/config"
	"github.com/GermanBogatov/user-service/pkg/logging"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type appHandler func(http.ResponseWriter, *http.Request) error

func appMiddleware(h appHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		routeContext := chi.RouteContext(r.Context())
		pattern := routeContext.RoutePattern()
		defer metrics.ObserveRequestDurationSeconds(method, pattern)()
		// todo remove/fix integration
		headerApiKey := r.Header.Get(config.XUserID)
		if headerApiKey == "" && routeContext.RoutePatterns[0] != integrationV2+"/*" {
			logging.Errorf("required X-USER-ID in header to path: %s", r.URL.Path)
			metrics.IncRequestTotal(metrics.FailStatus, method, pattern)
			response.RespondError(w, r, apperror.UnauthorizedError(apperror.ErrRequiredXUserID))

			return
		}

		err := h(w, r)
		if err != nil {
			metrics.IncRequestTotal(metrics.FailStatus, method, pattern)
			response.RespondError(w, r, err)
			return
		}

		metrics.IncRequestTotal(metrics.OkStatus, method, pattern)
	}
}
