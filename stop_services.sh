#!/bin/bash

# Stop promotion service
echo "Stopping promotion service..."
cd promotion
docker-compose down
cd ..

# Stop voucher service
echo "Stopping voucher service..."
cd voucher
docker-compose down
cd ..

# Stop ekyc service
echo "Stopping ekyc service..."
cd ekyc
docker-compose down
cd ..

echo "All services stopped successfully."
