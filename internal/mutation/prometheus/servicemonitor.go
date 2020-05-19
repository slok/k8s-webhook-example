package prometheus

import (
	"context"
	"time"

	monitoringv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
)

// ServiceMonitorSafer will ensure the service monitor has safe settings, and mutate them instead.
type ServiceMonitorSafer interface {
	EnsureSafety(ctx context.Context, sm *monitoringv1.ServiceMonitor) error
}

type serviceMonitorSafer struct {
	minScrapeInterval time.Duration
}

// NewServiceMonitorSafer returns a new ServiceMonitorSafer that will mutate
// the received service monitor in case it don0't have safe settings. Current checks:
// - Minimum scrape interval.
func NewServiceMonitorSafer(minScrapeInterval time.Duration) ServiceMonitorSafer {
	return serviceMonitorSafer{minScrapeInterval: minScrapeInterval}
}

func (s serviceMonitorSafer) EnsureSafety(_ context.Context, sm *monitoringv1.ServiceMonitor) error {
	endpoints := make([]monitoringv1.Endpoint, 0, len(sm.Spec.Endpoints))

	for _, e := range sm.Spec.Endpoints {
		// Set safe/correct scrape intervals if required.
		t, err := time.ParseDuration(e.Interval)
		if err != nil || t < s.minScrapeInterval {
			e.Interval = s.minScrapeInterval.String()
		}

		endpoints = append(endpoints, e)
	}

	sm.Spec.Endpoints = endpoints

	return nil
}

// DummyServiceMonitorSafer is a ServiceMonitorSafer that doesn't do anything.
const DummyServiceMonitorSafer = dummyServiceMonitorSafer(0)

type dummyServiceMonitorSafer int

var _ ServiceMonitorSafer = DummyServiceMonitorSafer

func (dummyServiceMonitorSafer) EnsureSafety(ctx context.Context, sm *monitoringv1.ServiceMonitor) error {
	return nil
}
