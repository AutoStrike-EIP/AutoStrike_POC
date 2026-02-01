# Build stage
FROM rust:1.75-alpine AS builder

WORKDIR /build

# Install dependencies
RUN apk add --no-cache musl-dev openssl-dev openssl-libs-static

# Copy Cargo files
COPY Cargo.toml Cargo.lock* ./

# Create dummy main.rs to cache dependencies
RUN mkdir src && echo "fn main() {}" > src/main.rs
RUN cargo build --release
RUN rm -rf src

# Copy actual source code
COPY src ./src

# Build the binary
RUN touch src/main.rs && cargo build --release

# Runtime stage - minimal for agent deployment
FROM alpine:3.19

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates

# Copy binary from builder
COPY --from=builder /build/target/release/autostrike-agent .

# Create non-root user
RUN adduser -D -u 1000 agent
USER agent

ENTRYPOINT ["./autostrike-agent"]
