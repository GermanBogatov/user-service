package metrics

import (
	"github.com/GermanBogatov/user-service/internal/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"time"
)

type (
	Status                   string
	DbRequestType            string
	Validation               string
	IntegrationServiceName   string
	IntegrationServiceMethod string
)

const (
	OkStatus   Status = "ok"
	FailStatus Status = "fail"

	Postgres DbRequestType = "postgres"
	Cache    DbRequestType = "cache"
)

var (
	requestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name:      "total_request",
		Namespace: config.Namespace,
		Help:      "Number of processed incoming requests",
	}, []string{"status", "method", "pattern"})

	requestDurationSeconds = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:      "duration_request",
		Namespace: config.Namespace,
		Help:      "Duration request processing time",
	}, []string{"method", "pattern"})

	requestsPsqlPathTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name:      "total_psql_requests",
		Namespace: config.Namespace,
		Help:      "Number of path requests to postgresql",
	}, []string{"method", "status"})

	requestsPsqlDurationsPerMethod = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:      "duration_psql_requests_per_method",
		Namespace: config.Namespace,
		Help:      "duration postgresql requests for a single method",
	}, []string{"database", "method"})
)

// IncRequestTotal - метод для увеличения счетчика запросов к сервису тегов
func IncRequestTotal(status Status, method, pattern string) {
	requestsTotal.With(prometheus.Labels{"status": string(status), "method": method, "pattern": pattern}).Inc()
}

// ObserveRequestDurationSeconds - метод для вычисления времени выполнения запроса к сервису toxic-message
func ObserveRequestDurationSeconds(method, pattern string) func() {
	start := time.Now()
	return func() {
		duration := time.Since(start).Seconds()
		requestDurationSeconds.With(prometheus.Labels{"method": method, "pattern": pattern}).Observe(duration)
	}
}

// IncRequestTotalDB - метод для увеличения счетчика запросов к базе данных
func IncRequestTotalDB(path DbRequestType, status Status) {
	requestsPsqlPathTotal.With(prometheus.Labels{"method": string(path), "status": string(status)}).Inc()
}

// ObserveRequestDurationPerMethodDB - метод для вычисления времени выполнения запроса к базе данных
func ObserveRequestDurationPerMethodDB(db, method DbRequestType) func() {
	start := time.Now()
	return func() {
		duration := time.Since(start).Seconds()
		requestsPsqlDurationsPerMethod.With(prometheus.Labels{"database": string(db), "method": string(method)}).Observe(duration)
	}
}
