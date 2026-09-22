#!/usr/bin/env bash
#
# HairHistory: create the PostgreSQL role and database (safe to re-run).
#
# The password is NEVER stored in this script. Supply it via the environment:
#
#   DB_PASSWORD='<strong password>' bash deploy/02-setup-db.sh
#
# Use the same password in DATABASE_URL inside /etc/hairhistory/api.env.
# Nothing here opens PostgreSQL to the network: the default pg_hba.conf
# (local socket + 127.0.0.1 only) is left untouched on purpose, and ufw keeps
# 5432 closed. The API reaches the database over localhost.

set -euo pipefail

DB_NAME="${DB_NAME:-hairhistory}"
DB_USER="${DB_USER:-hairhistory}"

log() { printf '\n==> %s\n' "$*"; }

if [[ -z "${DB_PASSWORD:-}" ]]; then
  cat >&2 <<'EOF'
ERROR: DB_PASSWORD is not set.

  DB_PASSWORD='<strong password>' bash deploy/02-setup-db.sh

Generate one with: openssl rand -base64 24
EOF
  exit 1
fi

if ! command -v psql >/dev/null 2>&1; then
  echo "ERROR: psql not found. Run deploy/01-setup-server.sh first." >&2
  exit 1
fi

psql_as_postgres() { sudo -u postgres psql -v ON_ERROR_STOP=1 "$@"; }

role_exists="$(psql_as_postgres -tAc \
  "SELECT 1 FROM pg_roles WHERE rolname = '${DB_USER}'")"

if [[ "${role_exists}" == "1" ]]; then
  log "Role ${DB_USER} exists; updating its password only"
else
  log "Creating role ${DB_USER}"
  psql_as_postgres -c "CREATE ROLE ${DB_USER} LOGIN"
fi

# ALTER ROLE ... PASSWORD is passed as a bound parameter so the password never
# reaches the shell history or the postgres log line.
psql_as_postgres -c "ALTER ROLE ${DB_USER} WITH LOGIN PASSWORD :'pw'" \
  -v "pw=${DB_PASSWORD}"

db_exists="$(psql_as_postgres -tAc \
  "SELECT 1 FROM pg_database WHERE datname = '${DB_NAME}'")"

if [[ "${db_exists}" == "1" ]]; then
  log "Database ${DB_NAME} already exists; leaving it untouched"
else
  log "Creating database ${DB_NAME} owned by ${DB_USER}"
  psql_as_postgres -c "CREATE DATABASE ${DB_NAME} OWNER ${DB_USER}"
fi

log "Verifying a localhost login works"
PGPASSWORD="${DB_PASSWORD}" psql -h 127.0.0.1 -U "${DB_USER}" -d "${DB_NAME}" \
  -tAc 'SELECT current_database()'

log "Done"
cat <<EOF
DATABASE_URL for /etc/hairhistory/api.env (substitute your password):

  DATABASE_URL=postgres://${DB_USER}:<password>@127.0.0.1:5432/${DB_NAME}?sslmode=disable

sslmode=disable is acceptable here because the connection never leaves the host.
URL-encode the password if it contains @ : / ? # or %.
EOF
