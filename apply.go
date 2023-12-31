package kcgo

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func (kc *Kc) ApplyFile(file, namespace string) error {
	yml, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("failed to read file %s : %w", file, err)
	}
	return kc.ApplyString(string(yml), namespace)
}

func (kc *Kc) ApplyString(yml, namespace string) error {
	resources := strings.Split(yml, "---")
	errs := make([]error, 0, len(resources))
	for _, r := range resources {
		if err := kc.apply(r, namespace); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("one or more resources failed : %v", errs)
	}
	return nil
}

func (kc *Kc) apply(yml, namespace string) error {
	dr, obj, err := kc.buildNamespaceableResourceInterface(yml, namespace)
	if err != nil {
		return fmt.Errorf("failed to build resource interface : %w", err)
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
