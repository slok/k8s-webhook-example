package ingress

import (
	"context"
	"fmt"

	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	networkingv1beta1 "k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SingleHostValidator checks if the ingress has a single host.
// It knows how to handle different ingress types.
// If the received object is not an ingress then will return `ErrNotIngress` error.
const SingleHostValidator = singleHostValidator(0)

type singleHostValidator int

func (s singleHostValidator) Validate(ctx context.Context, obj metav1.Object) error {
	var rulesLen int

	// Missing generics...
	switch ing := obj.(type) {
	case *extensionsv1beta1.Ingress:
		rulesLen = len(ing.Spec.Rules)
	case *networkingv1beta1.Ingress:
		rulesLen = len(ing.Spec.Rules)
	default:
		return ErrNotIngress
	}

	if rulesLen != 1 {
		return fmt.Errorf("ingress rules length should be 1, got: %d", rulesLen)
	}

	return nil
}

var _ Validator = SingleHostValidator
