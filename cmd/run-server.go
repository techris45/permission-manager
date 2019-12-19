package main

import (
	"log"
	"os"

	kubernetesResourcesService "github.com/sighupio/permission-manager/internal/app/kubernetes-resource"
)

func main() {
	cfg := config.New()

	clusterName := os.Getenv("CLUSTER_NAME")
	if clusterName == "" {
		log.Fatal("CLUSTER_NAME env cannot be empty")
	} else {
		cfg.ClusterName = clusterName
	}

	clusterControlPlaceAddress := os.Getenv("CONTROL_PLANE_ADDRESS")
	if clusterControlPlaceAddress == "" {
		log.Fatal("CONTROL_PLANE_ADDRESS env cannot be empty")
	} else {
		cfg.ClusterControlPlaceAddress = clusterControlPlaceAddress
	}

	kc := kubeclient.New()
	ks := kubernetesResourcesService.NewKubernetesResourcesService(kc)
	s := server.New(kc, cfg, ks)
	s.Logger.Fatal(s.Start(":4000"))
}
