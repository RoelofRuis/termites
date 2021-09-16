#!/usr/bin/env bash

set -e

docker image build \
    -t node-dev \
    .

docker run -it --rm \
    -v ${PWD}:/app \
    --network host \
    node-dev $@