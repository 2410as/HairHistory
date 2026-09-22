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
| ドメイン名（任意） | Route 53 など | EC2 の Public DNS でも動くが、HTTPS 化には実質必須 |

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

- **HTTP のみ** … セッション Cookie が平文で流れます。ドメインを用意して `certbot --nginx` で TLS を有効化するまでは、本番利用しないでください。
- **バックアップが自動化されていない** … `pg_dump` を cron に入れる、EBS スナップショットを取る、などを別途検討。
- **1 台構成** … このインスタンスが落ちるとサービス全体が落ちます（学習・ポートフォリオ用途としては妥当な割り切り）。
