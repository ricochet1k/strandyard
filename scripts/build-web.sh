#!/bin/bash
set -e

echo "Building dashboard..."
cd apps/dashboard
npm install
npm run build

echo "Building Go binary..."
cd ../..
go build -o strandyard main.go

echo "Done! Run './strandyard web' to start"
