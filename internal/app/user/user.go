package user

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"k8s.io/client-go/kubernetes"
)

// User is a yser saved inside ETCD as a Custom Resouce (CRD api)
type User struct {
	Name string `json:"name"`
}

type UserService interface {
	GetAll() []User
	DeleteUser(username string) User
	CreateUser(username string) User
}

const resourceUrl = "apis/permissionmanager.user/v1alpha1/permissionmanagerusers"

type userService struct {
	kubeclient kubernetes.Interface
}

func NewUserService(kc kubernetes.Interface) UserService {
	return &userService{
		kubeclient: kc,
	}
}

func (u *userService) GetAll() []User {
	users := []User{}

	/* generated from JSON response, most fields not used but usefull as documentation */
	type resType struct {
		APIVersion string `json:"apiVersion"`
		Items      []struct {
			APIVersion string `json:"apiVersion"`
			Kind       string `json:"kind"`
			Metadata   struct {
				Annotations struct {
					KubectlKubernetesIoLastAppliedConfiguration string `json:"kubectl.kubernetes.io/last-applied-configuration"`
				} `json:"annotations"`
				CreationTimestamp time.Time `json:"creationTimestamp"`
				Generation        int       `json:"generation"`
				Name              string    `json:"name"`
				ResourceVersion   string    `json:"resourceVersion"`
				SelfLink          string    `json:"selfLink"`
				UID               string    `json:"uid"`
			} `json:"metadata"`
			Spec struct {
				Name string `json:"name"`
			} `json:"spec"`
		} `json:"items"`
		Kind     string `json:"kind"`
		Metadata struct {
			Continue        string `json:"continue"`
			ResourceVersion string `json:"resourceVersion"`
			SelfLink        string `json:"selfLink"`
		} `json:"metadata"`
	}

	var res resType
	r, err := u.kubeclient.AppsV1().RESTClient().Get().AbsPath(resourceUrl).DoRaw()
	if err != nil {
		log.Print("Failed to get users from k8s CRUD api", err)
	}
	err = json.Unmarshal(r, &res)
	if err != nil {
		log.Print("Failed to decode users from k8s CRUD api", err)
	}

	for _, v := range res.Items {
		users = append(users, User{Name: v.Spec.Name})
	}

	return users
}

func (u *userService) CreateUser(username string) User {
	metadataName := "permissionmanager.user." + username
	jsonPayload := fmt.Sprintf(`{
		"apiVersion":"permissionmanager.user/v1alpha1",
		"kind":"Permissionmanageruser",
		"metadata":{
			"name": "%s"
		},
		"spec": {
			"name": "%s"
		}
	}`, metadataName, username)

	_, err := u.kubeclient.AppsV1().RESTClient().Post().AbsPath(resourceUrl).Body([]byte(jsonPayload)).DoRaw()
	if err != nil {
		log.Printf("Failed to create user:%s\n %v\n", username, err)
	}

	return User{Name: username}
}

func (u *userService) DeleteUser(username string) User {
	metadataName := "permissionmanager.user." + username
	_, err := u.kubeclient.AppsV1().RESTClient().Delete().AbsPath(resourceUrl + "/" + metadataName).DoRaw()
	if err != nil {
		log.Printf("Failed to delete user:%s\n %v\n", username, err)
	}

	return User{Name: username}
}
