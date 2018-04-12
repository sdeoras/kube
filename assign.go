package kube

import (
	"k8s.io/client-go/kubernetes"
)

// AssignClientSet sets clientset and namespace for each of the input coders
func AssignClientSet(clientset *kubernetes.Clientset, namespace string, coders ...Coder) {
	for _, coder := range coders {
		coder := coder

		// set clientset and namespace for coder
		coder.Clientset(clientset, namespace)

	}
}
