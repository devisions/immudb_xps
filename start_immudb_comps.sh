#!/bin/sh

## __________________________________________________________________
## 
## file: start_immudb_comps.sh
## info: This shell script starts the two main server side components
##       (immudb and immugw) as Docker containers.
## __________________________________________________________________
## 


NET_NAME=immudb_xps_net
VOL_NAME=immudb_xps_vol


echo
echo ">>> Creating prereqs: ${NET_NAME} network and ${VOL_NAME} volume ..."
echo


NETCOUNT=`docker network ls | grep ${NET_NAME} | wc -l`

if [ $NETCOUNT -gt 0 ]; then
    echo ">>> Network ${VOL_NAME} already exists. Skipped creating it."
else
    docker network create ${NET_NAME}
    echo ">>> Network ${NET_NAME} was created."
fi


VOLCOUNT=`docker volume ls | grep ${VOL_NAME} | wc -l`

if [ $VOLCOUNT -gt 0 ]; then
    echo ">>> Volume ${VOL_NAME} already exists. Skipped creating it."
else
    docker volume create ${VOL_NAME}
    echo ">>> Volume ${VOL_NAME} was created."
fi


echo
echo ">>> Starting immudb Server container ..."
echo

docker run -it -d --name immudb -p 3322:3322 -p 9497:9497   \
       -v ${VOL_NAME}:/var/lib/immudb                       \
       --network ${NET_NAME}                                \
       codenotary/immudb:latest

## Sleeping (to give time for the server to be available) 
## is one working approach. However, a better one is to
## periodically check for the readiness of the server.
sleep 3


echo
echo ">>> Starting immudb Gateway container ..."
echo

docker run -it -d -p 3323:3323 --name immugw    \
       --env IMMUGW_IMMUDB_ADDRESS=immudb       \
       --network ${NET_NAME}                    \
       codenotary/immugw:latest

