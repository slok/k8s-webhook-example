package webhook

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	monitoringv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	whhttp "github.com/slok/kubewebhook/pkg/http"
	mutatingwh "github.com/slok/kubewebhook/pkg/webhook/mutating"
	validatingwh "github.com/slok/kubewebhook/pkg/webhook/validating"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/slok/k8s-webhook-example/internal/log"
	"github.com/slok/k8s-webhook-example/internal/validation/ingress"
)

// allmark sets up the webhook handler for marking all kubernetes resources using Kubewebhook library.
func (h handler) allMark() (http.Handler, error) {
	mt := mutatingwh.MutatorFunc(func(ctx context.Context, obj metav1.Object) (bool, error) {
		err := h.marker.Mark(ctx, obj)
		if err != nil {
			return false, fmt.Errorf("could not mark the resource: %w", err)
		}

		return false, nil
	})

	logger := h.logger.WithKV(log.KV{"lib": "kubewebhook", "webhook": "allMark"})
	wh, err := mutatingwh.NewWebhook(mutatingwh.WebhookConfig{Name: "allMark"}, mt, nil, h.metrics, logger)
	if err != nil {
		return nil, fmt.Errorf("could not create webhook: %w", err)
	}
	whHandler, err := whhttp.HandlerFor(wh)
	if err != nil {
		return nil, fmt.Errorf("could not create handler from webhook: %w", err)
	}

	return whHandler, nil
}

// ingressValidation sets up the webhook handler for validating an ingress using a chain of validations.
// Thec validation chain will check first if the ingress has a single host, if not it will stop the
// validation chain, otherwirse it will check the nest ingress Validator that will try matching the host
// with allowed host.
func (h handler) ingressValidation() (http.Handler, error) {
	// Single host validator.
	vSingle := validatingwh.ValidatorFunc(func(ctx context.Context, obj metav1.Object) (bool, validatingwh.ValidatorResult, error) {
		err := h.ingSingleHostVal.Validate(ctx, obj)
		if err != nil {
			if errors.Is(err, ingress.ErrNotIngress) {
				h.logger.Warningf("received object is not an ingress")
				return false, validatingwh.ValidatorResult{Valid: true}, nil
			}

			// We want to stop the chain because we only check in this webhook hosts.
			return true, validatingwh.ValidatorResult{
				Message: fmt.Sprintf("ingress is invalid: %s", err),
				Valid:   false,
			}, nil
		}

		return false, validatingwh.ValidatorResult{Valid: true}, nil
	})

	// Host based on regex validator.
	vRegex := validatingwh.ValidatorFunc(func(ctx context.Context, obj metav1.Object) (bool, validatingwh.ValidatorResult, error) {
		err := h.ingRegexHostVal.Validate(ctx, obj)
		if err != nil {
			if errors.Is(err, ingress.ErrNotIngress) {
				h.logger.Warningf("received object is not an ingress")
				return false, validatingwh.ValidatorResult{Valid: true}, nil
			}

			return false, validatingwh.ValidatorResult{
				Message: fmt.Sprintf("ingress host is invalid: %s", err),
				Valid:   false,
			}, nil
		}

		return false, validatingwh.ValidatorResult{Valid: true}, nil
	})

	logger := h.logger.WithKV(log.KV{"lib": "kubewebhook", "webhook": "ingressValidation"})

	// Create a chain with both ingress validations and use these to create the webhook.
	v := validatingwh.NewChain(logger, vSingle, vRegex)
	wh, err := validatingwh.NewWebhook(validatingwh.WebhookConfig{Name: "ingressValidation"}, v, nil, h.metrics, logger)
	if err != nil {
		return nil, fmt.Errorf("could not create webhook: %w", err)
	}

	whHandler, err := whhttp.HandlerFor(wh)
	if err != nil {
		return nil, fmt.Errorf("could not create handler from webhook: %w", err)
	}

	return whHandler, nil
}

// safeServiceMonitor sets up the webhook handler to set safety Prometheus service monitor CR settings.
func (h handler) safeServiceMonitor() (http.Handler, error) {
	mt := mutatingwh.MutatorFunc(func(ctx context.Context, obj metav1.Object) (bool, error) {
		sm, ok := obj.(*monitoringv1.ServiceMonitor)
		if !ok {
			h.logger.Warningf("received object is not an monitoringv1.ServiceMonitor")
			return false, nil
		}

		err := h.servMonSafer.EnsureSafety(ctx, sm)
		if err != nil {
			return false, fmt.Errorf("could not set safety settings on service monitor: %w", err)
		}

		return false, nil
	})

	logger := h.logger.WithKV(log.KV{"lib": "kubewebhook", "webhook": "safeServiceMonitor"})

	// Create a static webhook, placing the specific object we are going to redeive, this is important
	// so we receive a CR instead of `runtume.Unstructured` on the mutator.
	wh, err := mutatingwh.NewWebhook(mutatingwh.WebhookConfig{
		Name: "safeServiceMonitor",
		Obj:  &monitoringv1.ServiceMonitor{},
	}, mt, nil, h.metrics, logger)
	if err != nil {
		return nil, fmt.Errorf("could not create webhook: %w", err)
	}

	whHandler, err := whhttp.HandlerFor(wh)
	if err != nil {
		return nil, fmt.Errorf("could not create handler from webhook: %w", err)
	}

	return whHandler, nil
}
