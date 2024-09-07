#!/bin/bash
set -e

# Function to read environment variables from env.json
read_env_json() {
    local env=${GO_ENV:-development}
    if [[ -f /app/config/env.json ]]; then
        eval $(jq -r ".$env | to_entries | map(\"export \(.key)='\(.value|tostring)'\") | .[]" /app/config/env.json)
        echo "Environment variables set from env.json for $env environment"
    else
        echo "env.json file not found"
    fi
}

# Read environment variables
read_env_json

# Check if we're running in the database container
if [ "$(id -u)" = '0' ] && [ -d "/var/lib/postgresql/data" ]; then
    # If running as root in the PostgreSQL container, execute the original entrypoint
    exec docker-entrypoint.sh "$@"
else
    # If running in the app container, start the application
    echo "Starting the application..."
    # Replace the following line with your actual application start command
    exec go run main.go
fi