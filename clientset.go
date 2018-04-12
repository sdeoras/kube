package kube

import (
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// GetDefaultClientSet returns k8s clientset based off of default config file
func GetDefaultClientSet() (*kubernetes.Clientset, error) {
	// kubernetes clientset init
	var clientset *kubernetes.Clientset
	kubeConfigFile := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	// use the current context in kubeconfig
	if kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigFile); err != nil {
		return nil, err
	} else {
		// create the clientset
		clientset, err = kubernetes.NewForConfig(kubeConfig)
		if err != nil {
			return nil, err
		}
	}

	return clientset, nil
}
