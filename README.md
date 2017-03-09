[![Build Status](https://travis-ci.org/uudashr/fibweb.svg?branch=master)](https://travis-ci.org/uudashr/fibweb)
[![Coverage Status](https://coveralls.io/repos/github/uudashr/fibweb/badge.svg?branch=master)](https://coveralls.io/github/uudashr/fibweb?branch=master)

# Fibweb
This are simple client that using fibgo service

## How to build
```shell
$ make install
```

## How to run
Ensure we have fibgo service run

```shell
$ make run-fibgo
$ make run
$ curl http://localhost:8080/api/fibonacci/numbers?limit=5
$ make stop-fibgo
```
