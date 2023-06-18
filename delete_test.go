package kcgo

import (
	"context"
	"testing"
)

func TestDelete(t *testing.T) {
	kc := NewKc(WithContext(context.Background()), WithKubeconfigFile("./testdata/config"))
	if err := kc.DeleteString(singleResource, "default"); err != nil {
		t.Error(err)
	}
	if err := kc.DeleteString(multipleResource, "default"); err != nil {
		t.Error(err)
	}
	if err := kc.DeleteFile("./testdata/dep.yaml", "default"); err != nil {
		t.Error(err)
	}
}
