#!/usr/bin/env bash
#
# HairHistory: Ubuntu server bootstrap (run once per server, safe to re-run).
#
#   bash deploy/01-setup-server.sh
#
# Installs: PostgreSQL 16, Go, Nginx, git. Configures ufw.
# Run as the "ubuntu" user; sudo is used only where root is required.

set -euo pipefail

# Match the Go toolchain used for local development (go.mod declares go 1.25.3).
GO_VERSION="${GO_VERSION:-1.25.3}"
GO_ROOT="/usr/local/go"

log() { printf '\n==> %s\n' "$*"; }

if [[ "$(id -u)" -eq 0 ]]; then
  echo "Run this as the ubuntu user, not as root." >&2
  exit 1
fi

export DEBIAN_FRONTEND=noninteractive

log "Updating apt package index"
sudo apt-get update -y
sudo apt-get upgrade -y

log "Installing base packages"
sudo apt-get install -y ca-certificates curl gnupg git ufw rsync

log "Installing PostgreSQL 16"
if command -v psql >/dev/null 2>&1 && psql --version | grep -q ' 16\.'; then
  echo "PostgreSQL 16 already installed; skipping."
else
  # Ubuntu 24.04 ships postgresql-16; older releases need the PGDG repository.
  if ! apt-cache policy postgresql-16 | grep -q 'Candidate: [0-9]'; then
    log "postgresql-16 not in the default repos; adding PGDG"
    sudo install -d -m 0755 /usr/share/postgresql-common/pgdg
    sudo curl -fsSL https://www.postgresql.org/media/keys/ACCC4CF8.asc \
      -o /usr/share/postgresql-common/pgdg/apt.postgresql.org.asc
    # shellcheck disable=SC1091
    echo "deb [signed-by=/usr/share/postgresql-common/pgdg/apt.postgresql.org.asc] https://apt.postgresql.org/pub/repos/apt $(. /etc/os-release && echo "${VERSION_CODENAME}")-pgdg main" \
      | sudo tee /etc/apt/sources.list.d/pgdg.list >/dev/null
    sudo apt-get update -y
  fi
  sudo apt-get install -y postgresql-16 postgresql-client-16
fi

sudo systemctl enable --now postgresql

log "Installing Go ${GO_VERSION}"
if [[ -x "${GO_ROOT}/bin/go" ]] && "${GO_ROOT}/bin/go" version | grep -q "go${GO_VERSION} "; then
  echo "Go ${GO_VERSION} already installed; skipping."
else
  arch="$(dpkg --print-architecture)"
  case "${arch}" in
    amd64 | arm64) ;;
    *)
      echo "Unsupported architecture for the official Go tarball: ${arch}" >&2
      exit 1
      ;;
  esac

  tarball="go${GO_VERSION}.linux-${arch}.tar.gz"
  tmpdir="$(mktemp -d)"
  trap 'rm -rf "${tmpdir}"' EXIT

  curl -fsSL "https://go.dev/dl/${tarball}" -o "${tmpdir}/${tarball}"
  # The Go release notes require removing the old tree rather than untarring over it.
  sudo rm -rf "${GO_ROOT}"
  sudo tar -C /usr/local -xzf "${tmpdir}/${tarball}"
fi

# $PATH and $HOME must stay literal so they expand when the profile is sourced.
# shellcheck disable=SC2016
printf 'export PATH=$PATH:%s/bin:$HOME/go/bin\n' "${GO_ROOT}" \
  | sudo tee /etc/profile.d/go.sh >/dev/null
sudo chmod 0644 /etc/profile.d/go.sh

log "Installing Node.js 22 (needed to build the frontend)"
if command -v node >/dev/null 2>&1 && node -v | grep -qE '^v(2[2-9]|[3-9][0-9])\.'; then
  echo "Node.js $(node -v) already installed; skipping."
else
  curl -fsSL https://deb.nodesource.com/setup_22.x | sudo -E bash -
  sudo apt-get install -y nodejs
fi

log "Installing Nginx"
sudo apt-get install -y nginx
sudo systemctl enable --now nginx

log "Configuring ufw"
# Allow SSH before enabling, otherwise the firewall locks out this session.
sudo ufw allow 22/tcp
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw default deny incoming
sudo ufw default allow outgoing
# Postgres stays unreachable from outside: the API talks to it over localhost only.
sudo ufw --force enable
sudo ufw status verbose

log "Done"
cat <<'EOF'
Next steps:
  1. Open a new shell (or: source /etc/profile.d/go.sh) so that go is on PATH.
  2. DB_PASSWORD='<strong password>' bash deploy/02-setup-db.sh
EOF
