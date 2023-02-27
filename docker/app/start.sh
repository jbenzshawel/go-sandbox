#!/bin/bash

while [[ $(curl -sL -w "%{http_code}\\n" "keycloak:8080" -o /dev/null) != "200" ]]
 do
    echo "Waiting for keycloak availability..."
    sleep 0.5
done

go run .