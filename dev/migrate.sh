#!/bin/bash

. vars.sh
cd ../src/migrations
go run main.go init
go run main.go up
