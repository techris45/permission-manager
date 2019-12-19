package kubernetesresources_test

import (
	"testing"

	kr "github.com/sighupio/permission-manager/internal/app/kubernetes-resources"
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/kubernetes/fake"
)

func TestGetNamespaces(t *testing.T) {
	kc := fake.NewSimpleClientset()
	svc := kr.NewKubernetesResourcesService(kc)

	names, err := svc.GetNamespaces()
	got := names
	want := []string{}
	if assert.NoError(t, err) {
		assert.ElementsMatch(t, want, got)
	}


	// svc.CreateUser("jaga")
	// svc.CreateUser("jacopo")

	// names, err = svc.GetNamespaces()
	// got = names
	// want = []string{"jaga", "jacopo"}
	// if assert.NoError(t, err) {
	// assert.ElementsMatch(t, want, got)
	// }
}
