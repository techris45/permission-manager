package main

import (
	"github.com/sighupio/permission-manager/kube"
	"github.com/sighupio/permission-manager/pkg/server"
)

func main() {
	kc := kube.NewKubeclient()
	s := server.New(kc)
	s.Logger.Fatal(s.Start(":4000"))
}
