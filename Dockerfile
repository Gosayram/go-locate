# Build stage v1.24.4-alpine3.22
FROM golang:1.24.4-alpine3.22@sha256:68932fa6d4d4059845c8f40ad7e654e626f3ebd3706eef7846f319293ab5cb7a AS builder

# Install git and ca-certificates (needed for go mod download)
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build arguments for version information
ARG VERSION=docker
ARG COMMIT=unknown
ARG DATE
ARG BUILT_BY=docker

# Set build date if not provided
RUN if [ -z "$DATE" ]; then DATE=$(date -u '+%Y-%m-%d_%H:%M:%S'); fi

# Build the binary with comprehensive version information
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w -X 'github.com/Gosayram/go-locate/internal/version.Version=${VERSION}' \
              -X 'github.com/Gosayram/go-locate/internal/version.Commit=${COMMIT}' \
              -X 'github.com/Gosayram/go-locate/internal/version.Date=${DATE}' \
              -X 'github.com/Gosayram/go-locate/internal/version.BuiltBy=${BUILT_BY}'" \
    -o glocate ./cmd/glocate

# Final stage
FROM scratch

# Copy ca-certificates from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy the binary
COPY --from=builder /app/glocate /usr/local/bin/glocate

# Set the entrypoint
ENTRYPOINT ["/usr/local/bin/glocate"]

# Default command
CMD ["--help"]

# Labels
LABEL org.opencontainers.image.title="glocate"
LABEL org.opencontainers.image.description="Modern file search tool - replacement for locate command"
LABEL org.opencontainers.image.source="https://github.com/Gosayram/go-locate"
LABEL org.opencontainers.image.licenses="MIT"
