package util

import (
	"fmt"
	"strings"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/cert"

	"github.com/openvswitch/ovn-kubernetes/go-controller/pkg/config"
)

// NewClientset creates a Kubernetes clientset from either a kubeconfig,
// TLS properties, or an apiserver URL
func NewClientset(conf *config.KubernetesConfig) (*kubernetes.Clientset, error) {
	var kconfig *rest.Config
	var err error

	if conf.Kubeconfig != "" {
		// uses the current context in kubeconfig
		kconfig, err = clientcmd.BuildConfigFromFlags("", conf.Kubeconfig)
	} else if strings.HasPrefix(conf.APIServer, "https") {
		if conf.APIServer == "" || conf.Token == "" || conf.CACert == "" {
			return nil, fmt.Errorf("TLS-secured apiservers require token and CA certificate")
		}

		if _, err := cert.NewPool(conf.CACert); err != nil {
			return nil, err
		}

		kconfig = &rest.Config{
			Host:            conf.APIServer,
			BearerToken:     conf.Token,
			TLSClientConfig: rest.TLSClientConfig{CAFile: conf.CACert},
		}
	} else if strings.HasPrefix(conf.APIServer, "http") {
		kconfig, err = clientcmd.BuildConfigFromFlags(conf.APIServer, "")
	} else {
		err = fmt.Errorf("a kubeconfig file or server/token/tls credentials are required")
	}
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(kconfig)
}
