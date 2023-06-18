package kcgo

import (
	"encoding/json"
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/kubernetes"
)

func (kc *Kc) Get(resource, namespace string) (*unstructured.Unstructured, error) {
	client, err := kubernetes.NewForConfig(kc.config)
	if err != nil {
		return nil, fmt.Errorf("failed to create k8s client : %w", err)
	}
	response := client.CoreV1().RESTClient().Get().Namespace(namespace).Resource(resource).Do(kc.ctx)
	if response.Error() != nil {
		return nil, fmt.Errorf("failed to GET requested %s : %w", resource, response.Error())
	}
	jsonResp, err := response.Raw()
	if err != nil {
		return nil, fmt.Errorf("failed to parse server response : %w", err)
	}
	unstructuredObj := &unstructured.Unstructured{}
	if err := json.Unmarshal(jsonResp, unstructuredObj); err != nil {
		return nil, fmt.Errorf("failed to parse server response : %w", err)
	}

	return unstructuredObj, nil
}
