package resources

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/kubernetes/fake"
)

func TestGetNamespaces(t *testing.T) {
	kc := fake.NewSimpleClientset()

	names, err := GetNamespaces(kc)
	if assert.NoError(t, err) {
		assert.Equal(t, []string{}, names)
	}
}
