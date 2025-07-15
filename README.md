# wedl - WeTransfer API Server

[![Test latest release](https://github.com/gnojus/wedl/actions/workflows/test.yml/badge.svg)](https://github.com/gnojus/wedl/actions/workflows/test.yml)

HTTP API server for downloading files from WeTransfer. Returns files as base64-encoded JSON responses.

Uses unofficial WeTransfer API used when downloading with a browser. Written in Go.

## Quick Start with Docker

### 1. Start the API Server

```bash
# Clone the repository
git clone https://github.com/gnojus/wedl.git
cd wedl

# Start the API server
docker-compose up -d
```

The API server will be available at `http://localhost:8080`

### 2. Test the API

```bash
# Check if the server is running
curl http://localhost:8080/health

# Expected response:
# {"status":"ok","timestamp":"2025-07-15T03:43:19Z"}
```

## API Usage

### Download Files

**Endpoint:** `POST /wetransfer`

**Request:**
```bash
curl --location 'http://localhost:8080/wetransfer' \
--header 'Content-Type: application/json' \
--data '{
    "wetransfer_url": "https://we.tl/your-wetransfer-url"
}'
```

**Response:**
```json
{
    "fileName": "document.pdf",
    "fileBinary": "JVBERi0xLjQKMSAwIG9iago8PAovVHlwZS..."
}
```

- `fileName`: The original filename from WeTransfer
- `fileBinary`: Base64-encoded file content

### Password-Protected Files

For password-protected transfers, include the password:

```bash
curl --location 'http://localhost:8080/wetransfer' \
--header 'Content-Type: application/json' \
--data '{
    "wetransfer_url": "https://we.tl/your-protected-url",
    "password": "your-password"
}'
```

### Get File Information

**Endpoint:** `POST /info`

Get file metadata without downloading:

```bash
curl --location 'http://localhost:8080/info' \
--header 'Content-Type: application/json' \
--data '{
    "wetransfer_url": "https://we.tl/your-wetransfer-url"
}'
```

**Response:**
```json
{
    "success": true,
    "filename": "document.pdf",
    "size": 1234567,
    "dl_url": "https://..."
}
```

## Docker Configuration

### Using docker-compose (Recommended)

```yaml
services:
  wedl-api:
    build: .
    container_name: wedl-api
    ports:
      - "8080:8080"
    environment:
      - TZ=UTC
    restart: unless-stopped
```

### Using Docker directly

```bash
# Build the image
docker build -t wedl .

# Run the container
docker run -d -p 8080:8080 --name wedl-api wedl
```

### Custom Port

To run on a different port:

```bash
# Using docker-compose (modify docker-compose.yml ports)
services:
  wedl-api:
    ports:
      - "3000:8080"  # External:Internal

# Using Docker directly
docker run -d -p 3000:8080 --name wedl-api wedl
```

## Development

### Local Development

```bash
# Install dependencies
go mod download

# Run API server locally
go run api-server.go -port 8080

# Or build and run
go build -o api-server api-server.go
./api-server -port 8080
```

### CLI Tool (Legacy)

The project also includes a CLI tool:

```bash
# Build CLI tool
go build -o wedl wedl.go

# Use CLI tool
./wedl https://we.tl/example-url
```

## Error Handling

The API returns appropriate HTTP status codes:

- `200 OK`: Success
- `400 Bad Request`: Invalid URL or request format
- `405 Method Not Allowed`: Wrong HTTP method
- `500 Internal Server Error`: Server error

Error responses include details:
```json
{
    "success": false,
    "error": "Failed to get download response: invalid URL"
}
```

## Notes

- Files are loaded entirely into memory before base64 encoding
- Large files may cause memory issues
- WeTransfer URLs typically expire after a certain time
- This uses the unofficial WeTransfer API and may break if they change their implementation