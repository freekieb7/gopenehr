# Multi-stage build for AQL Benchmark CLI
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o aql-bench ./cmd/bench/

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/aql-bench .
COPY --from=builder /app/cmd/bench/*.json ./suites/

# Make it easy to run
ENTRYPOINT ["./aql-bench"]
CMD ["--help"]