# HairHistory API (`apps/main`)

Go で実装した REST API。PostgreSQL に永続化し、HTTP ルーティングは [chi](https://github.com/go-chi/chi) v5 を使います。

## 進め方（初回セットアップ）

### 1. PostgreSQL を用意する

**A. Docker（手軽）**

```bash
cd apps/main
docker compose up -d
```

接続文字列（デフォルト・`DATABASE_URL` 未設定時も同じ）:

`postgres://postgres:postgres@127.0.0.1:5432/hairhistory?sslmode=disable`

**B. 既存の Postgres を使う**

`DATABASE_URL` を環境変数で渡す（例: `postgres://user:pass@host:5432/dbname?sslmode=disable`）。

### 2. スキーマを流し込む

マイグレーションツールは必須にしていない。**1 本の SQL** をそのまま実行する。

```bash
# DATABASE_URL を自分の環境に合わせる
export DATABASE_URL='postgres://postgres:postgres@127.0.0.1:5432/hairhistory?sslmode=disable'

psql "$DATABASE_URL" -f migrations/00001_init.sql
```

（`psql` が無い場合は GUI クライアントで `migrations/00001_init.sql` の内容を実行してもよい。）

### 3. API を起動する

```bash
cd apps/main
go run .
```

- ポート: 環境変数 `PORT` があればそれ（例: `8080` → `:8080`）、なければ `:8080`。
- `DATABASE_URL` が空なら、上記 Docker 用のデフォルト URL を使う（起動時にログに出る）。

**任意の環境変数**

| 変数 | 説明 |
|------|------|
| `HAIR_CORS_ORIGINS` | カンマ区切りのオリジン。未設定時は `http://localhost:3000` のみ許可（Next.js ローカル用）。その場合は起動時に標準ログへ警告が 1 行出る（本番では明示設定推奨）。値はあるが有効なオリジンが 1 つもない場合も同様にデフォルトへフォールバックし警告する。 |
| `HAIR_HEALTH_PING_TIMEOUT` | ヘルスチェックの DB `Ping` 期限（`time.ParseDuration` 形式、例: `2s`）。未設定・不正時は `2s`。 |

`GET /api/health` は DB プールに `Ping` する。Postgres に届かない場合は **503**（`{ "error": "database unavailable" }`）。

### 4. 動作確認の例

```bash
# ヘルス
curl -s localhost:8080/api/health

# ユーザー作成（レスポンスは { "ent": { "id": "..." } }）
UID=$(curl -s -X POST localhost:8080/api/users | jq -r '.ent.id')
echo "$UID"

# jq が無い場合は JSON を目で見て id をコピー

# 履歴作成（services は1件以上必須。salonName / stylistName は空文字可）
curl -s -X POST "localhost:8080/api/users/$UID/histories" \
  -H 'Content-Type: application/json' \
  -d '{"date":"2026-03-28T12:00:00Z","services":["color"],"salonName":"","stylistName":"","memo":"メモ"}'

# 一覧（各要素の id が履歴 UUID）
curl -s "localhost:8080/api/users/$UID/histories"

# 更新（履歴 id は一覧 JSON の list[0].id など）
HID=$(curl -s "localhost:8080/api/users/$UID/histories" | jq -r '.list[0].id')
curl -s -X PUT "localhost:8080/api/histories/$HID" \
  -H 'Content-Type: application/json' \
  -d '{"memo":"更新メモ"}'

# 削除
curl -s -X DELETE "localhost:8080/api/histories/$HID"
```

## レイヤ・ファイル構成

リポジトリルートの **`docs/HairHistory 基本設計（ドラフト）.md` §5** を参照。

## データモデル（DB）

- `users.id` … TEXT（UUID 文字列。匿名ユーザー発行時に生成）
- `hair_histories.id` … UUID（DB が `gen_random_uuid()`）
- `hair_histories.services` … JSONB（`["color","bleach"]` 形式）
