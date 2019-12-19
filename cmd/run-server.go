package main

import (
	"log"
	"os"

	"github.com/sighupio/permission-manager/internal/adapters/kubeclient"
	"github.com/sighupio/permission-manager/internal/app/user"
	"github.com/sighupio/permission-manager/internal/config"
	"github.com/sighupio/permission-manager/internal/server"
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
	userService := user.NewUserService(kc)
	s := server.New(kc, cfg, userService)
	s.Logger.Fatal(s.Start(":4000"))
}
