#!/usr/bin/env bash

SCRIPT_DIR=$(dirname ${BASH_SOURCE})

minikube_registry=$(podman container port minikube 5000)

podman build -t echo-func $SCRIPT_DIR/../images/echo-func
podman push echo-func $minikube_registry/echo-func:3 --tls-verify=false
