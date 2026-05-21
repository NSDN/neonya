#!/bin/sh
set -e

echo "=== Starting Podman database ==="
podman compose -f apps/server/podman-compose.yml up --detach

echo "=== Migrating database ==="
echo "Run: migrate -source file://apps/server/database/migrations -database \"postgres://forum_user:password@localhost:5432/forum?sslmode=disable\" up"

echo "=== Starting backend ==="
(cd apps/server && go run ./cmd/server) &

echo "=== Starting frontend ==="
pnpm dev
