#!/bin/bash
set -e

cd src/identity

jet -dsn=postgresql://postgres:postgres@localhost:5432/identity?sslmode=disable -schema=identity -ignore-tables=flyway_schema_history -path=./infrastructure/gen