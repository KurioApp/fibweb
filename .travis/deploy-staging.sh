#!/bin/bash

set -e

docker build -t asia.gcr.io/${PROJECT_NAME}/${DOCKER_IMAGE_NAME_STG}:$TRAVIS_COMMIT .

echo $GCLOUD_SERVICE_KEY | base64 --decode -i > ${HOME}/gcloud-service-key.json
gcloud auth activate-service-account --key-file ${HOME}/gcloud-service-key.json

gcloud --quiet config set project $PROJECT_NAME
gcloud --quiet config set container/cluster $CLUSTER_NAME_STG
gcloud --quiet config set compute/zone $COMPUTE_ZONE
gcloud --quiet container clusters get-credentials $CLUSTER_NAME_STG

gcloud docker -- push asia.gcr.io/$PROJECT_NAME/$DOCKER_IMAGE_NAME_STG

yes | gcloud beta container images add-tag asia.gcr.io/${PROJECT_NAME}/${DOCKER_IMAGE_NAME_STG}:$TRAVIS_COMMIT asia.gcr.io/${PROJECT_NAME}/${DOCKER_IMAGE_NAME_STG}:latest

kubectl config view
kubectl config current-context

kubectl set image deployment/${KUBE_DEPLOYMENT_NAME} ${KUBE_DEPLOYMENT_CONTAINER_NAME}=asia.gcr.io/${PROJECT_NAME}/${DOCKER_IMAGE_NAME_STG}:$TRAVIS_COMMIT

# Run e2e_test
DEP_IP=$(kubectl get service fibweb -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
DEP_PORT=$(kubectl get service fibweb -o jsonpath='{.spec.ports[0].port}')
curl --fail http://${DEP_IP}:${DEP_PORT} || exit 1
curl --fail http://${DEP_IP}:${DEP_PORT}/api/fibonacci/numbers || exit 1
