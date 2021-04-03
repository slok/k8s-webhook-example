package webhook

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	monitoringv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	kwhhttp "github.com/slok/kubewebhook/v2/pkg/http"
	kwhlog "github.com/slok/kubewebhook/v2/pkg/log"
	kwhmodel "github.com/slok/kubewebhook/v2/pkg/model"
	kwhwebhook "github.com/slok/kubewebhook/v2/pkg/webhook"
	kwhmutating "github.com/slok/kubewebhook/v2/pkg/webhook/mutating"
	kwhvalidating "github.com/slok/kubewebhook/v2/pkg/webhook/validating"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/slok/k8s-webhook-example/internal/log"
	"github.com/slok/k8s-webhook-example/internal/validation/ingress"
)

// kubewebhookLogger is a small proxy to use our logger with Kubewebhook.
type kubewebhookLogger struct {
	log.Logger
}

func (l kubewebhookLogger) WithValues(kv map[string]interface{}) kwhlog.Logger {
	return kubewebhookLogger{Logger: l.Logger.WithKV(kv)}
}
func (l kubewebhookLogger) WithCtxValues(ctx context.Context) kwhlog.Logger {
	return l.WithValues(kwhlog.ValuesFromCtx(ctx))
}
func (l kubewebhookLogger) SetValuesOnCtx(parent context.Context, values map[string]interface{}) context.Context {
	return kwhlog.CtxWithValues(parent, values)
}

// allmark sets up the webhook handler for marking all kubernetes resources using Kubewebhook library.
func (h handler) allMark() (http.Handler, error) {
	mt := kwhmutating.MutatorFunc(func(ctx context.Context, ar *kwhmodel.AdmissionReview, obj metav1.Object) (*kwhmutating.MutatorResult, error) {
		err := h.marker.Mark(ctx, obj)
		if err != nil {
			return nil, fmt.Errorf("could not mark the resource: %w", err)
		}

		return &kwhmutating.MutatorResult{
			MutatedObject: obj,
			Warnings:      []string{"Resource marked with custom labels"},
		}, nil
	})

	logger := kubewebhookLogger{Logger: h.logger.WithKV(log.KV{"lib": "kubewebhook", "webhook": "allMark"})}
	wh, err := kwhmutating.NewWebhook(kwhmutating.WebhookConfig{
		ID:      "allMark",
		Logger:  logger,
		Mutator: mt,
	})
	if err != nil {
		return nil, fmt.Errorf("could not create webhook: %w", err)
	}
	whHandler, err := kwhhttp.HandlerFor(kwhhttp.HandlerConfig{
		Webhook: kwhwebhook.NewMeasuredWebhook(h.metrics, wh),
		Logger:  logger,
	})
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
	vSingle := kwhvalidating.ValidatorFunc(func(ctx context.Context, ar *kwhmodel.AdmissionReview, obj metav1.Object) (*kwhvalidating.ValidatorResult, error) {
		err := h.ingSingleHostVal.Validate(ctx, obj)
		if err != nil {
			if errors.Is(err, ingress.ErrNotIngress) {
				h.logger.Warningf("received object is not an ingress")
				return &kwhvalidating.ValidatorResult{Valid: true}, nil
			}

			return &kwhvalidating.ValidatorResult{
				Message: fmt.Sprintf("ingress is invalid: %s", err),
				Valid:   false,
			}, nil
		}

		return &kwhvalidating.ValidatorResult{Valid: true}, nil
	})

	// Host based on regex validator.
	vRegex := kwhvalidating.ValidatorFunc(func(ctx context.Context, ar *kwhmodel.AdmissionReview, obj metav1.Object) (*kwhvalidating.ValidatorResult, error) {
		err := h.ingRegexHostVal.Validate(ctx, obj)
		if err != nil {
			if errors.Is(err, ingress.ErrNotIngress) {
				h.logger.Warningf("received object is not an ingress")
				return &kwhvalidating.ValidatorResult{Valid: true}, nil
			}

			return &kwhvalidating.ValidatorResult{
				Message: fmt.Sprintf("ingress host is invalid: %s", err),
				Valid:   false,
			}, nil
		}

		return &kwhvalidating.ValidatorResult{Valid: true}, nil
	})

	logger := kubewebhookLogger{Logger: h.logger.WithKV(log.KV{"lib": "kubewebhook", "webhook": "ingressValidation"})}

	// Create a chain with both ingress validations and use these to create the kwhwebhook.
	v := kwhvalidating.NewChain(logger, vSingle, vRegex)
	wh, err := kwhvalidating.NewWebhook(kwhvalidating.WebhookConfig{
		ID:        "ingressValidation",
		Validator: v,
		Logger:    logger,
	})
	if err != nil {
		return nil, fmt.Errorf("could not create webhook: %w", err)
	}

	whHandler, err := kwhhttp.HandlerFor(kwhhttp.HandlerConfig{
		Webhook: kwhwebhook.NewMeasuredWebhook(h.metrics, wh),
		Logger:  logger,
	})
	if err != nil {
		return nil, fmt.Errorf("could not create handler from webhook: %w", err)
	}

	return whHandler, nil
}

// safeServiceMonitor sets up the webhook handler to set safety Prometheus service monitor CR settings.
func (h handler) safeServiceMonitor() (http.Handler, error) {
	mt := kwhmutating.MutatorFunc(func(ctx context.Context, ar *kwhmodel.AdmissionReview, obj metav1.Object) (*kwhmutating.MutatorResult, error) {
		sm, ok := obj.(*monitoringv1.ServiceMonitor)
		if !ok {
			h.logger.Warningf("received object is not an monitoringv1.ServiceMonitor")
			return &kwhmutating.MutatorResult{}, nil
		}

		err := h.servMonSafer.EnsureSafety(ctx, sm)
		if err != nil {
			return nil, fmt.Errorf("could not set safety settings on service monitor: %w", err)
		}

		return &kwhmutating.MutatorResult{MutatedObject: sm}, nil
	})

	logger := kubewebhookLogger{Logger: h.logger.WithKV(log.KV{"lib": "kubewebhook", "webhook": "safeServiceMonitor"})}

	// Create a static webhook, placing the specific object we are going to redeive, this is important
	// so we receive a CR instead of `runtume.Unstructured` on the mutator.
	wh, err := kwhmutating.NewWebhook(kwhmutating.WebhookConfig{
		ID:      "safeServiceMonitor",
		Obj:     &monitoringv1.ServiceMonitor{},
		Mutator: mt,
		Logger:  logger,
	})
	if err != nil {
		return nil, fmt.Errorf("could not create webhook: %w", err)
	}

	whHandler, err := kwhhttp.HandlerFor(kwhhttp.HandlerConfig{
		Webhook: kwhwebhook.NewMeasuredWebhook(h.metrics, wh),
		Logger:  logger,
	})
	if err != nil {
		return nil, fmt.Errorf("could not create handler from webhook: %w", err)
	}

	return whHandler, nil
}
