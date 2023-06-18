package kcgo

import (
	"context"
	"fmt"
	"os"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Kc struct {
	ctx        context.Context
	kubeconfig []byte
	config     *rest.Config
}

type KcOption func(*Kc)

func NewKc(opts ...KcOption) *Kc {
	kc := &Kc{}
	for _, opt := range opts {
		opt(kc)
	}
	if len(kc.kubeconfig) > 0 {
		clientCfg, err := clientcmd.NewClientConfigFromBytes(kc.kubeconfig)
		if err != nil {
			panic(err)
		}
		restCfg, err := clientCfg.ClientConfig()
		if err != nil {
			panic(err)
		}
		kc.config = restCfg
	} else {
		//in cluster kubeconfig
		restCfg, err := rest.InClusterConfig()
		if err != nil {
			panic(err)
		}
		kc.config = restCfg
	}
	return kc
}

func WithContext(ctx context.Context) KcOption {
	return func(k *Kc) {
		k.ctx = ctx
	}
}

func WithKubeconfigString(config string) KcOption {
	return func(k *Kc) {
		k.kubeconfig = []byte(config)
	}
}

func WithKubeconfigFile(config string) KcOption {
	return func(k *Kc) {
		c, err := os.ReadFile(config)
		if err != nil {
			panic(fmt.Sprintf("failed to read kubeconfig file %s : %v", config, err))
		}
		k.kubeconfig = c
	}
}
