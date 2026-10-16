# HairHistory

美容院での施術内容を記録し、次に行くお店に共有リンクや QR コードで見せられる Web アプリケーションです。
Go 製の API サーバーと React の SPA を分離して実装し、AWS EC2 上に自前で構築したサーバーで動かしています。

**公開 URL: <https://hairhistory.duckdns.org>**

Google アカウントでログインして実際に試せます（HTTPS / Google ログインとも本番で動作中）。

## 解決したい課題

美容院で「前回どんな施術をしたか」を思い出せず、うまく伝えられないことがあります。

- カラー剤やパーマの内容は施術直後でないと覚えていられない
- 複数の美容院を使い分けていると、履歴がどこにもまとまらない
- 口頭やスマホのメモでの説明は、美容師側に正確に伝わりにくい

HairHistory は施術履歴を 1 か所に貯め、**閲覧専用のリンク 1 本**にまとめて美容師に渡せるようにします。
リンクには氏名・メールアドレスなどの個人情報を含めず、施術内容だけを見せます。

## 主な機能

| 機能 | 内容 |
| --- | --- |
| Google ログイン | Google Identity Services のボタンから受け取った ID Token をサーバー側で検証し、セッション Cookie を発行。本番で動作確認済み |
| 施術履歴の記録 | 施術日・施術項目・サロン名・メモ・金額を登録 / 一覧 / 編集 / 削除 |
| 施術項目の複数選択 | チェックボックス 10 種（カット・カラー・白髪染め・ブリーチ・ハイライト・パーマ・縮毛矯正・トリートメント・髪質改善・ヘッドスパ）＋ 自由入力。1 件につき最大 10 項目 |
| 共有リンクの発行 | 有効期限付きのトークンを発行。一覧表示と任意のタイミングでの無効化が可能 |
| 公開閲覧 | 発行されたリンクはログイン不要で閲覧できる。氏名・メールアドレスは返さない |
| QR コード生成 | 共有リンクの QR コードをブラウザ側で生成し、その場で見せられる |
| HTTPS 公開 | DuckDNS のサブドメイン＋ Let's Encrypt の証明書で TLS 化。HTTP は 443 へリダイレクト |

## 技術スタック

| 層 | 採用技術 |
| --- | --- |
| フロントエンド | Vite + React 19 + TypeScript（SPA）、Chakra UI、React Router、axios、qrcode.react |
| バックエンド | Go 1.25 + chi、pgx/v5、golang-migrate、`google.golang.org/api/idtoken` |
| データベース | PostgreSQL 16（ORM を使わず生 SQL） |
| インフラ | AWS EC2（Ubuntu）1 台構成、Nginx リバースプロキシ（TLS 終端・レート制限）、systemd、Let's Encrypt（certbot）、DuckDNS |
| 開発環境 | Docker Compose（PostgreSQL）、Makefile |

### なぜこの構成にしたか

**SSR を使わず SPA にした**
Next.js のようなフルスタックフレームワークを使うと、フロントとバックの境界が曖昧になります。
API を独立した Go サーバーに切り出し、フロントは HTTP でそれを叩くだけにすることで、
認証・CORS・Cookie といった通信の基礎を自分で組み立てる構成にしました。

**ORM を使わず生 SQL にした**
発行される SQL が読めない状態を避けるためです。
クエリは `internal/*/repository.go` にすべて直書きし、値は必ずプレースホルダ（`$1`）でバインドしています。
インデックス設計や `ON DELETE CASCADE` の効き方を、自分で確認しながら決められます。

**マネージドサービスではなく IaaS にした**
Vercel や Cloud Run であれば数分で公開できますが、その裏側がブラックボックスになります。
EC2 に OS から入って PostgreSQL・Nginx・systemd を自分で設定することで、
ビルド成果物の配置、リバースプロキシ、プロセスの常駐化までを一通り経験する構成にしました。

## アーキテクチャ

バックエンドは **handler → usecase → repository** のレイヤード構成です。
ドメイン（`auth` / `treatment` / `share`）ごとにパッケージを分け、横断的な処理は `httpx` に集約しています。

```
HTTP リクエスト
   │
   ├─ httpx: リクエスト ID / ログ / パニック復帰 / CORS
   │
   ▼
handler    … JSON のデコード、HTTP ステータスへの変換
   ▼
usecase    … 業務ルール（バリデーション、所有者チェック、有効期限判定）
   ▼
repository … SQL の実行。interface として usecase に注入される
   ▼
PostgreSQL
```

`usecase` は具体的な DB 実装ではなく `repository` の interface に依存しているため、
テストではインメモリのフェイクに差し替えられます（`internal/*/usecase_test.go`）。

## セキュリティ上の設計

個人の施術履歴を扱うため、以下は意図して作り込んでいます。

**セッション**
`hh_session` という HttpOnly Cookie で保持し、DB には生のトークンではなく SHA-256 ハッシュ
（`sessions.token_hash`）を保存しています（`internal/auth/auth.go` の `HashToken`）。
DB が流出してもセッションの乗っ取りには使えません。トークン自体は `crypto/rand` の 32 バイトです。

**共有トークン**
`crypto/rand` の 32 バイトを base64url で出力しています（`internal/share/share.go`）。
推測で当てられる長さではありません。

**公開閲覧で個人情報を返さない**
`GET /api/public/shares/{token}` のレスポンスは施術日・施術項目・サロン名・メモ・金額だけで、
氏名・メールアドレス・ユーザー ID は構造体の時点で含めていません（`internal/share/handler.go` の
`publicTreatmentResponse`）。

**トークンの存在を漏らさない**
共有リンクの削除（`DELETE /api/shares/{token}`）は「存在しないトークン」と「他人のトークン」を
区別せず 404 に統一しています。公開参照も、未発行・期限切れ・無効化済みをすべて同じ 404 で返します。
403 を返し分けると「そのトークンは存在する」という情報を与えてしまいます。

**共有トークンをログに残さない**
API 側のアクセスログとパニックログでトークン部分をマスクし（`internal/httpx/middleware.go`）、
Nginx 側でも `map $uri $logged_uri` で同じパスを `[REDACTED]` に置換しています。
Nginx は API とは独立にリクエストラインを記録するため、両方で塞ぐ必要がありました。

**本番で開発用ログインが無効になる**
`POST /api/auth/dev-login` のルート登録自体を `APP_ENV != production` の条件に入れているため、
本番ではエンドポイントが存在せず 404 になります（`cmd/api/main.go` から
`auth.NewHandler(..., !cfg.IsProduction())`）。あわせて `APP_ENV=production` では
Cookie に `Secure` が付き、`GOOGLE_CLIENT_ID` が未設定だと起動時にエラーで停止します。

**Nginx のレート制限とセキュリティヘッダ**
t3.micro 1 台構成なので、`/api/` への無制限なアクセスでサービス全体が落ちます。
用途ごとに 3 つのゾーンを張っています（`deploy/nginx-hairhistory.conf`）。

| 対象 | 制限 | 理由 |
| --- | --- | --- |
| `/api/auth/google`、`/api/auth/dev-login` | 10r/m（burst 5） | 認証情報を扱うエンドポイント。正常なログインは 1 POST |
| `/api/public/shares/` | 1r/s（burst 10） | 唯一の未認証エンドポイント。1 ページ閲覧＝1 リクエスト |
| `/api/` 全体 | 10r/s（burst 30） | 実測で画面を一通り触って約 33 リクエスト |

`/api/auth/me` と `/api/auth/logout` は意図的に認証系ゾーンから外しています。`/me` は全ページの
読み込みごとに走るため、ログイン相当の制限をかけるとアプリが自分で自分をロックアウトします。
超過時は 503 ではなく 429 を返し、SPA の axios インターセプタが読める
`{"error":{"code","message"}}` の形に揃えています。

セキュリティヘッダは `Strict-Transport-Security` / `X-Content-Type-Options` / `X-Frame-Options` /
`Content-Security-Policy`（`frame-ancestors 'none'`）/ `Referrer-Policy` を付けています。
`add_header` は階層をまたいで継承されないため、各 location に同じセットを書いています。
HSTS は `max-age=300` に抑えてあります。証明書の更新に失敗したときブラウザが HTTP に
フォールバックできず復旧不能になるのを避けるためで、certbot の自動更新が一度成功したら延ばす予定です。
CSP に `script-src` を入れていないのは、Google Identity Services が自前のスクリプトと iframe を
注入するため、ローカルで検証できないポリシーを入れるとログインが壊れるからです。

## ディレクトリ構成

```
HairHistory/
├── backend/
│   ├── cmd/api/main.go              # 依存の組み立てとサーバー起動
│   ├── internal/
│   │   ├── auth/                    # Google ID Token 検証、セッション管理
│   │   ├── treatment/               # 施術履歴の CRUD とバリデーション
│   │   ├── share/                   # 共有リンクの発行・無効化・公開参照
│   │   ├── httpx/                   # ミドルウェア、エラー整形、レスポンス
│   │   ├── config/                  # 環境変数の読み込みと検証
│   │   └── db/                      # pgxpool の接続
│   ├── migrations/                  # golang-migrate 用の SQL
│   ├── go.mod
│   └── go.sum
├── frontend/
│   ├── src/
│   │   ├── pages/                   # Home / Login / Dashboard / TreatmentList /
│   │   │                            #   TreatmentDetail / ShareLink / PublicShare / NotFound
│   │   ├── components/              # Layout, Navbar, TreatmentCard, ProtectedRoute ほか
│   │   ├── contexts/AuthContext.tsx # ログイン状態の保持
│   │   ├── hooks/                   # useTreatments, useShares
│   │   ├── lib/apiClient.ts         # axios インスタンス（Cookie 送信）
│   │   ├── types/                   # API のレスポンス型
│   │   ├── theme.ts / design.ts     # Chakra UI のテーマとデザイントークン
│   │   ├── App.tsx                  # ルーティング
│   │   └── main.tsx
│   ├── index.html
│   ├── package.json
│   └── vite.config.ts
├── deploy/                          # EC2 セットアップ〜デプロイのスクリプト一式
│   ├── 01-setup-server.sh           # PostgreSQL 16 / Go / Nginx / ufw の導入
│   ├── 02-setup-db.sh               # ロールとデータベースの作成
│   ├── 03-deploy-app.sh             # ビルドとリリース（2 回目以降もこれ）
│   ├── 04-setup-https.sh            # certbot で Let's Encrypt の証明書を取得し HTTPS 化
│   ├── 05-enable-production.sh      # APP_ENV=production に切り替え Google ログインを有効化
│   ├── hairhistory-api.service      # systemd ユニット
│   ├── nginx-hairhistory.conf       # リバースプロキシ / レート制限 / セキュリティヘッダ
│   ├── api.env.example
│   └── README.md
├── docs/
│   ├── API.md
│   └── DATABASE.md
├── docker-compose.yml               # ローカル開発用の PostgreSQL 16
├── Makefile
└── README.md
```

## 画面とルーティング

| パス | 内容 | 認証 |
| --- | --- | --- |
| `/` | トップページ | 不要 |
| `/login` | ログイン | 不要 |
| `/dashboard` | ダッシュボード（施術履歴の一覧） | 必要 |
| `/treatments` | 施術履歴の一覧 | 必要 |
| `/treatment/:id` | 施術履歴の詳細・編集・削除。`/treatment/new` は新規登録フォーム | 必要 |
| `/share` | 共有リンクの発行・QR コード表示・無効化 | 必要 |
| `/share/:token` | 共有された履歴の公開閲覧 | 不要 |

## ローカルセットアップ

### 前提ツール

- Go 1.25 以上
- Node.js 20 以上
- Docker（PostgreSQL の起動に使用）
- [golang-migrate](https://github.com/golang-migrate/migrate) の CLI（`migrate` コマンド）

### 1. データベースを起動する

```bash
make db-up
```

`docker compose` で PostgreSQL 16 を立ち上げ、接続可能になるまで待ちます。

**ポート 5432 が別の PostgreSQL で埋まっている場合**は、`POSTGRES_PORT` で変更できます。
`docker-compose.yml` と `Makefile` の `DATABASE_URL` の両方がこの値を参照します。

```bash
POSTGRES_PORT=5433 make db-up
POSTGRES_PORT=5433 make migrate-up
```

### 2. マイグレーションを適用する

```bash
make migrate-up
```

`users` / `treatments` / `share_links` / `sessions` の 4 テーブルが作られます。
スキーマの詳細は [`docs/DATABASE.md`](docs/DATABASE.md) を参照してください。

### 3. バックエンドを起動する

`backend/.env.example` を参考に環境変数を設定してから起動します。

```bash
APP_ENV=development \
DATABASE_URL='postgres://hairhistory:hairhistory@localhost:5432/hairhistory?sslmode=disable' \
make run
```

`http://localhost:8080/healthz` が `{"status":"ok"}` を返せば起動しています。

`APP_ENV=development` のときだけ `POST /api/auth/dev-login` が有効になり、
Google のクライアント ID なしでログインの動作確認ができます（本番では 404 になります）。

### 4. フロントエンドを起動する

```bash
cd frontend
npm install
npm run dev
```

`http://localhost:5173` で開きます。API のベース URL は `VITE_API_BASE_URL` で切り替えられ、
未設定なら `http://localhost:8080` を使います。

## 環境変数

### バックエンド（`backend/.env.example`）

| 変数 | 必須 | 既定値 | 説明 |
| --- | --- | --- | --- |
| `DATABASE_URL` | ○ | — | PostgreSQL の接続 URL。未設定だと起動時にエラーで停止する |
| `APP_ENV` | — | `development` | `development` または `production`。それ以外の値はエラー。`production` のとき Cookie に `Secure` が付き、開発用ログインが無効になる |
| `PORT` | — | `8080` | API サーバーの待ち受けポート |
| `GOOGLE_CLIENT_ID` | △ | 空 | Google ID Token の検証に使う。`APP_ENV=production` では必須 |
| `CORS_ALLOWED_ORIGINS` | — | `http://localhost:5173` | 許可するオリジンをカンマ区切りで指定。Cookie を送るため `*` は指定できない |
| `SESSION_TTL_HOURS` | — | `720` | セッションの有効時間（時間単位）。正の整数のみ |

### フロントエンド（`frontend/.env.example`）

| 変数 | 既定値 | 説明 |
| --- | --- | --- |
| `VITE_API_BASE_URL` | `http://localhost:8080` | API のベース URL |
| `VITE_GOOGLE_CLIENT_ID` | 空 | Google ログインボタンに渡すクライアント ID。空のときは Google ボタンを出さず、開発用ログインのフォームだけを表示する |

Vite の環境変数はビルド時に埋め込まれるため、値を変えたら再ビルドが必要です。
本番では `deploy/05-enable-production.sh` が `api.env` を書き換えたうえで `03-deploy-app.sh` を
呼び出し、この値を埋め込んで再ビルドします。本番の `VITE_API_BASE_URL` は空にしてあり、
SPA は自分と同じオリジンの `/api` を叩きます。同一オリジンなので CORS を経由せずに
セッション Cookie が送られます。

## API

詳細なリクエスト / レスポンス例、バリデーションルール、エラーコードは
[`docs/API.md`](docs/API.md) にまとめています。

| メソッド | パス | 認証 | 内容 |
| --- | --- | --- | --- |
| `GET` | `/healthz` | — | ヘルスチェック |
| `POST` | `/api/auth/google` | — | Google ID Token でログイン。`hh_session` Cookie を発行 |
| `POST` | `/api/auth/dev-login` | — | 開発用ログイン（`APP_ENV=development` のときのみ） |
| `GET` | `/api/auth/me` | ○ | ログイン中のユーザー情報 |
| `POST` | `/api/auth/logout` | — | セッションを削除し Cookie を失効させる |
| `GET` | `/api/treatments` | ○ | 施術履歴の一覧（施術日の降順） |
| `POST` | `/api/treatments` | ○ | 施術履歴の登録 |
| `GET` | `/api/treatments/{id}` | ○ | 施術履歴の取得 |
| `PUT` | `/api/treatments/{id}` | ○ | 施術履歴の更新（全項目置換） |
| `DELETE` | `/api/treatments/{id}` | ○ | 施術履歴の削除 |
| `POST` | `/api/shares` | ○ | 共有リンクの発行（`expiresInHours` は 1〜8760、既定 168） |
| `GET` | `/api/shares` | ○ | 有効な共有リンクの一覧 |
| `DELETE` | `/api/shares/{token}` | ○ | 共有リンクの無効化 |
| `GET` | `/api/public/shares/{token}` | — | 共有された施術履歴の公開参照。個人情報は返さない |

JSON はすべて camelCase、エラーは `{"error": {"code", "message", "details"}}` の形に統一しています。

## テスト

```bash
make test        # backend で go test ./...
```

Go のテストは 4 ファイル・18 テスト関数・約 1,100 行です。テーブルドリブンで書いており、
repository を interface として注入しているため DB を起動せずに実行できます。

| ファイル | 行数 | 対象 |
| --- | --- | --- |
| `internal/auth/usecase_test.go` | 321 | Google ログイン、開発用ログイン、セッション検証、ログアウト時のハッシュ照合、トークンの推測困難性 |
| `internal/share/usecase_test.go` | 360 | 有効期限の計算と境界、トークンの重複、公開参照、無効化、一覧の時刻依存 |
| `internal/treatment/usecase_test.go` | 330 | バリデーション、空文字の正規化、所有者チェック（他人のレコードを操作できないこと） |
| `internal/httpx/middleware_test.go` | 105 | 共有トークンがアクセスログ・パニックログに残らないこと |

**テストが無い箇所**（正直に書いておきます）

- handler 層（JSON のデコードと HTTP ステータスへの変換）のテストは書いていません
- 実際の PostgreSQL をつないだ統合テスト・repository のテストもありません。SQL の確認は手動と
  [`docs/API.md`](docs/API.md) 末尾の curl 手順で行っています
- フロントエンドのテストは 0 件です（型チェックと lint のみ）

その他のコマンド:

```bash
make lint        # gofmt と go vet
make build       # go build ./...
make help        # 全ターゲットの一覧
```

フロントエンドは `npm run build`（`tsc -b` を含む）と `npm run lint`（oxlint）で確認します。

## デプロイ

AWS EC2（Ubuntu）1 台の中に、Go API・PostgreSQL 16・Nginx をすべて配置する構成です。
<https://hairhistory.duckdns.org> がこの構成で動いています。

```
ブラウザ ──443(TLS)──▶ Nginx ──┬── /      → /var/www/hairhistory（React のビルド成果物）
         ──80──▶ 443 へリダイレクト └── /api/  → 127.0.0.1:8080（Go API / systemd）
                                            │
                                            └── 127.0.0.1:5432 PostgreSQL 16
```

PostgreSQL と API のポートは外部に公開せず（ufw で閉じ、PostgreSQL は localhost のみ待ち受け）、
必ず Nginx を経由させます。ドメインは DuckDNS の無料サブドメイン、証明書は certbot（Let's Encrypt）です。

スクリプトは番号順に実行します。すべて再実行しても壊れないように書いてあります。

| スクリプト | 役割 |
| --- | --- |
| `deploy/01-setup-server.sh` | PostgreSQL 16・Go・Nginx・git の導入と ufw の設定（サーバーごとに 1 回） |
| `deploy/02-setup-db.sh` | ロールとデータベースの作成。パスワードは環境変数で渡しスクリプトに残さない |
| `deploy/03-deploy-app.sh` | git fetch → Go のビルド → マイグレーション適用 → フロントのビルド → 配置 → systemd 再起動。2 回目以降のデプロイもこれ |
| `deploy/04-setup-https.sh` | certbot で証明書を取得し、Nginx を 443 待ち受け＋80 からのリダイレクトに書き換える |
| `deploy/05-enable-production.sh` | `api.env` を `APP_ENV=production` に書き換えて Google ログインを有効化。HTTPS が成立していないと何も変更せず停止する |

サーバーの初期セットアップからデプロイ、ロールバック、トラブルシュートまでの手順は
[`deploy/README.md`](deploy/README.md) にまとめています。

## 今後の課題・未実装

動いていないもの、手が回っていないものを正直に書きます。

**開発プロセス**

- **CI 未整備**。`.github/workflows/` は空のままで、テスト・lint・ビルドはすべて手元で実行しています。
  まず `make test` と `npm run build` を PR で回すところから入れたいところです
- **テストの範囲が狭い**。usecase 層と一部のミドルウェアだけで、handler 層・repository・
  フロントエンドのテストがありません（詳細は「テスト」の節）

**インフラ・運用**

- **IaC 未導入**。EC2・セキュリティグループ・DuckDNS はすべて手作業で、Terraform 等で
  コード化していません。サーバーを作り直す手順は `deploy/README.md` の文章に依存しています
- **DB バックアップ未自動化**。`pg_dump` を定期実行する仕組みが無く、今はインスタンスが
  壊れるとデータが戻りません
- **監視・ログ集約なし**。API は標準出力に JSON でログを吐き、systemd（journald）が受けるだけです。
  メトリクスもアラートもなく、落ちたことに気づく手段は自分でアクセスすることだけです
- **サーバー 1 台構成で冗長性なし**。API・PostgreSQL・Nginx が同じ t3.micro に載っているため、
  このインスタンスが落ちるとサービス全体が止まります

**アプリケーション**

- **ページネーション・検索なし**。`GET /api/treatments` は該当ユーザーの全件を `LIMIT` なしで
  返します。件数が増えたときに破綻します。施術項目や期間での絞り込みもありません
- **トランザクション未使用**。すべての SQL を単発で実行しており、複数テーブルにまたがる更新を
  原子的に行う箇所はまだありません（現状はそれで整合性が崩れる操作が無いためですが、
  機能追加時には必要になります）
- **バンドルサイズ**。コード分割をしておらず、JS が単一チャンク（約 630 KB、gzip 前）に
  なっています。`React.lazy` でのルート単位の分割が未対応です
- **共有リンクの有効期限が固定**。API は 1〜8760 時間を受け付けますが、
  フロントからは 168 時間（7 日）を決め打ちで送っています

今後やりたいこと: 写真の添付、施術履歴の検索、前回施術からの経過日数の通知。

## ライセンス

MIT License
