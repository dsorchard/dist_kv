# Dist KV
This is a tiny distributed key-value store implemented in Go.
Majority of the code is based on [MemberList](https://github.com/hashicorp/memberlist) and
[Consistent](https://github.com/buraksezer/consistent) Library to simplify the implementation.

### Features
- Replication and Quorum
- Redistribution of Shards on Node Join/Leave
- Membership List is maintained using SWIM Protocol
- Consistent Hashing for Shard Distribution


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
curl -Ss -XPOST "http://localhost:9001/kv/1/alex"
curl -Ss -XPOST "http://localhost:9001/kv/2/benny"
curl -Ss -XPOST "http://localhost:9001/kv/3/cassie"
```

```shell
curl -Ss -XGET "http://localhost:8001/shards"
curl -Ss -XGET "http://localhost:9001/shards"
curl -Ss -XGET "http://localhost:10001/shards"
```

```shell
curl -Ss -XGET "http://localhost:8001/kv/1" 
curl -Ss -XGET "http://localhost:9001/kv/2" 
curl -Ss -XGET "http://localhost:10001/kv/3" 
```


### TODO
- ~~Redistribution~~
- ~~Logging~~
- ~~Replication~~
- ~~Remove HTTP Port Number Hack~~


### Reference
- [Memberlist + KV](https://reintech.io/blog/implementing-distributed-key-value-store-go/)
- [Swim Ring](https://github.com/hungys/swimring)
- [Paypal JunoDB](https://github.com/paypal/junodb)
- [Uber RingPop](https://github.com/uber/ringpop-go)
- [OlricDB](https://github.com/buraksezer/olricdb/blob/f24016ca0379a2f0c652a1d38d04953f440d20e0/routing.go#L264)