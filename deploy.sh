#!/bin/bash

# WeTransfer Download API Deployment Script
set -e

echo "🚀 Deploying WeTransfer Download API..."

# Check if Docker and Docker Compose are installed
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed. Please install Docker first."
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose is not installed. Please install Docker Compose first."
    exit 1
fi

# Create .env file if it doesn't exist
if [ ! -f .env ]; then
    echo "📋 Creating .env file from template..."
    cp .env.example .env
    echo "✅ .env file created. You can edit it to customize your configuration."
fi

# Build and start the services
echo "🔨 Building Docker image..."
docker-compose build

echo "🏃 Starting services..."
docker-compose up -d

# Wait for the service to be healthy
echo "⏳ Waiting for service to be ready..."
timeout 60 bash -c 'until docker-compose exec wedl-api wget --quiet --tries=1 --spider http://localhost:8080/health 2>/dev/null; do sleep 2; done'

echo "✅ WeTransfer Download API is now running!"
echo "🌐 API URL: http://localhost:$(grep PORT .env | cut -d'=' -f2 || echo 8080)"
echo "🏥 Health check: http://localhost:$(grep PORT .env | cut -d'=' -f2 || echo 8080)/health"
echo ""
echo "📖 Usage:"
echo "  curl -X POST http://localhost:8080/info -H 'Content-Type: application/json' -d '{\"url\": \"https://we.tl/example\"}'"
echo "  curl -X POST http://localhost:8080/download -H 'Content-Type: application/json' -d '{\"url\": \"https://we.tl/example\"}' -o file.pdf"
echo ""
echo "🛑 To stop: docker-compose down"
echo "📊 To view logs: docker-compose logs -f"