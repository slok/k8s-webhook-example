package ingress

import (
	"context"
	"fmt"
	"regexp"

	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	networkingv1 "k8s.io/api/networking/v1"
	networkingv1beta1 "k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewHostRegexValidator returns a new validator that checks an ingress hosts match at least
// one of the received regexes.
// It knows how to handle different ingress types.
// If the received object is not an ingress then will return `ErrNotIngress` error.
func NewHostRegexValidator(hostRegexes []string) (Validator, error) {
	// Compile regex for the different hosts.
	regexes := make([]*regexp.Regexp, 0, len(hostRegexes))
	for _, r := range hostRegexes {
		rc, err := regexp.Compile(r)
		if err != nil {
			return nil, fmt.Errorf("the '%s' regex is not valid: %w", r, err)
		}
		regexes = append(regexes, rc)
	}

	return hostRegexValidator{
		allValid: len(regexes) == 0, // If no regexes then all valid.
		regexes:  regexes,
	}, nil
}

type hostRegexValidator struct {
	allValid bool
	regexes  []*regexp.Regexp
}

func (h hostRegexValidator) Validate(ctx context.Context, obj metav1.Object) error {
	hosts := []string{}

	// Missing generics...
	switch ing := obj.(type) {
	case *extensionsv1beta1.Ingress:
		for _, r := range ing.Spec.Rules {
			hosts = append(hosts, r.Host)
		}
	case *networkingv1beta1.Ingress:
		for _, r := range ing.Spec.Rules {
			hosts = append(hosts, r.Host)
		}
	case *networkingv1.Ingress:
		for _, r := range ing.Spec.Rules {
			hosts = append(hosts, r.Host)
		}
	default:
		return ErrNotIngress
	}

	if h.allValid {
		return nil
	}

	for _, host := range hosts {
		if !h.isValidHost(host) {
			return fmt.Errorf("host %s is not a valid host", host)
		}
	}

	return nil
}

func (h hostRegexValidator) isValidHost(host string) bool {
	for _, regex := range h.regexes {
		if regex.MatchString(host) {
			return true
		}
	}

	return false
}
