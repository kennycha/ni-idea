# Build stage
FROM golang:1.23-alpine AS builder

# Allow Go to download newer toolchain if needed
ENV GOTOOLCHAIN=auto

WORKDIR /app

# Install git for go mod download (some deps need it)
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o ni-server ./server

# Runtime stage
FROM alpine:3.19

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/ni-server .

# Create data directory
RUN mkdir -p /data

EXPOSE 8080

ENTRYPOINT ["./ni-server"]
CMD ["--port", "8080", "--data", "/data"]
