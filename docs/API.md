# API Documentation

HairHistory バックエンド API のリファレンス。

- ベース URL（ローカル）: `http://localhost:8080`
- すべてのパスは `/api` 配下（ヘルスチェックのみ例外）
- リクエスト / レスポンスの JSON はすべて **camelCase**
- 認証はセッション Cookie `hh_session`（HttpOnly / SameSite=Lax / Path=/）
- CORS は `CORS_ALLOWED_ORIGINS` に列挙したオリジンのみ許可。Cookie を使うため `AllowCredentials: true`、ワイルドカード `*` は設定不可

---

## 共通仕様

### エラーレスポンス

すべてのエラーは同じ形。内部エラーの詳細はレスポンスに含めない（サーバーログにのみ出力）。

```json
{
  "error": {
    "code": "invalid_argument",
    "message": "request contains invalid fields",
    "details": { "treatedOn": "must be formatted as YYYY-MM-DD" }
  }
}
```

`details` はバリデーションエラー時のみ付く。

| code | HTTP | 意味 |
| --- | --- | --- |
| `invalid_argument` | 400 | リクエストの内容が不正 |
| `unauthenticated` | 401 | 未ログイン / セッション無効・期限切れ |
| `permission_denied` | 403 | 他人のリソースへのアクセス |
| `not_found` | 404 | 対象が存在しない |
| `internal` | 500 | サーバー内部エラー（詳細は返さない） |

### リクエスト ID

すべてのレスポンスに `X-Request-Id` が付く。クライアントが同名ヘッダを送ればその値を引き継ぐ。

---

## ヘルスチェック

### `GET /healthz`

認証不要。

```json
{ "status": "ok" }
```

---

## 認証

### `POST /api/auth/google`

Google の ID Token を検証してログインする。ユーザーが未登録なら作成し、登録済みなら email / name を更新する（`google_sub` が一意キー）。成功時に `hh_session` Cookie をセットする。

リクエスト:

```json
{ "idToken": "eyJhbGciOi..." }
```

レスポンス `200 OK`:

```json
{
  "id": "01a0c4f6-0244-79bc-b976-6506066df7da",
  "email": "anna@example.com",
  "name": "Anna",
  "createdAt": "2026-09-21T17:14:11Z"
}
```

| ステータス | 条件 |
| --- | --- |
| 400 | `idToken` が空 |
| 401 | ID Token の検証に失敗 / `sub` が無い |

### `POST /api/auth/dev-login`（開発専用）

`APP_ENV=development` のときだけルーティングされる。`APP_ENV=production` では 404 を返す。Google のクライアント ID なしで動作確認するための入口。

リクエスト:

```json
{ "email": "anna@example.com", "name": "Anna" }
```

`name` を省略すると email が名前になる。ユーザーの `google_sub` は `dev|<email>` として保存されるため、本物の Google アカウントと衝突しない。

レスポンスは `POST /api/auth/google` と同じ `200 OK`。`email` が空なら 400。

### `GET /api/auth/me`

現在のログインユーザーを返す。要認証。

レスポンス `200 OK` は上記と同じ形。未ログインは 401。

### `POST /api/auth/logout`

セッションを DB から削除し、Cookie を失効させる。Cookie が無くても `204 No Content` を返す。

---

## 施術履歴

すべて要認証。他人のレコードに触れようとすると 403 `permission_denied`。

### `Treatment` オブジェクト

```json
{
  "id": "01a0c4f6-0287-7b9c-9e7b-c50a14e2a280",
  "userId": "01a0c4f6-0244-79bc-b976-6506066df7da",
  "treatedOn": "2026-09-10",
  "services": ["カラー"],
  "salonName": "Hair Salon ABC",
  "memo": "アッシュグレー",
  "cost": 8000,
  "createdAt": "2026-09-21T17:14:11Z",
  "updatedAt": "2026-09-21T17:14:11Z"
}
```

`salonName` / `memo` / `cost` は未設定なら `null`。

### バリデーション

| フィールド | 必須 | ルール |
| --- | --- | --- |
| `treatedOn` | ○ | `YYYY-MM-DD`。未来日も許可（予定の記録に使うため） |
| `services` | ○ | 1〜10 要素。各要素は空文字不可（前後の空白は除去される） |
| `salonName` | — | 最大 100 文字。空文字・空白のみは `null` として保存 |
| `memo` | — | 最大 1000 文字。空文字・空白のみは `null` として保存 |
| `cost` | — | 0 以上の整数 |

違反時は 400 で、`details` に各フィールドの理由が入る。

### `GET /api/treatments`

自分の履歴を `treatedOn` 降順（同日は `createdAt` 降順）で返す。

```json
{ "treatments": [ /* Treatment */ ] }
```

0 件でも `{"treatments": []}`（`null` ではない）。

### `GET /api/treatments/{id}`

`Treatment` を返す。`id` が UUID でなければ 400、存在しなければ 404、他人のものなら 403。

### `POST /api/treatments`

リクエスト:

```json
{
  "treatedOn": "2026-09-10",
  "services": ["カラー"],
  "salonName": "Hair Salon ABC",
  "memo": "アッシュグレー",
  "cost": 8000
}
```

`201 Created` で作成された `Treatment` を返す。

### `PUT /api/treatments/{id}`

全フィールド置換（部分更新ではない）。ボディは `POST` と同じ。`200 OK` で更新後の `Treatment` を返す。

### `DELETE /api/treatments/{id}`

`204 No Content`。

---

## 共有リンク

共有の単位は**ユーザー単位（履歴まるごと）**。トークンは `crypto/rand` の 32 バイトを base64url エンコードしたもの（43 文字）。

QR コードはフロント側（`qrcode.react`）で生成するため、サーバーに QR 用のエンドポイントは無い。

### `POST /api/shares`

要認証。

リクエスト（省略可）:

```json
{ "expiresInHours": 168 }
```

`expiresInHours` を省略すると 168（7 日）。範囲は 1〜8760（1 年）で、外れると 400。

レスポンス `201 Created`:

```json
{
  "id": "01a0c4f6-0366-7134-97e1-1635a42b5d2d",
  "token": "kY5Pk8t4vCQzVoYVHwdoNMmaAIddpoDwYUGtMWCAUmo",
  "expiresAt": "2026-09-28T17:14:12Z",
  "createdAt": "2026-09-21T17:14:12Z"
}
```

公開 URL はフロントが `token` から組み立てる（例: `https://<frontend>/share/<token>`）。

### `GET /api/shares`

要認証。自分の**有効な**リンク（未失効かつ未期限切れ）を `createdAt` 降順で返す。

```json
{ "shares": [ /* 上記と同じ形 */ ] }
```

### `DELETE /api/shares/{token}`

要認証。`revoked_at` をセットして無効化する。`204 No Content`。

自分のものでない、存在しない、すでに失効済みのいずれも 404（他人のリンクの存在を漏らさないため）。

### `GET /api/public/shares/{token}`

**認証不要。** 有効期限内かつ未失効のときだけ、そのユーザーの施術履歴を返す。

**ユーザーの氏名・メール・ID は一切返さない。** 施術レコードの `id` / `userId` も含まない。

```json
{
  "treatments": [
    {
      "treatedOn": "2026-09-11",
      "services": ["カラー", "カット"],
      "salonName": "Hair Salon ABC",
      "memo": "更新済み",
      "cost": 9500
    }
  ]
}
```

期限切れ・失効済み・存在しないトークンはすべて 404 `not_found`（区別しない）。

---

## curl での動作確認手順

```bash
# 1. DB 起動 + マイグレーション
make db-up
make migrate-up

# 2. サーバー起動（別ターミナル）
APP_ENV=development \
DATABASE_URL='postgres://hairhistory:hairhistory@localhost:5432/hairhistory?sslmode=disable' \
make run

# 3. 開発用ログイン（Cookie を jar に保存）
curl -s -c /tmp/hh.jar -X POST http://localhost:8080/api/auth/dev-login \
  -H 'Content-Type: application/json' \
  -d '{"email":"anna@example.com","name":"Anna"}'

# 4. ログイン中のユーザー
curl -s -b /tmp/hh.jar http://localhost:8080/api/auth/me

# 5. 施術を登録
curl -s -b /tmp/hh.jar -X POST http://localhost:8080/api/treatments \
  -H 'Content-Type: application/json' \
  -d '{"treatedOn":"2026-09-10","services":["カラー"],"salonName":"Hair Salon ABC","memo":"アッシュグレー","cost":8000}'

# 6. 一覧
curl -s -b /tmp/hh.jar http://localhost:8080/api/treatments

# 7. 共有リンクを発行 → 認証なしで閲覧
TOKEN=$(curl -s -b /tmp/hh.jar -X POST http://localhost:8080/api/shares \
  -H 'Content-Type: application/json' -d '{"expiresInHours":168}' \
  | python3 -c 'import json,sys;print(json.load(sys.stdin)["token"])')
curl -s http://localhost:8080/api/public/shares/$TOKEN

# 8. ログアウト
curl -s -b /tmp/hh.jar -c /tmp/hh.jar -X POST http://localhost:8080/api/auth/logout
```

---

## セキュリティ上の約束

- セッショントークンは平文で保存しない。SHA-256 ハッシュを `sessions.token_hash` に保存し、照合もハッシュで行う
- 共有トークンは 32 バイトの暗号論的乱数
- ID Token / セッショントークン / メールアドレスをログに出力しない
- SQL は必ずプレースホルダ（`$1`）でバインドする。文字列連結でクエリを組み立てない
- 内部エラーの詳細はレスポンスに載せず、サーバーログにのみ出力する
- Cookie の `Secure` 属性は `APP_ENV=production` のときのみ有効
