#!/usr/bin/env bash

docker build --force-rm -t battleship-client -f Dockerfile .
rm battleship