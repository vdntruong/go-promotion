#!/bin/bash

# Clean promotion service
echo "Cleaning promotion service..."
cd promotion
docker-compose down -v
cd ..

# Clean voucher service
echo "Cleaning voucher service..."
cd voucher
docker-compose down -v
cd ..

# Clean ekyc service
echo "Cleaning ekyc service..."
cd ekyc
docker-compose down -v
cd ..

echo "All services Cleaned successfully."
