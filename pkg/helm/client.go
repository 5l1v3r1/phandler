package helm

import (
	"fmt"
	"log"

	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/helm/pkg/helm"
	"k8s.io/helm/pkg/helm/portforwarder"
	"k8s.io/helm/pkg/kube"
)

const (
	tillerNamespace = "kube-system"
)

type Client struct {
	*kube.Tunnel
	*helm.Client
}

func NewHelmClient(clientset *kubernetes.Clientset, config *rest.Config) (*Client, error) {
	tillerTunnel, err := portforwarder.New(tillerNamespace, clientset, config)
	if err != nil {
		return nil, errors.Wrap(err, "[ERROR] Failed to portforward")
	}
	tillerHost := fmt.Sprintf("127.0.0.1:%d", tillerTunnel.Local)
	log.Printf(tillerHost)

	helmOpts := []helm.Option{
		helm.Host(tillerHost),
		helm.ConnectTimeout(30),
	}
	client := helm.NewClient(helmOpts...)

	return &Client{Tunnel: tillerTunnel, Client: client}, nil
}
