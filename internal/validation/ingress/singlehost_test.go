package ingress_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	networkingv1beta1 "k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/slok/k8s-webhook-example/internal/validation/ingress"
)

func TestSingleHostValidator(t *testing.T) {
	tests := map[string]struct {
		ingress metav1.Object
		expErr  bool
	}{
		"Having a non ingress should return an error.": {
			ingress: &extensionsv1beta1.Deployment{},
			expErr:  true,
		},

		"Having an ingress (extensions/v1beta1) with a single rule/host it should be validated as correct.": {
			ingress: &extensionsv1beta1.Ingress{
				Spec: extensionsv1beta1.IngressSpec{
					Rules: []extensionsv1beta1.IngressRule{
						{Host: "test1.slok.dev"},
					},
				},
			},
		},

		"Having an ingress with a single rule/host it should be validated as correct.": {
			ingress: &networkingv1beta1.Ingress{
				Spec: networkingv1beta1.IngressSpec{
					Rules: []networkingv1beta1.IngressRule{
						{Host: "test1.slok.dev"},
					},
				},
			},
		},

		"Having an ingress with a no rules/hosts it should be validated as incorrect.": {
			ingress: &networkingv1beta1.Ingress{
				Spec: networkingv1beta1.IngressSpec{},
			},
			expErr: true,
		},

		"Having an ingress with a 2 rules/hosts it should be validated as incorrect.": {
			ingress: &networkingv1beta1.Ingress{
				Spec: networkingv1beta1.IngressSpec{
					Rules: []networkingv1beta1.IngressRule{
						{Host: "test1.slok.dev"},
						{Host: "test2.slok.dev"},
					},
				},
			},
			expErr: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			err := ingress.SingleHostValidator.Validate(context.TODO(), test.ingress)

			if test.expErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
			}
		})
	}
}
