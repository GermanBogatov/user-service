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
	metricPath = "/metrics"
	healthPath = "/health"
	publicV1   = "/public/v1"
	authV1     = "/public/v1/auth"

	livePath       = "/live"
	readinessPath  = "/readiness"
	swaggerPattern = "/swagger-ui/*"
)

type Handler struct {
	userService service.IUser
	jwtService  service.IJWT
	cfg         *config.Config
}

func NewHandler(cfg *config.Config, userService service.IUser, jwtService service.IJWT) *Handler {
	return &Handler{
		userService: userService,
		jwtService:  jwtService,
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

	r.Route(authV1, func(r chi.Router) {
		r.Post("/sign-up", appMiddleware(h.SignUp))
		r.Post("/sign-in", appMiddleware(h.SignIn))
		r.Post("/refresh/{id}", appMiddleware(h.UpdateRefreshToken))
	})
	r.Route(publicV1, func(r chi.Router) {
		r.Post("/users/{id}", appMiddleware(h.GetUserByID))
	})

	return r
}
