package webhook

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	monitoringv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	whhttp "github.com/slok/kubewebhook/v2/pkg/http"
	whmodel "github.com/slok/kubewebhook/v2/pkg/model"
	"github.com/slok/kubewebhook/v2/pkg/webhook"
	whmutating "github.com/slok/kubewebhook/v2/pkg/webhook/mutating"
	whvalidating "github.com/slok/kubewebhook/v2/pkg/webhook/validating"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/slok/k8s-webhook-example/internal/log"
	"github.com/slok/k8s-webhook-example/internal/validation/ingress"
)

// allmark sets up the webhook handler for marking all kubernetes resources using Kubewebhook library.
func (h handler) allMark() (http.Handler, error) {
	mt := whmutating.MutatorFunc(func(ctx context.Context, ar *whmodel.AdmissionReview, obj metav1.Object) (*whmutating.MutatorResult, error) {
		err := h.marker.Mark(ctx, obj)
		if err != nil {
			return nil, fmt.Errorf("could not mark the resource: %w", err)
		}

		return &whmutating.MutatorResult{MutatedObject: obj}, nil
	})

	logger := h.logger.WithKV(log.KV{"lib": "kubewebhook", "webhook": "allMark"})
	wh, err := whmutating.NewWebhook(whmutating.WebhookConfig{
		ID:      "allMark",
		Logger:  logger,
		Mutator: mt,
	})
	if err != nil {
		return nil, fmt.Errorf("could not create webhook: %w", err)
	}
	whHandler, err := whhttp.HandlerFor(webhook.NewMeasuredWebhook(h.metrics, wh))
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
	vSingle := whvalidating.ValidatorFunc(func(ctx context.Context, ar *whmodel.AdmissionReview, obj metav1.Object) (*whvalidating.ValidatorResult, error) {
		err := h.ingSingleHostVal.Validate(ctx, obj)
		if err != nil {
			if errors.Is(err, ingress.ErrNotIngress) {
				h.logger.Warningf("received object is not an ingress")
				return &whvalidating.ValidatorResult{Valid: true}, nil
			}

			return &whvalidating.ValidatorResult{
				Message: fmt.Sprintf("ingress is invalid: %s", err),
				Valid:   false,
			}, nil
		}

		return &whvalidating.ValidatorResult{Valid: true}, nil
	})

	// Host based on regex validator.
	vRegex := whvalidating.ValidatorFunc(func(ctx context.Context, ar *whmodel.AdmissionReview, obj metav1.Object) (*whvalidating.ValidatorResult, error) {
		err := h.ingRegexHostVal.Validate(ctx, obj)
		if err != nil {
			if errors.Is(err, ingress.ErrNotIngress) {
				h.logger.Warningf("received object is not an ingress")
				return &whvalidating.ValidatorResult{Valid: true}, nil
			}

			return &whvalidating.ValidatorResult{
				Message: fmt.Sprintf("ingress host is invalid: %s", err),
				Valid:   false,
			}, nil
		}

		return &whvalidating.ValidatorResult{Valid: true}, nil
	})

	logger := h.logger.WithKV(log.KV{"lib": "kubewebhook", "webhook": "ingressValidation"})

	// Create a chain with both ingress validations and use these to create the webhook.
	v := whvalidating.NewChain(logger, vSingle, vRegex)
	wh, err := whvalidating.NewWebhook(whvalidating.WebhookConfig{
		ID:        "ingressValidation",
		Validator: v,
		Logger:    logger,
	})
	if err != nil {
		return nil, fmt.Errorf("could not create webhook: %w", err)
	}

	whHandler, err := whhttp.HandlerFor(webhook.NewMeasuredWebhook(h.metrics, wh))
	if err != nil {
		return nil, fmt.Errorf("could not create handler from webhook: %w", err)
	}

	return whHandler, nil
}

// safeServiceMonitor sets up the webhook handler to set safety Prometheus service monitor CR settings.
func (h handler) safeServiceMonitor() (http.Handler, error) {
	mt := whmutating.MutatorFunc(func(ctx context.Context, ar *whmodel.AdmissionReview, obj metav1.Object) (*whmutating.MutatorResult, error) {
		sm, ok := obj.(*monitoringv1.ServiceMonitor)
		if !ok {
			h.logger.Warningf("received object is not an monitoringv1.ServiceMonitor")
			return &whmutating.MutatorResult{}, nil
		}

		err := h.servMonSafer.EnsureSafety(ctx, sm)
		if err != nil {
			return nil, fmt.Errorf("could not set safety settings on service monitor: %w", err)
		}

		return &whmutating.MutatorResult{MutatedObject: sm}, nil
	})

	logger := h.logger.WithKV(log.KV{"lib": "kubewebhook", "webhook": "safeServiceMonitor"})

	// Create a static webhook, placing the specific object we are going to redeive, this is important
	// so we receive a CR instead of `runtume.Unstructured` on the mutator.
	wh, err := whmutating.NewWebhook(whmutating.WebhookConfig{
		ID:      "safeServiceMonitor",
		Obj:     &monitoringv1.ServiceMonitor{},
		Mutator: mt,
		Logger:  logger,
	})
	if err != nil {
		return nil, fmt.Errorf("could not create webhook: %w", err)
	}

	whHandler, err := whhttp.HandlerFor(webhook.NewMeasuredWebhook(h.metrics, wh))
	if err != nil {
		return nil, fmt.Errorf("could not create handler from webhook: %w", err)
	}

	return whHandler, nil
}
