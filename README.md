
### Reference
- [Memberlist + KV](https://reintech.io/blog/implementing-distributed-key-value-store-go/)
- [Swim Ring](https://github.com/hungys/swimring)
- [Paypal JunoDB](https://github.com/paypal/junodb)
- [Uber RingPop](https://github.com/uber/ringpop-go)

### Run

```sh
go build
```

```shell
./dist_kv
```


```shell
./dist_kv -gossip=9000 -http=9001
```

```shell
./dist_kv -gossip=10000 -http=10001
```

```shell
curl -Ss -XPOST "http://localhost:9001/put/1/alex"
curl -Ss -XPOST "http://localhost:9001/put/2/benny"
curl -Ss -XPOST "http://localhost:9001/put/3/cassie"
```

```shell
curl -Ss -XGET "http://localhost:8001/shards"
curl -Ss -XGET "http://localhost:9001/shards"
```

```shell
curl -Ss -XGET "http://localhost:8001/get/1" 
```

### TODO
- Redistribution