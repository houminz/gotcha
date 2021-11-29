## Compile

```
$ protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    helloworld/helloworld.proto
```

## Run

```server
$ go run server.go                                                                                                ─╯
2021/11/29 11:11:19 Received: world
```

```client
$ go run client.go                                                                                                                              ─╯
2021/11/29 11:11:19 Greeting: Hello world
```

## Reference

- https://grpc.io/docs/languages/go/quickstart/
