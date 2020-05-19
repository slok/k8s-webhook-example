package ingress

import (
	"context"
	"errors"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ErrNotIngress will be used when the validating object is not an ingress.
var ErrNotIngress = errors.New("object is not an ingress")

// Validator knows how to validate an ingress.
type Validator interface {
	Validate(ctx context.Context, obj metav1.Object) error
}

// DummyValidator is a Validator that doesn't do anything.
var DummyValidator Validator = dummyValidator(0)

type dummyValidator int

func (dummyValidator) Validate(_ context.Context, _ metav1.Object) error { return nil }
