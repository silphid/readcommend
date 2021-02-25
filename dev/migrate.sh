#!/bin/bash

export DB_URL="postgres://postgres:password123@localhost:5432/readcommend?sslmode=disable"

cd ../src/migrations
go run main.go init
go run main.go up
