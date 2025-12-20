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
    -trimpath \
    -ldflags="-s -w -extldflags=-static \
    -X github.com/freekieb7/gopenehr/internal/config.Version=${VERSION} \
    -X github.com/freekieb7/gopenehr/internal/config.GitCommit=${GIT_COMMIT} \
    -X github.com/freekieb7/gopenehr/internal/config.BuildTime=${BUILD_TIME}" \
    -o gopenehr ./cmd/main.go

# Verify binary is static
RUN file gopenehr && ldd gopenehr || true

# ---- Runtime stage ----
# Distroless static (Debian 12) - NO SHELL, NO PACKAGE MANAGER, NO ANYTHING
FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app

# Build args for labels
ARG VERSION=dev
ARG GIT_COMMIT=unknown
ARG BUILD_TIME=unknown

# Copy only what's needed
COPY --from=builder /app/gopenehr .

# Security labels
LABEL org.opencontainers.image.version="${VERSION}" \
      org.opencontainers.image.revision="${GIT_COMMIT}" \
      org.opencontainers.image.created="${BUILD_TIME}" \
      org.opencontainers.image.title="gopenehr" \
      org.opencontainers.image.description="OpenEHR server - distroless, no shell" \
      org.opencontainers.image.vendor="freekieb7" \
      org.opencontainers.image.licenses="MIT"

# Run as non-root user (UID 65532)
USER nonroot:nonroot

# No shell available - ENTRYPOINT only
ENTRYPOINT ["/app/gopenehr"]