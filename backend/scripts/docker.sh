#!/bin/bash

if [ "$VERSION" ]; then \
  docker build \
    --no-cache \
    -f Dockerfile \
    --tag "$IMAGE:$VERSION" .; \
else \
  echo "VERSION param required"; \
fi