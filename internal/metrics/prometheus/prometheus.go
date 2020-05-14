package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	gohttpmetrics "github.com/slok/go-http-metrics/metrics"
	gohttpmetricsprometheus "github.com/slok/go-http-metrics/metrics/prometheus"
	whmetrics "github.com/slok/kubewebhook/pkg/observability/metrics"

	"github.com/slok/k8s-webhook-example/internal/http/webhook"
)

// Types used to avoid collisions with the same interface naming.
type httpRecorder gohttpmetrics.Recorder
type webhookRecorder whmetrics.Recorder

// Recorder satisfies multiple metrics recording interfaces using a Prometheus backend.
type Recorder struct {
	httpRecorder
	webhookRecorder
}

// NewRecorder returns a new Prometheus Recorder.
func NewRecorder(reg prometheus.Registerer) Recorder {
	return Recorder{
		httpRecorder:    gohttpmetricsprometheus.NewRecorder(gohttpmetricsprometheus.Config{Registry: reg}),
		webhookRecorder: whmetrics.NewPrometheus(reg),
	}
}

// Interface assertion.
var _ webhook.MetricsRecorder = Recorder{}
