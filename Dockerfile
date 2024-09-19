# ---- Build Stage ----
    FROM golang:1.23-alpine AS builder

    # Set environment variables for Go build
    ENV GO111MODULE=on
    ENV CGO_ENABLED=0
    ENV GOOS=linux
    ENV GOARCH=amd64
    
    # Create app directory
    WORKDIR /app
    
    # Cache dependencies first by copying go.mod and go.sum separately
    COPY go.mod go.sum ./
    RUN go mod download
    
    # Copy the rest of the application source code
    COPY . .
    
    # Build the Go binary
    RUN go build -o /app/bin/myapp ./cmd/app
    
    # ---- Deploy Stage ----
    FROM alpine:3.18
    
    # Create a directory for the application binary
    WORKDIR /app
    
    # Copy the built binary from the build stage
    COPY --from=builder /app/bin/myapp /app/myapp
    
    # Expose application port 
    EXPOSE 8080
    
    # Set entrypoint for the container to run the binary
    ENTRYPOINT ["/app/myapp"]
    