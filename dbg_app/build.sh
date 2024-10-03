#!/usr/bin/env sh

docker-compose run --rm node sh -c "yarn mix build --production"

cp ./server/debugger.js ./../pkg/termites_dbg/debugger.js