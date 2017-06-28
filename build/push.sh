#!/usr/bin/env bash

./build-app.sh
./build-image.sh

docker login registry.nulana.com -u danil@nulana.com -p "@23[884#%*)38#6,"
docker push registry.nulana.com/bobrnor/battleship-client:latest
docker logout registry.nulana.com
