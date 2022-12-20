#!/bin/bash
set -e

find . -name go.mod -execdir go test ./... \;