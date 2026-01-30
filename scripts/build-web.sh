#!/bin/bash
set -e

echo "Building dashboard..."
cd apps/dashboard
npm install
npm run build

echo "Building Go binary..."
cd ../..
go build -o strand ./cmd/strand

echo "Done! Run './strand web' to start"
