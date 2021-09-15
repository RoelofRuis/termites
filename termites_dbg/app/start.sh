#!/usr/bin/env bash

set -e

docker image build \
    -t node-dev \
    .

docker run -it --rm \
    -v ${PWD}:/app \
    -v node_modules:/node_modules \
    --network host \
    node-dev $@