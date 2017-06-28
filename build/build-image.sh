#!/usr/bin/env bash

docker build --force-rm -t registry.nulana.com/bobrnor/battleship-client -f Dockerfile .
rm battleship
