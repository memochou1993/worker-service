# Worker Service

## Development

Start the server.

```BASH
go run server/main.go
```

Start the client.

```BASH
go run client/main.go
```

## API

### Server gRPC API

Running on port: http://localhost:8500

#### GetWorker

Dequeues a worker.

#### PutWorker

Enqueues a worker.

#### ListWorkers

Lists workers.

#### ShowWorker

Shows a worker.

### Server HTTP API

Running on port: http://localhost:8000

#### `GET` /worker

Dequeues a worker.

#### `PUT` /worker

Enqueues a worker.

#### `GET` /workers

Lists workers.

#### `GET` /workers/{n}

Shows a worker.

### Client HTTP API

Running on port: http://localhost:9000

#### `GET` /api/worker

Dequeues a worker.

#### `PUT` /api/worker

Enqueues a worker.

#### `GET` /api/workers

Lists workers.

#### `GET` /api/workers/{n}

Shows a worker.

#### `GET` /api/workers/summon/async/{a}/sync/{s}

Dequeues and enqueues workers.
