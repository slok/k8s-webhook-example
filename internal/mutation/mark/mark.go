package mark

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Marker knows how to mark Kubernetes resources.
type Marker interface {
	Mark(ctx context.Context, obj metav1.Object) error
}

// NewLabelMarker returns a new marker that will mark with labels.
func NewLabelMarker(marks map[string]string) Marker {
	return labelmarker{marks: marks}
}

type labelmarker struct {
	marks map[string]string
}

func (l labelmarker) Mark(_ context.Context, obj metav1.Object) error {
	labels := obj.GetLabels()
	if labels == nil {
		labels = map[string]string{}
	}

	for k, v := range l.marks {
		labels[k] = v
	}

	obj.SetLabels(labels)
	return nil
}

// DummyMarker is a marker that doesn't do anything.
var DummyMarker Marker = dummyMaker(0)

type dummyMaker int

func (dummyMaker) Mark(_ context.Context, _ metav1.Object) error { return nil }
