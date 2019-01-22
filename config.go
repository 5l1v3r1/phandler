package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type PhandlerConfig struct {
	kubernetesNamespace string
	kubernetesClientSet *kubernetes.Clientset
	kubernetesConfig    *rest.Config
}

func NewPhandlerConfig(kubernetesNamespace string) *PhandlerConfig {
	config, err := getOutClusterConfig()
	if err != nil {
		log.Print(err)
		return nil
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Print(err)
		return nil
	}
	return &PhandlerConfig{
		kubernetesNamespace: kubernetesNamespace,
		kubernetesConfig:    config,
		kubernetesClientSet: clientset,
	}

}

// Get kubernetes config
func getInClusterConfig() (*rest.Config, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	return config, nil
}

func getOutClusterConfig() (*rest.Config, error) {
	var kubeconfig string
	if home := homeDir(); home != "" {
		kubeconfig = *flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = *flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}
	return config, nil

}
func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
