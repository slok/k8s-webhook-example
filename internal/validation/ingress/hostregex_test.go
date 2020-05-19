package ingress_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	networkingv1beta1 "k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/slok/k8s-webhook-example/internal/validation/ingress"
)

func TestHostRegexValidator(t *testing.T) {
	tests := map[string]struct {
		hostRegexes []string
		ingress     metav1.Object
		expErr      bool
	}{
		"Having a non ingress should return an error.": {
			ingress: &extensionsv1beta1.Deployment{},
			expErr:  true,
		},

		"Having an ingress (extensions/v1beta1) and not specific regex hosts all ingresses should be valid.": {
			ingress: &extensionsv1beta1.Ingress{
				Spec: extensionsv1beta1.IngressSpec{
					Rules: []extensionsv1beta1.IngressRule{
						{Host: "test1.slok.dev"},
					},
				},
			},
		},

		"Having an ingress and not specific regex hosts all ingresses should be valid.": {
			ingress: &networkingv1beta1.Ingress{
				Spec: networkingv1beta1.IngressSpec{
					Rules: []networkingv1beta1.IngressRule{
						{Host: "test1.slok.dev"},
					},
				},
			},
		},

		"Having an ingress and an specific regex hosts that matches, validation should be valid.": {
			hostRegexes: []string{`^.*\.slok\.dev$`},
			ingress: &networkingv1beta1.Ingress{
				Spec: networkingv1beta1.IngressSpec{
					Rules: []networkingv1beta1.IngressRule{
						{Host: "test1.slok.dev"},
					},
				},
			},
		},

		"Having an ingress and an specific regex hosts that does not match, validation should be invalid.": {
			hostRegexes: []string{`^.*2\.slok\.dev$`},
			ingress: &networkingv1beta1.Ingress{
				Spec: networkingv1beta1.IngressSpec{
					Rules: []networkingv1beta1.IngressRule{
						{Host: "test1.slok.dev"},
					},
				},
			},
			expErr: true,
		},

		"Having an ingress with multiple rules and an specific regex hosts that does not match, validation should be invalid.": {
			hostRegexes: []string{`^.*\.slok\.dev$`},
			ingress: &networkingv1beta1.Ingress{
				Spec: networkingv1beta1.IngressSpec{
					Rules: []networkingv1beta1.IngressRule{
						{Host: "test1.slok.dev"},
						{Host: "test2.slok.dev"},
						{Host: "test1.slok.wrong"},
					},
				},
			},
			expErr: true,
		},

		"Having an ingress with multiple rules and an multiple regex hosts that does match, validation should be valid.": {
			hostRegexes: []string{
				`^.*\.slok\.dev$`,
				`^.*\.slok\.right$`,
			},
			ingress: &networkingv1beta1.Ingress{
				Spec: networkingv1beta1.IngressSpec{
					Rules: []networkingv1beta1.IngressRule{
						{Host: "test1.slok.dev"},
						{Host: "test2.slok.dev"},
						{Host: "test1.slok.right"},
					},
				},
			},
		},

		"Having an ingress with multiple rules and an multiple regex hosts that does not match, validation should be invalid.": {
			hostRegexes: []string{
				`^.*\.slok\.dev$`,
				`^.*\.slok\.right$`,
			},
			ingress: &networkingv1beta1.Ingress{
				Spec: networkingv1beta1.IngressSpec{
					Rules: []networkingv1beta1.IngressRule{
						{Host: "test1.slok.dev"},
						{Host: "test2.slok.dev"},
						{Host: "test1.slok.right"},
						{Host: "test1.slok.wrong"},
					},
				},
			},
			expErr: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			validator, err := ingress.NewHostRegexValidator(test.hostRegexes)
			require.NoError(err)

			err = validator.Validate(context.TODO(), test.ingress)

			if test.expErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
			}
		})
	}
}
