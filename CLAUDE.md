# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`wedl` is a command-line utility for downloading files from WeTransfer. It's written in Go and uses the unofficial WeTransfer API to facilitate downloads directly from the command line.

## Architecture

The codebase follows a clean, modular structure:

- **Main entry point**: `wedl.go` - Handles CLI argument parsing using docopt and delegates to the CLI package
- **CLI package** (`cli/`): 
  - `cli.go` - Core application logic, orchestrates the download process
  - `args.go` - Command-line argument parsing and validation
- **Transfer package** (`transfer/`):
  - `download.go` - WeTransfer API interaction, handles authentication and download URL retrieval
  - `write.go` - File writing operations with progress tracking

The application flow:
1. Parse CLI arguments → 2. Get download response from WeTransfer API → 3. Write file with optional progress bar

## Common Development Commands

### Building
```bash
# Build the application
go build

# Build with specific output name
go build -o wedl.exe wedl.go
```

### Running
```bash
# Run directly with go
go run . --help
go run . https://we.tl/example-url

# Run built binary
./wedl --help
```

### Testing
```bash
# Run all tests
go test ./...

# Run tests in current directory
go test

# Run tests with verbose output
go test -v ./...
```

### Development
```bash
# Format code
go fmt ./...

# Vet code for issues
go vet ./...

# Install dependencies
go mod download

# Update dependencies
go mod tidy
```

## Dependencies

The project uses minimal dependencies:
- `github.com/docopt/docopt-go` - CLI argument parsing
- `github.com/cheggaaa/pb` - Progress bar for downloads

## API Server

The application can run as an HTTP API server that streams downloaded files directly in the response.

### API Endpoints

**GET /health**
- Returns server health status

**POST /info**
- Get file metadata without downloading
- Request: `{"url": "https://we.tl/example", "password": "optional"}`
- Response: `{"success": true, "filename": "file.pdf", "size": 1234, "dl_url": "..."}`

**POST /download**
- Download and stream file directly
- Request: `{"url": "https://we.tl/example", "password": "optional"}`
- Response: File stream with proper headers

### Running API Server

```bash
# Run locally
go run api-server.go -port 8080

# Build and run
go build -o api-server api-server.go
./api-server -port 8080
```

### Usage Examples

```bash
# Get file info
curl -X POST http://localhost:8080/info \
  -H "Content-Type: application/json" \
  -d '{"url": "https://we.tl/example"}'

# Download file
curl -X POST http://localhost:8080/download \
  -H "Content-Type: application/json" \
  -d '{"url": "https://we.tl/example"}' \
  -o downloaded_file.pdf
```

## Docker Deployment

### Building and Running with Docker
```bash
# Build Docker image (includes API server)
docker build -t wedl .

# Run API server
docker run --rm -p 8080:8080 wedl

# Run with docker-compose
docker-compose up -d

# Test API
curl http://localhost:8080/health
```

### Docker Configuration
- Builds both CLI tool and API server
- API server runs on port 8080 by default
- Includes ca-certificates for HTTPS requests
- No volume mounting needed (files streamed directly)

## Key Implementation Details

- Uses regex parsing to extract transfer IDs and security hashes from WeTransfer URLs
- Implements a two-step download process: first gets download metadata, then streams the actual file
- Supports password-protected transfers via the `-P` flag
- Progress tracking is optional and can be disabled with `--silent` flag
- Output can be directed to stdout using `-o -` or to specific files/directories