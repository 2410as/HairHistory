# HairHistory デプロイ手順（AWS EC2 / Ubuntu）

サーバーに SSH で入って、上から順に実行する手順書です。
構成は **1 台完結**（EC2 の中に Go API・PostgreSQL 16・Nginx をすべて置く）。RDS などのマネージドサービスは使いません。

```
ブラウザ ──80──▶ Nginx ──┬── /        → /var/www/hairhistory（React のビルド成果物）
                          └── /api/    → 127.0.0.1:8080（Go API / systemd）
                                             │
                                             └── 127.0.0.1:5432 PostgreSQL 16
```

---

## 0. 事前に用意するもの

| もの | 取得元 | 無いとどうなるか |
|---|---|---|
| EC2 インスタンス（Ubuntu 22.04 / 24.04, t3.small 以上推奨） | AWS コンソール | — |
| SSH キーペア | AWS コンソール | ログインできない |
| DB パスワード | 自分で生成（`openssl rand -base64 24`） | 手順 2 が止まる |
| **Google OAuth クライアント ID** | Google Cloud Console（下の手順 4 参照） | `APP_ENV=production` では **API が起動しない** |
| ドメイン名 | **DuckDNS（無料）** または Route 53 など | EC2 の Public DNS でも HTTP では動くが、**HTTPS 化には必須**（手順 10） |

> t2.micro（メモリ 1GB）だと `npm run build`（TypeScript + Vite）がメモリ不足で落ちることがあります。落ちたらスワップを 2GB 足すか、インスタンスサイズを上げてください。

### セキュリティグループ（推奨設定）

| ポート | ソース | 用途 |
|---|---|---|
| 22 (SSH) | **自分のグローバル IP のみ**（`x.x.x.x/32`） | 運用作業 |
| 80 (HTTP) | 0.0.0.0/0 | 公開 |
| 443 (HTTPS) | 0.0.0.0/0 | 将来の TLS 用に開けておく |
| 5432 | **開けない** | DB は localhost からのみ接続する |

8080（API）も開けません。外部からは必ず Nginx を経由させます。

---

## 1. サーバーの初期セットアップ

```bash
ssh ubuntu@<サーバーのアドレス>
git clone https://github.com/2410as/HairHistory.git ~/HairHistory
cd ~/HairHistory
bash deploy/01-setup-server.sh
```

やること: apt 更新 / PostgreSQL 16 / Go 1.25.3（公式 tarball を `/usr/local/go`）/ Node.js 22 / Nginx / git / ufw（22・80・443 のみ許可、それ以外は拒否）。
すでに入っているものはスキップするので、何度実行しても壊れません。

終わったら **一度ログアウトして入り直す**（`/etc/profile.d/go.sh` の PATH を反映するため）。

```bash
go version   # go1.25.3 linux/... と出ればOK
```

---

## 2. データベースの作成

パスワードはスクリプトに書かず、環境変数で渡します。

```bash
cd ~/HairHistory
DB_PASSWORD='ここに生成したパスワード' bash deploy/02-setup-db.sh
```

やること: ロール `hairhistory` とデータベース `hairhistory` を作成（既にあればスキップ）し、localhost からログインできるか検証。
`pg_hba.conf` は触りません = **外部から PostgreSQL には接続できないまま**です。

パスワードはこのあと `DATABASE_URL` に使うので控えておいてください（シェル履歴に残したくない場合は、コマンドの先頭に半角スペースを入れて実行）。

---

## 3. 環境変数ファイルを置く

```bash
sudo install -d -m 0755 /etc/hairhistory
sudo cp ~/HairHistory/deploy/api.env.example /etc/hairhistory/api.env
sudo chown root:root /etc/hairhistory/api.env
sudo chmod 600 /etc/hairhistory/api.env
sudo nano /etc/hairhistory/api.env
```

`REPLACE_WITH_...` の 3 箇所を埋めます。

- `DATABASE_URL` … 手順 2 のパスワード（`@ : / ? # %` を含むなら URL エンコードする）
- `GOOGLE_CLIENT_ID` … 手順 4 で取得するもの
- `CORS_ALLOWED_ORIGINS` … `http://自分のドメイン` または `http://ec2-xx-xx.compute.amazonaws.com`（末尾スラッシュ無し）

> **`APP_ENV=production` は絶対に消さないこと。** これが `production` 以外だと、バックエンドが開発用ログイン `POST /api/auth/dev-login` を有効化します。このエンドポイントは Google 検証なしで任意のメールアドレスのセッションを発行するため、公開サーバーでは「誰でも誰にでもなりすませる」状態になります。

---

## 4. Google OAuth クライアント ID を取得する

`APP_ENV=production` のとき `GOOGLE_CLIENT_ID` が空だと **API は起動エラーで落ちます**（`GOOGLE_CLIENT_ID is required when APP_ENV=production`）。

1. [Google Cloud Console](https://console.cloud.google.com/) でプロジェクトを作成（既存でも可）
2. 「API とサービス」→「OAuth 同意画面」を設定
   - User Type: External、アプリ名・サポートメール・デベロッパー連絡先を入力
   - 公開ステータスがテストのうちは、ログインできるのは「テストユーザー」に登録したアカウントだけ
3. 「認証情報」→「認証情報を作成」→「OAuth クライアント ID」
   - アプリケーションの種類: **ウェブ アプリケーション**
4. **承認済みの JavaScript 生成元** に、ブラウザでアクセスする URL を**そのまま**追加
   - 例: `http://ec2-xx-xx-xx-xx.ap-northeast-1.compute.amazonaws.com`
   - 例: `http://your.domain` （HTTPS 化したら `https://your.domain` も追加）
   - 末尾スラッシュやパスは付けない。ポート番号まで含めて一致させる
5. 「承認済みのリダイレクト URI」は不要（このアプリは Google Identity Services のトークン方式で、リダイレクトを使いません）
6. 発行された `xxxxx.apps.googleusercontent.com` を `/etc/hairhistory/api.env` の `GOOGLE_CLIENT_ID` に貼る

クライアントシークレットはこのアプリでは使いません。**シークレットはサーバーに置かないでください。**

> `GOOGLE_CLIENT_ID` はフロントエンドのビルドにも使われます（`VITE_GOOGLE_CLIENT_ID`）。手順 5 のデプロイスクリプトが `/etc/hairhistory/api.env` から自動で読み取るので、二重に書く必要はありません。ID を変更したら **必ずデプロイを流し直す**（ビルド時に埋め込まれるため）。

---

## 5. systemd サービスを登録する

```bash
sudo cp ~/HairHistory/deploy/hairhistory-api.service /etc/systemd/system/hairhistory-api.service
sudo systemctl daemon-reload
sudo systemctl enable hairhistory-api
```

まだバイナリが無いので `start` はしません（次の手順のスクリプトが起動します）。

---

## 6. Nginx を設定する

```bash
sudo cp ~/HairHistory/deploy/nginx-hairhistory.conf /etc/nginx/sites-available/hairhistory
sudo nano /etc/nginx/sites-available/hairhistory   # server_name を自分のドメイン/Public DNS に書き換える
sudo ln -sfn /etc/nginx/sites-available/hairhistory /etc/nginx/sites-enabled/hairhistory
sudo rm -f /etc/nginx/sites-enabled/default
sudo nginx -t && sudo systemctl reload nginx
```

---

## 7. アプリをデプロイする

```bash
cd ~/HairHistory
bash deploy/03-deploy-app.sh main
```

引数はブランチ名（省略時 `main`）。やること:

1. `/opt/hairhistory` にリポジトリを clone / fetch して、指定ブランチに **`git reset --hard`**
2. `go build` → `/opt/hairhistory/bin/api`（旧バイナリは `bin/api.prev` に退避）
3. `golang-migrate` でマイグレーション適用（`schema_migrations` で管理されるので再実行しても二重適用されない）
4. `frontend/.env.production` を生成 → `npm ci && npm run build` → `dist/` を `/var/www/hairhistory` へ同期
5. `systemctl restart hairhistory-api` → `/healthz` を確認

> `/opt/hairhistory` の中身は毎回 `git reset --hard` で上書きされます。**サーバー上で直接ファイルを編集しないでください。**

---

## 8. 動作確認

```bash
sudo systemctl status hairhistory-api          # active (running) か
curl -i http://127.0.0.1:8080/healthz          # {"status":"ok"}
curl -i http://127.0.0.1/healthz               # Nginx 経由でも同じ
curl -I http://127.0.0.1/                      # 200 / text/html
```

最後にブラウザで `http://<ドメイン or Public DNS>` を開き、

- トップページが表示される
- Google ログインボタンからログインできる
- 施術記録の追加・一覧・共有リンク発行ができる

**開発用ログインが塞がれているかの確認（重要）:**

```bash
curl -i -X POST http://127.0.0.1/api/auth/dev-login \
  -H 'Content-Type: application/json' -d '{"email":"a@example.com","name":"a"}'
```

`404 not found` が返れば正しい状態です。`200` が返ったら `APP_ENV` が `production` になっていないので、すぐ直して再起動してください。

---

## 9. 更新（2 回目以降のデプロイ）

```bash
cd ~/HairHistory && git pull      # 手順書・スクリプト自体の更新を取り込む
bash deploy/03-deploy-app.sh main
```

---

## 10. HTTPS 化する（DuckDNS + Let's Encrypt / 無料）

HTTP のままだと `APP_ENV=production` に切り替えられません。production はセッション Cookie に `Secure` を付けるため、HTTPS でないとブラウザが Cookie を保存せず、ログインが成立しないからです。
Let's Encrypt の証明書は IP アドレスには発行できないので、まず**名前**が要ります。ドメインを買わずに済ませるなら DuckDNS の無料サブドメインを使います。

### 10-1. DuckDNS でサブドメインを取る

1. <https://www.duckdns.org/> を開き、GitHub / Google などでログイン
2. 好きな名前（例: `hairhistory-anna`）を入力して **add domain**
3. `current ip` の欄に **EC2 の Public IP** を入れて **update ip**
4. これで `hairhistory-anna.duckdns.org` がサーバーを指します

> ページに表示される **token は秘密情報**です。README・コミット・チャットに貼らないでください。IP を手で更新する運用なら token は使いません。
>
> EC2 を停止→起動すると Public IP が変わります（Elastic IP を割り当てていない場合）。変わったら DuckDNS の `current ip` を入れ直してください。証明書自体は取り直し不要です。

反映されたかの確認:

```bash
dig +short A hairhistory-anna.duckdns.org    # → EC2 の Public IP が返れば OK
```

### 10-2. セキュリティグループ

**443/tcp を 0.0.0.0/0 に開けます。80/tcp も開けたままにしてください**（Let's Encrypt の HTTP-01 認証が 80 番に来ます。閉じていると証明書が取れず、更新もできません）。

### 10-3. 証明書を取る

```bash
cd ~/HairHistory && git pull
CERTBOT_EMAIL=you@example.com \
  bash deploy/04-setup-https.sh hairhistory-anna.duckdns.org
```

やること:

1. このサーバーの Public IP を判定し、**引数のドメインがそこを指しているか確認**（不一致なら証明書取得を試みずに停止。Let's Encrypt のレート制限＝同一ドメイン 1 時間に 5 失敗、を無駄に消費しないため）
2. `certbot` と `python3-certbot-nginx` を apt で導入（導入済みならスキップ）
3. `/etc/nginx/sites-available/hairhistory` の `server_name` を引数のドメインに書き換え（バックアップあり・`nginx -t` 失敗時は自動で復元）
4. `certbot --nginx --redirect` で証明書取得 → 443 待ち受けと 80→443 リダイレクトを nginx に書き込み
5. `certbot.timer` を有効化し、`certbot renew --dry-run` で自動更新をテスト

`CERTBOT_EMAIL` は必須にしてあります。証明書は 90 日で失効するため、通知先が無いと更新失敗に気づけずサイトが丸ごと開かなくなるからです。どうしても登録したくない場合だけ `CERTBOT_ALLOW_NO_EMAIL=1` を付けてください。

確認:

```bash
curl -I https://hairhistory-anna.duckdns.org/          # 200
curl -I http://hairhistory-anna.duckdns.org/           # 301 → https://...
sudo certbot certificates                              # 有効期限
```

---

## 11. 本番モードに切り替える（Google ログインを有効化）

### 11-1. Google OAuth クライアント ID を用意する

手順 4 と同じ画面ですが、**承認済みの JavaScript 生成元に HTTPS の URL を登録**します。

1. <https://console.cloud.google.com/> → プロジェクトを選択（無ければ作成）
2. 「API とサービス」→「OAuth 同意画面」
   - User Type: External、アプリ名・サポートメール・デベロッパー連絡先を入力
   - 公開ステータスが「テスト」の間は、**テストユーザーに登録したアカウントしかログインできません**。自分の Gmail を追加しておくこと
3. 「認証情報」→「認証情報を作成」→「OAuth クライアント ID」→ 種類: **ウェブ アプリケーション**
4. **承認済みの JavaScript 生成元**に次を追加（末尾スラッシュ無し・完全一致）

   ```
   https://hairhistory-anna.duckdns.org
   ```

5. **承認済みのリダイレクト URI は空のままでよい**（Google Identity Services のトークン方式で、リダイレクトを使わないため）
6. **クライアントシークレットは使いません**。サーバーにも置かないでください
7. 発行された `xxxxx.apps.googleusercontent.com` を控える（次のコマンドで渡します）

> 既存のクライアント ID を使い回す場合も、生成元に `https://...` を追加し忘れると、ブラウザのコンソールに `The given origin is not allowed for the given client ID` が出てボタンが機能しません。

### 11-2. 切り替えスクリプトを実行する

```bash
cd ~/HairHistory && git pull
GOOGLE_CLIENT_ID=xxxxx.apps.googleusercontent.com \
  bash deploy/05-enable-production.sh hairhistory-anna.duckdns.org
```

やること:

1. `https://<domain>/healthz` が 200 を返すか確認（駄目なら **api.env を一切変更せずに停止**）
2. `/etc/hairhistory/api.env` をタイムスタンプ付きでバックアップ（`api.env.bak.YYYYmmddHHMMSS`、root:root 600）
3. `APP_ENV=production` / `GOOGLE_CLIENT_ID` / `CORS_ALLOWED_ORIGINS=https://<domain>` の 3 行だけ書き換え
   **`DATABASE_URL` は書き換え前後で一致することを検証してから書き戻す**ので、DB 接続情報が壊れることはありません
4. `deploy/03-deploy-app.sh` を呼んで再デプロイ（フロントは `VITE_GOOGLE_CLIENT_ID` を**ビルド時に**埋め込むので、再ビルドしないと Google ボタンは無効のままです）
5. `POST /api/auth/dev-login` が **404** を返すことを確認（404 以外なら警告を出して異常終了）

`SKIP_DEPLOY=1` を付けると api.env の更新だけ行い、再デプロイは自分のタイミングで実行できます。

### 11-3. 確認

```bash
curl -I https://hairhistory-anna.duckdns.org/
curl -i -X POST https://hairhistory-anna.duckdns.org/api/auth/dev-login \
  -H 'Content-Type: application/json' -d '{"email":"a@example.com","name":"a"}'   # → 404
sudo grep '^APP_ENV=' /etc/hairhistory/api.env                                     # → production
```

ブラウザで `https://hairhistory-anna.duckdns.org/login` を開き、「Google でログイン」ボタン（Google 公式のボタンが表示されます）からログインできれば完了です。

**元に戻したいとき:**

```bash
sudo cp -p /etc/hairhistory/api.env.bak.<タイムスタンプ> /etc/hairhistory/api.env
bash deploy/03-deploy-app.sh main
```

---

## トラブルシュート

```bash
# API のログをリアルタイムで見る（JSON 形式で出力される）
sudo journalctl -u hairhistory-api -f

# 直近の起動失敗の原因だけ見る
sudo journalctl -u hairhistory-api -n 50 --no-pager

# Nginx
sudo nginx -t
sudo tail -f /var/log/nginx/error.log
sudo tail -f /var/log/nginx/access.log

# PostgreSQL
sudo systemctl status postgresql
psql "postgres://hairhistory@127.0.0.1:5432/hairhistory" -c '\dt'
```

| 症状 | 原因の見当 |
|---|---|
| `GOOGLE_CLIENT_ID is required when APP_ENV=production` で起動しない | 手順 4 が未完。`/etc/hairhistory/api.env` を確認 |
| `DATABASE_URL is required` | `api.env` の読み込み失敗。`systemctl cat hairhistory-api` で `EnvironmentFile` のパスを確認 |
| `ping database: ... password authentication failed` | `DATABASE_URL` のパスワード違い、または記号の URL エンコード漏れ |
| ブラウザで 502 Bad Gateway | API が落ちている。`journalctl -u hairhistory-api` を見る |
| ブラウザで 404（ページ更新時だけ） | Nginx の SPA fallback が効いていない。`try_files` の行と `root` を確認 |
| Google ログインボタンが出ない / `origin is not allowed` | Google Console の「承認済み JavaScript 生成元」とアクセス URL が不一致 |
| API 呼び出しが CORS で落ちる | `CORS_ALLOWED_ORIGINS` が実際のアクセス URL と不一致（スキーム・末尾スラッシュに注意） |
| ログインしてもすぐログアウトされる | HTTPS 化後に `APP_ENV` が production でない等で Cookie の Secure 属性が不整合 |
| `npm run build` が Killed で終わる | メモリ不足。スワップを足すかインスタンスを大きくする |
| `04-setup-https.sh` が「DNS が一致しない」で止まる | DuckDNS の `current ip` が古い。EC2 の Public IP に更新して数分待つ |
| certbot が `Timeout during connect` で失敗 | セキュリティグループで 80/tcp が閉じている（HTTP-01 認証は 80 番に来る） |
| certbot が `too many failed authorizations` | Let's Encrypt のレート制限。1 時間待ってから、DNS を直した上で再実行 |
| HTTPS 化後に証明書が切れた | `sudo systemctl status certbot.timer` と `sudo certbot renew --dry-run` を確認 |

---

## ロールバック

### A. 直前のバイナリに戻す（一番速い）

```bash
sudo systemctl stop hairhistory-api
sudo -u ubuntu cp /opt/hairhistory/bin/api.prev /opt/hairhistory/bin/api
sudo systemctl start hairhistory-api
sudo systemctl status hairhistory-api
```

`api.prev` はデプロイのたびに「1 つ前」で上書きされます。2 世代前には戻れません。

### B. 特定のコミットに戻してビルドし直す

```bash
git -C /opt/hairhistory log --oneline -10     # 戻したいコミットを探す
bash ~/HairHistory/deploy/03-deploy-app.sh <戻したいブランチ名>
```

タグやコミット指定で戻したい場合は、手動で:

```bash
cd /opt/hairhistory
git checkout <commit-sha>
cd backend && go build -o /opt/hairhistory/bin/api ./cmd/api
cd ../frontend && npm ci && npm run build
sudo rsync -a --delete dist/ /var/www/hairhistory/
sudo systemctl restart hairhistory-api
```

### C. とにかく止める

```bash
sudo systemctl stop hairhistory-api
sudo systemctl disable hairhistory-api   # 再起動後も上がってこないようにする
```

### マイグレーションのロールバックについて

`migrate ... down` はテーブルを **DROP** します（`000001_init.down.sql` は全テーブル削除）。データが入っているなら、先に必ずバックアップを取ってください。

```bash
sudo -u postgres pg_dump hairhistory > ~/hairhistory-$(date +%Y%m%d-%H%M).sql
```

---

## 残っているリスク（把握しておくこと）

- **手順 10 を済ませるまでは HTTP のみ** … セッション Cookie が平文で流れます。`deploy/04-setup-https.sh` で TLS を有効化するまでは、本番利用しないでください。
- **DuckDNS 依存** … 無料サービスなので SLA はありません。停止すると名前が引けなくなり、証明書の更新も通りません。長期運用するなら独自ドメインへ移行してください。
- **Public IP が変わると名前が外れる** … Elastic IP を割り当てていない EC2 は停止→起動で IP が変わります。DuckDNS の `current ip` を更新するまでサイトは開けません。
- **バックアップが自動化されていない** … `pg_dump` を cron に入れる、EBS スナップショットを取る、などを別途検討。
- **1 台構成** … このインスタンスが落ちるとサービス全体が落ちます（学習・ポートフォリオ用途としては妥当な割り切り）。
