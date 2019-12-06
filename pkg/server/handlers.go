package server

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/sighupio/permission-manager/kube"
	"github.com/sighupio/permission-manager/users"
	v1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func listUsers(c echo.Context) error {
	ac := c.(*AppContext)
	return c.JSON(http.StatusOK, users.GetAll(ac.Kubeclient))
}

func createUser(c echo.Context) error {
	ac := c.(*AppContext)
	type Request struct {
		Name string `json:"name" validate:"required"`
	}
	type Reponse = users.User
	r := new(Request)
	if err := c.Bind(r); err != nil {
		panic(err)
	}
	if err := c.Validate(r); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorRes{err.Error()})
	}

	u := users.CreateUser(ac.Kubeclient, r.Name)

	return c.JSON(http.StatusOK, Reponse{Name: u.Name})
}

func listGroups(c echo.Context) error {
	type Group struct {
		Name string `json:"name"`
	}
	return c.JSON(http.StatusOK, []Group{})
}

func listNamespaces(c echo.Context) error {
	ac := c.(*AppContext)

	namespaces, err := ac.Kubeclient.CoreV1().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	type Response struct {
		Namespaces []v1.Namespace `json:"namespaces"`
	}

	return c.JSON(http.StatusOK, Response{
		Namespaces: namespaces.Items,
	})
}

func ListRbac(c echo.Context) error {
	ac := c.(*AppContext)
	type Response struct {
		ClusterRoles        []rbacv1.ClusterRole        `json:"clusterRoles"`
		ClusterRoleBindings []rbacv1.ClusterRoleBinding `json:"clusterRoleBindings"`
		Roles               []rbacv1.Role               `json:"roles"`
		RoleBindings        []rbacv1.RoleBinding        `json:"roleBindings"`
	}

	clusterRoles, err := ac.Kubeclient.RbacV1().ClusterRoles().List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	clusterRoleBindings, err := ac.Kubeclient.RbacV1().ClusterRoleBindings().List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	roles, err := ac.Kubeclient.RbacV1().Roles("").List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	roleBindings, err := ac.Kubeclient.RbacV1().RoleBindings("").List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	return c.JSON(http.StatusOK, Response{
		ClusterRoles:        clusterRoles.Items,
		ClusterRoleBindings: clusterRoleBindings.Items,
		Roles:               roles.Items,
		RoleBindings:        roleBindings.Items,
	})
}

func CreateClusterRole(c echo.Context) error {
	ac := c.(*AppContext)
	type Request struct {
		RoleName string              `json:"roleName" validate:"required"`
		Rules    []rbacv1.PolicyRule `json:"rules" validate:"required"`
	}
	r := new(Request)
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorRes{err.Error()})
	}

	type Response struct {
		Ok bool `json:"ok"`
	}

	ac.Kubeclient.RbacV1().ClusterRoles().Create(&rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: r.RoleName,
		},
		Rules: r.Rules,
	})

	return c.JSON(http.StatusOK, Response{Ok: true})
}

func CreateRolebinding(c echo.Context) error {
	ac := c.(*AppContext)
	type Request struct {
		RolebindingName string           `json:"rolebindingName" validate:"required"`
		Namespace       string           `json:"namespace" validate:"required"`
		Username        string           `json:"generated_for_user" validate:"required"`
		Subjects        []rbacv1.Subject `json:"subjects" validate:"required"`
		RoleKind        string           `json:"roleKind" validate:"required"`
		RoleName        string           `json:"roleName" validate:"required"`
	}
	r := new(Request)
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorRes{err.Error()})
	}

	type Response struct {
		Ok bool `json:"ok"`
	}

	ac.Kubeclient.RbacV1().RoleBindings(r.Namespace).Create(&rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      r.RolebindingName,
			Namespace: r.Namespace,
			Labels:    map[string]string{"generated_for_user": r.Username},
		},
		RoleRef: rbacv1.RoleRef{
			Kind:     r.RoleKind,
			Name:     r.RoleName,
			APIGroup: "rbac.authorization.k8s.io",
		},
		Subjects: r.Subjects,
	})

	return c.JSON(http.StatusOK, Response{Ok: true})
}

func createClusterRolebinding(c echo.Context) error {
	ac := c.(*AppContext)
	type Request struct {
		ClusterRolebindingName string           `json:"clusterRolebindingName"`
		Username               string           `json:"user"`
		Subjects               []rbacv1.Subject `json:"subjects"`
		RoleName               string           `json:"roleName"`
	}
	r := new(Request)
	if err := c.Bind(r); err != nil {
		return err
	}

	type Response struct {
		Ok bool `json:"ok"`
	}

	ac.Kubeclient.RbacV1().ClusterRoleBindings().Create(&rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:   r.ClusterRolebindingName,
			Labels: map[string]string{"generated_for_user": r.Username},
		},
		RoleRef: rbacv1.RoleRef{
			Kind:     "ClusterRole",
			Name:     r.RoleName,
			APIGroup: "rbac.authorization.k8s.io",
		},
		Subjects: r.Subjects,
	})

	return c.JSON(http.StatusOK, Response{Ok: true})
}

func deleteClusterRole(c echo.Context) error {
	ac := c.(*AppContext)
	type Request struct {
		RoleName string `json:"roleName" validate:"required"`
	}
	type Response struct {
		Ok bool `json:"ok"`
	}

	r := new(Request)
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorRes{err.Error()})
	}

	ac.Kubeclient.RbacV1().ClusterRoles().Delete(r.RoleName, nil)
	return c.JSON(http.StatusOK, Response{Ok: true})
}

func deleteClusterRolebinding(c echo.Context) error {
	ac := c.(*AppContext)
	type Request struct {
		RolebindingName string `json:"rolebindingName"`
	}
	type Response struct {
		Ok bool `json:"ok"`
	}

	r := new(Request)
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorRes{err.Error()})
	}

	ac.Kubeclient.RbacV1().ClusterRoleBindings().Delete(r.RolebindingName, nil)
	return c.JSON(http.StatusOK, Response{Ok: true})
}

func deleteRole(c echo.Context) error {
	ac := c.(*AppContext)

	type Request struct {
		RoleName  string `json:"roleName" validate:"required"`
		Namespace string `json:"namespace" validate:"required"`
	}
	type Response struct {
		Ok bool `json:"ok"`
	}

	r := new(Request)
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorRes{err.Error()})
	}

	ac.Kubeclient.RbacV1().Roles(r.Namespace).Delete(r.RoleName, nil)
	return c.JSON(http.StatusOK, Response{Ok: true})
}

func deleteRolebinding(c echo.Context) error {
	ac := c.(*AppContext)
	type Request struct {
		RolebindingName string `json:"rolebindingName" validate:"required"`
		Namespace       string `json:"namespace" validate:"required"`
	}
	type Response struct {
		Ok bool `json:"ok"`
	}

	r := new(Request)
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorRes{err.Error()})
	}

	ac.Kubeclient.RbacV1().RoleBindings(r.Namespace).Delete(r.RolebindingName, nil)
	return c.JSON(http.StatusOK, Response{Ok: true})
}

func createKubeconfig(c echo.Context) error {
	ac := c.(*AppContext)
	type Request struct {
		Username string `json:"username"`
	}
	r := new(Request)
	if err := c.Bind(r); err != nil {
		return err
	}

	kubeconfig := kube.CreateKubeconfigYAML(ac.Kubeclient, r.Username)

	type Response struct {
		Ok         bool   `json:"ok"`
		Kubeconfig string `json:"kubeconfig"`
	}
	return c.JSON(http.StatusOK, Response{Ok: true, Kubeconfig: kubeconfig})
}

type (
	// FallbackResponseWriter wraps an http.Requesthandler and surpresses
	// a 404 status code. In such case a given local file will be served.
	FallbackResponseWriter struct {
		WrappedResponseWriter http.ResponseWriter
		FileNotFound          bool
	}
)

// Header returns the header of the wrapped response writer
func (frw *FallbackResponseWriter) Header() http.Header {
	return frw.WrappedResponseWriter.Header()
}

// Write sends bytes to wrapped response writer, in case of FileNotFound
// It surpresses further writes (concealing the fact though)
func (frw *FallbackResponseWriter) Write(b []byte) (int, error) {
	if frw.FileNotFound {
		return len(b), nil
	}
	return frw.WrappedResponseWriter.Write(b)
}

// WriteHeader sends statusCode to wrapped response writer
func (frw *FallbackResponseWriter) WriteHeader(statusCode int) {
	if statusCode == http.StatusNotFound {
		frw.FileNotFound = true
		return
	}
	frw.WrappedResponseWriter.WriteHeader(statusCode)
}

// AddFallbackHandler wraps the handler func in another handler func covering authentication
func AddFallbackHandler(handler http.HandlerFunc, filename string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		frw := FallbackResponseWriter{
			WrappedResponseWriter: w,
			FileNotFound:          false,
		}
		handler(&frw, r)
		if frw.FileNotFound {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	}
}
