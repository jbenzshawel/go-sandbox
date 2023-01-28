#!/bin/bash

while true
 do
    if [[ $(curl -sL -w "%{http_code}\\n" "keycloak:8080" -o /dev/null) == "200" ]]; then
        break
    else
        echo "Waiting for keycloak availability..."
        sleep 0.5
    fi
done

go run .