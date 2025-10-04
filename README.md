# GopenEHR

A Go implementation of OpenEHR services with AQL (Archetype Query Language) support.

## Features

- **AQL Query Engine**: Full AQL parsing and execution support
- **REST API**: OpenEHR REST API endpoints
- **PostgreSQL Backend**: Robust data storage with PostgreSQL
- **Performance Benchmarking**: Comprehensive AQL query benchmarking tools
- **Protocol Buffers**: gRPC service definitions
- **Docker Support**: Easy deployment with Docker Compose

## Quick Start

### Prerequisites

- Go 1.24+
- PostgreSQL 12+
- Docker & Docker Compose (optional)

### Development Setup

1. Clone the repository:
```bash
git clone https://github.com/freekieb7/gopenehr.git
cd gopenehr
```

2. Start the database:
```bash
make up
```

3. Run the server:
```bash
make run
```

The API will be available at `http://localhost:8080`

### Docker Deployment

```bash
docker-compose up -d
```

## API Endpoints

### AQL Query Execution

```bash
POST /openehr/v1/query
Content-Type: application/json

{
  "aql": "SELECT e FROM EHR e LIMIT 10",
  "parameters": {}
}
```

### Health Check

```bash
GET /openehr/v1/system/health
```

## AQL Benchmarking

The project includes a comprehensive benchmarking CLI tool for performance testing AQL queries.

### Quick Benchmark

```bash
# Build the benchmark tool
make bench

# Run with default test suite
./aql-bench --url http://localhost:8080 --verbose

# Run a single query
./aql-bench --url http://localhost:8080 --aql "SELECT e FROM EHR e LIMIT 5"
```

### Production Benchmarking

```bash
# Run production test suite
make bench-production

# Run quick health checks only
make bench-quick

# Generate JUnit XML for CI/CD
make bench-junit
```

### Custom Test Suites

Create custom JSON test suites:

```json
{
  "name": "My Tests",
  "description": "Custom AQL tests",
  "tests": [
    {
      "name": "EHR Count",
      "aql": "SELECT COUNT(e) FROM EHR e",
      "expected_rows": 1,
      "tags": ["quick", "count"]
    }
  ]
}
```

Run with:
```bash
./aql-bench --suite my-tests.json --format junit --output results.xml
```

### Azure DevOps Integration

The benchmark tool outputs JUnit XML compatible with Azure Pipelines:

```yaml
- script: |
    go build -o aql-bench ./cmd/bench/
    ./aql-bench --url $(API_URL) --format junit --output test-results.xml
  displayName: 'Run AQL Benchmarks'

- task: PublishTestResults@2
  inputs:
    testResultsFormat: 'JUnit'
    testResultsFiles: 'test-results.xml'
    testRunTitle: 'AQL Performance Tests'
```

See [`cmd/bench/azure-pipelines.yml`](cmd/bench/azure-pipelines.yml) for a complete example.

## Project Structure

```
├── aql/                    # AQL parser and execution engine
│   ├── AQL.g4             # ANTLR grammar for AQL
│   ├── gen/               # Generated parser code
│   └── *.go               # AQL implementation
├── cmd/
│   ├── bench/             # Benchmark CLI tool
│   │   ├── bench.go       # Main benchmark implementation
│   │   ├── *.json         # Test suite examples
│   │   ├── Dockerfile     # Docker container for benchmarks
│   │   └── README.md      # Detailed benchmark documentation
│   └── main.go            # Server main entry point
├── rest/                  # REST API implementation
├── database/              # Database layer
├── model/                 # Data models
├── encoding/              # JSON encoding/decoding
└── proto/                 # Protocol buffer definitions
```

## Available Make Targets

```bash
make run                   # Run the development server
make up                    # Start PostgreSQL with Docker
make down                  # Stop Docker services
make air                   # Run with hot reload
make aql-gen              # Generate AQL parser from grammar

# Benchmark targets
make bench                # Build benchmark tool
make bench-simple         # Run simple test suite
make bench-production     # Run production test suite  
make bench-quick          # Run quick health check tests
make bench-junit          # Generate JUnit XML output
```

## Development

### AQL Grammar Changes

When modifying the AQL grammar (`aql/AQL.g4`), regenerate the parser:

```bash
make aql-gen
```

### Testing

Run the test suite:

```bash
go test ./...
```

Run benchmarks:

```bash
go test -bench=. ./...
```

### Hot Reload Development

Use Air for automatic reloading during development:

```bash
make air
```

## Production Deployment

### Environment Variables

- `DATABASE_URL`: PostgreSQL connection string
- `PORT`: Server port (default: 8080)
- `LOG_LEVEL`: Logging level (info, debug, warn, error)

### Docker Production Build

```bash
docker build -t gopenehr .
docker run -p 8080:8080 -e DATABASE_URL="postgres://..." gopenehr
```

## Performance Monitoring

The benchmark tool can be integrated into CI/CD pipelines for continuous performance monitoring:

1. **Quick Health Checks**: Fast queries to verify API availability
2. **Load Testing**: Medium-sized result sets for normal load simulation
3. **Stress Testing**: Large queries to test system limits
4. **Regression Detection**: Compare performance against baselines

Results are output in multiple formats:
- **JSON**: Detailed metrics for analysis
- **JUnit XML**: CI/CD integration
- **Console**: Human-readable summaries

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make changes and add tests
4. Run benchmarks to verify performance
5. Submit a pull request

## License

[Add your license information here]