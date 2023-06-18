package kcgo

import (
	"context"
	"testing"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestGet(t *testing.T) {
	kc := NewKc(WithContext(context.Background()), WithKubeconfigFile("./testdata/config"))
	if res, err := kc.Get("pods", "kube-system"); err != nil {
		t.Error(err)
	} else {
		var podList v1.PodList
		if err = runtime.DefaultUnstructuredConverter.FromUnstructured(res.Object, &podList); err != nil {
			t.Error(err)
		}
	}
}
