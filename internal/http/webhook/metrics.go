package webhook

import (
	gohttpmetrics "github.com/slok/go-http-metrics/metrics"
	whmetrics "github.com/slok/kubewebhook/pkg/observability/metrics"
)

// MetricsRecorder is the service used to record metrics in the internal HTTP webhook.
type MetricsRecorder interface {
	gohttpmetrics.Recorder
	whmetrics.Recorder
}

// Types used to avoid collisions with the same interface naming.
type httpRecorder gohttpmetrics.Recorder
type webhookRecorder whmetrics.Recorder

var dummyMetricsRecorder = struct {
	httpRecorder
	webhookRecorder
}{
	httpRecorder:    gohttpmetrics.Dummy,
	webhookRecorder: whmetrics.Dummy,
}
