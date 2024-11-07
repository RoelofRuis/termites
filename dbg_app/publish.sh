#!/usr/bin/env sh

docker-compose run --rm node sh -c "npm run build"

cp ./dist/index.html ./../pkg/termites_dbg/app/index.html