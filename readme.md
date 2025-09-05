Here’s a complete, polished **README.md** you can drop into your GitHub repo. It’s step-by-step, production-ready, and includes quick start, Docker, Kubernetes, gRPC-Gateway, protoc generation, API examples, and troubleshooting.

---

# Go gRPC Attendance Service

A simple yet robust **attendance management service** built with **Go**, **gRPC**, **MongoDB**, and a **REST JSON gateway**. It supports **check-in**, **check-out**, and **retrieval** of attendance records. The project demonstrates modern microservice best practices: clean architecture, containerization, and cloud-native deployment.

<p align="left">
  <a href="https://go.dev/"><img alt="Go" src="https://img.shields.io/badge/Go-1.20%2B-00ADD8?logo=go&logoColor=white"></a>
  <a href="https://grpc.io/"><img alt="gRPC" src="https://img.shields.io/badge/gRPC-Enabled-32CD32?logo=googlecloud&logoColor=white"></a>
  <a href="https://www.mongodb.com/"><img alt="MongoDB" src="https://img.shields.io/badge/MongoDB-6.x-47A248?logo=mongodb&logoColor=white"></a>
  <a href="https://www.docker.com/"><img alt="Docker" src="https://img.shields.io/badge/Docker-Ready-2496ED?logo=docker&logoColor=white"></a>
  <a href="https://kubernetes.io/"><img alt="Kubernetes" src="https://img.shields.io/badge/Kubernetes-Minikube-326CE5?logo=kubernetes&logoColor=white"></a>
</p>

---

## Table of Contents

* [Architecture](#architecture)
* [Project Structure](#project-structure)
* [Prerequisites](#prerequisites)
* [Quick Start](#quick-start)
* [Run Locally (no Docker)](#run-locally-no-docker)
* [Run with Docker Compose](#run-with-docker-compose)
* [Run on Kubernetes (Minikube)](#run-on-kubernetes-minikube)
* [Environment Variables](#environment-variables)
* [Generate Protobuf & Gateway Code](#generate-protobuf--gateway-code)
* [API Reference (REST via gRPC-Gateway)](#api-reference-rest-via-grpcgateway)
* [Inspecting MongoDB](#inspecting-mongodb)
* [Sample Dockerfile](#sample-dockerfile)
* [Sample docker-compose.yml](#sample-docker-composeyml)
* [Troubleshooting](#troubleshooting)
* [Development Notes](#development-notes)
* [License](#license)

---

## Architecture

```mermaid
flowchart LR
  C[Client / curl / Browser] -- REST JSON --> G[HTTP Gateway (grpc-gateway)]
  G -- gRPC --> S[gRPC Server (Go)]
  S -- MongoDB Driver --> DB[(MongoDB)]
```

* **Clients** send REST requests to the **HTTP Gateway**.
* The **Gateway** translates REST → gRPC and calls the **Go gRPC server**.
* The server applies business logic, stores data in **MongoDB**, and responds back.
* **Timezone:** Timestamps are **stored in UTC** in MongoDB and **rendered as IST** (Asia/Kolkata) in responses.

---

## Project Structure

```
GO_ATTENDANCE_SERVICE/
├── googleapis/                 # Google API protos (dependency for gRPC-Gateway)
├── mychart/                    # (Optional) Helm chart for Kubernetes deployment
├── proto/
│   ├── attendance.proto        # gRPC service definition with HTTP annotations
│   ├── attendance.pb.go        # Generated messages
│   ├── attendance_grpc.pb.go   # Generated gRPC client/server
│   └── attendance.pb.gw.go     # Generated REST <-> gRPC gateway
├── attendance-deployment.yaml  # K8s manifest for the attendance service
├── mongo-deployment.yaml       # K8s manifest for MongoDB
├── docker-compose.yml          # Local dev: MongoDB + service
├── Dockerfile                  # Image build for the Go service
├── go.mod                      # Go module definition
├── go.sum                      # Dependency lockfile
├── main.go                     # Starts gRPC server + HTTP gateway
└── README.md                   # You are here
```

> **Note:** If you change `proto/attendance.proto`, you **must** regenerate the `.pb.go` files (see [Generate Protobuf & Gateway Code](#generate-protobuf--gateway-code)).

---

## Prerequisites

* **Go** 1.20+
* **protoc** (Protocol Buffers compiler)
* **gRPC / Protobuf plugins for Go**

  ```bash
  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
  go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
  ```
* **Docker** & **Docker Compose**
* **kubectl** & **minikube** (for Kubernetes)
* **curl** (for REST testing)
* **git**

---

## Quick Start

Run everything with Docker Compose and test a check-in:

```bash
docker-compose up --build -d
curl -X POST http://localhost:8080/v1/checkin \
  -H "Content-Type: application/json" \
  -d '{"user_id":"emp_local_01", "username":"Aashish"}'
```

Expected response (IDs will be MongoDB ObjectIDs):

```json
{
  "id": "64f4e7b7a2d9f34a8b7b1e90",
  "userId": "emp_local_01",
  "username": "Aashish",
  "checkinTime": "2025-09-05 10:26:00",
  "statusMessage": "User checked in successfully."
}
```

---

## Run Locally (no Docker)

1. Start MongoDB:

```bash
docker run -d --name mongo-db -p 27017:27017 mongo:6.0
```

2. Install deps & run:

```bash
go mod tidy
go run main.go
```

3. Endpoints:

* gRPC: **localhost:50051**
* HTTP (REST): **localhost:8080**

---

## Run with Docker Compose

```bash
docker-compose up --build -d
docker-compose logs -f attendance-app
# stop
docker-compose down
```

---

## Run on Kubernetes (Minikube)

```bash
minikube start

# Deploy DB then app
kubectl apply -f mongo-deployment.yaml
kubectl apply -f attendance-deployment.yaml

# Check resources
kubectl get pods
kubectl get svc

# Port-forward HTTP gateway locally
kubectl port-forward svc/attendance-svc 8080:8080
```

Open another terminal and test with curl (see [API Reference](#api-reference-rest-via-grpcgateway)).

---

## Environment Variables

| Variable    | Default                     | Description               |
| ----------- | --------------------------- | ------------------------- |
| `MONGO_URI` | `mongodb://localhost:27017` | MongoDB connection string |

> In Docker Compose, the app uses `mongodb://mongo:27017` (service name `mongo`).

---

## Generate Protobuf & Gateway Code

1. Make sure `googleapis/` exists:

```bash
git clone https://github.com/googleapis/googleapis.git
```

2. From project root:

```bash
protoc \
  --proto_path=. \
  --proto_path=googleapis \
  --go_out=. \
  --go-grpc_out=. \
  --grpc-gateway_out=. \
  proto/attendance.proto
```

This regenerates:

* `proto/attendance.pb.go`
* `proto/attendance_grpc.pb.go`
* `proto/attendance.pb.gw.go`

---

## API Reference (REST via gRPC-Gateway)

Base URL: **[http://localhost:8080](http://localhost:8080)**

### 1) Check-In

* **POST** `/v1/checkin`

**Request**

```json
{
  "user_id": "emp_local_01",
  "username": "Aashish"
}
```

**Response**

```json
{
  "id": "64f4e7b7a2d9f34a8b7b1e90",
  "userId": "emp_local_01",
  "username": "Aashish",
  "checkinTime": "2025-09-05 10:26:00",
  "statusMessage": "User checked in successfully."
}
```

**curl**

```bash
curl -X POST -H "Content-Type: application/json" \
  -d '{"user_id":"emp_local_01","username":"Aashish"}' \
  http://localhost:8080/v1/checkin
```

---

### 2) Get Latest Attendance for a User

* **GET** `/v1/attendance/{user_id}`

**Response**

```json
{
  "id": "64f4e7b7a2d9f34a8b7b1e90",
  "userId": "emp_local_01",
  "username": "Aashish",
  "checkinTime": "2025-09-05 10:26:00",
  "checkoutTime": "2025-09-05 18:30:00",
  "statusMessage": "Record found."
}
```

**curl**

```bash
curl http://localhost:8080/v1/attendance/emp_local_01
```

---

### 3) Check-Out

* **PUT** `/v1/checkout/{record_id}`

**Request**

```json
{}
```

**Response**

```json
{
  "id": "64f4e7b7a2d9f34a8b7b1e90",
  "userId": "emp_local_01",
  "username": "Aashish",
  "checkinTime": "2025-09-05 10:26:00",
  "checkoutTime": "2025-09-05 18:30:00",
  "statusMessage": "User checked out successfully."
}
```

**curl**

```bash
curl -X PUT -H "Content-Type: application/json" \
  -d '{}' \
  http://localhost:8080/v1/checkout/64f4e7b7a2d9f34a8b7b1e90
```

> **Note:** `record_id` is the MongoDB ObjectID returned by **Check-In**.

---

### 4) Get All Attendance Records

* **GET** `/v1/attendance`

**Response**

```json
{
  "records": [
    {
      "id": "64f4e7b7a2d9f34a8b7b1e90",
      "userId": "emp_local_01",
      "username": "Aashish",
      "checkinTime": "2025-09-05 10:26:00",
      "checkoutTime": "2025-09-05 18:30:00",
      "statusMessage": "Record retrieved."
    }
  ]
}
```

**curl**

```bash
curl http://localhost:8080/v1/attendance
```

---

## Inspecting MongoDB

Open a shell into the running MongoDB container:

```bash
docker exec -it mongo-db mongosh
```

Then:

```javascript
use attendance_db
db.records.find().pretty()
```

---

## Sample Dockerfile

```dockerfile
# --- Build stage ---
FROM golang:1.22 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o attendance-service .

# --- Runtime stage ---
FROM debian:bullseye-slim
WORKDIR /root/

COPY --from=builder /app/attendance-service .

EXPOSE 50051 8080
CMD ["./attendance-service"]
```

---

## Sample docker-compose.yml

```yaml
version: "3.9"
services:
  mongo:
    image: mongo:6.0
    container_name: mongo-db
    ports:
      - "27017:27017"
    volumes:
      - mongo_data:/data/db

  attendance-app:
    build: .
    container_name: attendance-app
    depends_on:
      - mongo
    environment:
      - MONGO_URI=mongodb://mongo:27017
    ports:
      - "50051:50051"
      - "8080:8080"

volumes:
  mongo_data:
```

---

## Troubleshooting

**Container name conflict**

```bash
docker rm -f attendance-app
```

**Binary not found in container**

* Ensure Dockerfile copies the built binary and `CMD ["./attendance-service"]` matches the path.
* Use the provided Dockerfile.

**Mongo connection refused**

* Verify `MONGO_URI`:

  * Local run: `mongodb://localhost:27017`
  * Docker Compose: `mongodb://mongo:27017`
* Confirm Mongo is up: `docker ps`, logs: `docker logs mongo-db`.

**Wrong time shown**

* DB stores **UTC**; responses convert to **IST** (`Asia/Kolkata`). Ensure you format times with `.In(ist)` on response.

**Protoc plugins not found**

* Ensure `$GOPATH/bin` (or Go bin path) is on your `PATH`.
* Re-run the `go install` commands in [Prerequisites](#prerequisites).

**bson.D unkeyed struct error**

* Use keyed fields:

  ```go
  opts := options.FindOne().SetSort(bson.D{{Key: "checkin_time", Value: -1}})
  ```

---

## Development Notes

* **Code style:** store timestamps as `time.Now().UTC()`; present as IST in responses.
* **Indexing:** consider indexes for faster lookups:

  ```js
  db.records.createIndex({ user_id: 1, checkin_time: -1 })
  ```
* **Health checks:** add a simple health RPC/REST if needed.
* **gRPC tooling:** you can test gRPC directly with `grpcurl` or `evans`.

---

## License

MIT — feel free to use, modify, and share.

---

**Happy shipping! 🚀**
If you want, I can also add a small **architecture diagram image** or refine the **Helm chart** in `mychart/` to match this setup.
