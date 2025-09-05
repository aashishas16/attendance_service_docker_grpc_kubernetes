Go gRPC Attendance Service
A simple yet robust attendance management service built with Go, gRPC, MongoDB, and a RESTful JSON gateway. This service allows for checking in, checking out, and retrieving attendance records for users. It serves as a comprehensive example of a modern microservice, demonstrating best practices for creating a scalable and maintainable backend system.

🏛️ System Architecture
The service operates with a simple yet powerful architecture. All external REST API calls from clients (like curl or a web browser) are received by the HTTP Gateway. The gateway translates these requests into the gRPC format and forwards them to the main Go application server. The server contains the core business logic, interacts with the MongoDB database, and then sends the response back through the same path.

📂 Project Structure
A brief overview of the key files and directories in this project.

GO_ATTENDANCE_SERVICE/
├── googleapis/                 # Google API proto definitions (dependency for gRPC-Gateway)
├── mychart/                    # Helm chart for Kubernetes deployment
├── proto/                      # Protobuf files and generated Go code
│   ├── attendance.proto        # The gRPC service definition
│   ├── attendance.pb.go        # Generated Go code for messages
│   ├── attendance_grpc.pb.go   # Generated Go code for the gRPC service client/server
│   └── attendance.pb.gw.go     # Generated Go code for the REST <-> gRPC gateway
├── attendance-deployment.yaml  # Kubernetes manifest for the attendance service
├── mongo-deployment.yaml       # Kubernetes manifest for the MongoDB instance
├── docker-compose.yml          # Docker Compose to run MongoDB + service locally
├── Dockerfile                  # Instructions to build the attendance service Docker image
├── go.mod                      # Go module definition file
├── go.sum                      # Go dependencies lock file
├── main.go                     # Main application: starts the gRPC server and HTTP gateway
└── readme.md                   # This README file

✨ Features
User Check-In & Check-Out: Record attendance with precise timestamps.

Numeric IDs: Uses auto-incrementing integer IDs for records (limited to 1-999).

Timezone Handling: Stores all data in UTC and displays it in Indian Standard Time (IST).

Dual API Support: Exposes a high-performance gRPC service and a user-friendly JSON REST gateway.

Database Integration: Persists all data in a MongoDB database.

Flexible Data Retrieval: Fetch the latest record for a specific user or a complete list of all records.

⚙️ Prerequisites
Before you begin, ensure you have the following installed on your system:

Go: Version 1.18 or higher

Protobuf Compiler (protoc): Installation Guide

Go gRPC Plugins:

go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install [github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest](https://github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest)

Docker & Docker Compose

Git

curl for API testing

🚀 Running the Project
You can run this project locally, with Docker Compose, or on Kubernetes.

1. Run Locally (Without Docker)
Start MongoDB (either locally or via a simple Docker command):

docker run -d --name mongo-db -p 27017:27017 mongo:6.0

Run the Go Service:

go mod tidy
go run main.go

The service is now available:

gRPC Server: localhost:50051

HTTP Gateway: localhost:8080

2. Run with Docker Compose
This is the simplest way to run both the application and the database together.

Build and Start Services:

docker-compose up --build -d

View Logs:

docker-compose logs -f attendance-app

Stop Services:

docker-compose down

3. Run on Kubernetes (using Minikube)
Apply Deployments:

kubectl apply -f mongo-deployment.yaml
kubectl apply -f attendance-deployment.yaml

Port-Forward for Local Access:

kubectl port-forward svc/attendance-svc 8080:8080

📡 API Usage (curl Commands)
All curl commands should be directed to the HTTP gateway on localhost:8080.

1. Check-In an Employee
Method: POST

Endpoint: /v1/checkin

curl -X POST -H "Content-Type: application/json" \
  -d '{"user_id": "emp_local_01", "username": "Aashish"}' \
  http://localhost:8080/v1/checkin

Sample Response:

{
  "id": "1",
  "userId": "emp_local_01",
  "username": "Aashish",
  "checkinTime": "2025-09-05 10:26:00 IST",
  "statusMessage": "User checked in successfully."
}

2. Get a Single User's Latest Attendance
Method: GET

Endpoint: /v1/attendance/{user_id}

# Replace 'emp_local_01' with the user ID you want to find
curl -X GET http://localhost:8080/v1/attendance/emp_local_01

Sample Response:

{
  "id": "1",
  "userId": "emp_local_01",
  "username": "Aashish",
  "checkinTime": "2025-09-05 10:26:00 IST",
  "statusMessage": "Record found."
}

3. Check-Out an Employee
Method: PUT

Endpoint: /v1/checkout/{record_id}

# Replace '1' with the actual ID of the record you want to check out
curl -X PUT -H "Content-Type: application/json" \
  -d '{}' \
  http://localhost:8080/v1/checkout/1

Sample Response:

{
  "id": "1",
  "userId": "emp_local_01",
  "username": "Aashish",
  "checkinTime": "2025-09-05 10:26:00 IST",
  "checkoutTime": "2025-09-05 18:30:00 IST",
  "statusMessage": "User checked out successfully."
}

4. Get All Attendance Records
Method: GET

Endpoint: /v1/attendance

curl -X GET http://localhost:8080/v1/attendance

Sample Response:

{
  "records": [
    {
      "id": "1",
      "userId": "emp_local_01",
      "username": "Aashish",
      "checkinTime": "2025-09-05 10:26:00 IST",
      "checkoutTime": "2025-09-05 18:30:00 IST",
      "statusMessage": "Record retrieved."
    }
  ]
}

💻 Development
Regenerating Protobuf Code
This is a critical step. If you make any changes to proto/attendance.proto, you must regenerate the Go client, server, and gateway files. Run the following command from the project's root directory:

protoc \
  --proto_path=. \
  --proto_path=googleapis \
  --go_out=. \
  --go-grpc_out=. \
  --grpc-gateway_out=. \
  proto/attendance.proto

🗄️ Inspecting the Database
To manually check the data in MongoDB, you can use the Mongo Shell.

Connect to the Container:

docker exec -it mongo-db mongosh

Run Queries:

// Switch to the correct database
use attendance_db

// View all records in a readable format
db.records.find().pretty()







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



Go gRPC Attendance Service
A simple yet robust attendance management service built with Go, gRPC, MongoDB, and a RESTful JSON gateway. This service allows for checking in, checking out, and retrieving attendance records for users.

Features
User Check-In & Check-Out: Record attendance with precise timestamps.

Numeric IDs: Uses auto-incrementing integer IDs for records (limited to 1-999).

Timezone Handling: Stores all data in UTC and displays it in Indian Standard Time (IST).

gRPC & REST Support: Exposes a high-performance gRPC service and a user-friendly JSON REST gateway.

Database Integration: Persists all data in a MongoDB database.

Get Individual Records: Fetch the latest attendance record for a specific user.

Get All Records: Retrieve a complete list of all attendance records in the database.

Prerequisites
Before you begin, ensure you have the following installed on your system:

Go: Version 1.18 or higher.

Protobuf Compiler (protoc): For generating Go code from .proto files.

Go gRPC Plugins: For protoc.

MongoDB: A running instance (local or remote).

Git: For cloning repositories.

curl: For testing the REST API.

🛠️ Setup and Installation
Follow these steps to get the project running locally.

1. Clone the Project
Clone this repository to your local machine:

git clone <your-repository-url>
cd GO_ATTENDANCE_SERVICE

2. Initialize Go Module
If you haven't already, initialize the Go module. The go.mod file is crucial for managing dependencies.

go mod init GO_ATTENDANCE_SERVICE
go mod tidy

3. Get googleapis Dependency
The gRPC-Gateway requires Google's API proto files to function. These must be present in your project directory.

git clone [https://github.com/googleapis/googleapis.git](https://github.com/googleapis/googleapis.git)

Your file structure should now look like this:

GO_ATTENDANCE_SERVICE/
├── googleapis/
├── proto/
│   └── attendance.proto
└── main.go
...

4. Generate Protobuf Code
This is a critical step. Run the protoc command from the root of your project (GO_ATTENDANCE_SERVICE) to generate the necessary Go files from your .proto definition.

protoc \
  --proto_path=. \
  --proto_path=googleapis \
  --go_out=. \
  --go-grpc_out=. \
  --grpc-gateway_out=. \
  proto/attendance.proto

5. Run the Service
Start the Go server. This will launch both the gRPC service on port 50051 and the HTTP gateway on port 8080.

go run main.go

You should see the following output in your terminal, indicating that the servers are running:

✅ gRPC Service is listening on 50051
✅ HTTP Gateway is listening on port 8080

🚀 API Usage (curl Commands)
All curl commands should be directed to the HTTP gateway on port 8080.

1. Check-In an Employee
Method: POST

Endpoint: /v1/checkin

curl -X POST -H "Content-Type: application/json" \
  -d '{"user_id": "emp_local_01", "username": "Aashish"}' \
  http://localhost:8080/v1/checkin

Sample Response:

{
  "id": "1",
  "userId": "emp_local_01",
  "username": "Aashish",
  "checkinTime": "2025-09-05 10:23:00 IST",
  "statusMessage": "User checked in successfully."
}

2. Get a Single User's Latest Attendance
Method: GET

Endpoint: /v1/attendance/{user_id}

# Replace 'emp_local_01' with the user ID you want to find
curl -X GET http://localhost:8080/v1/attendance/emp_local_01

Sample Response:

{
  "id": "1",
  "userId": "emp_local_01",
  "username": "Aashish",
  "checkinTime": "2025-09-05 10:23:00 IST",
  "statusMessage": "Record found."
}

3. Check-Out an Employee
Method: PUT

Endpoint: /v1/checkout/{record_id}

# Replace '1' with the actual ID of the record you want to check out
curl -X PUT -H "Content-Type: application/json" \
  -d '{}' \
  http://localhost:8080/v1/checkout/1

Sample Response:

{
  "id": "1",
  "userId": "emp_local_01",
  "username": "Aashish",
  "checkinTime": "2025-09-05 10:23:00 IST",
  "checkoutTime": "2025-09-05 18:30:00 IST",
  "statusMessage": "User checked out successfully."
}

4. Get All Attendance Records
Method: GET

Endpoint: /v1/attendance

curl -X GET http://localhost:8080/v1/attendance

Sample Response:

{
  "records": [
    {
      "id": "1",
      "userId": "emp_local_01",
      "username": "Aashish",
      "checkinTime": "2025-09-05 10:23:00 IST",
      "checkoutTime": "2025-09-05 18:30:00 IST",
      "statusMessage": "Record retrieved."
    },
    {
      "id": "2",
      "userId": "emp_local_02",
      "username": "Priya",
      "checkinTime": "2025-09-05 10:25:15 IST",
      "statusMessage": "Record retrieved."
    }
  ]
}







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
