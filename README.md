# Fetch - Receipt Processor Challenge

`receipt-processor` is a service application which facilitates the processing of receipts, and the awarding of points associated with those receipts.

This Dockerized application is written by **ashyrae** (Ashelyn) for Go 1.23.4, utilizing Protocol Buffers ([`protobuf`](https://github.com/protocolbuffers/protobuf)), [`grpc-go`](https://github.com/grpc/grpc-go), and [`grpc-gateway`](https://github.com/grpc-ecosystem/grpc-gateway).

## Running the Service

You can run this application with a local installation of the Go runtime, or via Docker Compose.

### Go Runtime (v1.23.4+)

```shell
go mod verify
go run main.go
```

### Docker Compose

```shell
docker compose up
```

## Using the Service

By default, the server is bound to `localhost:8081`.

### Endpoints

The Receipt Processor service has two `HTTP` endpoints, `ProcessReceipt` and `AwardPoints`.

These can be accessed through via `curl`, Postman, or similar tools, conforming to the Challenge API Spec.

```shell
curl -X POST localhost:8081/receipts/process -H "content-type: application/json" -d @receipt-processor/api/challenge-api-spec/simple-receipt.json
    
curl localhost:8081/receipts/{your-receipt-id}/points
```

## Rationale & Post-mortem

### Why Golang?

The bulk of my professional experience is in Golang - and frankly, I'd written too many *Scala 3 & ZIO* HTTP services recently.

This project served as a good way to de-rust & refamiliarize myself with Go.

### Why use `gRPC`?

In an environment demanding scalability, creating services that can easily & performantly communicate with each other is crucial.
`gRPC` provides a solid foundation for this, designed for long-lived connections & real-time bidirectional streams, with built-in load balancing.

While `receipt-processor` does not leverage these capabilities (to conform with the challenge API spec),
`gRPC` presents a wide variety of possibilities as to future refactors & new features that could be proposed.

For example, a future refactor of `ProcessReceipt` could enable the processing of multiple receipts concurrently from a single request,
returning a stream of responses containing `ID`s for each `Receipt` - and preserve the order in which they are returned.

The same logic follows for a potential refactor of `AwardPoints`, where multiple receipt IDs could be provided,
returning a stream of responses, containing point award amounts for the provided receipt IDs.

### Why not use `HTTP` & configure this as a REST API application?

REST APIs are definitely simpler, and sometimes more ideal for CRUD (Create/Retrieve/Update/Delete) applications.

However, `HTTP` services come with the cost of higher latency comparative to `gRPC`, and primarily use `JSON` & `XML` for serialization.
With `gRPC` & Protocol Buffers for binary serialization, payload size is lower & transmission over `HTTP/2` is faster.

A platform like Fetch, designed for large volumes of users, would benefit significantly in a `gRPC` environment.

### What challenges did you encounter?

There are two primary pain points with `gRPC` services, both of which I encountered here.

The first is compiling Protocol Buffers (`*.proto` files), & generating the appropriate code. `protoc` does not generate Go code out of the box, requiring the `protoc-gen-go` plugin at bare minimum. Sometimes, additional plugins are also necessary (like for generating `grpc-gateway` code).

The second is serialization. In a `gRPC` server, communication from a `gRPC` client using protobufs is preferable.
Ensuring that `JSON` input maps appropriately to `protobuf` messages can be tedious, particularly if learning `*.proto` files for the first time.

While protobufs do support `JSON` for serialization, it's for compatibility reasons only. `gRPC` will always perform better with `protobuf`.

### What would you do differently next time?

Firstly, I'd implement a proper database, rather than my in-memory solution. In a production environment, this data would need to not only be persistent,
but also updated more frequently and faster. A read-write mutex on a string-keyed map of `Receipt` objects is just simply too inefficient to work at scale.

Secondly, the error handling & logging in this project is not up-to-par with production standards, & needs to be revisited.

### What if you could change the API spec?

As mentioned previously, leveraging streams with `gRPC` would drastically increase efficiency in processing receipts & awarding points for those receipts.

I would implement both endpoints as bidirectional streaming RPCs, instead of unary (single request, single response) RPCs - allowing multiple requests to be provided at once, and responses to be received, preserving the order in which the requests were sent.
