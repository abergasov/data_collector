#!/usr/bin/env bash

clean() {
  echo "stop containers";
  docker container stop gcr_db.mf
  echo "drop containers"
  docker rm -v gcr_db.mf
}

clean

FILE_HASH=$(git rev-parse HEAD)
export GIT_HASH=$FILE_HASH

echo "RUN docker-compose-dev.yml "
serviceList="dc_db"
echo "RUNNING SERVICES: $serviceList"
docker-compose -f docker-compose.yml pull
docker-compose -f docker-compose.yml up --build ${serviceList}