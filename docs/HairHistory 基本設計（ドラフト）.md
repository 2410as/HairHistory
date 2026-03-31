# HairHistory 基本設計（ドラフト）

要件定義（`HairHistory 要件定義.md`）をもとにした、画面・データ・API のざっくり設計。  
このファイルを「実装に近い詳細の一次情報」として扱う。

目的は **「髪の履歴をメモして、ダメージを美容師が把握しやすくする」** こと。

- --

## 1. 画面構成（MVP）

1. **トップ / 初回アクセス画面**

- 説明テキスト（このアプリの目的）

- 「はじめる」ボタン → ユーザーID（URL用ID）を自動発行し、自分用の履歴ページへ遷移

2. **マイ履歴一覧画面（/h/:userId）**

- ユーザーURLでアクセスするメイン画面

- 表示内容（例）:

- ヘッダー: プロジェクト名、簡単な説明（「美容師さんと一緒にダメージを確認できます」など）

- 「新しい履歴を追加」ボタン

- 履歴一覧（新しい順）

- 日付

- 施術内容（例: カラー / ブリーチ / 縮毛矯正）

- サロン名・スタイリスト名（簡易表示）

- ダメージが分かるような短いサマリー（例: 「ブリーチ2回目 / 毛先ダメージ強め」など）

3. **履歴入力・編集画面（モーダル or 別ページ）**

- 1件の履歴を登録・編集する画面

- 入力項目（要件から）:

- 日付

- 施術内容

- チェックボックス例: カラー / ブリーチ / 縮毛矯正 / トリートメント / その他（自由入力）

- サロン名（フリーテキスト）

- スタイリスト名（フリーテキスト）

- メモ

- ダメージに関するメモを書きやすいように、プレースホルダー例:

- 「例: 毛先がかなりダメージ / ブリーチ〇回目 / 薬剤がしみた など」

4. **共有ビュー（美容師さんが見る想定の画面）**

- 基本はマイ履歴一覧画面と同一URL（/h/:userId）

- スマホを見せる / QR を送る前提なので、

- 直近3〜5件の履歴を上にまとめて表示

- 各履歴に「施術内容」と「ダメージ関連メモ」が一目で分かるように

- --

## 2. データ設計（詳細）

### 2.1 テーブル: users（匿名ユーザー）

| カラム名 | 型 | NULL | 説明 |
|---------|----|------|------|

| id | TEXT | ❌ | ユーザーID（URL用ID、nanoid等） |

| name | TEXT | ✅ | 表示名（「〇〇さんの履歴」用） |

| email | TEXT | ✅ | メールアドレス（将来ログイン用） |

| password_hash | TEXT | ✅ | パスワードハッシュ（将来ログイン用） |

| last_login_at | TIMESTAMPTZ | ✅ | 最終ログイン日時 |

| is_deactivated | BOOLEAN | ❌ | 無効化フラグ（退会・BAN時に true、DEFAULT false） |

| created_at | TIMESTAMPTZ | ❌ | 作成日時 |

※ MVP では id / created_at のみ使用。name は任意入力。email 以降はログイン導入時に利用。

### 2.2 テーブル: hair_histories（髪の履歴）

PostgreSQL 想定の型で記載。

| カラム名 | 型 | NOT NULL | 説明 |
|---------|----|----------|------|

| id | UUID | ✅ | 履歴ID（サーバー側で生成） |

| user_id | UUID | ✅ | users.id への外部キー |

| date | DATE | ✅ | 施術日 |

| services | TEXT | ✅ | 施術内容。JSON 文字列（例: `["color","bleach"]`）として保存 |

| salon_name | TEXT | ❌ | サロン名（フリーテキスト） |

| stylist_name | TEXT | ❌ | スタイリスト名（フリーテキスト） |

| memo | TEXT | ❌ | メモ（ダメージ状況・薬剤・注意点など） |

| created_at | TIMESTAMPTZ | ✅ | 作成日時（サーバー側で設定） |

| updated_at | TIMESTAMPTZ | ✅ | 更新日時（サーバー側で設定） |

※ `services` は Go 側では `[]ServiceType` として扱い、DB には JSON 文字列で保存する想定。

- --

## 2.3 Go の型（domain / DTO）

### 2.3.1 施術内容の種類

```go

type ServiceType string

const (

ServiceTypeColor        ServiceType = "color"

ServiceTypeBleach       ServiceType = "bleach"

ServiceTypeStraightPerm ServiceType = "straight_perm" *// 縮毛矯正*

ServiceTypeTreatment    ServiceType = "treatment"

ServiceTypeOther        ServiceType = "other"

)

```

### 2.3.2 ドメインモデル

```go

type User struct {

ID           string     `json:"id"`

Name         string     `json:"name,omitempty"`

Email        string     `json:"email,omitempty"`

LastLoginAt   *time.Time `json:"lastLoginAt,omitempty"`

IsDeactivated bool       `json:"isDeactivated"`

CreatedAt    time.Time  `json:"createdAt"`

}

type HairHistory struct {

ID          string        `json:"id"`

UserID      string        `json:"userId"`

Date        time.Time     `json:"date"`

Services    []ServiceType `json:"services"`

SalonName   string        `json:"salonName"`

StylistName string        `json:"stylistName"`

Memo        string        `json:"memo"`

CreatedAt   time.Time     `json:"createdAt"`

UpdatedAt   time.Time     `json:"updatedAt"`

}

```

### 2.3.3 リクエスト / レスポンス用 DTO（概念）


```go

*// POST /api/users レスポンス*

type CreateUserResponse struct {

UserID string `json:"userId"`

}

*// GET /api/users/{userId}/histories レスポンス*

type ListHistoriesResponse struct {

Histories []HairHistory `json:"histories"`

}

*// POST /api/users/{userId}/histories リクエスト*

type CreateHistoryRequest struct {

Date        time.Time     `json:"date"`

Services    []ServiceType `json:"services"`

SalonName   string        `json:"salonName"`

StylistName string        `json:"stylistName"`

Memo        string        `json:"memo"`

}

*// POST /api/users/{userId}/histories レスポンス*

type CreateHistoryResponse struct {

History HairHistory `json:"history"`

}

*// PUT /api/histories/{historyId} リクエスト*

type UpdateHistoryRequest struct {

Date        *time.Time     `json:"date,omitempty"`

Services    *[]ServiceType `json:"services,omitempty"`

SalonName   *string        `json:"salonName,omitempty"`

StylistName *string        `json:"stylistName,omitempty"`

Memo        *string        `json:"memo,omitempty"`

}

*// PUT /api/histories/{historyId} レスポンス*

type UpdateHistoryResponse struct {

History HairHistory `json:"history"`

}

```

- --

## 3. API 設計（MVP 想定）

ベースURL: `/api`

0. **ヘルスチェック**

- `GET /api/health` … 監視・デプロイ確認用（例: `{ "ok": true }`）

1. **ユーザー作成**

- `POST /api/users`

- 内容: 新しい匿名ユーザーIDを発行

- レスポンス: 実装では `ent.id` 形式（§2.3.3 注記・§5）

2. **履歴一覧取得**

- `GET /api/users/{userId}/histories`

- 役割: 指定ユーザーの履歴一覧を取得（新しい順）

- レスポンス: 実装では `list` 配列（§5）

3. **履歴作成**

- `POST /api/users/{userId}/histories`

- ボディ: `CreateHistoryRequest` 相当（日付・施術・サロン名等）

- レスポンス: 実装では `ent`（§5）

4. **履歴更新**

- `PUT /api/histories/{historyId}`

- ボディ: `UpdateHistoryRequest` 相当

- レスポンス: 実装では `ent`（§5）

5. **履歴削除**

- `DELETE /api/histories/{historyId}`

- レスポンス: `{ "ok": true }`（実装済み想定）

- --

## 4. 「ダメージが分かりやすい」ためのUIポイントメモ

- 履歴一覧で、**直近のブリーチ・縮毛矯正の回数**がざっくり分かるようにしたい

- 例: 一覧の1行に「カラー / ブリーチ（通算3回目） / 縮毛矯正（前回から6ヶ月）」など

- メモ欄のプレースホルダーやラベルで、「ダメージ」「しみた」「切れ毛」などを書きやすくする

- 直近3〜5件をカード風に大きめ表示 → 美容師さんがパッと見て判断しやすい

- --

## 5. バックエンド実装構成（`apps/main`）

モノレポ内の Go API。機能ごとの縦割り、usecase の request/response、domain service によるドメイン操作の集約を前提とする。**リポジトリのポートは `domain` パッケージ直下のインターフェース**（`UserRepository` / `HairHistoryRepository`）。`domain/repository` のような別パッケージには切り出さない。

### 5.1 レイヤの流れ

```
HTTP リクエスト
  → controller（ハンドラ・ルーティング登録）
  → usecase（アプリケーションフロー・request 組み立て・response 組み立て）
  → domain/service（ドメイン操作・将来のルール集約）
  → domain の Repository インターフェース
  → infra（PostgreSQL 等の実装）
```

- **entity** … 永続化モデルと値オブジェクト（JSON タグは基本付けない。API 形は response が担当）
- **domain** … 上記 Repository インターフェース（ポート）
- **domain/service** … ユースケースが呼ぶドメインサービス（現状はリポジトリへの薄い委譲から始め、ルールが増えたらここに集約）
- **usecase/request** … `*http.Request` から入力を取り出し、検証可能な構造体へ
- **usecase/response** … JSON 用 DTO と `New…` コンストラクタ（例: `list` / `ent`）
- **controller/render** … JSON 書き出し・エラー応答

### 5.2 ルーティング（[chi](https://github.com/go-chi/chi) v5）

`controller` で `chi.NewRouter()` を組み立て、`/api` 配下にルートを登録する。

| メソッド | パス（`chi` パターン） | 処理 |
|---------|------------------------|------|
| GET | `/api/health` | ヘルスチェック |
| POST | `/api/users` | ユーザー作成 |
| GET | `/api/users/{userId}/histories` | 履歴一覧 |
| POST | `/api/users/{userId}/histories` | 履歴作成 |
| PUT | `/api/histories/{historyId}` | 履歴更新 |
| DELETE | `/api/histories/{historyId}` | 履歴削除 |

`{userId}` / `{historyId}` は `usecase/request` で `chi.URLParam` により取得する。

### 5.3 ファイル分割の方針

**原則: 役割・操作ごとに 1 ファイル。** 機能名のプレフィックスで揃える（例: `hair_history_list.go`）。

**`app/domain`**

| ファイル | 内容 |
|---------|------|
| `user_repository.go` | `UserRepository` |
| `hair_history_repository.go` | `HairHistoryRepository` |

**`app/domain/entity`**

モデル・パラメータ・列挙を**種類ごと**に分割（例: `user.go`, `hair_history.go`, `service_type.go`, `create_hair_history_params.go`）。`Created()` / `Updated()` などのライフサイクル用メソッドを載せる。

**`app/domain/service/*`**

`Service` インターフェースと `NewService` を `service.go` に置き、操作は `create.go` / `finder_*.go` など**ファイル単位**で分割する。

**`app/infra`**

| ファイル | 内容 |
|---------|------|
| `user_repository_pg.go` | `UserRepository` の PostgreSQL 実装（現状スタブ可） |
| `hair_history_repository_pg.go` | `HairHistoryRepository` の実装 |

**`app/controller`**

| 種別 | 例 |
|------|-----|
| ルーティング | `router.go`（chi 組み立て・`/api`）、`router_health.go`、`router_users.go`、`router_hair_history.go` |
| 依存の型 | `deps_types.go`（`Deps`） |
| ハンドラ | `health_type.go` + `health_get.go`、`users_type.go` + `users_create.go`、`hair_history_type.go` + `hair_history_list.go` など**操作ごと** |
| JSON | `render/json.go`, `render/error.go` |

**`app/usecase`**

| 種別 | 例 |
|------|-----|
| 集約型 | `user_type.go`, `hair_history_type.go`（コンストラクタのみ） |
| メソッド | `user_create.go`, `hair_history_list.go`, `hair_history_create.go`, … |

**`app/usecase/request`**

`hair_history_list.go`, `hair_history_create.go`, … , `user_create.go` のように**エンドポイント入力ごと**。

**`app/usecase/response`**

- 履歴 API: `history_ent.go`（内部用マッパ）+ `hair_history_list.go` 等**操作ごと**
- ユーザー: `user_ent.go`, `user_create.go`

**`app/utility`**

例: `service_type_marshal.go` / `service_type_unmarshal.go` のように**処理単位**。

**エントリ（`apps/main`）**

| ファイル | 内容 |
|---------|------|
| `main.go` | DB 接続 → repository → domain service → usecase → `controller.NewRouter` → `ListenAndServe`（グレースフル shutdown） |

### 5.4 新機能を足すとき（例: 別リソース）

1. `domain/entity` にモデルをファイル追加  
2. `domain` に `*Repository` インターフェースをファイル追加  
3. `domain/service/<name>` を追加  
4. `infra/*_repository_pg.go`  
5. `usecase/*` / `request` / `response`  
6. `controller/*` と `router_*` にルート登録  
7. `main.go` に repository / service / usecase の配線を追加  

フロント（例: Next.js）からは `/api/*` をこの Go API に向けて呼び出す想定。

- --

この設計はドラフトです。画面項目や API の細部は、実装しながら調整する。

