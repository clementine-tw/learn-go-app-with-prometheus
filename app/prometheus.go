package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// API total request counter
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "API total request",
		},
		[]string{"method", "endpoint", "status"},
	)
	httpRequestDurationSeconds = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP Request Latency Distribution (s)",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)
)

type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *ResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func prometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrappedWriter := &ResponseWriter{
			ResponseWriter: w,
		}
		next.ServeHTTP(wrappedWriter, r)

		httpRequestsTotal.WithLabelValues(
			r.Method,
			r.Pattern, // 用"Pattern"來避免用到用戶ID這類動態數值作為標籤，降低基數
			strconv.Itoa(wrappedWriter.statusCode),
		).Inc()
		httpRequestDurationSeconds.WithLabelValues(
			r.Method,
			r.Pattern,
			strconv.Itoa(wrappedWriter.statusCode),
		).Observe(time.Since(start).Seconds())
	})
}
