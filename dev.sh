#!/usr/bin/env bash

clean() {
  echo "stop containers";
  docker container stop dc_tnt.mf dc_db.mf collector.mf
  echo "drop containers"
  docker rm -v dc_tnt.mf dc_db.mf collector.mf
}

clean

FILE_HASH=$(git rev-parse HEAD)
export GIT_HASH=$FILE_HASH

echo "RUN docker-compose-dev.yml "
#serviceList="dc_db dc_tarantool collector_app"
serviceList="dc_db dc_tarantool"
echo "RUNNING SERVICES: $serviceList"
docker-compose -f docker-compose.yml pull
docker-compose -f docker-compose.yml up --build ${serviceList}