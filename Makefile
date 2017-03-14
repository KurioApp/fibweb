IMAGE_NAME=fibweb
GCLOUD_REGISTRY_PREFIX=asia.gcr.io/kurio-dev
GCLOUD_CLUSTER_NAME=fibgo-cluster-stg
GCLOUD_PROJECT_NAME=kurio-dev
GCLOUD_ZONE=asia-east1-a

install:
	@go install -v ./...

prepare-install:
	@go get -v ./...

run:
	@fibgoweb -fibgo-addr http://localhost:9090

docker-net-init:
	@docker network create fibnet_back-tier

docker-net-clean:
	@docker network rm fibnet_back-tier

docker-run-fibgo:
	@docker run -d --name fibgo-server --network fibnet_back-tier uudashr/fibgo

docker-stop-fibgo:
	@docker stop fibgo-server
	@docker rm -v fibgo-server

docker-build:
	@docker build -t $(IMAGE_NAME) .

setup-docker-run: docker-net-init docker-run-fibgo

teardown-docker-run: docker-stop-fibgo docker-net-clean

docker-run:
	@docker run -it --rm -p 8080:8080 --network fibnet_back-tier -e FIBGO_ADDR=fibgo-server:8080 $(IMAGE_NAME)

docker-console:
	@docker run -it --rm -p 8080:8080 --network fibnet_back-tier -e FIBGO_ADDR=fibgo-server:8080 $(IMAGE_NAME) /bin/sh

test:
	@go test

test-cover:
	@go test

check:
	@gometalinter --deadline=15s

prepare-check:
	@go get -u github.com/alecthomas/gometalinter
	@gometalinter --install

init-config:
	@gcloud auth activate-service-account --key-file kurio-gcp.json
	@gcloud --quiet config set project $(GCLOUD_PROJECT_NAME)
	@gcloud --quiet config set compute/zone $(GCLOUD_ZONE)

init-cluster-stg:
	@gcloud container --project "$(GCLOUD_PROJECT_NAME)" clusters create "$(GCLOUD_CLUSTER_NAME)" --zone "$(GCLOUD_ZONE)" --machine-type "n1-standard-1" --image-type "GCI" --disk-size "100" --scopes "https://www.googleapis.com/auth/compute","https://www.googleapis.com/auth/devstorage.read_only","https://www.googleapis.com/auth/logging.write","https://www.googleapis.com/auth/monitoring","https://www.googleapis.com/auth/servicecontrol","https://www.googleapis.com/auth/service.management.readonly","https://www.googleapis.com/auth/trace.append" --num-nodes "1" --network "default" --enable-cloud-logging --enable-cloud-monitoring

init-image-stg:
	@docker tag $(IMAGE_NAME) $(GCLOUD_REGISTRY_PREFIX)/fibweb-stg
	@gcloud docker -- push $(GCLOUD_REGISTRY_PREFIX)/fibweb-stg

init-deployment-stg:
	@gcloud --quiet config set container/cluster $(GCLOUD_CLUSTER_NAME)
	@gcloud --quiet container clusters get-credentials $(GCLOUD_CLUSTER_NAME)
	@kubectl create -f fibweb-app-stg.yaml --record

init-cloud: init-config docker-build init-image-stg init-deployment-stg

init-cloud-all: init-config init-cluster-stg docker-build init-image-stg init-deployment-stg
