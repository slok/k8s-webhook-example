package webhook

import (
	"context"
	"fmt"
	"net/http"

	"github.com/slok/k8s-webhook-example/internal/log"
	whhttp "github.com/slok/kubewebhook/pkg/http"
	mutatingwh "github.com/slok/kubewebhook/pkg/webhook/mutating"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
