clean:
	rm -rf permission-manager
	rm -rf statik
	rm -rf ./web-client/build
	rm -rf ./web-client/node_modules

go-build:
	go mod download
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o permission-manager ./cmd/run-server.go

ui-build:
	npm install --prefix ./web-client
	npm run build --prefix ./web-client
	statik -src=./web-client/build

build: ui-build go-build

run: build
	./permission-manager

dev:
	CLUSTER_NAME=minikube CONTROL_PLANE_ADDRESS=https://192.168.64.33:8443 BASIC_AUTH_PASSWORD=secret gomon cmd/run-server.go

delete-users:
	kubectl delete -f ./crd/user-crd-definition.yml && kubectl apply -f ./crd/user-crd-definition.yml

release-image:
	docker build -t reg.sighup.io/sighup-products/permission-manager:1.0.0 .
	docker push reg.sighup.io/sighup-products/permission-manager:1.0.0

forward-service:
	kubectl port-forward svc/permission-manager-service 4000 --namespace permission-manager

