#!/bin/bash

docker network create promotion-bridge-network

# Start ekyc service
echo "Starting ekyc service..."
cd ekyc
docker-compose up -d
cd ..

# Start voucher service
echo "Starting voucher service..."
cd voucher
docker-compose up -d
cd ..

# Start promotion service
echo "Starting promotion service..."
cd promotion
docker-compose up -d
cd ..

echo "All services started successfully."
