#!/usr/bin/env bash

for i in `seq 1 8`;
do
    docker stop battleship-client-$i | true
    docker rm battleship-client-$i | true
    docker run -d --name=battleship-client-$i \
		--volume /private/var/battleship/client:/external \
		--network=bridge \
		--network=battleship-network \
		battleship-client
done
