#!/bin/sh

## __________________________________________________________________
## 
## file: start_immudb_comps.sh
## info: This shell script starts the two main server side components
##       (immudb and immugw) as Docker containers.
## __________________________________________________________________
## 

## TODO: Consider checking first if not already exists
## using something like: docker network ls | grep immudb_net | wc -l

echo
echo ">>> Creating prereqs: immudb_xps_net network and immudb_xps_data volume ..."
echo

docker network create immudb_xps_net

docker volume create immudb_xps_data

echo
echo ">>> Starting immudb Server container ..."
echo

docker run -it -d --name immudb -p 3322:3322 -p 9497:9497   \
       -v immudb_xps_data:/var/lib/immudb                   \
       --network immudb_xps_net                             \
       codenotary/immudb:latest

sleep 3

echo
echo ">>> Starting immudb Gateway container ..."
echo

docker run -it -d -p 3323:3323 --name immugw    \
       --env IMMUGW_IMMUDB_ADDRESS=immudb       \
       --network immudb_xps_net                 \
       codenotary/immugw:latest

