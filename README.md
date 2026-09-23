# HairHistory

美容院での施術内容を記録し、次に行くお店に共有リンクや QR コードで見せられる Web アプリケーションです。
Go 製の API サーバーと React の SPA を分離して実装し、AWS EC2 上に自前で構築したサーバーで動かしています。

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
| 施術履歴の記録 | 施術日・施術項目・サロン名・メモ・金額を登録 / 一覧 / 編集 / 削除 |
| 共有リンクの発行 | 有効期限付きのトークンを発行。一覧表示と任意のタイミングでの無効化が可能 |
| 公開閲覧 | 発行されたリンクはログイン不要で閲覧できる。個人情報は返さない |
| QR コード生成 | 共有リンクの QR コードをブラウザ側で生成し、その場で見せられる |
| Google ログイン | Google ID Token を検証してセッション Cookie を発行（下記「未実装」も参照） |

## 技術スタック

| 層 | 採用技術 |
| --- | --- |
| フロントエンド | Vite + React 19 + TypeScript（SPA）、Chakra UI、React Router、axios、qrcode.react |
| バックエンド | Go 1.25 + chi、pgx/v5、golang-migrate、`google.golang.org/api/idtoken` |
| データベース | PostgreSQL 16（ORM を使わず生 SQL） |
| インフラ | AWS EC2（Ubuntu）1 台構成、Nginx リバースプロキシ、systemd |
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

セッションは `hh_session` という HttpOnly Cookie で保持し、DB には生のトークンではなく
SHA-256 ハッシュ（`sessions.token_hash`）を保存しています。

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
│   ├── 01-setup-server.sh
│   ├── 02-setup-db.sh
│   ├── 03-deploy-app.sh
│   ├── hairhistory-api.service
│   ├── nginx-hairhistory.conf
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
| `/dashboard` | ダッシュボード | 必要 |
| `/treatments` | 施術履歴の一覧 | 必要 |
| `/treatment/:id` | 施術履歴の詳細・編集・削除 | 必要 |
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
| `VITE_GOOGLE_CLIENT_ID` | 空 | Google ログインボタンに渡すクライアント ID。空だとボタンは無効のまま |

Vite の環境変数はビルド時に埋め込まれるため、値を変えたら再ビルドが必要です。

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

`auth` / `treatment` / `share` の usecase 層に、テーブルドリブンのユニットテストがあります。
repository を interface として注入しているため、DB を起動せずに実行できます。

その他のコマンド:

```bash
make lint        # gofmt と go vet
make build       # go build ./...
make help        # 全ターゲットの一覧
```

フロントエンドは `npm run build`（`tsc -b` を含む）と `npm run lint`（oxlint）で確認します。

## デプロイ

AWS EC2（Ubuntu）1 台の中に、Go API・PostgreSQL 16・Nginx をすべて配置する構成です。

```
ブラウザ ──80──▶ Nginx ──┬── /      → /var/www/hairhistory（React のビルド成果物）
                          └── /api/  → 127.0.0.1:8080（Go API / systemd）
                                            │
                                            └── 127.0.0.1:5432 PostgreSQL 16
```

PostgreSQL と API のポートは外部に公開せず、必ず Nginx を経由させます。
サーバーの初期セットアップからデプロイ、ロールバック、トラブルシュートまでの手順は
[`deploy/README.md`](deploy/README.md) にまとめています。

## 今後の課題・未実装

正直なところ、以下はまだ動いていません。

- **本番の Google ログインが未接続**。Google OAuth のクライアント ID をまだ取得していないため、
  ログインボタンは無効状態です。動作確認は開発用ログイン（`POST /api/auth/dev-login`）で代替しています
- **HTTPS 未対応**。現状は HTTP のみで、セッション Cookie が平文で流れます。
  ドメインを用意して `certbot` で TLS を有効化するまでは本番利用できません
- **施術項目の複数選択 UI が未実装**。API と DB（`text[]`）は複数項目に対応済みですが、
  画面からは 1 項目しか送っていません
- **CI 未整備**。`.github/workflows/` が空で、テストとビルドは手元で実行しています
- **共有リンクの有効期限が固定**。API は 1〜8760 時間を受け付けますが、
  フロントからは 168 時間（7 日）を決め打ちで送っています
- **バックアップ未自動化・単一障害点**。1 台構成のため、このインスタンスが落ちるとサービス全体が止まります

今後やりたいこと: 写真の添付、施術履歴の検索、前回施術からの経過日数の通知。

## ライセンス

MIT License
