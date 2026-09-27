#!/usr/bin/env bash
#
# HairHistory: switch the running server to APP_ENV=production (safe to re-run).
#
#   GOOGLE_CLIENT_ID=xxxxx.apps.googleusercontent.com \
#     bash deploy/05-enable-production.sh <domain>
#
# Run deploy/04-setup-https.sh first: with APP_ENV=production the session cookie
# gets the Secure flag, so login only works over HTTPS.
#
# Rewrites in /etc/hairhistory/api.env (a timestamped backup is taken first):
#   APP_ENV=production
#   GOOGLE_CLIENT_ID=<GOOGLE_CLIENT_ID>
#   CORS_ALLOWED_ORIGINS=https://<domain>
# Every other line, including DATABASE_URL, is left byte-for-byte untouched.
#
# Optional environment variables:
#   SKIP_DEPLOY=1   Update api.env only; do not run deploy/03-deploy-app.sh.
#   BRANCH          Branch passed to 03-deploy-app.sh (default main).
#   ENV_FILE        Path to api.env (default below).

set -euo pipefail
umask 077

DOMAIN="${1:-}"
ENV_FILE="${ENV_FILE:-/etc/hairhistory/api.env}"
BRANCH="${BRANCH:-main}"
SKIP_DEPLOY="${SKIP_DEPLOY:-0}"
SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"

log() { printf '\n==> %s\n' "$*"; }
warn() { printf '\n!!! %s\n' "$*" >&2; }
die() { printf '\nERROR: %s\n' "$*" >&2; exit 1; }

if [[ -z "${DOMAIN}" ]]; then
  cat >&2 <<'EOF'
使い方:
  GOOGLE_CLIENT_ID=xxxxx.apps.googleusercontent.com \
    bash deploy/05-enable-production.sh <domain>

  <domain> … HTTPS で公開しているドメイン（例: hairhistory-anna.duckdns.org）

事前条件:
  - deploy/04-setup-https.sh が完了し https://<domain> が開けること
  - Google Cloud Console の「承認済みの JavaScript 生成元」に
    https://<domain> を登録済みであること
EOF
  exit 1
fi

DOMAIN="${DOMAIN#https://}"
DOMAIN="${DOMAIN#http://}"
DOMAIN="${DOMAIN%/}"

CLIENT_ID="${GOOGLE_CLIENT_ID:-}"
[[ -n "${CLIENT_ID}" ]] || die "環境変数 GOOGLE_CLIENT_ID が未設定です。APP_ENV=production では必須です（API が起動しません）。"

# 値そのものはログに出さず、形だけ検証する。
if [[ "${CLIENT_ID}" != *.apps.googleusercontent.com ]]; then
  die "GOOGLE_CLIENT_ID が .apps.googleusercontent.com で終わっていません。Google Cloud Console のウェブアプリ用クライアント ID を渡してください。"
fi

[[ -f "${ENV_FILE}" ]] || die "${ENV_FILE} がありません。先に deploy/README.md の「環境変数ファイルを置く」を実施してください。"

# ------------------------------------------------------------- HTTPS の事前確認
# production は Secure Cookie を付けるので、HTTPS が成立していない状態で切り替える
# とログインできなくなります。先に確認して、駄目なら何も変更せず止めます。
log "Checking https://${DOMAIN}/healthz"
status="$(curl -sS -o /dev/null -w '%{http_code}' --max-time 15 "https://${DOMAIN}/healthz" || true)"
echo "HTTP status: ${status}"
if [[ "${status}" != "200" ]]; then
  cat >&2 <<EOF

ERROR: https://${DOMAIN}/healthz が 200 を返しません（received: ${status:-none}）。

  APP_ENV=production はセッション Cookie に Secure を付けるため、HTTPS が
  成立していない状態で切り替えるとログインできなくなります。
  ${ENV_FILE} は変更していません。

  先に実行してください:
    CERTBOT_EMAIL=you@example.com bash deploy/04-setup-https.sh ${DOMAIN}
EOF
  exit 1
fi

# ------------------------------------------------------------------ バックアップ
backup="${ENV_FILE}.bak.$(date +%Y%m%d%H%M%S)"
log "Backing up ${ENV_FILE} to ${backup}"
# cp -p なので所有者 root:root と mode 600 がそのまま引き継がれる。
sudo cp -p "${ENV_FILE}" "${backup}"
sudo ls -l "${backup}"

# -------------------------------------------------------------- api.env の更新
work_dir="$(mktemp -d)"
trap 'rm -rf "${work_dir}"' EXIT
staged="${work_dir}/api.env"

# shellcheck disable=SC2024  # 書き込み先は自分の mktemp -d 配下なので root 権限は不要
sudo cat "${ENV_FILE}" >"${staged}" ||
  die "${ENV_FILE} が読めません。sudo 権限のあるユーザーで実行してください（sudo -v を先に）。"

# 値は ENVIRON 経由で awk に渡す。-v やコマンド引数だと ps で他ユーザーに見える。
set_env_var() {
  HH_KEY="$1" HH_VALUE="$2" awk '
    BEGIN { key = ENVIRON["HH_KEY"]; value = ENVIRON["HH_VALUE"]; done = 0 }
    $0 ~ "^[[:space:]]*"key"=" {
      if (!done) { print key "=" value; done = 1 }
      next
    }
    { print }
    END { if (!done) print key "=" value }
  ' "${staged}" >"${staged}.new"
  mv "${staged}.new" "${staged}"
}

log "Rewriting APP_ENV / GOOGLE_CLIENT_ID / CORS_ALLOWED_ORIGINS"
set_env_var APP_ENV production
set_env_var GOOGLE_CLIENT_ID "${CLIENT_ID}"
set_env_var CORS_ALLOWED_ORIGINS "https://${DOMAIN}"

# DATABASE_URL を壊していないことを、書き戻す前に必ず確認する。
original_db_line="$(sudo grep -m1 -E '^[[:space:]]*DATABASE_URL=' "${ENV_FILE}" || true)"
staged_db_line="$(grep -m1 -E '^[[:space:]]*DATABASE_URL=' "${staged}" || true)"
[[ -n "${staged_db_line}" ]] || die "更新後の api.env に DATABASE_URL がありません。書き戻しを中止しました。"
[[ "${original_db_line}" == "${staged_db_line}" ]] || die "DATABASE_URL が変化しています。書き戻しを中止しました（${backup} が元の内容です）。"
unset original_db_line staged_db_line

grep -qxF 'APP_ENV=production' "${staged}" || die "APP_ENV=production を書き込めませんでした。"
grep -qxF "CORS_ALLOWED_ORIGINS=https://${DOMAIN}" "${staged}" || die "CORS_ALLOWED_ORIGINS を書き込めませんでした。"

log "Installing the updated ${ENV_FILE} (root:root 600)"
sudo install -o root -g root -m 600 "${staged}" "${ENV_FILE}"
sudo ls -l "${ENV_FILE}"
# クライアント ID の実値は出さず、設定されたことだけ示す。
sudo sed -E 's/^([[:space:]]*(DATABASE_URL|GOOGLE_CLIENT_ID)=).*/\1<redacted>/' "${ENV_FILE}"

# ---------------------------------------------------------------- 再デプロイ
if [[ "${SKIP_DEPLOY}" == "1" ]]; then
  cat <<EOF

SKIP_DEPLOY=1 のため再デプロイしていません。
フロントエンドは VITE_GOOGLE_CLIENT_ID をビルド時に埋め込むので、
次を実行するまで Google ボタンは無効のままです:

    bash ${SCRIPT_DIR}/03-deploy-app.sh ${BRANCH}
EOF
  exit 0
fi

log "Re-deploying so the frontend is rebuilt with VITE_GOOGLE_CLIENT_ID"
bash "${SCRIPT_DIR}/03-deploy-app.sh" "${BRANCH}"

# ------------------------------------------------------- dev-login が塞がれたか
log "Verifying that POST /api/auth/dev-login is gone"
dev_status="$(curl -sS -o /dev/null -w '%{http_code}' --max-time 15 \
  -X POST "https://${DOMAIN}/api/auth/dev-login" \
  -H 'Content-Type: application/json' \
  -d '{"email":"probe@example.com","name":"probe"}' || true)"
echo "POST /api/auth/dev-login -> ${dev_status:-none}"

if [[ "${dev_status}" == "404" ]]; then
  echo "OK: 開発用ログインは無効化されています。"
else
  cat >&2 <<EOF

!!! 警告: POST /api/auth/dev-login が 404 ではなく ${dev_status:-none} を返しました。

  404 以外ということは開発用ログインが生きている可能性があります。この状態では
  Google 検証なしで誰でも任意のメールアドレスとしてログインできます。
  公開を止めて、次を確認してください:

    sudo grep '^APP_ENV=' ${ENV_FILE}
    sudo systemctl status hairhistory-api
    sudo journalctl -u hairhistory-api -n 50
EOF
  exit 1
fi

cat <<EOF

================================================================================
本番モードに切り替えました: https://${DOMAIN}

  APP_ENV=production          … 開発用ログインは無効 / Cookie に Secure が付く
  CORS_ALLOWED_ORIGINS        … https://${DOMAIN}
  GOOGLE_CLIENT_ID            … 設定済み（フロントにも埋め込み済み）

ブラウザで https://${DOMAIN}/login を開き、「Google でログイン」から
ログインできることを確認してください。

ボタンが出ない / 押しても失敗する場合:
  - Google Cloud Console の「承認済みの JavaScript 生成元」が
    https://${DOMAIN} と完全一致しているか（末尾スラッシュ無し）
  - OAuth 同意画面が「テスト」状態なら、ログインするアカウントを
    テストユーザーに追加しているか
  - sudo journalctl -u hairhistory-api -n 50

元の設定に戻す場合:
  sudo cp -p ${backup} ${ENV_FILE}
  bash ${SCRIPT_DIR}/03-deploy-app.sh ${BRANCH}
================================================================================
EOF
