[![Build Status](https://travis-ci.org/uudashr/fibweb.svg?branch=master)](https://travis-ci.org/uudashr/fibweb)
[![Coverage Status](https://coveralls.io/repos/github/uudashr/fibweb/badge.svg?branch=master)](https://coveralls.io/github/uudashr/fibweb?branch=master)

# Fibweb
This are simple client that using fibgo service

## How to build
```shell
$ make install
```

## How to run
Using docker compose
```shell
$ docker-compose up -d
```

Then open browser http://localhost:8080

## Cloud Deployment
Set GCLOUD_SERVICE_KEY on travis environment variable
```shell
$ base64 gcloud-service-key.json
```

Make sure we create the cluster
```shell
$ make init-cloud-all
```

or the cluster might be already created, then do
```shell
$ make init-cloud
```
This will do the initialization without creating the cluster

Then every successful Travis build, it will be automatically deploy the service.
