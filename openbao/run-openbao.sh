#!/bin/bash

# OpenBao Docker Runner Script
# This script helps you run OpenBao in development mode using Docker

set -e

echo "🔐 Starting OpenBao in development mode..."

# Create config directory if it doesn't exist
mkdir -p config

# Check if Docker is running
if ! docker info >/dev/null 2>&1; then
	echo "❌ Docker is not running. Please start Docker first."
	exit 1
fi

# Start OpenBao using docker-compose
docker compose up -d

# Wait for OpenBao to be ready
echo "⏳ Waiting for OpenBao to start..."
sleep 5

# Check if OpenBao is accessible
if curl -s http://localhost:8200/v1/sys/health >/dev/null 2>&1; then
	echo "✅ OpenBao is running!"
	echo ""
	echo "🌐 Web UI: http://localhost:8200"
	echo "🔑 Root Token: myroot"
	echo "📋 API Endpoint: http://localhost:8200"
	echo ""
	echo "To stop OpenBao, run: docker compose down"
else
	echo "❌ OpenBao failed to start. Check logs with: docker compose logs"
	exit 1
fi
