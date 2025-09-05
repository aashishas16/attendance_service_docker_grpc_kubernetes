Perfect 🙌 Thanks for sharing your project structure.
I’ll prepare a **`README.md`** for you that:

1. Explains **purpose of each file/folder** in your structure.
2. Gives **setup, build, and run instructions** (local + Docker + Kubernetes).
3. Includes **examples with `curl`** for testing.

Here’s your README code ⬇️

---

```markdown
# Attendance Management Service (Go + gRPC + MongoDB)

This project is a simple **Attendance Management Service** written in Go, using **gRPC** for inter-service communication and **MongoDB** as the database.  
It also includes a lightweight **HTTP wrapper** so you can test APIs easily with `curl` or Postman.

---

## 📂 Project Structure

```

GO\_ATTENDANCE\_SERVICE/
├── googleapis/                 # Google API proto definitions (for gRPC-Gateway if needed)
├── mychart/                    # Helm chart for Kubernetes deployment
├── proto/                      # Protobuf files and generated Go code
│   ├── attendance.proto        # gRPC service definition
│   ├── attendance.pb.go        # Generated Go code for messages
│   ├── attendance\_grpc.pb.go   # Generated Go code for gRPC service
│   ├── attendance.pb.gw\.go     # gRPC-Gateway generated file (REST <-> gRPC bridge)
├── attendance-deployment.yaml  # Kubernetes deployment manifest for attendance service
├── mongo-deployment.yaml       # Kubernetes deployment manifest for MongoDB
├── docker-compose.yml          # Docker Compose to run MongoDB + service locally
├── Dockerfile                  # Build attendance service image
├── go.mod                      # Go module definition
├── go.sum                      # Go dependencies lock file
├── main.go                     # Entry point: starts gRPC + HTTP server
├── readme.md                   # This README file

````

---

## ⚙️ Prerequisites

- [Go](https://go.dev/) >= 1.20
- [protoc](https://grpc.io/docs/protoc-installation/) (Protocol Buffers compiler)
- gRPC Go plugins:
  ```bash
  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
  go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
````

* [Docker](https://www.docker.com/)
* [kubectl](https://kubernetes.io/docs/tasks/tools/) + [minikube](https://minikube.sigs.k8s.io/docs/) (if running on Kubernetes)

---

## 🚀 Running the Project

### 1️⃣ Run Locally (without Docker)

1. Start MongoDB (either installed locally or via Docker):

   ```bash
   docker run -d --name mongo-db -p 27017:27017 mongo:6.0
   ```

2. Run the service:

   ```bash
   go mod tidy
   go run main.go
   ```

3. Service endpoints:

   * gRPC: `localhost:50051`
   * HTTP (REST wrapper): `localhost:8080`

---

### 2️⃣ Run with Docker Compose

1. Build and start services:

   ```bash
   docker-compose up --build -d
   ```

2. Check running containers:

   ```bash
   docker ps
   ```

3. Logs:

   ```bash
   docker-compose logs -f attendance-app
   ```

---

### 3️⃣ Run on Kubernetes

1. Apply MongoDB deployment:

   ```bash
   kubectl apply -f mongo-deployment.yaml
   ```

2. Apply Attendance service:

   ```bash
   kubectl apply -f attendance-deployment.yaml
   ```

3. Check pods:

   ```bash
   kubectl get pods
   ```

4. Port-forward for local access:

   ```bash
   kubectl port-forward svc/attendance-svc 8080:8080
   ```

---

## 📡 API Usage (via HTTP wrapper)

### ✅ Check-in

```bash
curl -X POST http://localhost:8080/checkin \
  -H "Content-Type: application/json" \
  -d '{"user_id":"u1","username":"Alice","address":"Office"}'
```

### ✅ Check-out

```bash
curl -X POST http://localhost:8080/checkout \
  -H "Content-Type: application/json" \
  -d '{"record_id":"<paste_record_id_here>"}'
```

### ✅ Get Attendance

```bash
curl "http://localhost:8080/attendance?user_id=u1"
```

---

## 🗄 Inspect MongoDB

Enter Mongo shell:

```bash
docker exec -it mongo-db mongosh
```

Inside:

```js
use attendance_db
db.records.find().pretty()
```

---

## 📌 Notes

* `proto/attendance.proto` defines gRPC services and messages.
* gRPC clients (Go, Python, Node.js, etc.) can call directly on port `50051`.
* REST (via HTTP wrapper or gRPC-Gateway) is available on port `8080`.
* For Kubernetes, service discovery is handled via `attendance-svc` and `mongo-svc`.

---

## ✅ Quick Start

One-liner to run everything with Docker:

```bash
docker-compose up --build -d
curl -X POST http://localhost:8080/checkin -H "Content-Type: application/json" -d '{"user_id":"u1","username":"Alice","address":"Office"}'
```

---

```

---

👉 This README gives **file explanations + setup instructions + curl examples**.  

Do you want me to also add a **diagram of system architecture** (gRPC <-> HTTP <-> MongoDB) inside your README?
```
