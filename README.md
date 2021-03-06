# Worker Service

## Development

Start the server.

```BASH
go run ./server/main.go
```

Start the client.

```BASH
go run ./client/main.go
```

## Docker

Start the server and client.

```
git checkout docker
docker-compose up -d
```

## API

### Server - gRPC API

Running on port: <http://localhost:8500>

Method|Description
-|-
GetWorker|Dequeues a worker.
PutWorker|Enqueues a worker.
ListWorkers|Lists workers.
ShowWorker|Shows a worker.

### Server - HTTP API

Running on port: <http://localhost:8000>

Method|Path|Description
-|-|-
GET|/worker|Dequeues a worker.
PUT|/worker|Enqueues a worker.
GET|/workers|Lists workers.
GET|/workers/{n}|Shows a worker.

### Client - HTTP API

Running on port: <http://localhost:9000>

Method|Path|Description
-|-|-
GET|/api/worker|Dequeues a worker.
PUT|/api/worker|Enqueues a worker.
GET|/api/workers|Lists workers.
GET|/api/workers/{n}|Shows a worker.
GET|/api/workers/summon/async/{a}/sync/{s}|Dequeues and enqueues workers.

## Compiling

Compile the server.

```BASH
cd server
go build
```

Compile the client.

```BASH
cd client
packr build
```

## Testing

Test the server.

```BASH
go test ./server/app
```

Test the client.

```BASH
go run server/main.go
go test ./client/handler
```
