#!/bin/bash

GVM_PATH=

if [ $GVM_V ]; then
  GVM_PATH=~/.gvm/gos/$GVM_V/bin/;
fi

if [ -z $SERVICE ]; then
  echo "Error: Please declared which CMD to call during execution! the argument is SERVICE";
  exit 1;
fi

${GVM_PATH}go mod tidy;

GOOS=linux GOARCH=amd64 ${GVM_PATH}go build \
  -o ./bin/$SERVICE \
  ./cmd/$SERVICE/main.go;