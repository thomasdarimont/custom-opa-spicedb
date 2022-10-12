#!/usr/bin/env sh

echo ENDPOINT="$ENDPOINT"
echo TOKEN="$TOKEN"

echo "Wait for Authzed..."
sleep 3

echo "Setup authzed context..."
zed context set first-dev-context "$ENDPOINT" "$TOKEN" --insecure

echo "Importing schema and data..."
zed import --schema=true file:///opt/schema-and-data.yml --insecure

echo "Done."