#!/usr/bin/env bash

set -e

./build-app.sh
./build-image.sh
./run-container.sh