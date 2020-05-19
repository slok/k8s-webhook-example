package mark_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/slok/k8s-webhook-example/internal/mutation/mark"
)

func TestLabelMarkerMark(t *testing.T) {
	tests := map[string]struct {
		marks  map[string]string
		obj    metav1.Object
		expObj metav1.Object
	}{
		"Having a pod, the labels should be mutated.": {
			marks: map[string]string{
				"test1": "value1",
				"test2": "value2",
			},
			obj: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
					Labels: map[string]string{
						"test2": "old-value2",
						"test3": "value3",
					},
				},
			},
			expObj: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
					Labels: map[string]string{
						"test1": "value1",
						"test2": "value2",
						"test3": "value3",
					},
				},
			},
		},

		"Having a service, the labels should be mutated.": {
			marks: map[string]string{
				"test1": "value1",
				"test2": "value2",
			},
			obj: &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
				},
			},
			expObj: &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
					Labels: map[string]string{
						"test1": "value1",
						"test2": "value2",
					},
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			m := mark.NewLabelMarker(test.marks)

			err := m.Mark(context.TODO(), test.obj)
			require.NoError(err)

			assert.Equal(test.expObj, test.obj)
		})
	}
}
