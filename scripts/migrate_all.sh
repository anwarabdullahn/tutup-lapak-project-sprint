#!/usr/bin/env bash
set -euo pipefail

# Usage: ./scripts/migrate_all.sh [DATABASE_URL]
# If DATABASE_URL not passed, tries per-service .env files.

ROOT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
DB_URL="${1:-}"

services=(auth-service profile-service product-service purchase-service)

log(){ printf "[%s] %s\n" "$1" "$2"; }
fail(){ log ERROR "$1"; exit 1; }

for svc in "${services[@]}"; do
  pushd "$ROOT_DIR/$svc" >/dev/null
  if [ -z "$DB_URL" ]; then
    if [ -f .env ]; then
      DB_URL_LOCAL=$(grep -E '^DATABASE_URL=' .env | cut -d '=' -f2- || true)
    else
      DB_URL_LOCAL=""
    fi
  else
    DB_URL_LOCAL="$DB_URL"
  fi
  if [ -z "$DB_URL_LOCAL" ]; then
    log WARN "Skipping $svc (no DATABASE_URL)"
    popd >/dev/null
    continue
  fi
  if ! command -v migrate >/dev/null 2>&1; then
    fail "migrate CLI not found in PATH. Install with: go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest"
  fi
  log INFO "Migrating $svc"
  migrate -path migrations/db -database "$DB_URL_LOCAL" up
  popd >/dev/null
done

log INFO "All done"
