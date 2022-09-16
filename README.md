# loadbalancer-in-go

Building a loadbalancer from scratch. DO NOT use this in production :p (just in case)

uses round robin algorithm.

## Running it

Start multiple instances of `app.go`

```bash

# spawn multiple servers
go run application/app.go --server="server-$i" --port=500$i
```

Start the loadbalancer

```bash
go run loadbalancer.go
```

Bombard the loadbalancer with requests

```bash
for i in {1..20}; do curl 127.0.0.1:8000; done
```
