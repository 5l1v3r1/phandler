package helm

import (
	"fmt"

	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/helm/pkg/helm"
	"k8s.io/helm/pkg/helm/portforwarder"
)

const (
	tillerNamespace = "kube-system"
)

func NewHelmClient(clientset *kubernetes.Clientset, config *rest.Config) (*helm.Client, error) {
	tillerTunnel, err := portforwarder.New(tillerNamespace, clientset, config)
	if err != nil {
		return nil, errors.Wrap(err, "[ERROR] Failed to portforward")
	}
	tillerHost := fmt.Sprint("127.0.0.1:%d", tillerTunnel.Local)

	helmOpts := []helm.Option{
		helm.Host(tillerHost),
		helm.ConnectTimeout(10),
	}
	client := helm.NewClient(helmOpts...)

	return client, nil
}
