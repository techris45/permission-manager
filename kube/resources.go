package kube

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetNamespaces(kubeclient kubernetes.Interface) (names []string) {
	namespaces, err := kubeclient.CoreV1().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	for _, ns := range namespaces.Items {
		names = append(names, ns.Name)
	}

	return
}
