package main_test

import (
	"net/http"
	"net/http/httptest"

	"testing"

	"github.com/labstack/echo"
	"github.com/sighupio/permission-manager/pkg/server"
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/kubernetes/fake"
)

func TestMain(t *testing.T) {
	kc := fake.NewSimpleClientset()
	s := server.New(kc)
	req := httptest.NewRequest(echo.GET, "/list-users", nil)
	res := httptest.NewRecorder()
	c := s.NewContext(req, res)

	appContext := &server.AppContext{Context: c, Kubeclient: kc}

	if assert.NoError(t, server.ListRbac(appContext)) {
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, `123`, res.Body.String())
	}
}
