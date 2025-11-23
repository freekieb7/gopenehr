# ---- Build stage ----
FROM golang:1.25.4-alpine AS builder
WORKDIR /app

# Build args
ARG VERSION=dev
ARG GIT_COMMIT=unknown
ARG BUILD_TIME=unknown

# Copy go.mod and go.sum first (for verification)
COPY go.mod go.sum ./

# Copy vendor directory (contains all dependencies)
COPY vendor/ ./vendor/

# Copy source code
COPY . .

# Verify vendor directory is in sync (fails if vendor is outdated)
RUN go mod verify

# Build static binary using vendored dependencies (no network needed)
RUN CGO_ENABLED=0 GOOS=linux go build \
    -mod=vendor \
    -ldflags="-s -w \
    -X github.com/freekieb7/gopenehr/internal/config.Version=${VERSION} \
    -X github.com/freekieb7/gopenehr/internal/config.GitCommit=${GIT_COMMIT} \
    -X github.com/freekieb7/gopenehr/internal/config.BuildTime=${BUILD_TIME}" \
    -o gopenehr ./cmd/main.go

# ---- Runtime stage ----
# Distroless for Kubernetes (has certificates + minimal attack surface)
FROM gcr.io/distroless/static:nonroot
WORKDIR /app

# Build args
ARG VERSION=dev
ARG GIT_COMMIT=unknown
ARG BUILD_TIME=unknown

# Copy binary
COPY --from=builder /app/gopenehr .

# Labels for metadata
LABEL org.opencontainers.image.version="${VERSION}"
LABEL org.opencontainers.image.revision="${GIT_COMMIT}"
LABEL org.opencontainers.image.created="${BUILD_TIME}"

# Expose
EXPOSE 3000

# Run as non-root by default
USER nonroot:nonroot

# No shell - must use ENTRYPOINT
ENTRYPOINT ["/app/gopenehr"]