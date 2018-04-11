package main

import (
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func newClientset() (*kubernetes.Clientset, error) {
	kubeConfigFile := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	// kubernetes clientset init
	var clientset *kubernetes.Clientset
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
