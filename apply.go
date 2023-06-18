package kcgo

import (
	"encoding/json"
	"fmt"

	"k8s.io/apimachinery/pkg/api/meta"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/restmapper"
)

func (kc *Kc) ApplyString(yml string, namespace string) error {
	// Deconstruct yaml and define group/version/kind
	decUnstructured := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	obj := &unstructured.Unstructured{}
	_, gvk, err := decUnstructured.Decode([]byte(yml), nil, obj)
	if err != nil {
		return fmt.Errorf("failed to decode provided string : %w", err)
	}

	// Prepare a RESTMapper and find GVR
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(kc.config)
	if err != nil {
		return fmt.Errorf("failed to create discovery client : %w", err)
	}

	mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(discoveryClient))
	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return fmt.Errorf("failed to map group/version/kind of provided string : %w", err)
	}

	// Prepare the dynamic client
	dyn, err := dynamic.NewForConfig(kc.config)
	if err != nil {
		return fmt.Errorf("failed to create new dynamic client : %w", err)
	}

	// Obtain REST interface for the GVR
	dr := dyn.Resource(mapping.Resource)
	// Set namespace if resource is namespace scoped
	if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
		dr = dr.Namespace(namespace).(dynamic.NamespaceableResourceInterface)
	}

	data, err := json.Marshal(obj)
	if err != nil {
		return fmt.Errorf("failed to marshal provided string : %w", err)
	}

	// Apply string
	if _, err = dr.Patch(kc.ctx, obj.GetName(), types.ApplyPatchType, data, v1.PatchOptions{FieldManager: "kc-go"}); err != nil {
		return fmt.Errorf("failed to apply provided string : %w", err)
	}

	return nil
}
