package kubernetesresources

import (
	"k8s.io/client-go/kubernetes"
)

// User is a yser saved inside ETCD as a Custom Resouce (CRD api)

type KubernetesResourcesService interface {
	UserService
	GetNamespaces() (names []string, err error)
}

type kubernetesResourcesService struct {
	kubeclient kubernetes.Interface
}

func NewKubernetesResourcesService(kc kubernetes.Interface) KubernetesResourcesService {
	return &kubernetesResourcesService{
		kubeclient: kc,
	}
}
