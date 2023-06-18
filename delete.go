package kcgo

import (
	"fmt"
	"os"
	"strings"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (kc *Kc) DeleteFile(file, namespace string) error {
	yml, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("failed to read file %s : %w", file, err)
	}
	return kc.DeleteString(string(yml), namespace)
}

func (kc *Kc) DeleteString(yml, namespace string) error {
	resources := strings.Split(yml, "---")
	errs := make([]error, 0, len(resources))
	for _, r := range resources {
		if err := kc.delete(r, namespace); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("one or more resources failed : %v", errs)
	}
	return nil
}

func (kc *Kc) delete(yml, namespace string) error {
	dr, obj, err := kc.buildNamespaceableResourceInterface(yml, namespace)
	if err != nil {
		return fmt.Errorf("failed to build resource interface : %w", err)
	}

	// Apply string
	if err = dr.Delete(kc.ctx, obj.GetName(), v1.DeleteOptions{}); err != nil {
		return fmt.Errorf("failed to apply provided string : %w", err)
	}

	return nil
}
