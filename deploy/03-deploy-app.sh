#!/usr/bin/env bash
#
# HairHistory: build and release the app (safe to re-run).
#
#   bash deploy/03-deploy-app.sh [branch]      # branch defaults to main
#
# Requires /etc/hairhistory/api.env to exist (see deploy/api.env.example).

set -euo pipefail

BRANCH="${1:-main}"
REPO_URL="${REPO_URL:-https://github.com/2410as/HairHistory.git}"
APP_DIR="${APP_DIR:-/opt/hairhistory}"
WEB_ROOT="${WEB_ROOT:-/var/www/hairhistory}"
ENV_FILE="${ENV_FILE:-/etc/hairhistory/api.env}"
SERVICE="hairhistory-api"

export PATH="/usr/local/go/bin:${HOME}/go/bin:${PATH}"

log() { printf '\n==> %s\n' "$*"; }

for cmd in go node npm git psql; do
  command -v "${cmd}" >/dev/null 2>&1 || {
    echo "ERROR: ${cmd} not found. Run deploy/01-setup-server.sh first." >&2
    exit 1
  }
done

if [[ ! -f "${ENV_FILE}" ]]; then
  echo "ERROR: ${ENV_FILE} is missing. Copy deploy/api.env.example there, fill it in, and chmod 600 it." >&2
  exit 1
fi

# ---------------------------------------------------------------- source code
log "Syncing ${REPO_URL} (${BRANCH}) into ${APP_DIR}"
if [[ ! -d "${APP_DIR}/.git" ]]; then
  sudo install -d -o "$(id -un)" -g "$(id -gn)" "${APP_DIR}"
  git clone "${REPO_URL}" "${APP_DIR}"
fi

git -C "${APP_DIR}" fetch --prune origin
# Discarding local state keeps the deployed tree an exact mirror of the branch;
# never edit files under ${APP_DIR} directly.
git -C "${APP_DIR}" checkout -B "${BRANCH}" "origin/${BRANCH}"
git -C "${APP_DIR}" reset --hard "origin/${BRANCH}"
git -C "${APP_DIR}" rev-parse --short HEAD

# ------------------------------------------------------------------ variables
# The build needs DATABASE_URL and GOOGLE_CLIENT_ID; read them from the same
# file systemd uses so the values are defined in exactly one place.
#
# ${ENV_FILE} holds the DB password, so it stays root-owned and chmod 600 - this
# script runs as ubuntu and therefore reads it through sudo instead of sourcing
# it directly. systemd reads the same file as root, so the permissions must not
# be relaxed. The contents travel through a variable and a here-string, never
# through a command line argument, so the password cannot show up in ps; xtrace
# is suspended around the read so it cannot show up in a `bash -x` log either.
xtrace_was_on=0
if [[ $- == *x* ]]; then
  xtrace_was_on=1
  set +x
fi

env_contents=""
if ! env_contents="$(sudo cat "${ENV_FILE}")"; then
  cat >&2 <<EOF
ERROR: ${ENV_FILE} が読めません。sudo 権限が必要です。
       Could not read ${ENV_FILE} (root-owned, mode 600).
       Run this script as a user with sudo, for example:
         sudo -v && bash deploy/03-deploy-app.sh ${BRANCH}
       Do NOT chmod the file to make it world-readable; it contains the DB password.
EOF
  exit 1
fi

set -a
# shellcheck disable=SC1090,SC1091
source /dev/stdin <<<"${env_contents}"
set +a

unset env_contents

: "${DATABASE_URL:?DATABASE_URL is not set in ${ENV_FILE}}"
: "${GOOGLE_CLIENT_ID:?GOOGLE_CLIENT_ID is not set in ${ENV_FILE}}"

if [[ "${xtrace_was_on}" -eq 1 ]]; then
  set -x
fi
unset xtrace_was_on

# --------------------------------------------------------------- backend build
log "Building the API binary"
mkdir -p "${APP_DIR}/bin"
BIN="${APP_DIR}/bin/api"

if [[ -f "${BIN}" ]]; then
  cp -p "${BIN}" "${BIN}.prev"
  echo "Previous binary kept at ${BIN}.prev (used by the rollback procedure)."
fi

(cd "${APP_DIR}/backend" && go build -o "${BIN}.new" ./cmd/api)
mv "${BIN}.new" "${BIN}"

# --------------------------------------------------------------- migrations
# golang-migrate is used rather than looping psql -f because it records the
# applied version in schema_migrations, which makes re-running this script a no-op.
log "Applying database migrations"
if ! command -v migrate >/dev/null 2>&1; then
  log "Installing golang-migrate"
  go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
  hash -r
fi

if ! command -v migrate >/dev/null 2>&1; then
  echo "ERROR: migrate が PATH にありません（go install 後も見つかりません）。" >&2
  echo "       Expected it at $(go env GOPATH)/bin/migrate - add that directory to PATH." >&2
  exit 1
fi

migrate -path "${APP_DIR}/backend/migrations" -database "${DATABASE_URL}" up

# -------------------------------------------------------------- frontend build
log "Building the frontend"
# Vite inlines these at build time, so they must be present before npm run build.
# An empty VITE_API_BASE_URL makes the SPA call /api on its own origin, which is
# what nginx proxies - same-origin means the session cookie is sent without CORS.
cat > "${APP_DIR}/frontend/.env.production" <<EOF
VITE_API_BASE_URL=
VITE_GOOGLE_CLIENT_ID=${GOOGLE_CLIENT_ID}
EOF

(cd "${APP_DIR}/frontend" && npm ci && npm run build)

log "Publishing static files to ${WEB_ROOT}"
sudo install -d -o www-data -g www-data "${WEB_ROOT}"
# --delete keeps hashed assets from older releases from piling up.
sudo rsync -a --delete "${APP_DIR}/frontend/dist/" "${WEB_ROOT}/"
sudo chown -R www-data:www-data "${WEB_ROOT}"

# -------------------------------------------------------------------- restart
log "Restarting ${SERVICE}"
sudo systemctl restart "${SERVICE}"
sleep 2
sudo systemctl --no-pager --full status "${SERVICE}" || true

log "Health check"
curl -fsS "http://127.0.0.1:${PORT:-8080}/healthz" && echo

log "Deployed ${BRANCH} at $(git -C "${APP_DIR}" rev-parse --short HEAD)"
