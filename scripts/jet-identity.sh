#!/bin/bash
set -e

cd src/identity

jet -dsn=postgresql://postgres:postgres@localhost:5432/identity?sslmode=disable -schema=identity -path=./infrastructure/gen