#!/bin/sh

## __________________________________________________________________
## 
## file: start_immudb_comps.sh
## meta: 2020-10-02 | 2nd version: working now
## info: This shell script starts the two main server side components
##       (immudb and immugw) as Docker containers
## __________________________________________________________________
## 

## TODO: Consider checking first if not already exists
## using something like: docker network ls | grep immudb_net | wc -l

echo
echo ">>> Creating immudb_net network ..."
echo
docker network create immudb_net

echo
echo ">>> Starting immudb Server ..."
echo

docker run -it -d --name immudb -p 3322:3322 -p 9497:9497             \
       -v "$(pwd)/.immudb_data:/var/lib/immudb" --network immudb_net \
       codenotary/immudb:latest

sleep 3

echo
echo ">>> Starting immudb Gateway ..."
echo

docker run -it -d -p 3323:3323 --name immugw                   \
       --env IMMUGW_IMMUDB_ADDRESS=immudb --network immudb_net \
       codenotary/immugw:latest

