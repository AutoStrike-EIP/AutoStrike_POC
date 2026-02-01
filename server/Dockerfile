# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /build

# Install dependencies
RUN apk add --no-cache gcc musl-dev sqlite-dev

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=1 GOOS=linux go build -a -ldflags '-linkmode external -extldflags "-static"' -o autostrike ./cmd/autostrike

# Runtime stage
FROM alpine:3.19

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Copy binary from builder
COPY --from=builder /build/autostrike .

# Create directories
RUN mkdir -p /app/config /app/data /app/certs

# Copy default config
COPY config/config.yaml /app/config/

# Create non-root user
RUN adduser -D -u 1000 autostrike
RUN chown -R autostrike:autostrike /app
USER autostrike

EXPOSE 8443

ENTRYPOINT ["./autostrike"]
