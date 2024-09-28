package http

import (
	_ "github.com/GermanBogatov/user-service/docs"
	"github.com/GermanBogatov/user-service/internal/config"
	"github.com/GermanBogatov/user-service/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
)

const (
	metricPath     = "/metrics"
	healthPath     = "/health"
	integrationV2  = "/integration/v2"
	privateV2      = "/private/v2"
	publicV2       = "/public/v2"
	livePath       = "/live"
	readinessPath  = "/readiness"
	swaggerPattern = "/swagger-ui/*"
)

type Handler struct {
	userService service.IUser
	cfg         *config.Config
}

func NewHandler(cfg *config.Config, userService service.IUser) *Handler {
	return &Handler{
		userService: userService,
		cfg:         cfg,
	}
}

// InitRoutes - инициализация роутера приложения
func (h *Handler) InitRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	r.Handle(metricPath, promhttp.Handler())

	r.Route(healthPath, func(r chi.Router) {
		r.Get(livePath, live)
		r.Get(readinessPath, readiness)
	})
	r.Get(swaggerPattern, httpSwagger.Handler())

	r.Route(publicV2, func(r chi.Router) {
		r.Post("/user", appMiddleware(h.CreateUser))
	})
	return r
}
