# Database Schema

- RDBMS: **PostgreSQL 16**（ローカルは `docker-compose.yml` の `db` サービス）
- ドライバ: `github.com/jackc/pgx/v5`（pgxpool）。ORM は使わず生 SQL
- マイグレーション: `golang-migrate`（`backend/migrations/`）
- 主キーは **UUID v7**（`github.com/google/uuid` の `NewV7`）。時系列順に並ぶため、生成順＝挿入順になりインデックスが断片化しにくい
- 命名は **snake_case**（API の JSON は camelCase なので、リポジトリ層で変換する）
- 日付だけを持つ `treated_on` は `date`、それ以外の時刻は `timestamptz`

---

## ER 図

```mermaid
erDiagram
    users ||--o{ treatments  : "has"
    users ||--o{ share_links : "issues"
    users ||--o{ sessions    : "owns"

    users {
        uuid        id PK
        text        google_sub UK
        text        email
        text        name
        timestamptz created_at
        timestamptz updated_at
    }
    treatments {
        uuid        id PK
        uuid        user_id FK
        date        treated_on
        text_array  services
        text        salon_name
        text        memo
        integer     cost
        timestamptz created_at
        timestamptz updated_at
    }
    share_links {
        uuid        id PK
        uuid        user_id FK
        text        token UK
        timestamptz expires_at
        timestamptz revoked_at
        timestamptz created_at
    }
    sessions {
        uuid        id PK
        uuid        user_id FK
        text        token_hash UK
        timestamptz expires_at
        timestamptz created_at
    }
```

すべての子テーブルの外部キーは `ON DELETE CASCADE`。ユーザーを削除すれば履歴・共有リンク・セッションも消える。

---

## `users`

Google アカウントでログインしたユーザー。

| カラム | 型 | 制約 | 説明 |
| --- | --- | --- | --- |
| `id` | `uuid` | PK | UUID v7 |
| `google_sub` | `text` | NOT NULL, UNIQUE | Google ID Token の `sub`。ログインの同一性判定に使う一意キー |
| `email` | `text` | NOT NULL | ログインのたびに更新される |
| `name` | `text` | NOT NULL | 同上。Google 側に名前が無ければ email が入る |
| `created_at` | `timestamptz` | NOT NULL DEFAULT now() | |
| `updated_at` | `timestamptz` | NOT NULL DEFAULT now() | upsert 時に `now()` で更新 |

ログインは `google_sub` を衝突キーにした upsert で行う。

開発用の `POST /api/auth/dev-login` で作られるユーザーは `google_sub` が `dev|<email>` になるため、本物の Google アカウント（数字列の `sub`）と衝突しない。

---

## `treatments`

施術履歴。アプリの中心テーブル。

| カラム | 型 | 制約 | 説明 |
| --- | --- | --- | --- |
| `id` | `uuid` | PK | UUID v7 |
| `user_id` | `uuid` | NOT NULL, FK → `users(id)` ON DELETE CASCADE | 所有者 |
| `treated_on` | `date` | NOT NULL | 施術日。時刻を持たない。未来日も許可（予定の記録に使うため） |
| `services` | `text[]` | NOT NULL, `cardinality(services) > 0` | 施術内容。「カット＋カラー」のような複数施術に対応するため配列 |
| `salon_name` | `text` | NULL 可 | 美容院名（最大 100 文字はアプリ側で検証） |
| `memo` | `text` | NULL 可 | メモ（最大 1000 文字はアプリ側で検証） |
| `cost` | `integer` | NULL 可, `cost IS NULL OR cost >= 0` | 料金（円） |
| `created_at` | `timestamptz` | NOT NULL DEFAULT now() | |
| `updated_at` | `timestamptz` | NOT NULL DEFAULT now() | UPDATE 時に `now()` で更新 |

インデックス:

- `treatments_user_id_treated_on_idx (user_id, treated_on DESC)` — 一覧取得は必ず「自分の履歴を施術日降順」で引くため、この複合インデックスだけでソートまで賄える

`services` を配列にした理由: 1 回の来店で複数の施術を受けるのが普通で、1 レコードに複数の施術内容を持たせたいため。フロントは当面 1 要素の配列を送る。

---

## `share_links`

履歴を第三者に見せるための共有リンク。共有の単位は**ユーザー単位（履歴まるごと）**で、施術単位ではない。

| カラム | 型 | 制約 | 説明 |
| --- | --- | --- | --- |
| `id` | `uuid` | PK | UUID v7 |
| `user_id` | `uuid` | NOT NULL, FK → `users(id)` ON DELETE CASCADE | 発行者 |
| `token` | `text` | NOT NULL, UNIQUE | 32 バイトの暗号論的乱数を base64url エンコードした 43 文字 |
| `expires_at` | `timestamptz` | NOT NULL | 既定は発行から 168 時間（7 日） |
| `revoked_at` | `timestamptz` | NULL 可 | 手動で無効化した時刻。NULL なら未失効 |
| `created_at` | `timestamptz` | NOT NULL DEFAULT now() | |

インデックス: `share_links_user_id_idx (user_id)`

**有効判定**: `revoked_at IS NULL AND expires_at > now()`。境界は「期限ちょうどは無効」。

`token` はセッショントークンと違い**平文で保存する**。所有者が `GET /api/shares` で自分のリンクを再表示・再配布できる必要があるため（ハッシュだと元の値を復元できない）。漏洩時の影響は「その 1 人の施術履歴が読まれる」に限定され、氏名・メールは公開レスポンスに含まれない。

---

## `sessions`

Cookie ベースのログインセッション。

| カラム | 型 | 制約 | 説明 |
| --- | --- | --- | --- |
| `id` | `uuid` | PK | UUID v7 |
| `user_id` | `uuid` | NOT NULL, FK → `users(id)` ON DELETE CASCADE | |
| `token_hash` | `text` | NOT NULL, UNIQUE | セッショントークンの **SHA-256 ハッシュ（16 進 64 文字）** |
| `expires_at` | `timestamptz` | NOT NULL | 既定は `SESSION_TTL_HOURS`（720 時間 = 30 日） |
| `created_at` | `timestamptz` | NOT NULL DEFAULT now() | |

インデックス: `sessions_user_id_idx (user_id)`, `sessions_expires_at_idx (expires_at)`

**生のトークンは DB に保存しない。** Cookie に入れた値の SHA-256 だけを保存し、認証時も受け取った値をハッシュ化して照合する。DB が漏れてもセッションを乗っ取られない。

期限切れセッションはサーバー起動時に `DELETE FROM sessions WHERE expires_at <= now()` で掃除する。

---

## マイグレーション

ファイル: `backend/migrations/000001_init.up.sql` / `000001_init.down.sql`

```bash
make db-up        # Postgres コンテナ起動（pg_isready まで待つ）
make migrate-up   # すべて適用
make migrate-down # 直近 1 つをロールバック
make db-down      # コンテナ停止
```

接続文字列は Makefile の `DATABASE_URL` を上書きできる。

```
postgres://hairhistory:hairhistory@localhost:5432/hairhistory?sslmode=disable
```

ホストの 5432 が別の Postgres に使われている場合は `POSTGRES_PORT` を変える。

```bash
POSTGRES_PORT=5442 make db-up
POSTGRES_PORT=5442 make migrate-up
```

データは名前付きボリューム `hairhistory-db-data` に永続化される。完全に作り直すときは:

```bash
docker compose down -v
```

---

## クエリの方針

- SQL は必ずプレースホルダ（`$1`, `$2`, ...）でバインドする。**文字列連結でクエリを組み立てない**
- 所有者チェックはユースケース層で行う。加えて `share_links` の無効化は `WHERE token = $2 AND user_id = $1` と SQL 側でも所有者を絞る
- 公開エンドポイント（`GET /api/public/shares/{token}`）は `treatments` しか読まず、`users` には触れない
