package kcgo

import (
	"context"
	"testing"
)

var (
	singleResource = `apiVersion: v1
kind: Pod
metadata:
  name: nginx
spec:
  containers:
  - name: nginx
    image: nginx:1.14.2
    ports:
    - containerPort: 80`

	multipleResource = `apiVersion: v1
kind: Pod
metadata:
  name: nginx
spec:
  containers:
  - name: nginx
    image: nginx:1.14.2
    ports:
    - containerPort: 80
---
apiVersion: v1
kind: Pod
metadata:
  name: nginx
spec:
  containers:
  - name: nginx
    image: nginx:1.14.2
    ports:
    - containerPort: 80`
)

func TestApplyString(t *testing.T) {
	kc := NewKc(WithContext(context.Background()), WithKubeconfigFile("./testdata/config"))
	if err := kc.ApplyString(singleResource, "default"); err != nil {
		t.Error(err)
	}
	if err := kc.ApplyString(multipleResource, "default"); err != nil {
		t.Error(err)
	}
}
