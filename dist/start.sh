#!/bin/bash

# Wait for DB to be ready
./wait-for-it.sh --timeout=30 $DB_HOST_PORT

# Apply migrations, if any
./migrate init
./migrate up

# Start server
./server