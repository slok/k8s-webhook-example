package webhook

import (
	gohttpmetrics "github.com/slok/go-http-metrics/metrics"
	"github.com/slok/kubewebhook/v2/pkg/webhook"
)

// MetricsRecorder is the service used to record metrics in the internal HTTP webhook.
type MetricsRecorder interface {
	gohttpmetrics.Recorder
	webhook.MetricsRecorder
}

// Types used to avoid collisions with the same interface naming.
type httpRecorder gohttpmetrics.Recorder
type webhookRecorder webhook.MetricsRecorder

var dummyMetricsRecorder = struct {
	httpRecorder
	webhookRecorder
}{
	httpRecorder:    gohttpmetrics.Dummy,
	webhookRecorder: webhook.NoopMetricsRecorder,
}
