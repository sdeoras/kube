package kube

import (
	"k8s.io/client-go/kubernetes"
)

func AssignClientSet(clientset *kubernetes.Clientset, namespace string, coders ...Coder) {
	for _, coder := range coders {
		coder := coder

		// set clientset and namespace for coder
		coder.Clientset(clientset, namespace)

	}
}
