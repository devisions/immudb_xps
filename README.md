# immudb Experiments

Experiments using [immudb](https://codenotary.io/technologies/immudb/), a lightweight and high-speed immutable database.

<br/>

## Starting the components

Use `start_immudb_comps.sh` script that starts the server (`immudb`) and the gateway (`immugw`).

<br/>

## Interacting with components

### Using CLI client

Using the CLI client (`immuclient`), you can talk directly with the server (`immudb`):

```shell
$ immuclient login immudb
Password:
Successfully logged in.
immudb user has the default password: please change it to ensure proper security
$
$ immuclient status              # pinging to check if server connection is alive
Health check OK
$
$ immuclient current             # getting the last merkle tree root and index stored locally
immudb is empty
$
$ immuclient set k0 v0           # setting a value for key k0
index:		0
key:		k0
value:		v0
hash:		10bb348095d289f0731c7ef115e5abc5dc438ecb3ea784aac3c9d7354fcc0fc8
time:		2020-10-02 15:34:35 +0300 EEST

$ immuclient set k0 v1           # setting another value for key k0
index:		1
key:		k0
value:		v1
hash:		a02fd8cbc78003ad1f86000bbf936f25636d0f48be21572c66eee0bb843b8c8f
time:		2020-10-02 15:34:45 +0300 EEST

$ immuclient set k1 v2           # setting a value for key k1
index:		2
key:		k1
value:		v2
hash:		87dd1398817dc7d87f93f919bc5327a0d53872a98ac6323a69c81e0e6acd8e74
time:		2020-10-02 15:34:59 +0300 EEST

$ immuclient get k0              # getting the value for key k0
index:		1
key:		k0
value:		v1
hash:		a02fd8cbc78003ad1f86000bbf936f25636d0f48be21572c66eee0bb843b8c8f
time:		2020-10-02 15:34:45 +0300 EEST

$ immuclient current             # getting the last merkle tree root and index stored locally
index:		2
hash:		f5f84dbfe5dcfb8cfd4512aff1f30f813b83f1a472ae16b4779055ed5c9b1b51

$
```
