package webhook

import (
	"net/http"

	"github.com/slok/go-http-metrics/middleware"
)

// measuredHandler wraps a handler and measures the request handled
// by this handler.
func (h handler) measuredHandler(next http.Handler) http.Handler {
	mdlw := middleware.New(middleware.Config{Recorder: h.metrics})
	return mdlw.Handler("", next)
}
