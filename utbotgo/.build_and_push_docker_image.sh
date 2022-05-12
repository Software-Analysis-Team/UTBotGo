#!/bin/bash
source .env
make .docker_image > /dev/null || ( \
    echo "$1" | docker login ghcr.io -u sergeykuz1001 --password-stdin; \
    docker build -t "$UTBOTGO_DOCKER_IMAGE" .; \
    docker push "$UTBOTGO_DOCKER_IMAGE"; \
    make .docker_image > /dev/null \
)
