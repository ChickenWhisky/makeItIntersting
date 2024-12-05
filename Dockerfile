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

    COPY . .
    RUN go build -o /app/bin/myapp ./cmd/app

    # ---- Deploy Stage ----

    FROM alpine:3.18

    WORKDIR /app
    COPY --from=builder /app/bin/myapp /app/myapp
    EXPOSE 8080
    ENTRYPOINT ["/app/myapp"]
    