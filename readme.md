# Unity Back End Engineer Assignment

## Folder Structure

- example - keeps hub and client usage examples
- pkg - keeps hub, client and protocol packages

## Data Structure

A hub message consist of fields below

- **Command** _IDENTITY, LIST, RELAY_
- **Id** _Client id of message sender_
- **ConnectedClientIds** _Array of other connected client ids._
- **RelayTo** _Array of relay message reciever ids._
- **BodyType** _PLAIN TEXT, JSON, ERROR_
- **Body** _Byte array_

## Prerequisites

Google Protocol Buffers is used for serializing structured data to consume consumes minimal amount resources.

You should download and install protocol buffer compiler from [here](https://github.com/protocolbuffers/protobuf)

## Usage

- Download and install google protocol buffer compiler
- Run the hub in a new terminal

```
go run example/hub/main.go
```

- Run the client as much as you want in the new terminals

```
go run example/client/main.go
```

- Send some test commands from one of client terminals. Available test commands:
  - **identity** _Hub will answer the id of client_
  - **list** _Hub will answer the other connected client ids_
  - **relay** _Hub will send data the other connected client ids. Client ids are hardcoded in the file for demo purposes_

### Running the tests

```
go test -v ./pkg/... -cover
```

## Notes

- Files in the protocol folder are autogenerated with [protoc](https://github.com/golang/protobuf/tree/master/protoc-gen-go), hence didn't write tests.
-
