#!/bin/bash
set -e

for f in $(find . -name go.mod)
    do (cd $(dirname $f); go mod vendor && go mod tidy)
done