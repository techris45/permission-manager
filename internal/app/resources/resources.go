package resources

import (
	"k8s.io/client-go/kubernetes"
)

// User is a yser saved inside ETCD as a Custom Resouce (CRD api)

type ResourcesService interface {
	UserService
	GetNamespaces() (names []string, err error)
}

type resourcesService struct {
	kubeclient kubernetes.Interface
}

func NewResourcesService(kc kubernetes.Interface) ResourcesService {
	return &kubernetesResourcesService{
		kubeclient: kc,
	}
}
