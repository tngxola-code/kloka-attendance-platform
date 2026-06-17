package metrics

import (
    "net/http"
    "strconv"
    "time"
)

type statusRecorder struct {
    http.ResponseWriter
    status int
}

func (r *statusRecorder) WriteHeader(code int) {
    r.status = code
    r.ResponseWriter.WriteHeader(code)
}

func RequestMetrics(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        recorder := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
        next.ServeHTTP(recorder, r)

        route := r.URL.Path
        method := r.Method
        status := strconv.Itoa(recorder.status)

        RequestsTotal.WithLabelValues(method, route, status).Inc()
        RequestDuration.WithLabelValues(method, route).Observe(time.Since(start).Seconds())
    })
}
