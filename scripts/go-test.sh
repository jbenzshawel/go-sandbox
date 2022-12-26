#!/bin/bash
set -e

if [ "$(docker ps -a | grep postgres)" ]; then
    echo "Running tests with postgres and in memory db"
    IDENTITY_POSTGRES='host=localhost port=5432 user=identity password=dockerdbpw dbname=identity sslmode=disable' bash -c 'find . -name go.mod -execdir go test ./... \;'
else 
    echo "Running tests with in memory db"
    find . -name go.mod -execdir go test ./... \;
fi
