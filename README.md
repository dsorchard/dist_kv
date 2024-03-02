
### Reference
- [Memberlist + KV](https://reintech.io/blog/implementing-distributed-key-value-store-go/)

### Run

```sh
go build
./pkg
```

```shell
./pkg -gossip=9000 -http=9001
```

```shell
curl -Ss -XPOST "http://localhost:9001/put/1/habibi"
```

```shell
curl -Ss -XGET "http://localhost:8001/get/1" 
```
