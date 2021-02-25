## immudb 0.9.x Experiments

Included:
- A *walk* through each entry of each transaction.

### Walk

Usage example: `go run walk/main.go -pass {password} -db {dbName}`.

Note that the arguments' default values are the same as the server ones.<br/>
Therefore, you can just run the sample without any additional argument,<br/>
just like this: `go run walk/main.go`

To get the complete usage help, pass `-h` argument:
```shell
$ go run walk/main.go -h                                                 
Usage of /tmp/go-build751243515/b001/exe/main:
  -db string
        name of the database to use (default "defaultdb")
  -host string
        hostname or IP adddress of immudb's listening endpoint (default "127.0.0.1")
  -pass string
        password to authenticate (default "immudb")
  -port int
        port of immudb's listening endpoint (default 3322)
  -user string
        username to authenticate (default "immudb")
$ 
```
