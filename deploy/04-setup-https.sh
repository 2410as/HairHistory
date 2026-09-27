#!/usr/bin/env bash
#
# HairHistory: issue a Let's Encrypt certificate and turn on HTTPS (safe to re-run).
#
#   CERTBOT_EMAIL=you@example.com bash deploy/04-setup-https.sh <domain>
#
# <domain> is the name that already resolves to this server, for example a free
# DuckDNS subdomain: hairhistory-anna.duckdns.org
#
# Optional environment variables:
#   CERTBOT_EMAIL          Address Let's Encrypt uses for expiry warnings. Required
#                          unless CERTBOT_ALLOW_NO_EMAIL=1 is set.
#   CERTBOT_ALLOW_NO_EMAIL Set to 1 to register without an address
#                          (--register-unsafely-without-email).
#   EXPECTED_IP            Public IPv4 of this server. Detected automatically when unset.
#   NGINX_SITE             Path to the nginx site file (default below).

set -euo pipefail

DOMAIN="${1:-}"
NGINX_SITE="${NGINX_SITE:-/etc/nginx/sites-available/hairhistory}"
CERTBOT_EMAIL="${CERTBOT_EMAIL:-}"
CERTBOT_ALLOW_NO_EMAIL="${CERTBOT_ALLOW_NO_EMAIL:-0}"
EXPECTED_IP="${EXPECTED_IP:-}"

log() { printf '\n==> %s\n' "$*"; }
warn() { printf '\n!!! %s\n' "$*" >&2; }
die() { printf '\nERROR: %s\n' "$*" >&2; exit 1; }

usage() {
  cat >&2 <<'EOF'
使い方:
  CERTBOT_EMAIL=you@example.com bash deploy/04-setup-https.sh <domain>

  <domain> … このサーバーの IP を指している名前（例: hairhistory-anna.duckdns.org）

事前条件:
  - DuckDNS などで <domain> がこのサーバーの Public IP を指していること
  - セキュリティグループで 80/tcp と 443/tcp が 0.0.0.0/0 に開いていること
    （Let's Encrypt の HTTP-01 認証は外部から 80 番に来ます）
EOF
}

if [[ -z "${DOMAIN}" ]]; then
  warn "ドメイン名を第1引数で指定してください。"
  usage
  exit 1
fi

if [[ ! "${DOMAIN}" =~ ^[a-zA-Z0-9]([a-zA-Z0-9-]*[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9-]*[a-zA-Z0-9])?)+$ ]]; then
  die "\"${DOMAIN}\" はドメイン名として不正です（http:// やスラッシュは付けない）。"
fi

# Let's Encrypt はメールアドレス無しでも登録できますが、そのアカウントには証明書
# 期限切れの警告が一切届きません。90 日で切れる証明書を通知なしで運用するのは
# 事故の元なので、既定ではアドレスを必須にし、明示的に opt-out させます。
CERTBOT_CONTACT_ARGS=()
if [[ -n "${CERTBOT_EMAIL}" ]]; then
  CERTBOT_CONTACT_ARGS=(-m "${CERTBOT_EMAIL}")
elif [[ "${CERTBOT_ALLOW_NO_EMAIL}" == "1" ]]; then
  warn "CERTBOT_EMAIL が未設定です。期限切れの通知メールは届きません（CERTBOT_ALLOW_NO_EMAIL=1 のため続行）。"
  CERTBOT_CONTACT_ARGS=(--register-unsafely-without-email)
else
  cat >&2 <<'EOF'

ERROR: CERTBOT_EMAIL が未設定です。

  Let's Encrypt の証明書は 90 日で失効します。メールアドレスを登録しておかないと
  更新に失敗しても気づけず、ある日突然サイト全体が開けなくなります。

  実行例:
    CERTBOT_EMAIL=you@example.com bash deploy/04-setup-https.sh <domain>

  どうしても登録したくない場合のみ:
    CERTBOT_ALLOW_NO_EMAIL=1 bash deploy/04-setup-https.sh <domain>
EOF
  exit 1
fi

command -v nginx >/dev/null 2>&1 || die "nginx が見つかりません。先に deploy/01-setup-server.sh を実行してください。"
[[ -f "${NGINX_SITE}" ]] || die "${NGINX_SITE} がありません。先に deploy/README.md の「Nginx を設定する」を実施してください。"

# --------------------------------------------------------------- DNS の事前確認
# Let's Encrypt には厳しいレート制限（同一ドメインで 1 時間あたり 5 回の失敗）が
# あります。DNS が向いていない状態で certbot を叩くと確実に失敗して回数を消費する
# ため、先に名前解決を突き合わせて、一致しない限り証明書取得に進みません。
log "Detecting this server's public IPv4"
if [[ -z "${EXPECTED_IP}" ]]; then
  imds_token=""
  if imds_token="$(curl -fsS --max-time 3 -X PUT http://169.254.169.254/latest/api/token \
    -H 'X-aws-ec2-metadata-token-ttl-seconds: 60' 2>/dev/null)"; then
    EXPECTED_IP="$(curl -fsS --max-time 3 -H "X-aws-ec2-metadata-token: ${imds_token}" \
      http://169.254.169.254/latest/meta-data/public-ipv4 2>/dev/null || true)"
  fi
  unset imds_token
fi

if [[ -z "${EXPECTED_IP}" ]]; then
  EXPECTED_IP="$(curl -fsS --max-time 5 https://checkip.amazonaws.com 2>/dev/null | tr -d '[:space:]' || true)"
fi

[[ -n "${EXPECTED_IP}" ]] || die "このサーバーの Public IP を自動判定できませんでした。EXPECTED_IP=<IP> を指定して再実行してください。"
echo "このサーバーの Public IP: ${EXPECTED_IP}"

log "Checking that ${DOMAIN} resolves to ${EXPECTED_IP}"
resolved=()
if command -v dig >/dev/null 2>&1; then
  mapfile -t resolved < <(dig +short A "${DOMAIN}" | grep -E '^[0-9.]+$' || true)
else
  # getent は /etc/hosts も見るので dig が使えるならそちらを優先する。
  mapfile -t resolved < <(getent hosts "${DOMAIN}" | awk '{print $1}' | grep -E '^[0-9.]+$' || true)
fi

if [[ "${#resolved[@]}" -eq 0 ]]; then
  die "${DOMAIN} の A レコードが引けません。DuckDNS の設定と反映（数分かかります）を確認してください。"
fi

echo "${DOMAIN} の A レコード:"
printf '  %s\n' "${resolved[@]}"

if ! printf '%s\n' "${resolved[@]}" | grep -qxF "${EXPECTED_IP}"; then
  cat >&2 <<EOF

ERROR: ${DOMAIN} が このサーバー (${EXPECTED_IP}) を指していません。

  証明書取得は試行しません（Let's Encrypt のレート制限を無駄に消費するため）。
  DuckDNS の管理画面で current ip を ${EXPECTED_IP} に更新し、数分待ってから
  もう一度このスクリプトを実行してください。
EOF
  exit 1
fi

# ------------------------------------------------------------------- certbot
if command -v certbot >/dev/null 2>&1 && dpkg -s python3-certbot-nginx >/dev/null 2>&1; then
  log "certbot is already installed - skipping apt"
else
  log "Installing certbot and python3-certbot-nginx"
  sudo apt-get update
  sudo DEBIAN_FRONTEND=noninteractive apt-get install -y certbot python3-certbot-nginx
fi

# --------------------------------------------------------------- server_name
# 初期状態の server_name は IP（または REPLACE_WITH_...）がハードコードされて
# いるため、certbot が対象の server ブロックを見つけられません。全ての
# server_name 行を引数のドメインに揃えます。既にドメインでも結果は同じ（冪等）。
log "Pointing server_name at ${DOMAIN} in ${NGINX_SITE}"
site_backup="${NGINX_SITE}.bak.$(date +%Y%m%d%H%M%S)"
sudo cp -p "${NGINX_SITE}" "${site_backup}"
sudo sed -i -E "s|^([[:space:]]*)server_name[[:space:]]+[^;]*;|\1server_name ${DOMAIN};|" "${NGINX_SITE}"
sudo grep -nE '^[[:space:]]*server_name' "${NGINX_SITE}"

if ! sudo nginx -t; then
  warn "nginx -t が失敗したため ${site_backup} から復元します。"
  sudo cp -p "${site_backup}" "${NGINX_SITE}"
  die "nginx の設定が壊れています。上のエラーを確認してください。"
fi
sudo systemctl reload nginx

# ------------------------------------------------------------------ 証明書取得
log "Requesting a certificate for ${DOMAIN}"
# --keep-until-expiring: 有効な証明書が残っているうちは再取得しないので、この
# スクリプトを何度流してもレート制限を消費しません。
# --redirect: 80 番を 443 へ 301 する設定を nginx に書き込みます。
sudo certbot --nginx \
  -d "${DOMAIN}" \
  --redirect \
  --agree-tos \
  --non-interactive \
  --keep-until-expiring \
  "${CERTBOT_CONTACT_ARGS[@]}"

sudo nginx -t
sudo systemctl reload nginx

# ------------------------------------------------------------------- 自動更新
log "Checking the certbot renewal timer"
if systemctl list-unit-files 'certbot.timer' | grep -q certbot.timer; then
  sudo systemctl enable --now certbot.timer
  systemctl --no-pager list-timers certbot.timer || true
else
  warn "certbot.timer が見つかりません。snap 版 certbot か cron での更新設定を確認してください。"
fi

log "Dry-running the renewal"
sudo certbot renew --dry-run

# ----------------------------------------------------------------- 動作確認
log "Verifying HTTPS"
curl -fsS -o /dev/null -w 'https://%{http_code} %{url_effective}\n' "https://${DOMAIN}/healthz"
curl -fsS -o /dev/null -w 'http -> %{http_code} %{redirect_url}\n' "http://${DOMAIN}/healthz" || true

cat <<EOF

================================================================================
HTTPS が有効になりました: https://${DOMAIN}

次にやること:

1. Google Cloud Console で OAuth クライアント ID を作る（または既存のものを編集）
   「承認済みの JavaScript 生成元」に次を追加:
       https://${DOMAIN}
   （末尾スラッシュ無し。リダイレクト URI とクライアントシークレットは不要）

2. 本番モードへ切り替える（APP_ENV=production / CORS / クライアント ID の反映）
       GOOGLE_CLIENT_ID=xxxxx.apps.googleusercontent.com \\
         bash deploy/05-enable-production.sh ${DOMAIN}

   このスクリプトが /etc/hairhistory/api.env を更新し、
   deploy/03-deploy-app.sh を呼んでフロントを再ビルドします。

3. 切り替え後、開発用ログインが塞がれたことを確認する
       curl -i -X POST https://${DOMAIN}/api/auth/dev-login \\
         -H 'Content-Type: application/json' -d '{"email":"a@example.com","name":"a"}'
   → 404 が返るのが正しい状態です（05 のスクリプトも自動で確認します）。
================================================================================
EOF
