package kubernetesresources_test

import (
	"testing"

	"github.com/sighupio/permission-manager/internal/app/resources"
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/kubernetes/fake"
)

func TestGetNamespaces(t *testing.T) {
	kc := fake.NewSimpleClientset()
	svc := resources.NewResourcesService(kc)

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
