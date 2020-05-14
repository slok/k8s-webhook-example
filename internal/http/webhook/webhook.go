package webhook

import (
	"fmt"
	"net/http"

	"github.com/slok/k8s-webhook-example/internal/log"
	"github.com/slok/k8s-webhook-example/internal/mark"
)

// Config is the handler configuration.
type Config struct {
	MetricsRecorder MetricsRecorder
	Marker          mark.Marker
	Logger          log.Logger
}

func (c *Config) defaults() error {
	if c.Marker == nil {
		return fmt.Errorf("marker is required")
	}

	if c.MetricsRecorder == nil {
		c.MetricsRecorder = dummyMetricsRecorder
	}

	if c.Logger == nil {
		c.Logger = log.Dummy
	}

	return nil
}

type handler struct {
	handler http.Handler
	marker  mark.Marker
	metrics MetricsRecorder
	logger  log.Logger
}

// New returns a new webhook handler.
func New(config Config) (http.Handler, error) {
	err := config.defaults()
	if err != nil {
		return nil, fmt.Errorf("handler configuration is not valid: %w", err)
	}

	mux := http.NewServeMux()

	h := handler{
		handler: mux,
		marker:  config.Marker,
		metrics: config.MetricsRecorder,
		logger:  config.Logger.WithKV(log.KV{"service": "webhook-handler"}),
	}

	// Register all the routes with our router.
	err = h.routes(mux)
	if err != nil {
		return nil, fmt.Errorf("could not register routes on handler: %w", err)
	}

	// Register middlware.
	h.registerRootMiddleware()

	return h, nil
}

func (h handler) registerRootMiddleware() {
	h.handler = h.measuredHandler(h.handler) // Add metrics middleware.
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.handler.ServeHTTP(w, r)
}
