# --- Build Stage ---
    FROM alpine:latest AS builder

    # Set the working directory inside the container
    WORKDIR /app
    
    # Copy go.mod and go.sum first to leverage Docker layer caching
    COPY go.mod go.sum ./
    RUN go mod download
    
    # Copy the rest of the source code
    COPY . .
    
    # Build the Go app (static binary for portability)
    RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o attendance-service .
    
    # --- Final Stage ---
    FROM alpine:latest
    
    # Install CA certificates (needed for MongoDB TLS connections, HTTP calls, etc.)
    RUN apk --no-cache add ca-certificates
    
    WORKDIR /app
    
    # Copy the compiled binary from builder
    COPY --from=builder /app/attendance-service .
    
    # Expose ports: 50051 for gRPC, 8080 for HTTP-Gateway
    EXPOSE 50051
    EXPOSE 8080
    
    # Run the service
    CMD ["./attendance-service"]
    


# # --- Build Stage ---
# # Use the official Go image as a builder
# FROM golang:1.25-alpine AS builder

# # Set the working directory inside the container
# WORKDIR /app

# # Copy the go.mod and go.sum files to download dependencies
# COPY go.mod go.sum ./
# RUN go mod download

# # Copy the rest of the source code
# COPY . .
# # Build the Go app
# RUN go build -o attendance-service .

# # Build the application, creating a statically linked binary
# # CGO_ENABLED=0 is important for creating a truly portable binary
# # -o /app/server creates the compiled program named 'server'
# RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/server .

# # --- Final Stage ---
# # Use a minimal base image for the final container
# FROM alpine:latest

# # Set the working directory
# WORKDIR /app

# # Copy only the compiled binary from the builder stage
# COPY --from=builder /app/attendance-service .

# # Expose the ports the application will use
# # 50051 for gRPC and 8080 for the HTTP gateway
# EXPOSE 50051
# EXPOSE 8080

# # Command to run the application when the container starts
# CMD ["./attendance-service"]