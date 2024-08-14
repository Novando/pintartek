#!/bin/bash

if [ -z $SERVICE ]; then
  echo "Error: Please declared which CMD to call during execution! the argument is SERVICE";
  exit 1;
fi

go run ./cmd/$SERVICE/main.go