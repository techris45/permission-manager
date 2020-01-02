clean:
	rm -rf permission-manager
	rm -rf statik
	rm -rf ./web-client/build
	rm -rf ./web-client/node_modules

dependencies:
	go mod download
	go get github.com/rakyll/statik
	npm install --prefix ./web-client

init: clean dependencies

ui-build:
	npm run build --prefix ./web-client
	statik -src=./web-client/build

go-build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o permission-manager ./cmd/run-server.go

build: dependencies ui-build go-build

release-image: build
	docker build -t quay.io/sighup/permission-manager:$${VERSION:-local} .
	# docker push quay.io/sighup/permission-manager:$${VERSION:-local}

run: build
	./permission-manager

dev:
	CLUSTER_NAME=minikube CONTROL_PLANE_ADDRESS=https://192.168.64.33:8443 BASIC_AUTH_PASSWORD=secret gomon cmd/run-server.go

delete-users:
	kubectl delete -f ./crd/user-crd-definition.yml && kubectl apply -f ./crd/user-crd-definition.yml

forward-service:
	kubectl port-forward svc/permission-manager-service 4000 --namespace permission-manager
